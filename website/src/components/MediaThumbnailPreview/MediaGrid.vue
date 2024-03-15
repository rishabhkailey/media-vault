<script setup lang="ts">
import { computed, onBeforeUnmount } from "vue";
import MonthMediaGrid from "./MonthMediaGrid.vue";
import { getMonthlyMedia } from "@/js/date";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { storeToRefs } from "pinia";
import { useErrorsStore } from "@/piniaStore/errors";

const props = defineProps<{
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: LoadMoreMedia;
  mediaDateGetter: (media: Media) => Date;
}>();

const emits = defineEmits<{
  (e: "thumbnailClick", mediaID: number): void;
}>();

const monthlyMediaList = computed<Array<MonthlyMedia>>(() =>
  getMonthlyMedia(props.mediaList, props.mediaDateGetter),
);

const mediaSelectionStore = useMediaSelectionStore();
const { selectedMediaMap } = storeToRefs(mediaSelectionStore);

const { reset: resetMediaSelection } = useMediaSelectionStore();

function updateMediaSelection(mediaID: number, value: boolean) {
  mediaSelectionStore.updateSelection(mediaID, value);
}

const { appendError } = useErrorsStore();

/**
 * Fetches all media until the provided date
 * @returns {Promise<boolean>} A Promise that resolves with:
 *   - `true` if more media was successfully loaded.
 *   - Rejects with an error if media loading fails.
 */
async function loadAllMediaUntil(date: Date): Promise<boolean> {
  let lastMediaDate = props.mediaDateGetter(
    props.mediaList[props.mediaList.length - 1],
  );
  while (
    date.getDate() === lastMediaDate.getDate() &&
    date.getFullYear() === lastMediaDate.getFullYear() &&
    date.getMonth() === lastMediaDate.getMonth() &&
    !props.allMediaLoaded
  ) {
    await props.loadMoreMedia();
    lastMediaDate = props.mediaList[props.mediaList.length - 1].date;
  }
  return true;
}

/**
 * Fetches and select all media for a date
 */
function selectMediaForDate(
  date: number,
  month: number,
  year: number,
  value: boolean,
) {
  loadAllMediaUntil(new Date(`${year}-${month + 1}-${date}`))
    .then(() => {
      monthlyMediaList.value
        .find(
          (monthlyMediaList) =>
            monthlyMediaList.month === month && monthlyMediaList.year === year,
        )
        ?.media.filter((media) => {
          const mediaDate = props.mediaDateGetter(media);
          return (
            mediaDate.getDate() === date &&
            mediaDate.getMonth() === month &&
            mediaDate.getFullYear() === year
          );
        })
        ?.forEach((media) => {
          updateMediaSelection(media.id, value);
        });
    })
    .catch((err) => {
      appendError("failed to select all day's media", err, -1);
    });
}

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
            v-for="monthlyMedia in monthlyMediaList"
            :key="`${monthlyMedia.month}_${monthlyMedia.year}`"
          >
            <MonthMediaGrid
              :month="monthlyMedia.month"
              :year="monthlyMedia.year"
              :media-list="monthlyMedia.media"
              :load-all-media-until="loadAllMediaUntil"
              :media-date-getter="props.mediaDateGetter"
              :selected-media-map="selectedMediaMap"
              @select-media="updateMediaSelection"
              @select-media-for-date="selectMediaForDate"
              @thumbnail-click="
                (clickedMediaID) => emits('thumbnailClick', clickedMediaID)
              "
            />
            <v-divider />
          </template>
        </template>
      </v-infinite-scroll>
    </div>
  </div>
</template>
