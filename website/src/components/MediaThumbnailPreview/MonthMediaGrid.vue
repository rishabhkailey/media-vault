<script setup lang="ts">
import { computed } from "vue";
import { getDailyMediaIndex, monthLong } from "@/js/date";
import DayMediaGrid from "./DayMediaGrid.vue";
const props = defineProps<{
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  indexOffset: number;
  loadAllMediaUntil: (date: Date) => Promise<any>;
  // mediaDateGetter returns the date of media using which the media is sorted. it can be media date, upload date or any other date.
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

const dailyMediaList = computed<Array<DailyMedia>>(() =>
  getDailyMediaIndex(props.indexMediaList, props.mediaDateGetter),
);
</script>

<template>
  <v-card class="bg-secondary-background">
    <v-card-title>{{ `${monthLong[props.month]} ${props.year}` }}</v-card-title>
    <v-divider class="mb-2" :thickness="2" />
    <div class="d-flex flex-row flex-wrap">
      <DayMediaGrid
        v-for="(dailyMedia, index) in dailyMediaList"
        :key="index"
        :month="dailyMedia.month"
        :day="dailyMedia.day"
        :year="dailyMedia.year"
        :date="dailyMedia.date"
        :index-media-list="dailyMedia.media"
        :load-all-media-until="props.loadAllMediaUntil"
        @thumbnail-click="
          (clickedMediaID, clickedIndex, clickLocation) =>
            emits('thumbnailClick', clickedMediaID, clickedIndex, clickLocation)
        "
      />
    </div>
  </v-card>
</template>
@/js/date
