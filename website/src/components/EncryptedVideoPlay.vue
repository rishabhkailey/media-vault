<script setup lang="ts">
import { ref, onMounted } from "vue";
const video = ref<HTMLVideoElement | undefined>(undefined);
var source = new ReadableStream({
  start(controller) {
    // Fetch the video as a stream
    fetch(
      "http://localhost:5173/v1/testGetVideoWithRange?file=SampleVideo_1280x720_30mb.mp4"
    )
      .then((response) => {
        return response.body;
      })
      .then((stream) => {
        if (!stream) {
          console.log("invalid response");
          return;
        }
        // Read chunks from the stream and push them to the controller
        var reader = stream.getReader();

        function push() {
          reader.read().then((result) => {
            if (result.done) {
              controller.close();
            } else {
              controller.enqueue(result.value);
              push();
            }
          });
        }
        push();
      });
  },
});

onMounted(() => {
  source.pipeTo(new WritableStream()).then((stream) => {
    // we will probably end up using a lot of memory
    // video.value?.src = URL.createObjectURL(stream);
  });

  video.value?.play();
});
</script>

<template>
  <video ref="video"></video>
</template>
// custom stream source 
// https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/srcObject

