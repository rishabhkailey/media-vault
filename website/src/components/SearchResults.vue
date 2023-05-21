<script setup lang="ts">
import LazyMediaThumbnailsPreview from "@/components/LazyMediaThumbnailsPreview.vue";
import { onBeforeMount } from "vue";
import { useRoute } from "vue-router";
import { useSearchStore } from "@/piniaStore/search";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
const route = useRoute();
const searchQuery = Array.isArray(route.params.query)
  ? route.params.query[0]
  : route.params.query;

const { accessToken } = storeToRefs(useAuthStore());
const searchStore = useSearchStore();
const { mediaList, allMediaLoaded } = storeToRefs(searchStore);
const { loadMoreSearchResults, setQuery } = searchStore;
onBeforeMount(() => {
  // as we are using global store for search results, it can still have results of old media search
  // this will ensure to update search query and results in store
  setQuery(searchQuery);
});

function loadAllMediaOfDate(date: Date): Promise<boolean> {
  return new Promise((resolve, reject) => {
    let lastMedia = mediaList.value[mediaList.value.length - 1];
    while (date < lastMedia.date && !allMediaLoaded.value) {
      loadMoreSearchResults(accessToken.value, searchQuery)
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
</script>

<template>
  <LazyMediaThumbnailsPreview
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="() => loadMoreSearchResults(accessToken, searchQuery)"
    :load-all-media-of-date="loadAllMediaOfDate"
  />
  <!-- :load-more-media="() => loadMoreSearchResults(accessToken, searchQuery)" -->
</template>
