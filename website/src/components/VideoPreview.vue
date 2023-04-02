<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  src: { src: string; type: string };
}>();
const sources = computed(() => [props.src]);
</script>

<template>
  <div>
    <video-player
      style="max-height: 90vh; max-width: 90vw"
      :sources="sources"
      :controls="true"
      :autoplay="true"
      :loop="true"
      :volume="0.6"
    />
  </div>
</template>
