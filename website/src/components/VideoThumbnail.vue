<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
import { generateThumbnail } from "@/utils/thumbnail";

const uploadFiles = ref<Array<File>>([]);
const downloadButton = ref<HTMLAnchorElement | undefined>(undefined);
const SubmitHandler: (e: SubmitEventPromise) => any = async (
  e: SubmitEventPromise
) => {
  if (uploadFiles.value.length === 0) {
    // todo error message
    return;
  }
  let file = uploadFiles.value[0];
  console.log(file);
  generateThumbnail(file, {
    maxHeightWidth: 300,
  })
    .then((thumbnail) => {
      console.log("thumbnail generated");
      if (downloadButton.value !== undefined) {
        downloadButton.value.href = URL.createObjectURL(thumbnail);
        downloadButton.value.download = "thumbnail.jpeg";
      }
    })
    .catch((err) => {
      console.log(err);
    });
  // console.log(URL.createObjectURL(file));
  // this doesn't use much memory but to dataURL uses memory, checked it with 10GB file
  // video events
  // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video#events
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
