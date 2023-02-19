<script setup lang="ts">
import { ref, onMounted } from "vue";
import decryptWorker from "@/worker/decrypt?url";
import streamSaver from "streamsaver";
import streamSaverMITM from "@/worker/mitm.html?url";
import axios from "axios";
import sodium, { type StateAddress } from "libsodium-wrappers";
import { transform } from "@vue/compiler-core";

const fileName = ref<string>("file_example_AVI_1920_2_3MG.avi");
onMounted(async () => {
  streamSaver.mitm = streamSaverMITM;
  // streamSaver.WritableStream = WritableStream;
  // streamSaver.TransformStream = TransformStream;

  if ("serviceWorker" in navigator) {
    navigator.serviceWorker
      .register(decryptWorker, {
        scope: "./",
        type: "module",
      })
      .then((registration) => {
        if (registration.installing) {
          console.log("Service worker installing");
        } else if (registration.waiting) {
          console.log("Service worker installed");
        } else if (registration.active) {
          console.log("Service worker active");
        }
      })
      .catch((error) => {
        console.error(`Registration failed with ${error}`);
      });
  }
});
interface IRange {
  unit: string;
  start: number;
  end: number;
  size: number;
}
// bytes 200-1000/67589
// todo add suport for * size
function parseRangeHeader(range: string): IRange {
  let parts = range.split(" ");
  if (parts.length != 2) {
    throw new Error("Invalid range header " + range);
  }
  let unit = parts[0];
  // rangeSizePart = 200-1000/67589
  let rangeSizePart = parts[1];
  parts = rangeSizePart.split("/");
  if (parts.length != 2) {
    throw new Error("Invalid range header " + range);
  }
  let size = Number(parts[1]);
  // rangePart = 200-1000
  let rangePart = parts[0];
  if (rangePart.split("-").length != 2) {
    throw new Error("Invalid range header " + range);
  }
  let start = Number(rangePart.split("-")[0]);
  let end = Number(rangePart.split("-")[1]);
  return {
    unit,
    start,
    end,
    size,
  };
}
const download = async () => {
  // // complete file download (causing memory issues)
  // const url = "http://localhost:5173/v1/testDownload?file=10GB.dat";
  // const fileStream = streamSaver.createWriteStream("cat.mp4");
  // fetch(url).then((res) => {
  //   const readableStream = res.body;
  //   // more optimized
  //   if (window.WritableStream && readableStream?.pipeTo) {
  //     return readableStream
  //       .pipeTo(fileStream)
  //       .then(() => console.log("done writing"));
  //   }
  //   let writer = fileStream.getWriter();
  //   const reader = res.body.getReader();
  //   const pump = () =>
  //     reader
  //       .read()
  //       .then((res) =>
  //         res.done ? writer.close() : writer.write(res.value).then(pump)
  //       );
  //   pump();
  // });

  // file download using range request
  // const url = "http://localhost:5173/v1/testGetVideoWithRange?file=10GB.dat";
  // const url = "http://localhost:5173/v1/testGetVideoWithRange?file=test.dat";
  const url = "http://localhost:5173/v1/testGetVideoWithRange?file=100MB.dat";
  whileRangeDownloadWithDecrypt(url, "test.dat")
    .then((done) => {
      console.log("done " + done);
    })
    .catch((err) => {
      console.log("error ", err);
    });
  // new Response("StreamSaver is awesome").body?.pipeTo(fileStream);
  // const link = document.createElement("a");
  // link.href = `/v1/testGetEncryptedVideoWithRange?file=${fileName.value}`;
  // link.download = fileName.value;
  // link.click();
  // let url = `/v1/testGetEncryptedVideoWithRange?file=${fileName.value}`;
  // let url = `/v1/testGetVideoWithRange?file=${fileName.value}`;
  // let response = await fetch(url, {
  //   headers: {
  //     Range: "bytes=0-",
  //   },
  // });
  // let acceptRanges = response.headers.get("accept-ranges");
  // let supportRangeRequest: boolean = false;
  // if (acceptRanges && acceptRanges === "bytes") {
  //   supportRangeRequest = true;
  // }
  // let writer = fileStream.getWriter();
  // if (!supportRangeRequest) {
  //   const readableStream = response.body;
  //   if (!readableStream) {
  //     throw new Error("unable to get response stream");
  //   }
  //   readableStream
  //     .pipeTo(fileStream)
  //     .then(() => console.log("done writing"))
  //     .catch((err) => {
  //       throw new Error("error writing the response to file " + err);
  //     });
  //   return;
  //   // todo not all browsers support pipeTo
  //   // check the fallback method here https://github.com/jimmywarting/StreamSaver.js/blob/master/examples/fetch.html
  // }
  // // support range requests
  // let supportedRange = response.headers.get("content-range");
  // if (!supportedRange) {
  //   throw new Error("Expected content-range not present");
  // }
  // let range = parseRangeHeader(supportedRange);
  // let length = Number(response.headers.get("content-length"));
  // let index = 0;
  // if (length !== 0) {
  //   // write the current request data to file and then we will loop over the range
  //   const readableStream = response.body;
  //   if (!readableStream) {
  //     throw new Error("unable to get response stream");
  //   }
  //   await readableStream.pipeTo(fileStream);
  //   index += length;
  // }
  // // 1mb
  // let idealRangeSize = 1000000;
  // while (index < range.size) {
  //   response = await fetch(url, {
  //     headers: {
  //       Range: `bytes=${index}-${index + idealRangeSize}`,
  //     },
  //   });
  //   length = Number(response.headers.get("content-length"));
  //   if (length == 0) {
  //     throw new Error("empty response received");
  //   }
  //   // write the current request data to file and then we will loop over the range
  //   const readableStream = response.body;
  //   if (!readableStream) {
  //     throw new Error("unable to get response stream");
  //   }
  //   await readableStream.pipeTo(fileStream);
  //   index += length;
  //   // fetch(url).then((res) => {
  //   //   const readableStream = res.body;
  //   //   // more optimized
  //   //   if (window.WritableStream && readableStream?.pipeTo) {
  //   //     return readableStream.pipeTo(fileStream).then(() => console.log("done writing"));
  //   //   }
  //   //   const reader = res.body?.getReader();
  //   //   if (!reader) {
  //   //     throw new Error("unable to get response reader");
  //   //   }
  //   //   const pump = () =>
  //   //     reader
  //   //       .read()
  //   //       .then((res) =>
  //   //         res.done ? writer.close() : writer.write(res.value).then(pump)
  //   //       );
  //   //   pump();
  //   // });
  // }
  // writer.close();
};

