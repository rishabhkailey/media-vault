<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
import axios from "axios";

const uploadFiles = ref<Array<File>>([]);

const SubmitHandler: (e: SubmitEventPromise) => any = (e: SubmitEventPromise) => {
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
  // /v1/testEncryptedUpload
  // /v1/testNormalUpload
  axios
    .post("/v1/testEncryptedUpload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
    .then(function () {
      console.log("SUCCESS!!");
    })
    .catch(function () {
      console.log("FAILURE!!");
    });
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
