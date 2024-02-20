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
import { getSingleMediaById } from "@/js/api/media";
import { useErrorsStore } from "@/piniaStore/errors";

const router = useRouter();
const route = useRoute();
const { appendError } = useErrorsStore();

// params
const index = ref(-1);
const mediaID = ref(-1);
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

const loading = ref<boolean>(false);
function initSingleMediaPreviewRefsAndStore() {
  loading.value = true;
  getSingleMediaById(mediaID.value)
    .then((media) => {
      allMediaLoaded.value = true;
      mediaList.value = [media];
      index.value = 0;
      loadMoreMedia = () =>
        new Promise<LoadMoreMediaStatus>((resolve) => resolve("empty"));
    })
    .catch((err) => {
      appendError(
        "failed to get media info from server",
        `error message - ${err}`,
        -1,
      );
    })
    .finally(() => {
      loading.value = false;
    });
}

function updateIndex(newIndex: number) {
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
    console.error(err);
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
    :loading="loading"
    :index="index"
    @update:index="updateIndex"
    :media-list="mediaList"
    :load-more-media="loadMoreMedia"
    :all-media-loaded="allMediaLoaded"
    route-name="MediaPreview"
    :animation-origin-selector="`#thumbnail_${mediaList[index]?.id}`"
    @close="
      () => {
        router.push(searchRoute(query));
      }
    "
  />
</template>
