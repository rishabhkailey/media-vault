<script setup lang="ts">
import { computed, ref } from "vue";

const props = defineProps<{
  album: Album;
  padding: number;
  aspectRatio: number;
}>();
const emit = defineEmits<{
  (e: "click"): void;
}>();

const containerElement = ref<HTMLElement | undefined>(undefined);
const contianerSize = computed<{
  height: number;
  width: number;
}>(() => {
  if (containerElement.value === undefined) {
    return {
      height: 150,
      width: 150,
    };
  }
  return containerElement.value.getBoundingClientRect();
});

const imageElement = ref<HTMLElement | undefined>(undefined);
</script>
<template>
  <div ref="containerElement">
    <v-hover v-slot="{ isHovering, props: hoverProps }">
      <v-card
        class="flex-grow-1"
        @click="
          () => {
            emit('click');
          }
        "
        :style="`padding: ${props.padding}px`"
        v-bind="hoverProps"
        :elevation="isHovering ? 12 : 6"
      >
        <v-img
          v-if="props.album.thumbnail_url.length > 0"
          :src="props.album.thumbnail_url"
          class="w-100"
          :aspect-ratio="props.aspectRatio"
          ref="imageElement"
          transition="scale"
          cover
        >
          <template v-slot:error>
            <div class="d-flex align-center justify-center fill-height">
              <v-icon
                icon="mdi-image-broken-variant"
                :style="`aspect-ratio: 1; font-size: ${contianerSize.width}px`"
              />
            </div>
          </template>
        </v-img>
        <div
          v-else
          class="d-flex align-center justify-center fill-height w-100"
          :style="`aspect-ratio: ${props.aspectRatio}; overflow: hidden;`"
        >
          <v-icon
            icon="mdi-image-album"
            :style="{
              fontSize: `${contianerSize.height}px`,
            }"
          />
        </div>
        <v-card-subtitle>
          <div class="d-flex flex-column align-start">
            <v-chip size="x-small" variant="text" class="ma-0 pa-0"
              >{{ album.media_count ?? "0" }} items</v-chip
            >
            <div>
              {{ album.name.length != 0 ? album.name : "!" }}
            </div>
          </div>
        </v-card-subtitle>
      </v-card>
    </v-hover>
  </div>
</template>
