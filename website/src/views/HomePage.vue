<script setup lang="ts">
import MediaGrid from "@/components/MediaThumbnailPreview/MediaGrid.vue";
import { base64UrlEncode } from "@/js/utils";
import { useMediaStore } from "@/piniaStore/media";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";

const router = useRouter();
const mediaStore = useMediaStore();
const { mediaList, allMediaLoaded } = storeToRefs(mediaStore);
const { loadMoreMedia, getMediaDateAccordingToOrderBy } = mediaStore;

async function loadAllMediaUntil(date: Date): Promise<boolean> {
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
  <MediaGrid
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreMedia()"
    :load-all-media-until="loadAllMediaUntil"
    :media-date-getter="getMediaDateAccordingToOrderBy"
    @thumbnail-click="
      (clickedMediaID, clickedIndex, thumbnailClickLocation) => {
        router.push({
          name: `MediaPreview`,
          params: {
            index: clickedIndex,
            media_id: clickedMediaID,
          },
          hash: `#${base64UrlEncode(thumbnailClickLocation)}`,
        });
      }
    "
  />
  <Teleport to="body">
    <div style="position: absolute; top: 0px; left: 0px; z-index: 9999999">
      <RouterView />
    </div>
  </Teleport>
</template>
