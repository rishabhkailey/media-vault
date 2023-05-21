import { Chacha20 } from "ts-chacha20";
import { parseRequestRangeHeader } from "@/utils/request";
import type { IRequestRange } from "@/utils/request";
// stream saver worker code + our custom code for decryption

declare let self: ServiceWorkerGlobalScope;
// const sw = self as ServiceWorkerGlobalScope & typeof globalThis;

self.addEventListener("install", () => {
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});

const map = new Map();

// This should be called once per download
// Each event has a dataChannel that the data will be piped through
self.onmessage = (event) => {
  // We send a heartbeat every x second to keep the
  // service worker alive if a transferable stream is not sent
  if (event.data === "ping") {
    return;
  }

  const data = event.data;
  const downloadUrl =
    data.url ||
    self.registration.scope +
      Math.random() +
      "/" +
      (typeof data === "string" ? data : data.filename);
  const port = event.ports[0];
  const metadata = new Array(3); // [stream, data, port]

  metadata[1] = data;
  metadata[2] = port;

  // Note to self:
  // old streamsaver v1.2.0 might still use `readableStream`...
  // but v2.0.0 will always transfer the stream through MessageChannel #94
  if (event.data.readableStream) {
    metadata[0] = event.data.readableStream;
  } else if (event.data.transferringReadable) {
    port.onmessage = (evt) => {
      port.onmessage = null;
      metadata[0] = evt.data.readableStream;
    };
  } else {
    metadata[0] = createStream(port);
  }

  map.set(downloadUrl, metadata);
  port.postMessage({ download: downloadUrl });
};

function createStream(port: MessagePort) {
  // ReadableStream is only supported by chrome 52
  return new ReadableStream({
    start(controller) {
      // When we receive data on the messageChannel, we write
      port.onmessage = ({ data }) => {
        if (data === "end") {
          return controller.close();
        }

        if (data === "abort") {
          controller.error("Aborted the download");
          return;
        }

        controller.enqueue(data);
      };
    },
    cancel(reason) {
      console.log("user aborted", reason);
      port.postMessage({ abort: true });
    },
  });
}

const newDecryptTransformer: (
  decryptor: Chacha20
) => TransformStream<Uint8Array, Uint8Array> = (decryptor) =>
  new TransformStream<Uint8Array, Uint8Array>({
    start() {},
    transform(chunk, controller) {
      if (!chunk) {
        console.log("undefined chunk");
      }
      // console.log("encrypted ", new TextDecoder().decode(chunk));
      const decryptedChunk = decryptChunk(chunk, decryptor);
      // console.log("decrypted ", new TextDecoder().decode(decryptedChunk));
      controller.enqueue(decryptedChunk);
    },
    flush() {},
  });

const decryptChunk: (input: Uint8Array, decryptor: Chacha20) => Uint8Array = (
  input,
  decryptor
) => {
  try {
    const decrypted = decryptor.decrypt(input);
    if (!decrypted?.length || decrypted.length !== input.length) {
      console.log(decrypted);
    }
    // console.log(new TextDecoder().decode(decrypted));
    return decrypted;
  } catch (err) {
    console.log(err);
  }
  return new Uint8Array(0);
};

async function internalFetch(req: Request) {
  const rangeHeader = req.headers.get("Range");
  let range: IRequestRange | undefined;
  if (rangeHeader !== null && rangeHeader.length !== 0) {
    range = parseRequestRangeHeader(rangeHeader);
  } else {
    // todo
  }
  const password = "01234567890123456789012345678901";
  const nonce = "012345678901";
  const useless = new TextEncoder().encode(
    "00000000000000000000000000000000000000000000000000000000000000000"
  );
  let i = 0;
  if (range?.start !== undefined) {
    i = range.start;
  }
  const counter = Math.floor(i / 64);
  const byteCounter = i % 64;

  const textEncoder = new TextEncoder();
  const decryptor = new Chacha20(
    textEncoder.encode(password),
    textEncoder.encode(nonce),
    counter
  );
  // set the internal byte counter
  if (byteCounter !== 0) {
    decryptor.decrypt(useless.slice(0, byteCounter));
  }
  // console.log(req, self.clients.matchAll());
  // await new Promise((r) => {
  //   setTimeout(r, 1000);
  // });
  try {
    const res = await fetch(req);
    // return;
    const encryptedStream = res.body;
    if (encryptedStream === null) {
      return new Response("res");
    }
    const decryptedStream = encryptedStream.pipeThrough(
      newDecryptTransformer(decryptor)
    );
    // console.log(rangeHeader, range);
    // console.log(req);
    // const blob = await new Response(decryptedStream).blob();
    // console.log(await blob.text());
    // console.log(res.headers.get("Content-Type"));
    return new Response(decryptedStream, {
      headers: res.headers,
      status: res.status,
      statusText: res.statusText,
    });
  } catch (err) {
    console.log(err);
    return new Response(undefined, {
      status: 500,
    });
  }
}

