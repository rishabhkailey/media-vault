import axios from "axios";
import { Chacha20 } from "ts-chacha20";

type ProgressCallback = (percentage: number) => void;
// todo update promise<any> to request info or something

interface ChunkUploadInfo {
  requestID: string;
  uploadedBytes: number;
}

export function chunkUpload(
  file: File,
  callback: ProgressCallback
): Promise<ChunkUploadInfo> {
  return new Promise((resolve, reject) => {
    initChunkUpload(file)
      .then((chunkRequestInfo) => {
        uploadFileChunks(file, chunkRequestInfo, callback)
          .then((uploadedBytes) => {
            console.log("uploaded " + uploadedBytes + " bytes of " + file.name);
            finishChunkUpload(chunkRequestInfo)
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
function initChunkUpload(file: File): Promise<ChunkRequestInfo> {
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
          },
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
            await encryptAndUploadChunk(
              requestID,
              bytesUploaded,
              buffer.slice(0, bufferIndex),
              encryptor
            );
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

        readBytes += value.length;
        while (bufferIndex + value.length >= chunkSize) {
          buffer.set(value.slice(0, chunkSize - bufferIndex), bufferIndex);
          value = value.slice(chunkSize - bufferIndex);
          bufferIndex = chunkSize;
          // upload buffer
          await encryptAndUploadChunk(
            requestID,
            bytesUploaded,
            buffer,
            encryptor
          );
          console.log("upload buffer", bytesUploaded);
          // if buffer upload successful
          bufferIndex = 0;
          // do something here
          bytesUploaded += chunkSize;
          // update progress
          if (callback !== undefined) {
            callback((bytesUploaded * 100) / file.size);
          }
          console.log(`chunk of length ${value.length}`);
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
  encryptor: Chacha20
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
  const response = await axios.post(
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
    }
  );
  console.log(response);
  if (response.status !== 200) {
    throw new Error(
      "upload chuck request failed with status" + response.status
    );
  }
}

function finishChunkUpload(
  chunkRequestInfo: ChunkRequestInfo
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
