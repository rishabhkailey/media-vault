<script setup lang="ts">
const props = defineProps<{
  media: Media;
  padding: number;
  aspectRatio: number;
  width: number;
  height: number;
}>();
const emit = defineEmits<{
  (e: "click"): void;
}>();

const getIcon = (mediaType: string) => {
  switch (mediaType.split("/")[0]) {
    case "image":
      return "mdi-image-outline";
    case "video":
      return "mdi-video-outline";
    default:
      return "mdi-file-outline";
  }
};
</script>
<template>
  <v-container
    class="pa-0 ma-0"
    :style="`aspect-ratio: 1; width: ${props.width}px`"
  >
    <v-hover v-slot="{ isHovering, props: hoverProps }">
      <v-card
        @click="
          () => {
            emit('click');
          }
        "
        v-bind="hoverProps"
        :elevation="isHovering ? 6 : 0"
      >
        <v-img
          v-if="props.media.thumbnail"
          :src="props.media.thumbnail_url"
          :width="props.width"
          :height="props.height"
          :aspect-ratio="props.aspectRatio"
          transition="scale"
          cover
        >
          <template v-slot:error>
            <div class="d-flex align-center justify-center fill-height">
              <v-icon
                :style="`aspect-ratio: 1; font-size: ${props.width}px`"
                :icon="getIcon(props.media.type)"
              />
            </div>
          </template>
        </v-img>
        <v-icon
          v-else
          :style="`aspect-ratio: 1; font-size: ${props.width}px`"
          :icon="getIcon(props.media.type)"
        />
      </v-card>
    </v-hover>
  </v-container>
</template>
