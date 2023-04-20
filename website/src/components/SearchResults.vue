<script setup lang="ts">
import LazyMediaThumbnailsPreview from "@/components/LazyMediaThumbnailsPreview.vue";
import {
  LOAD_MORE_SEARCH_RESULTS_ACTION,
  SEARCH_RESULTS_GETTER,
  ALL_SEARCH_RESULTS_LOADED_GETTER,
} from "@/store/modules/search";
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useStore } from "vuex";
import type { LoadSearchResults } from "@/store/modules/search";
const route = useRoute();
const searchQuery = Array.isArray(route.params.query)
  ? route.params.query[0]
  : route.params.query;

const store = useStore();
// todo create conts for nested modules searchModule/mediaList
const mediaList = computed<Array<Media>>(
  () => store.getters[SEARCH_RESULTS_GETTER]
);
const allMediaLoaded = computed<boolean>(
  () => store.getters[ALL_SEARCH_RESULTS_LOADED_GETTER]
);

const payload: LoadSearchResults = {
  query: searchQuery,
};
const loadMoreMedia = () =>
  store.dispatch(LOAD_MORE_SEARCH_RESULTS_ACTION, payload);

onMounted(() => {
  // as we are using global store for search results, it can still have results of old media search
  // this will ensure to update search query and results in store
  store.dispatch(LOAD_MORE_SEARCH_RESULTS_ACTION, payload);
});
</script>

<template>
  <LazyMediaThumbnailsPreview
    :media-list="mediaList"
    :all-media-loaded="allMediaLoaded"
    :load-more-media="loadMoreMedia"
  />
</template>
