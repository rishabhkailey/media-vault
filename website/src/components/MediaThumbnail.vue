<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  media: Media;
}>();
const emit = defineEmits<{
  (e: "click"): void;
}>();

const imageElement = ref<HTMLElement | undefined>(undefined);
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
  <v-hover v-slot="{ isHovering, props: hoverProps }">
    <v-card
      class="flex-grow-1 d-flex justify-content-strech"
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
        width="150"
        aspect-ratio="1"
        ref="imageElement"
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
              :icon="getIcon(props.media.type)"
              style="font-size: 150px"
            />
          </div>
        </template>
      </v-img>
      <div v-else>
        <div class="d-flex align-center justify-center fill-height">
          <v-icon :icon="getIcon(props.media.type)" style="font-size: 150px" />
        </div>
      </div>
    </v-card>
  </v-hover>
</template>
