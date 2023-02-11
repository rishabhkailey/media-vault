<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
import axios from "axios";

const uploadFiles = ref<Array<File>>([]);

const SubmitHandler: (e: SubmitEventPromise) => any = async (e: SubmitEventPromise) => {
  if (uploadFiles.value.length === 0) {
    // todo error message
    return;
  }
  let file = uploadFiles.value[0];
  console.log(file);

  let formData = new FormData();
  formData.append("name", file.name);
  formData.append("size", file.size.toString());
  formData.append("file", file);
  let stream = file.stream();
  let reader = stream.getReader();
  let readBytes = 0;
  let requestID: string;
  try {
    let response = await axios.post(
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
    );
    if (response.status !== 200) {
      throw new Error("init equest failed with " + response.status);
    }
    if (!response.data?.requestID || typeof response.data.requestID !== "string") {
      throw new Error("invalid response for init request " + response);
    }
    requestID = response.data.requestID;
    console.log(response);
  } catch (err) {
    console.log(err);
    return;
  }
  await new Promise((resolve, reject) => {
    reader
      .read()
      .then(async function uploadChunk({ done, value }) {
        if (done) {
          console.log(`${readBytes} read`);
          resolve(readBytes);
          return;
        }
        if (value == undefined) {
          throw new Error("empty chunk received");
        }
        let chunkBlob = new Blob([value]);
        // let formData = new FormData();
        // formData.append("requestID", requestID);
        // formData.append("index", readBytes);
        // formData.append("chunkSize", value.length);
        // formData.append("chunkData", Blob);
        let response = await axios.post(
          "/v1/uploadChunk",
          {
            requestID: requestID,
            index: readBytes,
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
        readBytes += value.length;
        console.log(`chunk of length ${value.length}`);
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

  axios
    .post(
      "/v1/finishChunkUpload",
      {
        requestID: requestID,
        checksum: "file.size",
      },
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    )
    .then((res) => {
      console.log(res);
    })
    .catch((err) => {
      console.log(err);
    });
  // /v1/testEncryptedUpload
  // /v1/testNormalUpload
  // /v1/testVideoUploadWithThumbnail
  // /v1/testStreamVideoUploadWithThumbnail
  // /v1/testEncryptedFileSave
  // axios
  //   .post("/v1/testNormalUpload", formData, {
  //     headers: {
  //       "Content-Type": "multipart/form-data",
  //     },
  //   })
  //   .then(function () {
  //     console.log("SUCCESS!!");
  //   })
  //   .catch(function () {
  //     console.log("FAILURE!!");
  //   });
};
</script>

<template>
  <v-form min-width="400px" ref="form" @submit.prevent="SubmitHandler">
    <v-file-input
      style="min-width: 400px"
      v-model="uploadFiles"
      show-size
      label="File input"
    />
    <v-btn type="submit">upload</v-btn>
  </v-form>
</template>
