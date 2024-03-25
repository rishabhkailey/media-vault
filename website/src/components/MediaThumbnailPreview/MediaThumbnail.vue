<script setup lang="ts">
import { isPdf } from "@/js/files/type";

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
  >
    <template #error>
      <v-icon
        icon="mdi-image-broken-variant"
        :style="`font-size: ${width}px;`"
      />
    </template>
  </v-img>
  <div
    v-else
    :style="`height: ${height}px; width: ${width}px;`"
    :id="`thumbnail_${media.id}`"
    :class="props.class"
  >
    <v-icon
      v-if="isPdf(media.type)"
      icon="mdi-file-pdf-box"
      :style="`font-size: ${width}px;`"
    />
    <v-icon
      v-else
      icon="mdi-file-document-outline"
      :style="`font-size: ${width}px;`"
    />
  </div>
</template>
