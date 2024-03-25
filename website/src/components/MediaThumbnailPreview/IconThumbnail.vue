<script setup lang="ts">
const props = defineProps<{
  icon: string;
  width: number;
  height: number;
  fileName: string;
}>();
</script>

<template>
  <div
    class="d-flex flex-column justify-center align-center w-100 h-100 bg-secondary-background"
    style="overflow: hidden"
  >
    <div
      class="d-flex flex-column justify-center align-center border-radius: 5em bg-background w-100 mb-1"
      style="flex: 1; border-radius: 1em"
    >
      <v-icon :icon="props.icon" :style="`font-size: ${props.width - 50}px;`" />
    </div>
    <span class="text-no-wrap" style="width: 100%; font-size: 0.8em">{{
      props.fileName
    }}</span>
  </div>
</template>
