<script setup lang="ts">
import LazyMediaThumbnailsPreview from "@/components/MediaThumbnailPreview/LazyMediaThumbnailsPreview.vue";
import { useMediaStore } from "@/piniaStore/media";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
const mediaStore = useMediaStore();
const { mediaList, allMediaLoaded } = storeToRefs(mediaStore);
console.log(mediaStore);
const { loadMoreMedia } = mediaStore;
const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

async function loadAllMediaOfDate(date: Date): Promise<boolean> {
  let lastMediaDate = mediaList.value[mediaList.value.length - 1].date;
  console.log(date, lastMediaDate);
  while (
    date.getDate() === lastMediaDate.getDate() &&
    date.getFullYear() === lastMediaDate.getFullYear() &&
    date.getMonth() === lastMediaDate.getMonth() &&
    !allMediaLoaded.value
  ) {
    await loadMoreMedia(accessToken.value);
    lastMediaDate = mediaList.value[mediaList.value.length - 1].date;
  }
  return true;
}
</script>

<template>
  <LazyMediaThumbnailsPreview
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreMedia(accessToken)"
    :load-all-media-of-date="loadAllMediaOfDate"
  />
</template>
