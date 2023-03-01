<script setup lang="ts">
import { ref, onMounted } from "vue";
import streamSaverMITM from "@/worker/mitm.html?url";
import decryptWorker from "@/worker/decrypt?url";
const video = ref<HTMLVideoElement | undefined>(undefined);

const url = ref<string>("");
const loading = ref<boolean>(true);

onMounted(async () => {
  // streamSaver.WritableStream = WritableStream;
  // streamSaver.TransformStream = TransformStream;

  if ("serviceWorker" in navigator) {
    // todo try catch
    await navigator.serviceWorker.register(decryptWorker, {
      scope: "./",
      type: "module",
    });
    await navigator.serviceWorker.ready;
    console.log("Service worker active");
    url.value =
      "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4";
    loading.value = false;
    fetch(
      "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4",
      {
        headers: {
          Range: `bytes=10-100`,
        },
      }
    )
      .then((res) => {
        res.body
          ?.getReader()
          .read()
          .then((data) => {
            console.log(new TextDecoder().decode(data.value));
          })
          .catch((err) => {
            console.log(err);
          });
        console.log(res);
      })
      .catch((err) => {
        console.log(err);
      });
  }
});
</script>

<template>
  <div v-if="loading">loading...</div>
  <video-player
    v-else
    :src="url"
    :controls="true"
    :autoplay="true"
    :loop="true"
    :volume="0.6"
  />
  <!-- <video
    src="http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4"
  ></video> -->
</template>
// custom stream source //
https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/srcObject
