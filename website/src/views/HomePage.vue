<script setup lang="ts">
import MediaGrid from "@/components/MediaThumbnailPreview/MediaGrid.vue";
import { useMediaStore } from "@/piniaStore/media";
import { mediaPreviewRoute } from "@/router/routesConstants";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { MEDIA_PREVIEW_CONTAINER_Z_INDEX } from "@/js/constants/z-index";

const router = useRouter();
const mediaStore = useMediaStore();
const { mediaList, allMediaLoaded } = storeToRefs(mediaStore);
const { loadMoreMedia, getMediaDateAccordingToOrderBy } = mediaStore;

async function loadAllMediaUntil(date: Date): Promise<boolean> {
  let lastMediaDate = getMediaDateAccordingToOrderBy(
    mediaList.value[mediaList.value.length - 1],
  );
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

function handleThumbnailClick(clickedMediaID: number) {
  try {
    const clickedIndex = mediaList.value.findIndex(
      (m) => m.id === clickedMediaID,
    );
    router.push(mediaPreviewRoute(clickedIndex, clickedMediaID));
  } catch (err) {
    // todo error page?
    console.error("error in homepage", err);
  }
}
</script>

<template>
  <MediaGrid
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreMedia()"
    :load-all-media-until="loadAllMediaUntil"
    :media-date-getter="getMediaDateAccordingToOrderBy"
    @thumbnail-click="handleThumbnailClick"
  />
  <Teleport to="body">
    <div class="media-preview-container">
      <RouterView />
    </div>
  </Teleport>
</template>

<style scoped>
.media-preview-container {
  position: absolute;
  top: 0px;
  left: 0px;
  z-index: v-bind(MEDIA_PREVIEW_CONTAINER_Z_INDEX);
}
</style>
