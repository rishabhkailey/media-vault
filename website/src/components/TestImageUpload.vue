<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
import { generateThumbnail } from "@/utils/image";
import axios from "axios";

const uploadFiles = ref<Array<File>>([]);
const downloadButton = ref<HTMLAnchorElement|undefined>(undefined);

const SubmitHandler: (e: SubmitEventPromise) => any = (e: SubmitEventPromise) => {
  if (uploadFiles.value.length === 0) {
    // todo error message
    return;
  }
  let file = uploadFiles.value[0];
  console.log(file);

  if (
    !file.name.toLowerCase().endsWith(".png") &&
    !file.name.toLowerCase().endsWith(".jpg") &&
    !file.name.toLowerCase().endsWith(".jpeg")
  ) {
    console.log("not a image");
    return;
  }
  let formData = new FormData();
  formData.append("name", file.name);
  formData.append("size", file.size.toString());
  formData.append("file", file);
  generateThumbnail(
    file,
    { maxHeightWidth: 300 },
    (file, thumbnail) => {
      if (downloadButton.value !== undefined){
        downloadButton.value.href = URL.createObjectURL(thumbnail);
        downloadButton.value.download = "thumbnail.jpeg";
      }
      console.log("file", file, "thumbnail", thumbnail);
    },
    (err) => {
      console.log("failed to generate thumbnail", err);
    }
  );

  // scaleImage(
  //   file,
  //   { height: 300, width: 300 },
  //   (file, thumbnail) => {
  //     if (downloadButton.value !== undefined){
  //       downloadButton.value.href = URL.createObjectURL(thumbnail);
  //       downloadButton.value.download = "thumbnail.jpeg";
  //     }
  //     console.log("file", file, "thumbnail", thumbnail);
  //   },
  //   (err) => {
  //     console.log("failed to generate thumbnail", err);
  //   }
  // );
  
  // /v1/testEncryptedUpload
  // /v1/testNormalUpload
  // axios
  //   .post("/v1/testEncryptedUpload", formData, {
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
    <a ref="downloadButton"> Download thumbnail</a>
  </v-form>
</template>
