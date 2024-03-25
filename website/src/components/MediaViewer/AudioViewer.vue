<script setup lang="ts">
import { getQueryParamNumberValue } from "@/js/utils";
import { ref, watch } from "vue";
import { useRoute } from "vue-router";

const props = defineProps<{
  src: { src: string; type: string };
  mediaId: number;
}>();

const playerRef = ref<HTMLVideoElement | undefined>(undefined);
const route = useRoute();
watch(
  () => getQueryParamNumberValue(route.params, "media_id"),
  () => {
    if (props.mediaId !== getQueryParamNumberValue(route.params, "media_id")) {
      console.log("pausing audio");
      playerRef.value?.pause();
    }
  },
);
</script>

<template>
  <video
    ref="playerRef"
    style="width: 100vw; max-height: 100vh"
    controls
    autoplay
  >
    <source :src="props.src.src" :type="props.src.type" />
  </video>
</template>
