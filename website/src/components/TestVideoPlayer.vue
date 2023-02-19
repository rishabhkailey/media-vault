<script setup lang="ts">
import { ref, onMounted } from "vue";
import { VideoPlayer } from "@videojs-player/vue";
import "video.js/dist/video-js.css";
import decryptWorker from "@/worker/decrypt?url";
const fileName = ref<string>("file_example_AVI_1920_2_3MG.avi");
onMounted(() => {
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
  }
});

</script>

<template>
  <v-col>
    <v-row>
      <v-text-field
        v-model="fileName"
        :counter="10"
        label="File Input"
        required
      ></v-text-field>
    </v-row>
    <!-- :src="`/v1/testGetEncryptedVideo?file=${file}`" -->
    <!-- src="/v1/testGetVideoWithRange/test.mp4" -->
    <!-- src="/c1/testGetEncryptedVideoWithRange" -->
    <!-- type dynamic type="video/mp4" -->
    <v-row>
      <video-player
        :src="`/v1/testGetEncryptedVideoWithRange?file=${fileName}`"
        :controls="true"
        :autoplay="true"
        :loop="true"
        :volume="0.6"
      />
    </v-row>
  </v-col>
  <!-- <video-player src="https://vjs.zencdn.net/v/oceans.mp4" /> -->
</template>
