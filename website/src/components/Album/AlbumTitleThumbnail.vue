<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  album: Album;
  height: number;
  width: number;
  padding: number;
  aspectRatio: number;
}>();
const emit = defineEmits<{
  (e: "click"): void;
}>();

const imageElement = ref<HTMLElement | undefined>(undefined);
</script>
<template>
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
      :elevation="isHovering ? 6 : 0"
    >
      <v-img
        v-if="props.album.thumbnail_url.length > 0"
        :src="props.album.thumbnail_url"
        :width="props.width - 2 * props.padding"
        :height="props.height - 2 * props.padding"
        :aspect-ratio="props.aspectRatio"
        ref="imageElement"
        transition="scale"
        cover
      >
        <!-- <template v-slot:placeholder>
        <div class="d-flex align-center justify-center fill-height">
          <v-progress-circular
            color="grey-lighten-4"
            indeterminate
          ></v-progress-circular>
        </div>
      </template> -->

        <template v-slot:error>
          <div class="d-flex align-center justify-center fill-height">
            <v-icon
              icon="mdi-image-broken-variant"
              :style="`font-size: ${props.width - 2 * props.padding}px`"
            />
          </div>
        </template>
      </v-img>
      <div v-else>
        <div class="d-flex align-center justify-center fill-height">
          <v-icon
            icon="mdi-image-album"
            :style="`font-size: ${props.width - 2 * props.padding}px`"
          />
        </div>
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
</template>
