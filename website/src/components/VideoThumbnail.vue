<script setup lang="ts">
import { ref } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";
import { proccessThumbnailConstraints } from "@/utils/image";

const uploadFiles = ref<Array<File>>([]);

const SubmitHandler: (e: SubmitEventPromise) => any = async (
  e: SubmitEventPromise
) => {
  if (uploadFiles.value.length === 0) {
    // todo error message
    return;
  }
  let file = uploadFiles.value[0];
  console.log(file);
  await getVideoThumbnail(file);
  // console.log(URL.createObjectURL(file));
  // this doesn't use much memory but to dataURL uses memory, checked it with 10GB file
  // video events
  // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video#events
};

const getVideoThumbnail = async (file: File) => {
  return new Promise((resolve, reject) => {
    let video = document.createElement("video");

    video.src = URL.createObjectURL(file);
    video.onerror = (e) => {
      reject(e);
    };
    video.onloadedmetadata = (event) => {
      console.log("metadata", event);
      console.log(video.duration);
      // console.log(video.fastSeek);
    };
    let requestVideoFrameCallbackCalled = false;
    let thumbnailGenerated = false;
    video.oncanplay = async (event) => {
      if (requestVideoFrameCallbackCalled) {
        return;
      }
      requestVideoFrameCallbackCalled = true;
      console.log("canplay", event);
      // we want the thumbnail to be in the first minute
      // for video longer than 1 minute = 30, for video less than 1 minute = duration/2
      let thumbnailTime = Math.min(30, Math.floor(video.duration / 2));
      // video.fastSeek(thumbnailTime);
      video.currentTime = thumbnailTime;
      video.volume = 0;
      await video.play();
      await video.pause();
    };
    video.requestVideoFrameCallback(async (now, metadata) => {
      console.log(now, metadata);
      thumbnailGenerated = true;
      requestVideoFrameCallbackCalled = true;
      let thumbnail: Blob | undefined;
      // thumbnailGenerated = true;
      try {
        // video.height, width will be the html element width, height
        // which will always be 0, 0 as we are not rendering it on display
        thumbnail = await generateThumbnail(video, {
          width: metadata.width,
          height: metadata.height,
        });
      } catch (err) {
        reject(err);
        video.remove();
        return;
      }
      console.log("thumbnail", thumbnail);
    });
  });
};
// not sure why vue is not taking types from type.d.ts?
type WidthHeight = {
  width: number;
  height: number;
};
const generateThumbnail: (
  video: HTMLVideoElement,
  resolution: WidthHeight
) => Promise<Blob> = (video: HTMLVideoElement, resolution: WidthHeight) => {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement("canvas");
    const ctx: CanvasRenderingContext2D | null = canvas.getContext("2d");
    if (!ctx) {
      reject(new Error("canvas.GetContext returned null`"));
      return;
    }
    const { destinationResolution, sourceResolution, offset } =
      proccessThumbnailConstraints(
        { maxHeightWidth: 300 },
        {
          width: resolution.width,
          height: resolution.height,
        }
      );
    canvas.width = destinationResolution.width;
    canvas.height = destinationResolution.height;
    ctx.drawImage(
      video,
      offset.x,
      offset.y,
      sourceResolution.width,
      sourceResolution.height,
      0,
      0,
      destinationResolution.width,
      destinationResolution.height
    );
    document.body.append(canvas);
    canvas.toBlob((thumbnail) => {
      if (thumbnail === null) {
        reject(new Error("got null"));
        return;
      }
      resolve(thumbnail);
      return;
    }, "image/jpeg");
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
