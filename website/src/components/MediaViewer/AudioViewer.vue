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
  <div class="d-flex flex-column justify-center-align-center h-100 w-100">
    <div class="d-flex justify-center align-center" style="flex: 1">
      <v-icon icon="mdi-file-music" style="font-size: 15em" />
    </div>
    <audio
      class="align-end"
      ref="playerRef"
      style="width: 100vw; align-self: flex-end"
      controls
      autoplay
    >
      <source :src="props.src.src" :type="props.src.type" />
    </audio>
  </div>
</template>
