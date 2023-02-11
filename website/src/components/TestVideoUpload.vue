<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
// import { generateVideoThumbnail } from "@/utils/video";
import axios from "axios";
// import { createFFmpeg, fetchFile } from "@ffmpeg/ffmpeg";
const uploadFiles = ref<Array<File>>([]);
const downloadButton = ref<HTMLAnchorElement | undefined>(undefined);
const thumbnail = ref<HTMLImageElement | undefined>(undefined);
// const ffmpeg = createFFmpeg({
//   log: true,
// });

const SubmitHandler: (e: SubmitEventPromise) => any = (
  e: SubmitEventPromise
) => {
  if (uploadFiles.value.length === 0) {
    // todo error message
    return;
  }
  let file = uploadFiles.value[0];
  console.log(file);

  if (!file.name.toLowerCase().endsWith(".mp4")) {
    console.log("not a video");
    return;
  }
  let formData = new FormData();
  formData.append("name", file.name);
  formData.append("size", file.size.toString());
  formData.append("file", file);
  // generateVideoThumbnail(
  //   file,
  //   () => {},
  //   () => {}
  // );
  // requires some cors headers will try later if no other solution works
  // ffmpeg.load().then(() => {
  //   fetchFile(file).then((fileBuffer) => {
  //     ffmpeg.FS("writeFile", file.name, fileBuffer);
  //     ffmpeg
  //       .run("-ss", "00:00:01", "-i", file.name, "thumbnail.png")
  //       .then(() => {
  //         let thumbnailBuffer = ffmpeg.FS("readFile", file.name);
  //         if (thumbnail.value === undefined) {
  //           console.log("thumbnail not intialized");
  //           return;
  //         }
  //         thumbnail.value.src = URL.createObjectURL(
  //           new Blob([thumbnailBuffer.buffer], { type: "image/png" })
  //         );
  //       });
  //   });
  // });
  // generateThumbnail(
  //   file,
  //   { maxHeightWidth: 300 },
  //   (file, thumbnail) => {
  //     if (downloadButton.value !== undefined) {
  //       downloadButton.value.href = URL.createObjectURL(thumbnail);
  //       downloadButton.value.download = "thumbnail.jpeg";
  //     }
  //     console.log("file", file, "thumbnail", thumbnail);
  //   },
  //   (err) => {
  //     console.log("failed to generate thumbnail", err);
  //   }
  // );

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
    <img ref="thumbnail" />
  </v-form>
</template>
