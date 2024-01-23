<script setup lang="ts">
import axios from "axios";
import { computed, watch, onBeforeUnmount } from "vue";
import { useAuthStore } from "@/piniaStore/auth";
import MonthMediaGrid from "./MonthMediaGrid.vue";
import InfiniteScroller from "@/components/InfiniteScroller/InfiniteScroller.vue";
import { getMonthlyMediaIndex } from "@/js/date";
import { storeToRefs } from "pinia";
import { useLoadingStore } from "@/piniaStore/loading";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";

const props = defineProps<{
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: () => Promise<any>;
  loadAllMediaUntil: (date: Date) => Promise<any>;
  mediaDateGetter: (media: Media) => Date;
}>();

const emits = defineEmits<{
  (
    e: "thumbnailClick",
    mediaID: number,
    index: number,
    clickLocation: ThumbnailClickLocation | undefined,
  ): void;
}>();

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
console.log(authStore);
const monthlyMediaList = computed<Array<MonthlyMedia>>(() =>
  getMonthlyMediaIndex(props.mediaList, props.mediaDateGetter),
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
      },
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
        <MonthMediaGrid
          :month="monthlyMedia.month"
          :year="monthlyMedia.year"
          :index-media-list="monthlyMedia.media"
          :index-offset="0"
          :load-all-media-until="props.loadAllMediaUntil"
          :media-date-getter="props.mediaDateGetter"
          @thumbnail-click="
            (clickedMediaID, clickedIndex, clickLocation) =>
              emits(
                'thumbnailClick',
                clickedMediaID,
                clickedIndex,
                clickLocation,
              )
          "
        />
        <v-divider />
      </div>
      <InfiniteScroller
        v-if="allMediaLoaded === false"
        :on-threshold-reach="loadMoreMedia"
        :threshold="0.1"
        :min-height="100"
        :min-width="100"
        :root-margin="10"
      ></InfiniteScroller>
    </div>
  </div>
</template>
@/js/date
