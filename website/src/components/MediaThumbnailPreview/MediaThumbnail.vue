<script setup lang="ts">
const props = defineProps<{
  media: Media;
  class: any;
  width: number;
  height: number;
  transition: string | boolean;
}>();
</script>

<template>
  <v-img
    v-if="media.thumbnail"
    :id="`thumbnail_${media.id}`"
    :src="media.thumbnail_url"
    :width="width"
    :height="height"
    :transition="transition"
    :class="props.class"
    cover
  />
  <div
    v-else
    :style="`height: ${height}px; width: ${width}px;`"
    :class="props.class"
  >
    <v-icon
      icon="mdi-file-document-outline"
      :style="`font-size: ${width}px;`"
    />
  </div>
</template>
