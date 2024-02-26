<script setup lang="ts">
import { computed } from "vue";
import { getDailyMedia, monthLong } from "@/js/date";
import DayMediaGrid from "./DayMediaGrid.vue";

const props = defineProps<{
  month: number;
  year: number;
  mediaList: Array<Media>;
  loadAllMediaUntil: (date: Date) => Promise<any>;
  mediaDateGetter: (media: Media) => Date;
  selectedMediaMap: Map<number, boolean>;
}>();

const emits = defineEmits<{
  (e: "thumbnailClick", mediaID: number): void;
  (e: "selectMedia", mediaID: number, value: boolean): void;
  (
    e: "selectMediaForDate",
    date: number,
    month: number,
    year: number,
    value: boolean,
  ): void;
}>();

const dailyMediaList = computed<Array<DailyMedia>>(() =>
  getDailyMedia(props.mediaList, props.mediaDateGetter),
);
</script>

<template>
  <v-card class="bg-secondary-background mt-2">
    <v-card-title style="font-size: 1.75rem; font-weight: 400">
      {{ `${monthLong[props.month]} ${props.year}` }}
    </v-card-title>
    <v-divider class="mb-2" :thickness="2" />
    <div class="d-flex flex-row flex-wrap">
      <DayMediaGrid
        v-for="(dailyMedia, index) in dailyMediaList"
        :key="`${dailyMedia.date}_${dailyMedia.month}_${dailyMedia.year}_${dailyMedia.media.length}_${index}`"
        :selected-media-map="props.selectedMediaMap"
        :month="dailyMedia.month"
        :day="dailyMedia.day"
        :year="dailyMedia.year"
        :date="dailyMedia.date"
        :media-list="dailyMedia.media"
        :load-all-media-until="props.loadAllMediaUntil"
        @thumbnail-click="
          (clickedMediaID) => emits('thumbnailClick', clickedMediaID)
        "
        @select-media="(mediaID, value) => emits('selectMedia', mediaID, value)"
        @select-all-media="
          (value) =>
            emits(
              'selectMediaForDate',
              dailyMedia.date,
              dailyMedia.month,
              dailyMedia.year,
              value,
            )
        "
      />
    </div>
  </v-card>
</template>
