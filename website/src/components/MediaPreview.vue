<script setup lang="ts">
import { ref } from "vue";
import ImagePreview from "./ImagePreview.vue";
import VideoPreview from "./VideoPreview.vue";
const props = defineProps<{
  index: number;
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: () => Promise<any>; // emit always returns void so we cannot use emit
}>();
const internalIndex = ref(props.index);
const media = ref(props.mediaList[props.index]);
const loadingMoreMedia = ref(false);
const next = () => {
  if (internalIndex.value <= props.mediaList.length - 2) {
    media.value = props.mediaList[++internalIndex.value];
    return;
  }
  if (props.allMediaLoaded) {
    return;
  }
  loadingMoreMedia.value = true;
  props
    .loadMoreMedia()
    .then(() => {
      loadingMoreMedia.value = false;
      if (internalIndex.value <= props.mediaList.length - 2) {
        media.value = props.mediaList[++internalIndex.value];
      }
    })
    .catch((err) => {
      loadingMoreMedia.value = false;
      console.log("error loaing more media ", err);
    });
};
</script>

<template>
  <v-carousel style="width: 100vw; height: 100vh" :cycle="false">
    <template v-slot:next>
      <!--  not sure why v-btn icon is not working, that's why we are using default slot -->
      <v-btn
        icon="mdi-chevron-right"
        @click="next"
        :disabled="
          (internalIndex === props.mediaList.length && props.allMediaLoaded) ||
          loadingMoreMedia
        "
      >
        <template v-slot:default>
          <v-progress-circular v-if="loadingMoreMedia" indeterminate />
          <v-icon v-else icon="mdi-chevron-right" size="x-large" />
        </template>
      </v-btn>
    </template>
    <template v-slot:prev>
      <!--  not sure why v-btn icon is not working, that's why we are using default slot -->
      <v-btn
        icon="mdi-chevron-left"
        :disabled="internalIndex === 0"
        @click="
          () => {
            media = props.mediaList[--internalIndex];
          }
        "
      >
        <template v-slot:default>
          <v-icon icon="mdi-chevron-left" size="x-large" />
        </template>
      </v-btn>
    </template>
    <v-carousel-item class="d-flex justify-center align-center">
      <template v-slot:default>
        <div class="d-flex justify-center align-center">
          <div v-if="media?.type.startsWith('image')">
            <ImagePreview :src="media.url" />
          </div>
          <div v-else-if="media?.type.startsWith('video')">
            <VideoPreview :src="{ src: media.url, type: media.type }" />
          </div>
          <div v-else>unknown</div>
          <div>{{ props.index }}</div>
        </div>
      </template>
    </v-carousel-item>
  </v-carousel>
</template>
