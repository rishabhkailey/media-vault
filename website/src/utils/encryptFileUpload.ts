import axios, { type AxiosResponse } from "axios";
import { Chacha20 } from "ts-chacha20";

type ProgressCallback = (percentage: number) => void;
// todo update promise<any> to request info or something

interface ChunkUploadInfo {
  requestID: string;
  uploadedBytes: number;
}

export function chunkUpload(
  file: File,
  bearerToken: string,
  controller: AbortController,
  callback: ProgressCallback
): Promise<ChunkUploadInfo> {
  return new Promise((resolve, reject) => {
    initChunkUpload(file, bearerToken, controller)
      .then((chunkRequestInfo) => {
        uploadFileChunks(file, chunkRequestInfo, controller, callback)
          .then((uploadedBytes) => {
            console.log("uploaded " + uploadedBytes + " bytes of " + file.name);
            finishChunkUpload(chunkRequestInfo, controller)
              .then((success) => {
                if (success) {
                  // todo
                  resolve({
                    requestID: chunkRequestInfo.requestID,
                    uploadedBytes: uploadedBytes,
                  });
                  return;
                } else {
                  reject(new Error("finish chunk upload request failed"));
                  return;
                }
              })
              .catch((err) => {
                reject(err);
                return;
              });
          })
          .catch((err) => {
            reject(err);
            return;
          });
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
}

interface ChunkRequestInfo {
  requestID: string;
}
function initChunkUpload(
  file: File,
  bearerToken: string,
  controller: AbortController
): Promise<ChunkRequestInfo> {
  return new Promise((resolve, reject) => {
    axios
      .post(
        "/v1/initChunkUpload",
        {
          fileName: file.name,
          size: file.size,
          fileType: "txt",
        },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${bearerToken}`,
          },
          signal: controller.signal,
        }
      )
      .then((response) => {
        if (response.status !== 200) {
          throw new Error("init equest failed with " + response.status);
        }
        if (
          !response.data?.requestID ||
          typeof response.data.requestID !== "string"
        ) {
          throw new Error("invalid response for init request " + response);
        }
        resolve({
          requestID: response.data.requestID,
        });
        return;
      })
      .catch((err) => {
        reject(new Error("init request failed: " + err));
      });
  });
}
// test cases
// read bytes usually around 2Mi > chunksize (then for each read multiple chunks uploaded)
// read bytes < chunk size (multipile read to create 1 chunk)
// file size < defaultChunkSize
// 25 Mi
const defaultChunkSize = 25 * 1024 * 1024;
// todo buffer reading, once buffer full then send upload request
function uploadFileChunks(
  file: File,
  chunkRequestInfo: ChunkRequestInfo,
  controller: AbortController,
  callback: ProgressCallback
): Promise<number> {
  const stream = file.stream();
  const reader = stream.getReader();
  // change this to uploaded bytes
  let readBytes = 0;
  let bytesUploaded = 0;
  const requestID = chunkRequestInfo.requestID;
  // init upload
  const password = "01234567890123456789012345678901";
  const nonce = "012345678901";
  const textEncoder = new TextEncoder();
  const encryptor = new Chacha20(
    textEncoder.encode(password),
    textEncoder.encode(nonce)
  );
  let chunkSize = defaultChunkSize;
  if (chunkSize > file.size) {
    chunkSize = file.size;
  }
  console.log("chunk size", chunkSize);
  const buffer = new Uint8Array(chunkSize);
  let bufferIndex = 0;
  return new Promise((resolve, reject) => {
    reader
      .read()
      .then(async function uploadChunk({ done, value }) {
        if (done) {
          // upload remaining
          if (bufferIndex != 0) {
            try {
              await encryptAndUploadChunk(
                requestID,
                bytesUploaded,
                buffer.slice(0, bufferIndex),
                encryptor,
                controller
              );
            } catch (err) {
              reject(err);
              return;
            }
            console.log("upload buffer", bytesUploaded);
            bytesUploaded += bufferIndex;
          }
          console.log(`${bytesUploaded} uploaded`);
          console.log(`${readBytes} read`);
          resolve(bytesUploaded);
          return;
        }

        if (value == undefined) {
          throw new Error("empty chunk received");
        }

        console.log(`read chunk of length ${value.length}`);
        readBytes += value.length;
        while (bufferIndex + value.length >= chunkSize) {
          buffer.set(value.slice(0, chunkSize - bufferIndex), bufferIndex);
          value = value.slice(chunkSize - bufferIndex);
          bufferIndex = chunkSize;
          // upload buffer
          try {
            await encryptAndUploadChunk(
              requestID,
              bytesUploaded,
              buffer,
              encryptor,
              controller
            );
          } catch (err) {
            reject(err);
            return;
          }
          console.log("upload buffer", bytesUploaded);
          // if buffer upload successful
          bufferIndex = 0;
          // do something here
          bytesUploaded += chunkSize;
          // update progress
          if (callback !== undefined) {
            callback((bytesUploaded * 100) / file.size);
          }
        }
        // new chunk is not satisfying the buffer size so we will save it in the buffer and read more
        if (value.length > 0) {
          buffer.set(value, bufferIndex);
          bufferIndex += value.length;
        }
        reader
          .read()
          .then(uploadChunk)
          .catch((err) => {
            throw err;
          });
      })
      .catch((err) => {
        reject(err);
      });
  });
}

async function encryptAndUploadChunk(
  requestID: string,
  index: number,
  value: Uint8Array,
  encryptor: Chacha20,
  controller: AbortController
) {
  let chunkBlob: Blob;
  try {
    // todo add encryption again
    const encryptedData = encryptor.encrypt(value);
    // console.log(new TextDecoder().decode(value));
    // console.log(new TextDecoder().decode(encryptedData));
    chunkBlob = new Blob([encryptedData]);
  } catch (err) {
    throw new Error("Encryptiong failed " + err);
  }
  let response: AxiosResponse<any>;
  try {
    response = await axios.post(
      "/v1/uploadChunk",
      {
        requestID: requestID,
        index: index,
        chunkSize: value.length,
        chunkData: chunkBlob,
      },
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        signal: controller.signal,
      }
    );
  } catch (err) {
    console.log("Upload chunk failed with error " + err);
    throw new Error("Upload chunk failed with error " + err);
  }
  console.log(response);
  if (response.status !== 200) {
    throw new Error(
      "upload chunk request failed with status" + response.status
    );
  }
}

function finishChunkUpload(
  chunkRequestInfo: ChunkRequestInfo,
  controller: AbortController
): Promise<boolean> {
  return new Promise((resolve, reject) => {
    // finish upload
    axios
      .post(
        "/v1/finishChunkUpload",
        {
          requestID: chunkRequestInfo.requestID,
          checksum: "file.size",
        },
        {
          headers: {
            "Content-Type": "application/json",
          },
          signal: controller.signal,
        }
      )
      .then((res) => {
        if (res.status !== 200) {
          throw new Error(
            "upload chuck request failed with status" + res.status
          );
        }
        resolve(true);
        return;
      })
      .catch((err) => {
        reject(err);
      });
  });
}
