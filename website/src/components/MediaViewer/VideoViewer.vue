<script setup lang="ts">
import { getQueryParamNumberValue } from "@/js/utils";
import { ref, watch } from "vue";
import { useRoute } from "vue-router";

const props = defineProps<{
  mediaId: number;
  src: { src: string; type: string };
}>();

const playerRef = ref<HTMLVideoElement | undefined>(undefined);
const route = useRoute();
watch(
  () => getQueryParamNumberValue(route.params, "media_id"),
  () => {
    if (props.mediaId !== getQueryParamNumberValue(route.params, "media_id")) {
      console.debug("pausing video");
      playerRef.value?.pause();
    }
  },
);
</script>

<template>
  <video
    ref="playerRef"
    style="max-width: 100vw; max-height: 100vh"
    controls
    autoplay
  >
    <source :src="props.src.src" :type="props.src.type" />
    <!-- if the specified video type is not supported then try using video/mp4 -->
    <source :src="props.src.src" type="video/mp4" />
    <!-- in the end let the browser guess the type -->
    <source :src="props.src.src" />
  </video>
</template>
