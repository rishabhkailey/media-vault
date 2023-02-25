<script setup lang="ts">
import { ref, onMounted } from "vue";
import streamSaverMITM from "@/worker/mitm.html?url";
import decryptWorker from "@/worker/decrypt?url";
const video = ref<HTMLVideoElement | undefined>(undefined);

onMounted(() => {
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
    fetch(
      "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4",
      {
        headers: {
          Range: `bytes=0-100`,
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
  <video-player
    :src="`/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4`"
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
