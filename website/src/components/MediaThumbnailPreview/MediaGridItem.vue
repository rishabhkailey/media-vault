<script setup lang="ts">
import { computed, ref } from "vue";

const props = defineProps<{
  media: Media;
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
  console.log(containerElement.value);
  return containerElement.value.getBoundingClientRect();
});
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
  <v-container ref="containerElement" class="pa-0 ma-0">
    <v-hover v-slot="{ isHovering, props: hoverProps }">
      <v-card
        @click="
          () => {
            emit('click');
          }
        "
        class="w-100 h-100"
        :style="[`padding: ${props.padding}px`, 'aspec-ratio: 1']"
        v-bind="hoverProps"
        :elevation="isHovering ? 6 : 0"
      >
        <v-img
          v-if="props.media.thumbnail"
          :src="props.media.thumbnail_url"
          class="w-100"
          :aspect-ratio="props.aspectRatio"
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
                :style="`aspect-ratio: 1; font-size: ${contianerSize.height}px`"
                :icon="getIcon(props.media.type)"
              />
            </div>
          </template>
        </v-img>
        <v-container
          v-else
          class="pa-0 ma-0 w-100 h-100 d-flex align-center justify-center fill-height"
        >
          <v-icon
            :style="`aspect-ratio: 1; font-size: ${contianerSize.height}px`"
            :icon="getIcon(props.media.type)"
          />
        </v-container>
      </v-card>
    </v-hover>
  </v-container>
</template>
