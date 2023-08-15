<script setup lang="ts">
import axios from "axios";
import { computed, watch, provide, onBeforeUnmount } from "vue";
import { useAuthStore } from "@/piniaStore/auth";
import {
  loadMoreMediaKey,
  allMediaLoadedKey,
  mediaListKey,
  mediaDateGetterKey,
} from "@/symbols/injectionSymbols";
import MonthlyThumbnailPreview from "./MonthlyThumbnailPreview.vue";
import LazyLoading from "@/components/LazyLoading/LazyLoading.vue";
import { getMonthlyMediaIndex } from "@/js/date";
import { storeToRefs } from "pinia";
import { useLoadingStore } from "@/piniaStore/loading";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";

const props = defineProps<{
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: () => Promise<any>;
  loadAllMediaOfDate: (date: Date) => Promise<any>;
  mediaDateGetter: (media: Media) => Date;
}>();

provide(
  mediaListKey,
  computed(() => props.mediaList)
);
provide(
  allMediaLoadedKey,
  computed(() => props.allMediaLoaded)
);
provide(loadMoreMediaKey, props.loadMoreMedia);
provide(mediaDateGetterKey, props.mediaDateGetter);

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
console.log(authStore);
const monthlyMediaList = computed<Array<MonthlyMedia>>(() =>
  getMonthlyMediaIndex(props.mediaList, props.mediaDateGetter)
);

const { initializing } = storeToRefs(useLoadingStore());
watch(initializing, async (newValue, oldValue) => {
  console.log("initializing changed to ", newValue);
  if (newValue === oldValue) {
    return;
  }
  if (!newValue) {
    let response = await axios.post(
      "/v1/refreshSession",
      {},
      {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      }
    );
    console.log(response);
  }
});

const { reset: resetMediaSelection } = useMediaSelectionStore();

onBeforeUnmount(() => {
  resetMediaSelection();
});
</script>

<template>
  <div class="pa-0 ma-0">
    <div v-if="initializing">loading...</div>
    <div v-else class="d-flex flex-column align-stretch">
      <div v-for="(monthlyMedia, index) in monthlyMediaList" :key="index">
        <MonthlyThumbnailPreview
          :month="monthlyMedia.month"
          :year="monthlyMedia.year"
          :index-media-list="monthlyMedia.media"
          :index-offset="0"
          :load-all-media-of-date="props.loadAllMediaOfDate"
        />
        <v-divider />
      </div>
      <LazyLoading
        v-if="allMediaLoaded === false"
        :on-threshold-reach="loadMoreMedia"
        :threshold="0.1"
        :min-height="100"
        :min-width="100"
        :root-margin="10"
      ></LazyLoading>
    </div>
  </div>
</template>
@/js/date