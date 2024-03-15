<script setup lang="ts">
// import { computed } from "vue";
import { useDisplay } from "vuetify";

const props = defineProps<{
  src: { src: string; type: string };
}>();
// const sources = computed(() => [props.src]);
const display = useDisplay();
</script>

<template>
  <div>
    <video
      style="max-width: 100vw; max-height: 100vh"
      controls
    >
      <source :src="props.src.src" :type="props.src.type" />
    </video>
    <!-- <video-player
      style="max-height: 90vh; max-width: 90vw"
      :sources="sources"
      :controls="true"
      :autoplay="true"
      :loop="true"
      :volume="0.6"
      :responsive="true"
      :liveui="true"
    /> -->
    <!-- :fluid="true" -->
  </div>
</template>