self.onfetch = (event) => {
  const url = event.request.url;

  // this only works for Firefox
  if (url.endsWith("/ping")) {
    return event.respondWith(new Response("pong"));
  }
  const urlObj = new URL(url);
  if (
    event?.request?.method === "GET" &&
    typeof event?.request?.url === "string" &&
    (urlObj.pathname.startsWith("/v1/media/") ||
      urlObj.pathname.startsWith("/v1/thumbnail/"))
    // new URL(event.request.url).pathname.startsWith("/v1/testGetVideoWithRange")
  ) {
    return event.respondWith(internalFetch(event.request));

    // fetch(event.request)
    //   .then((res) => {
    //     const encryptedStream = res.body;
    //     if (encryptedStream === null) {
    //       return event.respondWith(res);
    //     }
    //     const decryptedStream = encryptedStream.pipeThrough(
    //       newDecryptTransformer()
    //     );
    //     return event.respondWith(
    //       new Response(decryptedStream, {
    //         headers: res.headers,
    //         status: res.status,
    //         statusText: res.statusText,
    //       })
    //     );
    //   })
    //   .catch((err) => {
    //     console.log(err);
    //     return event.respondWith(
    //       new Response("decryption failed", {
    //         status: 500,
    //       })
    //     );
    //   });
    // console.log(rangeHeader, range);
    // console.log(event?.request);
  }
  const hijacke = map.get(url);

  if (!hijacke) {
    return null;
  }

  const [stream, data, port] = hijacke;

  map.delete(url);

  // Not comfortable letting any user control all headers
  // so we only copy over the length & disposition
  const responseHeaders = new Headers({
    "Content-Type": "application/octet-stream; charset=utf-8",

    // To be on the safe side, The link can be opened in a iframe.
    // but octet-stream should stop it.
    "Content-Security-Policy": "default-src 'none'",
    "X-Content-Security-Policy": "default-src 'none'",
    "X-WebKit-CSP": "default-src 'none'",
    "X-XSS-Protection": "1; mode=block",
  });

  const headers = new Headers(data.headers || {});
  {
    const contentLength = headers.get("Content-Length");
    if (contentLength != null) {
      responseHeaders.set("Content-Length", contentLength);
    }
  }

  {
    const contentDisposition = headers.get("Content-Disposition");
    if (contentDisposition != null) {
      responseHeaders.set("Content-Disposition", contentDisposition);
    }
  }

  // data, data.filename and size should not be used anymore
  if (data.size) {
    console.warn("Depricated");
    responseHeaders.set("Content-Length", data.size);
  }

  let fileName = typeof data === "string" ? data : data.filename;
  if (fileName) {
    console.warn("Depricated");
    // Make filename RFC5987 compatible
    fileName = encodeURIComponent(fileName)
      .replace(/['()]/g, escape)
      .replace(/\*/g, "%2A");
    responseHeaders.set(
      "Content-Disposition",
      "attachment; filename*=UTF-8''" + fileName
    );
  }

  event.respondWith(new Response(stream, { headers: responseHeaders }));

  port.postMessage({ debug: "Download started" });
};

// // export const onmessage = function (event) {
// //   console.log("message received", event)
// // };
// import sodium from "libsodium-wrappers";

// self.addEventListener("message", (message) => {
//   console.log("message received", message);
//   // postMessage("sending response" + message.data);
// });

// self.addEventListener("install", function (event) {
//   console.log("install");
// });

// self.addEventListener("activate", function (event) {
//   console.log("Claiming control");
//   return self.clients.claim();
// });

// let password = "berrysikure";
// const decrypt: (input: Uint8Array) => Uint8Array = (input) => {
//   if (password.length == 0) {
//     throw new Error("password length zero");
//   }
//   while (password.length < 32) {
//     password += password;
//   }
//   password = password.substring(0, 32);
//   const encoder = new TextEncoder();
//   const key = encoder.encode(password);
//   console.log(sodium);
//   const res = sodium.crypto_secretstream_xchacha20poly1305_init_push(key);
//   const state_in = sodium.crypto_secretstream_xchacha20poly1305_init_pull(
//     res.header,
//     key
//   );
//   const output = sodium.crypto_secretstream_xchacha20poly1305_pull(
//     state_in,
//     input
//   );
//   return output.message;
// };
// self.addEventListener("fetch", (event: any) => {
//   if (
//     event?.request?.method !== "GET" ||
//     (typeof event?.request?.url === "string" &&
//       !new URL(event.request.url).pathname.startsWith("/v1/testGetEncryptedVideoWithRange"))
//   ) {
//     return;
//   }
//   console.log("fetch event", event);
//   event.respondWith(
//     fetch(event.request).then(async function (response) {
//       while (password.length < 32) {
//         password += password;
//       }
//       password = password.substring(0, 32);
//       // console.log(response);
//       // const reader = response.body?.getReader();
//       // if (!reader) {
//       //   console.log("unable to get response reader");
//       //   return new Response(new Blob([]), {
//       //     headers: response.headers,
//       //     status: 500,
//       //   });
//       // }
//       try {
//         const encryptedBytes = new Uint8Array(await response.arrayBuffer());
//         await sodium.ready;
//         const decryptedBytes = decrypt(encryptedBytes);
//         return new Response(new Blob([decryptedBytes]), {
//           headers: response.headers,
//           status: response.status,
//         });
//       } catch (err) {
//         console.log("error decyrpting the response", err);
//         return new Response(new Blob([]), {
//           headers: response.headers,
//           status: 500,
//         });
//       }

//       // new Response()
//       // return response;
//       // if(response.url.match(".mp4")){
//       //   console.log(event);
//       //   responseCloned = response.clone();
//       //   responseCloned.arrayBuffer().then(
//       //     buffer =>{
//       //       let length = 100
//       //       view = new Uint8Array(buffer,0,length)
//       //       for(i=0,j=length - 1; i<j; i++,j--){
//       //           view[i] = view[i] ^ view[j]
//       //           view[j] = view[j] ^ view[i]
//       //           view[i] = view[i] ^ view[j]
//       //       }
//       //     }
//       //   )
//       // }
//       // return responseCloned;
//     })
//   );
// });
