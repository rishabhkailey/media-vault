<script setup lang="ts">
import { ref, onMounted } from "vue";
import streamSaverMITM from "@/worker/mitm.html?url";
import decryptWorker from "@/worker/decrypt?url";
import axios from "axios";
import { useStore } from "vuex";
const store = useStore();
const video = ref<HTMLVideoElement | undefined>(undefined);

const url = ref<string>("");
const sources = ref<Array<{ src: string; type: string }>>([]);
const loading = ref<boolean>(true);
const accessToken = store.getters.accessToken;
console.log(accessToken);
onMounted(async () => {
  let response = await axios.post(
    "/v1/refreshSession",
    {},
    {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    }
  );
  console.log(response);
  response = await axios.get("/v1/mediaList", {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  // console.log(response);
  // response = await axios.get("/v1/media/6b809789-5d19-46d0-b087-7ab46ec8c0dd", {
  //   headers: {
  //     Authorization: `Bearer ${accessToken}`,
  //   },
  // });
  // console.log(response, undefined);

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
    url.value = "/v1/media/6ad05e28-13e3-488f-8503-06fd367f9106";
    sources.value = [
      {
        src: "/v1/media/6ad05e28-13e3-488f-8503-06fd367f9106",
        type: "video/mp4",
      },
    ];
    // url.value =
    //   "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4";
    loading.value = false;
    // fetch(
    //   "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4",
    //   {
    //     headers: {
    //       Range: `bytes=10-100`,
    //     },
    //   }
    // )
    //   .then((res) => {
    //     res.body
    //       ?.getReader()
    //       .read()
    //       .then((data) => {
    //         console.log(new TextDecoder().decode(data.value));
    //       })
    //       .catch((err) => {
    //         console.log(err);
    //       });
    //     console.log(res);
    //   })
    //   .catch((err) => {
    //     console.log(err);
    //   });
  }
});
</script>

<template>
  <v-container>
    <div v-if="loading">loading...</div>
    <video-player
      v-else
      :sources="sources"
      :controls="true"
      :autoplay="true"
      :loop="true"
      :volume="0.6"
    />
  </v-container>
  <!-- <video
    src="http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4"
  ></video> -->
</template>
// custom stream source //
https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/srcObject
