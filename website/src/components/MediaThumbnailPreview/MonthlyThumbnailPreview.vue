<script setup lang="ts">
import { computed, inject } from "vue";
import { getDailyMediaIndex, monthLong } from "@/utils/date";
import DailyThumbnailPreview from "./DailyThumbnailPreview.vue";
import { mediaDateGetterKey } from "@/symbols/injectionSymbols";
const props = defineProps<{
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  indexOffset: number;
  loadAllMediaOfDate: (date: Date) => Promise<any>;
}>();
const mediaDateGetter = inject<(media: Media) => Date>(mediaDateGetterKey);
if (mediaDateGetter === undefined) {
  // todo error message
  throw new Error();
}
const dailyMediaList = computed<Array<DailyMedia>>(() =>
  getDailyMediaIndex(props.indexMediaList, mediaDateGetter)
);
</script>

<template>
  <v-card class="bg-secondary-background">
    <v-card-title>{{ `${monthLong[props.month]} ${props.year}` }}</v-card-title>
    <div class="d-flex flex-row flex-wrap">
      <DailyThumbnailPreview
        v-for="(dailyMedia, index) in dailyMediaList"
        :key="index"
        :month="dailyMedia.month"
        :day="dailyMedia.day"
        :year="dailyMedia.year"
        :date="dailyMedia.date"
        :index-media-list="dailyMedia.media"
        :load-all-media-of-date="props.loadAllMediaOfDate"
      />
    </div>
  </v-card>
</template>
