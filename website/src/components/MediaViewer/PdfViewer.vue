<script setup lang="ts">
const props = defineProps<{
  src: string;
}>();
</script>

<template>
  <div class="pdf-viewer-container">
    <iframe
      :src="props.src"
      style="width: 100%; height: 100%; border: none"
      type="application/pdf"
    >
    </iframe>
  </div>
</template>
<style scoped>
.pdf-viewer-container {
  /* equals to media viewer header */
  top: 4em;
  width: 100%;
  height: 100%;
}
</style>
