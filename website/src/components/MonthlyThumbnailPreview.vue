<script setup lang="ts">
import { computed } from "vue";
import { getDailyMediaIndex, monthLong } from "@/utils/date";
import DailyThumbnailPreview from "./DailyThumbnailPreview.vue";
const props = defineProps<{
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  indexOffset: number;
}>();

const dailyMediaList = computed<Array<DailyMedia>>(() =>
  getDailyMediaIndex(props.indexMediaList)
);
</script>

<template>
  <v-card :title="`${monthLong[props.month]} ${props.year}`">
    <div class="d-flex flex-row flex-wrap">
      <DailyThumbnailPreview
        v-for="(dailyMedia, index) in dailyMediaList"
        :key="index"
        :month="dailyMedia.month"
        :day="dailyMedia.day"
        :year="dailyMedia.year"
        :date="dailyMedia.date"
        :index-media-list="dailyMedia.media"
      />
    </div>
  </v-card>
</template>