const rangeDownload = async function (url: string, fileName: string) {
  let index = 0;
  // 100 mb
  let idealRangeSize = 100_000_000;
  let response = await fetch(url, {
    headers: {
      Range: `bytes=${index}-${index + idealRangeSize}`,
    },
  });
  if (response.status !== 206) {
    throw new Error(`non 206 status from server. status = ${response.status}`);
  }
  let acceptRanges = response.headers.get("accept-ranges");
  if (acceptRanges && acceptRanges !== "bytes") {
    throw new Error(`server doesn't support range requests`);
  }
  let rangeHeader = response.headers.get("content-range");
  if (!rangeHeader) {
    throw new Error("Expected content-range not present");
  }

  let bufferSize = 10_000_000;
  const fileStream = streamSaver.createWriteStream(fileName, {
    writableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
    readableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
  });
  let range = parseRangeHeader(rangeHeader);
  let length = Number(response.headers.get("content-length"));
  if (length !== 0) {
    // write the current request data to file and then we will loop over the range
    const readableStream = response.body;
    if (!readableStream) {
      throw new Error("unable to get response stream");
    }
    await readableStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
  }
  // let writer = fileStream.getWriter();
  while (index < range.size) {
    response = await fetch(url, {
      headers: {
        Range: `bytes=${index}-${index + idealRangeSize}`,
      },
    });
    if (response.status !== 206) {
      throw new Error(
        `non 200 status from server. status = ${response.status}`
      );
    }
    length = Number(response.headers.get("content-length"));
    if (length == 0) {
      throw new Error("empty response received");
    }
    // write the current request data to file and then we will loop over the range
    const readableStream = response.body;
    if (!readableStream) {
      throw new Error("unable to get response stream");
    }
    await readableStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
    console.log(index);
  }
  fileStream.close();
};

// currently working on this
const whileRangeDownload = async function (url: string, fileName: string) {
  let bufferSize = 10_000_000;
  const fileStream = streamSaver.createWriteStream(fileName, {
    writableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
    readableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
  });

  let range: IRange | undefined = undefined;
  let length: number,
    response: Response,
    index = 0,
    idealRangeSize = 100_000_000;
  let fileSize = 1; /*fileSize is set to 1 to send atleast 1 request, it will be updated after first response*/
  // 100 mb
  // let writer = fileStream.getWriter();
  while (index < fileSize) {
    response = await fetch(url, {
      headers: {
        Range: `bytes=${index}-${index + idealRangeSize}`,
      },
    });
    if (response.status !== 206) {
      throw new Error(
        `non 206 status from server. status = ${response.status}`
      );
    }
    let acceptRanges = response.headers.get("accept-ranges");
    if (acceptRanges && acceptRanges !== "bytes") {
      throw new Error(`server doesn't support range requests`);
    }
    let rangeHeader = response.headers.get("content-range");
    if (!rangeHeader) {
      throw new Error("Expected content-range not present");
    }
    if (range === undefined) {
      range = parseRangeHeader(rangeHeader);
      fileSize = range.size;
    }
    length = Number(response.headers.get("content-length"));
    if (length == 0) {
      throw new Error("empty response received");
    }
    // write the current request data to file and then we will loop over the range
    const readableStream = response.body;
    if (!readableStream) {
      throw new Error("unable to get response stream");
    }

    // working but using a lot of memory
    await readableStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
    console.log(index);
  }
  fileStream.close();
};

