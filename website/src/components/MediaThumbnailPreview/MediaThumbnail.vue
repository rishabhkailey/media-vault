<script setup lang="ts">
import { isAudio, isPdf } from "@/js/files/type";
import IconThumbnail from "./IconThumbnail.vue";

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
      <IconThumbnail
        icon="mdi-image-broken-variant"
        :file-name="props.media.name"
        :width="props.width"
        :height="props.height"
      />
    </template>
  </v-img>
  <div
    v-else
    :style="`height: ${height}px; width: ${width}px;`"
    :id="`thumbnail_${media.id}`"
    :class="props.class"
  >
    <IconThumbnail
      v-if="isPdf(media.type)"
      icon="mdi-file-pdf-box"
      :file-name="props.media.name"
      :width="props.width"
      :height="props.height"
    />
    <IconThumbnail
      v-else-if="isAudio(media.type)"
      icon="mdi-file-music"
      :file-name="props.media.name"
      :width="props.width"
      :height="props.height"
    />
    <IconThumbnail
      v-else
      icon="mdi-file-document-outline"
      :file-name="props.media.name"
      :width="props.width"
      :height="props.height"
    />
  </div>
</template>
