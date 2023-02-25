// stream saver worker code + our custom code for decryption

/* global self ReadableStream Response */
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

function createStream(port) {
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

self.onfetch = (event) => {
  const url = event.request.url;

  // this only works for Firefox
  if (url.endsWith("/ping")) {
    return event.respondWith(new Response("pong"));
  }

  const hijacke = map.get(url);

  if (!hijacke) return null;

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

  if (headers.has("Content-Length")) {
    responseHeaders.set("Content-Length", headers.get("Content-Length"));
  }

  if (headers.has("Content-Disposition")) {
    responseHeaders.set(
      "Content-Disposition",
      headers.get("Content-Disposition")
    );
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