// currently working on this
const whileRangeDownloadWithDecrypt = async function (
  url: string,
  fileName: string
) {
  await sodium.ready;
  let bufferSize = 10_000_000;
  const fileStream = streamSaver.createWriteStream(fileName, {
    writableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
    readableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
  });

  let range: IRange | undefined = undefined;
  let length: number,
    response: Response,
    index = 0,
    idealRangeSize = 100_000_000;
  let fileSize = 1; /*fileSize is set to 1 to send atleast 1 request, it will be updated after first response*/
  // 100 mb
  // let writer = fileStream.getWriter();
  while (index < fileSize) {
    response = await fetch(url, {
      headers: {
        Range: `bytes=${index}-${index + idealRangeSize}`,
      },
    });
    if (response.status !== 206) {
      throw new Error(
        `non 206 status from server. status = ${response.status}`
      );
    }
    let acceptRanges = response.headers.get("accept-ranges");
    if (acceptRanges && acceptRanges !== "bytes") {
      throw new Error(`server doesn't support range requests`);
    }
    let rangeHeader = response.headers.get("content-range");
    if (!rangeHeader) {
      throw new Error("Expected content-range not present");
    }
    if (range === undefined) {
      range = parseRangeHeader(rangeHeader);
      fileSize = range.size;
    }
    length = Number(response.headers.get("content-length"));
    if (length == 0) {
      throw new Error("empty response received");
    }
    // write the current request data to file and then we will loop over the range
    const readableStream = response.body;
    if (!readableStream) {
      throw new Error("unable to get response stream");
    }

    let password = "berrysikure";
    _global_state_in = DecryptInit(password);
    let decryptedStream = readableStream.pipeThrough(decryptTransformer);
    // let decryptedStream = decryptStream(readableStream);
    console.log("got decrypted stream", decryptedStream);
    // working but using a lot of memory
    await decryptedStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
    console.log(index);
  }
  fileStream.close();
};
const decryptTransformer = new TransformStream<Uint8Array, Uint8Array>({
  start() {},
  transform(chunk, controller) {
    if (!chunk) {
      console.log("undefined chunk");
    }
    controller.enqueue(decryptChunk(chunk));
  },
  flush() {},
});

const decryptStream: (
  input: ReadableStream<Uint8Array>
) => ReadableStream<Uint8Array> = (input) => {
  let reader = input.getReader();
  let idealRangeSize = 100_000_000;
  return new ReadableStream(
    {
      start(controller) {
        // The following function handles each data chunk
        function push() {
          // "done" is a Boolean and value a "Uint8Array"
          reader.read().then(({ done, value }) => {
            // If there is no more data to read
            if (done) {
              console.log("done", done);
              controller.close();
              return;
            }
            console.log(done);
            // Get the data and send it to the browser via the controller
            controller.enqueue(decryptChunk(value));
            // Check chunks by logging to the console
            console.log(done, value);
            push();
          });
        }
        push();
      },
    },
    new ByteLengthQueuingStrategy({
      highWaterMark: idealRangeSize,
    })
  );
};

// let _global_res: {
//   state: StateAddress;
//   header: Uint8Array;
// };
// let _global_key: Uint8Array;
let _global_state_in: StateAddress;
const DecryptInit: (password: string) => StateAddress = (password) => {
  if (password.length == 0) {
    throw new Error("password length zero");
  }
  while (password.length < 32) {
    password += password;
  }
  password = password.substring(0, 32);
  const encoder = new TextEncoder();
  let key = encoder.encode(password);
  console.log(sodium);
  let res = sodium.crypto_secretstream_xchacha20poly1305_init_push(key);
  const state_in = sodium.crypto_secretstream_xchacha20poly1305_init_pull(
    res.header,
    key
  );
  return state_in;
};

const decryptChunk: (input: Uint8Array) => Uint8Array = (input) => {
  let output = sodium.crypto_secretstream_xchacha20poly1305_pull(
    _global_state_in,
    input
  );
  if (output === false || !output.message) {
    console.log(input);
    throw new Error("decryption failed");
  }
  return output.message;
};
</script>

<template>
  <v-col>
    <v-row>
      <v-text-field
        v-model="fileName"
        :counter="10"
        label="File Input"
        required
        :style="{ width: '300px' }"
      ></v-text-field>
    </v-row>
    <!-- :src="`/v1/testGetEncryptedVideo?file=${file}`" -->
    <!-- src="/v1/testGetVideoWithRange/test.mp4" -->
    <!-- src="/c1/testGetEncryptedVideoWithRange" -->
    <!-- type dynamic type="video/mp4" -->
    <v-row>
      <!-- <video-player
        :src="`/v1/testGetEncryptedVideoWithRange?file=${fileName}`"
        :controls="true"
        :autoplay="true"
        :loop="true"
        :volume="0.6"
      /> -->
    </v-row>
    <!-- <v-btn :href="`/v1/testGetEncryptedVideoWithRange?file=${fileName}`" download
      >download</v-btn
    > -->
    <v-btn @click="download">download using stream saver</v-btn>
  </v-col>
  <!-- <video-player src="https://vjs.zencdn.net/v/oceans.mp4" /> -->
</template>
//
https://stackoverflow.com/questions/65984220/node-js-high-memory-usage-when-using-createreadstream-and-createwritestream
// https://github.com/jimmywarting/StreamSaver.js/issues/133
