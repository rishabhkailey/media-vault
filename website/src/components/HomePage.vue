<script setup lang="ts">
import LazyMediaThumbnailsPreview from "@/components/MediaThumbnailPreview/LazyMediaThumbnailsPreview.vue";
import { useMediaStore } from "@/piniaStore/media";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";

const router = useRouter();
const mediaStore = useMediaStore();
const { mediaList, allMediaLoaded } = storeToRefs(mediaStore);
console.log(mediaStore);
const { loadMoreMedia, getMediaDateAccordingToOrderBy } = mediaStore;

async function loadAllMediaOfDate(date: Date): Promise<boolean> {
  let lastMediaDate = mediaList.value[mediaList.value.length - 1].date;
  console.log(date, lastMediaDate);
  while (
    date.getDate() === lastMediaDate.getDate() &&
    date.getFullYear() === lastMediaDate.getFullYear() &&
    date.getMonth() === lastMediaDate.getMonth() &&
    !allMediaLoaded.value
  ) {
    await loadMoreMedia();
    lastMediaDate = mediaList.value[mediaList.value.length - 1].date;
  }
  return true;
}
</script>

<template>
  <LazyMediaThumbnailsPreview
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreMedia()"
    :load-all-media-of-date="loadAllMediaOfDate"
    :media-date-getter="getMediaDateAccordingToOrderBy"
    @thumbnail-click="
      (clickedMediaID, clickedIndex) => {
        router.push({
          name: `MediaPreview`,
          params: {
            index: clickedIndex,
            media_id: clickedMediaID,
          },
        });
      }
    "
  />
</template>
