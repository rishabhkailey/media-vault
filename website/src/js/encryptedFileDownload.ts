import * as streamSaver from "streamsaver";
import { parseResponseRangeHeader, type IResponseRange } from "./request";
// import { Chacha20 } from "ts-chacha20";
// const streamSaver = require("streamsaver");
export function download(url: string, fileName: string) {
  return new Promise<boolean>((resolve, reject) => {
    whileRangeDownloadWithDecrypt(url, fileName)
      .then((done) => {
        console.log("done " + done);
        resolve(true);
      })
      .catch((err) => {
        console.log("error ", err);
        reject(err);
      });
  });
}

const bufferSize = 10_000_000;
const idealRangeSize = 100_000_000;

// currently working on this
const whileRangeDownloadWithDecrypt = async function (
  url: string,
  fileName: string,
) {
  const fileStream = streamSaver.createWriteStream(fileName, {
    writableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
    readableStrategy: new ByteLengthQueuingStrategy({
      highWaterMark: bufferSize,
    }),
  });

  let range: IResponseRange | undefined = undefined;
  let length: number,
    response: Response,
    index = 0;
  let fileSize = 1; /*fileSize is set to 1 to send atleast 1 request, it will be updated after first response*/
  // 100 mb
  // let writer = fileStream.getWriter();
  // const password = "01234567890123456789012345678901";
  // const nonce = "012345678901";
  // const textEncoder = new TextEncoder();
  // const decryptor = new Chacha20(
  //   textEncoder.encode(password),
  //   textEncoder.encode(nonce)
  // );
  while (index < fileSize) {
    response = await fetch(url, {
      headers: {
        // range start and end is inclusive
        Range: `bytes=${index}-${index + idealRangeSize - 1}`,
      },
    });
    if (response.status !== 206) {
      throw new Error(
        `non 206 status from server. status = ${response.status}`
      );
    }
    const acceptRanges = response.headers.get("accept-ranges");
    if (acceptRanges && acceptRanges !== "bytes") {
      throw new Error(`server doesn't support range requests`);
    }
    const rangeHeader = response.headers.get("content-range");
    if (!rangeHeader) {
      throw new Error("Expected content-range not present");
    }
    if (range === undefined) {
      range = parseResponseRangeHeader(rangeHeader);
      fileSize = range.size;
    }
    length = range.end - range.start + 1;
    // length = Number(response.headers.get("content-length"));
    // if (length == 0) {
    //   throw new Error("empty response received");
    // }
    // write the current request data to file and then we will loop over the range
    const readableStream = response.body;
    if (!readableStream) {
      throw new Error("unable to get response stream");
    }
    console.log(readableStream);
    // let decryptedStream = readableStream;
    // no need to decrypt this twice we are decrypting in service worker
    // const decryptedStream = readableStream.pipeThrough(
    //   newDecryptTransformer(decryptor)
    // );
    // let decryptedStream = decryptStream(readableStream);
    console.log("got decrypted stream", readableStream);
    // working but using a lot of memory
    await readableStream.pipeTo(fileStream, {
      preventClose: true,
    });
    index += length;
    console.log(index);
    // await new Promise(r => setTimeout(r, 5000));
  }
  fileStream.close();
};

// export function newDecryptTransformer(decryptor: Chacha20) {
//   return new TransformStream<Uint8Array, Uint8Array>({
//     start() {},
//     transform(chunk, controller) {
//       if (!chunk) {
//         console.log("undefined chunk");
//       }
//       controller.enqueue(decryptChunk(decryptor, chunk));
//     },
//     flush() {},
//   });
// }

// export function decryptChunk(decryptor: Chacha20, input: Uint8Array) {
//   try {
//     const decrypted = decryptor.decrypt(input);
//     if (!decrypted?.length || decrypted.length !== input.length) {
//       console.log(decrypted);
//     }
//     return decrypted;
//   } catch (err) {
//     console.log(err);
//   }
//   return new Uint8Array(0);
// }
