<script setup lang="ts">
import type { StyleValue } from "vue";
import ImageViewer from "@/components/MediaViewer/ImageViewer.vue";
import VideoViewer from "@/components/MediaViewer/VideoViewer.vue";

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
      v-if="props.media.type.startsWith('image')"
      :src="props.media.url"
      :low-res-src="props.media.thumbnail_url"
    />
    <VideoViewer
      v-else-if="props.media.type.startsWith('video')"
      :src="{
        src: props.media.url,
        type: props.media.type,
      }"
    />
    <div v-else>unknown</div>
  </div>
</template>
