<script setup lang="ts">
import { ref, watch } from "vue";
const container = ref<HTMLElement | undefined>(undefined);
const height = ref(0);
const width = ref(0);

watch(container, () => {
  if (container.value === undefined) {
    return;
  }
  const dimensions = container.value.getBoundingClientRect();
  width.value = dimensions.width;
  height.value = dimensions.height;
});
</script>
<template>
  <div ref="container">
    <slot :height="height" :width="width"></slot>
  </div>
</template>
