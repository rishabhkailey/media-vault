<script setup lang="ts">
import { ref, onMounted } from "vue";
import decryptWorker from "@/worker/decrypt?url";
import streamSaver from "streamsaver";
import streamSaverMITM from "@/worker/mitm.html?url";
import { Chacha20 } from "ts-chacha20";
import axios from "axios";
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
  // let file = "text.txt";
  // let file = "100MB.dat";
  let file = "1GB.dat";
  const url = `http://localhost:5173/v1/testGetVideoWithRange?file=${file}`;
  whileRangeDownloadWithDecrypt(url, file)
    .then((done) => {
      console.log("done " + done);
    })
    .catch((err) => {
      console.log("error ", err);
    });
};
// currently working on this
const whileRangeDownloadWithDecrypt = async function (
  url: string,
  fileName: string
) {
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
  let password = "01234567890123456789012345678901";
  let nonce = "012345678901";
  const textEncoder = new TextEncoder();
  _global_decryptor = new Chacha20(
    textEncoder.encode(password),
    textEncoder.encode(nonce)
  );
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
    console.log(readableStream);
    // let decryptedStream = readableStream;
    let decryptedStream = readableStream.pipeThrough(newDecryptTransformer());
    // let decryptedStream = decryptStream(readableStream);
    console.log("got decrypted stream", decryptedStream);
    // working but using a lot of memory
    await decryptedStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
    console.log(index);
    // await new Promise(r => setTimeout(r, 5000));
  }
  fileStream.close();
};
const newDecryptTransformer: () => TransformStream<
  Uint8Array,
  Uint8Array
> = () =>
  new TransformStream<Uint8Array, Uint8Array>({
    start() {},
    transform(chunk, controller) {
      if (!chunk) {
        console.log("undefined chunk");
      }
      controller.enqueue(decryptChunk(chunk));
    },
    flush() {},
  });

let _global_decryptor: Chacha20;

const decryptChunk: (input: Uint8Array) => Uint8Array = (input) => {
  try {
    let decrypted = _global_decryptor.decrypt(input);
    if (!decrypted?.length || decrypted.length !== input.length) {
      console.log(decrypted);
    }
    return decrypted;
  } catch (err) {
    console.log(err);
  }
  return new Uint8Array(0);
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
