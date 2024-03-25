<script setup lang="ts">
import type { StyleValue } from "vue";
import ImageViewer from "@/components/MediaViewer/ImageViewer.vue";
import VideoViewer from "@/components/MediaViewer/VideoViewer.vue";
import PdfViewer from "@/components/MediaViewer/PdfViewer.vue";
import UnknownMediaViewer from "@/components/MediaViewer/UnknownMediaViewer.vue";
import { isImage, isPdf, isVideo } from "@/js/files/type";

const props = defineProps<{
  media: Media;
  class?: any;
  style?: StyleValue;
}>();
</script>

<template>
  <div
    :class="props.class"
    :style="props.style"
    :id="`media_window_${props.media.id}`"
  >
    <ImageViewer
      v-if="isImage(props.media.type)"
      :src="props.media.url"
      :low-res-src="props.media.thumbnail_url"
    />
    <VideoViewer
      v-else-if="isVideo(props.media.type)"
      :src="{
        src: props.media.url,
        type: props.media.type,
      }"
    />
    <PdfViewer v-else-if="isPdf(props.media.type)" :src="props.media.url" />
    <UnknownMediaViewer :media="props.media" v-else />
  </div>
</template>
