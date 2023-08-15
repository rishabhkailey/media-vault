<script setup lang="ts">
import { computed } from "vue";
import { getDailyMediaIndex, monthLong } from "@/js/date";
import DailyThumbnailPreview from "./DailyThumbnailPreview.vue";
const props = defineProps<{
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  indexOffset: number;
  loadAllMediaOfDate: (date: Date) => Promise<any>;
  mediaDateGetter: (media: Media) => Date;
}>();

const emits = defineEmits<{
  (e: "thumbnailClick", mediaID: number, index: number): void;
}>();

const dailyMediaList = computed<Array<DailyMedia>>(() =>
  getDailyMediaIndex(props.indexMediaList, props.mediaDateGetter)
);
</script>

<template>
  <v-card class="bg-secondary-background">
    <v-card-title>{{ `${monthLong[props.month]} ${props.year}` }}</v-card-title>
    <v-divider :thickness="2" />
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
        @thumbnail-click="
          (clickedMediaID, clickedIndex) =>
            emits('thumbnailClick', clickedMediaID, clickedIndex)
        "
      />
    </div>
  </v-card>
</template>
@/js/date
