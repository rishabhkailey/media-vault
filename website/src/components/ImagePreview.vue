<script setup lang="ts">
const props = defineProps<{
  src: string;
}>();
</script>

<template>
  <div>
    <img style="max-width: 100vw; max-height: 100vh" :src="props.src" />
  </div>
</template>
