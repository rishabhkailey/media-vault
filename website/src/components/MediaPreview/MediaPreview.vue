<script setup lang="ts">
import ImagePreview from "@/components/ImagePreview.vue";
import VideoPreview from "@/components/VideoPreview.vue";
import { download } from "@/js/encryptedFileDownload";
import { computed } from "vue";

const props = defineProps<{
  index: number;
  mediaList: Array<Media>;
  loadMoreMedia: () => Promise<Boolean>;
  allMediaLoaded: boolean;
  routeName: string;
}>();

const emits = defineEmits<{
  (e: "close"): void;
  (e: "next"): void;
  (e: "previous"): void;
  (e: "update:index", value: number): void;
}>();

const media = computed<Media>(() => {
  return props.mediaList[props.index];
});
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
        style="height: 100%; width: 100%"
        :model-value="index"
        @update:model-value="(value) => emits('update:index', value)"
        :continuous="false"
        :show-arrows="true"
        touch
        direction="horizontal"
      >
        <v-window-item
          style="height: 100%; width: 100%"
          v-for="(media, index) in mediaList"
          :key="index"
        >
          <template v-slot:default>
            <div
              style="
                width: 100%;
                height: 100%;
                display: flex;
                justify-content: center;
                align-items: center;
              "
            >
              <div v-if="media.type.startsWith('image')">
                <ImagePreview
                  :src="media.url"
                  :low-res-src="media.thumbnail_url"
                />
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
      </v-window>
    </v-row>
  </v-col>
</template>
