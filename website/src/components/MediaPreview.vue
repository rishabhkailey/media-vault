<script setup lang="ts">
import { ref } from "vue";
import ImagePreview from "./ImagePreview.vue";
import VideoPreview from "./VideoPreview.vue";
import { download } from "@/utils/encryptedFileDownload";
const props = defineProps<{
  index: number;
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: () => Promise<any>; // emit always returns void so we cannot use emit
}>();

const emits = defineEmits<{
  (e: "close"): void;
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
  <v-col
    class="d-flex flex-column justify-center align-stretch"
    style="width: 100vw; height: 100vh; background-color: black"
  >
    <!-- header -->
    <div
      class="pt-1 pr-4 d-flex justify-end align-center"
      style="
        background-color: rgba(50, 50, 50, 0.3);
        box-shadow: 0px 5px 15px rgb(50, 50, 50, 0.3);
        position: absolute;
        width: 100vw;
        top: 0;
        z-index: 5000;
      "
    >
      <v-btn
        icon="mdi-download"
        @click.stop="
          () => {
            download(media.url, media.name);
          }
        "
        style="background: none; border: none; box-shadow: none"
      >
        <v-icon color="white">mdi-download</v-icon>
      </v-btn>
      <v-btn
        icon="mdi-close"
        @click.stop="
          () => {
            emits('close');
          }
        "
        style="background: none; border: none; box-shadow: none"
      >
        <v-icon color="white">mdi-close</v-icon>
      </v-btn>
    </div>
    <v-row class="flex-grow-1 d-flex flex-column justify-center align-center">
      <v-window
        :continuous="false"
        :show-arrows="true"
        v-model="internalIndex"
        touch
      >
        <template v-slot:next>
          <!--  not sure why v-btn icon is not working, that's why we are using default slot -->
          <v-btn
            icon="mdi-chevron-right"
            @click="next"
            :disabled="
              (internalIndex === props.mediaList.length &&
                props.allMediaLoaded) ||
              loadingMoreMedia
            "
          >
            <template v-slot:default>
              <v-progress-circular v-if="loadingMoreMedia" indeterminate />
              <v-icon v-else icon="mdi-chevron-right" />
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
        <!-- previous required for left arrow -->
        <v-window-item v-if="internalIndex !== 0" :value="internalIndex - 1">
          <template v-slot:default>
            <div />
          </template>
        </v-window-item>
        <!-- current -->
        <v-window-item
          style="width: 100vw"
          class="d-flex justify-center align-center"
          :value="internalIndex"
        >
          <template v-slot:default>
            <div>
              <div v-if="media.type.startsWith('image')">
                <ImagePreview :src="media.url" />
              </div>
              <div v-else-if="media.type.startsWith('video')">
                <VideoPreview
                  :src="{
                    src: media.url,
                    type: media.type,
                  }"
                />
              </div>
              <div v-else>unknown</div>
            </div>
          </template>
        </v-window-item>
        <!-- next required for left arrow -->
        <v-window-item
          v-if="internalIndex !== props.mediaList.length - 1"
          :value="internalIndex + 1"
        >
          <template v-slot:default>
            <div></div>
          </template>
        </v-window-item>
      </v-window>
    </v-row>
    <!-- footer
    <v-row class="d-flex justify-center align-center flex-grow-1"> </v-row> -->
  </v-col>
</template>
