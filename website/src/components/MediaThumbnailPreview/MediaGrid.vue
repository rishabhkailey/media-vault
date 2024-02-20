<script setup lang="ts">
import { computed, onBeforeUnmount } from "vue";
import MonthMediaGrid from "./MonthMediaGrid.vue";
import { getMonthlyMediaIndex } from "@/js/date";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";

const props = defineProps<{
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: LoadMoreMedia;
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

const monthlyMediaList = computed<Array<MonthlyMedia>>(() =>
  getMonthlyMediaIndex(props.mediaList, props.mediaDateGetter),
);

const { reset: resetMediaSelection } = useMediaSelectionStore();

onBeforeUnmount(() => {
  resetMediaSelection();
});
</script>

<template>
  <div class="pa-0 ma-0">
    <div class="d-flex flex-column align-stretch">
      <v-infinite-scroll
        style="overflow-y: hidden"
        class="bg-secondary-background"
        :items="monthlyMediaList"
        side="end"
        @load="
          ({ done }) => {
            loadMoreMedia()
              .then((status) => {
                done(status);
              })
              .catch((_) => done('error'));
          }
        "
      >
        <template #error> failed to load data from server </template>
        <template #default>
          <template
            v-for="(monthlyMedia, index) in monthlyMediaList"
            :key="index"
          >
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
          </template>
        </template>
      </v-infinite-scroll>
    </div>
  </div>
</template>
