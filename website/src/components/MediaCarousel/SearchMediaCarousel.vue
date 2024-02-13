<script setup lang="ts">
import { useSearchStore } from "@/piniaStore/search";
import MediaCarousel from "@/components/MediaCarousel/MediaCarousel.vue";
import { storeToRefs } from "pinia";
import { onBeforeMount, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/piniaStore/auth";
import {
  errorScreenRoute,
  searchMediaPreviewRoute,
  searchRoute,
} from "@/router/routesConstants";
import { getQueryParamNumberValue, getQueryParamStringValue } from "@/js/utils";

const router = useRouter();
const route = useRoute();

// params
const index = ref(0);
const mediaID = ref(0);
const query = ref("");

function initParams() {
  // media_id
  let mediaIdParam = getQueryParamNumberValue(route.params, "media_id");
  if (mediaIdParam === undefined) {
    throw new Error("invalid media id param");
  }
  mediaID.value = Number(mediaIdParam);

  // media index
  let indexParam = getQueryParamNumberValue(route.params, "index");
  if (indexParam === undefined) {
    throw new Error("invalid index param");
  }
  index.value = indexParam;

  // query
  let queryParam = getQueryParamStringValue(route.params, "query");
  if (queryParam === undefined || queryParam.length === 0) {
    throw new Error("invalid search param");
  }
  query.value = queryParam;
}

let allMediaLoaded = ref(true);
let mediaList = ref<Array<Media>>([]);
let loadMoreMedia: LoadMoreMedia;

function initMediaPreviewRefsAndStore() {
  const searchStore = useSearchStore();
  if (searchStore.query !== query.value) {
    searchStore.setQuery(query.value);
  }
  if (
    searchStore.mediaList.findIndex((m) => m.id === mediaID.value) !==
    index.value
  ) {
    initSingleMediaPreviewRefsAndStore();
    return;
  }
  ({ allMediaLoaded, mediaList } = storeToRefs(searchStore));
  loadMoreMedia = () =>
    searchStore.loadMoreSearchResults(useAuthStore().accessToken, query.value);
}

function initSingleMediaPreviewRefsAndStore() {
  allMediaLoaded.value = true;
  mediaList.value = [];
  loadMoreMedia = () =>
    new Promise<LoadMoreMediaStatus>((resolve) => resolve("empty"));
}

function updateIndex(newIndex: number) {
  console.log(newIndex);
  index.value = newIndex;
  router.push(
    searchMediaPreviewRoute(
      newIndex,
      mediaList.value[newIndex].id,
      query.value,
    ),
  );
}

onBeforeMount(() => {
  try {
    initParams();
    initMediaPreviewRefsAndStore();
  } catch (err) {
    console.log(err);
    router.push(
      errorScreenRoute(
        "AlbumMediaCarousel component intialization failed",
        `error message = "${err}"`,
      ),
    );
  }
});
</script>
<template>
  <MediaCarousel
    :index="index"
    @update:index="updateIndex"
    :media-list="mediaList"
    :load-more-media="loadMoreMedia"
    :all-media-loaded="allMediaLoaded"
    route-name="MediaPreview"
    :animation-origin-selector="`#thumbnail_${mediaList[index].id}`"
    @close="
      () => {
        router.push(searchRoute(query));
      }
    "
  />
</template>
