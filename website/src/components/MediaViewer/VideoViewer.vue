<script setup lang="ts">
// import { computed } from "vue";
// import { useDisplay } from "vuetify";

const props = defineProps<{
  src: { src: string; type: string };
}>();
// const sources = computed(() => [props.src]);
// const display = useDisplay();
</script>

<template>
  <video style="max-width: 100vw; max-height: 100vh" controls autoplay>
    <source :src="props.src.src" :type="props.src.type" />
  </video>
</template>
