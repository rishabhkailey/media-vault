<script setup lang="ts">
import MediaGrid from "@/components/MediaThumbnailPreview/MediaGrid.vue";
import { onBeforeMount, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useSearchStore } from "@/piniaStore/search";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
import { getQueryParamStringValue } from "@/js/utils";
import { searchMediaPreviewRoute } from "@/router/routesConstants";
import { MEDIA_PREVIEW_CONTAINER_Z_INDEX } from "@/js/constants/z-index";
import ErrorMessage from "../Error/ErrorMessage.vue";
const errorMessage = ref("");
const router = useRouter();

const route = useRoute();
function getSearchQuery(): string {
  const queryParam = getQueryParamStringValue(route.params, "query");
  if (queryParam === undefined) {
    errorMessage.value = "Invalid search query";
    return "";
  }
  return queryParam;
}

const { accessToken } = storeToRefs(useAuthStore());
const searchStore = useSearchStore();
const { mediaList, allMediaLoaded } = storeToRefs(searchStore);
const { loadMoreSearchResults, setQuery, getMediaDateAccordingToOrderBy } =
  searchStore;
const searchQuery = ref("");

onBeforeMount(async () => {
  await router.isReady();
  // as we are using global store for search results, it can still have results of old media search
  // this will ensure to update search query and results in store
  searchQuery.value = getSearchQuery();
  setQuery(searchQuery.value);
});

function loadAllMediaUntil(date: Date): Promise<boolean> {
  return new Promise((resolve, reject) => {
    let lastMedia = mediaList.value[mediaList.value.length - 1];
    while (date < lastMedia.date && !allMediaLoaded.value) {
      loadMoreSearchResults(accessToken.value, searchQuery.value)
        .then(() => {
          lastMedia = mediaList.value[mediaList.value.length - 1];
        })
        .catch((err) => {
          reject(err);
          return;
        });
    }
    resolve(true);
    return;
  });
}

function handleThumbnailClick(clickedMediaID: number) {
  try {
    const clickedIndex = mediaList.value.findIndex(
      (m) => m.id === clickedMediaID,
    );
    router.push(
      searchMediaPreviewRoute(clickedIndex, clickedMediaID, searchQuery.value),
    );
  } catch (err) {
    // todo error page?
    console.error("error in homepage", err);
  }
}
</script>

<template>
  <ErrorMessage
    v-if="errorMessage.length > 0"
    :message="errorMessage"
    title="Something went wrong"
    :expanded-by-default="true"
    @close="() => (errorMessage = '')"
  />
  <MediaGrid
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreSearchResults(accessToken, searchQuery)"
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
