<script setup lang="ts">
import { computed } from "vue";
import { useDisplay } from "vuetify";
const props = defineProps<{
  thumbnailsMetadata: Array<{
    url: string;
    title: string;
    subTitle: string;
  }>;
  availableWidth: number;
  aspectRatio: number;
  xs: number;
  sm: number;
  md: number;
  lg: number;
  xl: number;
  xxl: number;
  cols: number;
}>();
const display = useDisplay();
const width = computed<number>(() => {
  switch (display.name.value) {
    case "xs":
      return display.width.value / props.xs;
    case "sm":
      return display.width.value / props.sm;
    case "md":
      return display.width.value / props.md;
    case "lg":
      return display.width.value / props.lg;
    case "xl":
      return display.width.value / props.xl;
    case "xxl":
      return display.width.value / props.xxl;
    default:
      return display.width.value / props.cols;
  }
});
</script>
<template>
  <div class="d-flex flex-row flew-wrap">
    <div
      class="d-flex flex-column"
      v-for="({ url, title, subTitle }, index) in thumbnailsMetadata"
      :key="index"
    >
      <slot :title="title" :subTitle="subTitle" name="top"></slot>

      <v-img :width="width" :src="url" :aspect-ratio="props.aspectRatio" />
      <slot :title="title" :subTitle="subTitle" name="bottom"></slot>
    </div>
  </div>
</template>
