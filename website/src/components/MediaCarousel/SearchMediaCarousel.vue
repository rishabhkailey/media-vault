<script setup lang="ts">
import { useSearchStore } from "@/piniaStore/search";
import MediaCarousel from "@/components/MediaCarousel/MediaCarousel.vue";
import { storeToRefs } from "pinia";
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/piniaStore/auth";

const router = useRouter();
const route = useRoute();

// params
const index = ref(0);
const mediaID = ref(0);
const query = ref("");

initParams();
function initParams() {
  // media_id
  let mediaIdParam = Array.isArray(route.params.media_id)
    ? route.params.media_id[0]
    : route.params.media_id;
  if (Number.isNaN(mediaIdParam)) {
    router.replace({
      name: "errorscreen",
      query: {
        title: "Invalid Media ID",
        message: `got media id "${mediaIdParam}", expected a number.`,
      },
    });
    return;
  }
  mediaID.value = Number(mediaIdParam);

  // query
  let queryParam = Array.isArray(route.params.query)
    ? route.params.query[0]
    : route.params.query;
  if (queryParam.length === 0) {
    router.replace({
      name: "errorscreen",
      query: {
        title: "Invalid Search Query",
        message: `got search query "${queryParam}", expected a number.`,
      },
    });
    return;
  }
  query.value = queryParam;

  // media index
  let indexParam = Array.isArray(route.params.index)
    ? route.params.index[0]
    : route.params.index;
  if (!Number.isNaN(indexParam)) {
    index.value = Number(indexParam);
  }
}

let allMediaLoaded = ref(true);
let mediaList = ref<Array<Media>>([]);
let loadMoreMedia: () => Promise<boolean>;
initMediaPreviewRefsAndStore();

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
  loadMoreMedia = () => new Promise<boolean>((resolve) => resolve(true));
}

function updateIndex(newIndex: number) {
  console.log(newIndex);
  index.value = newIndex;
  router.push({
    name: `SearchMediaPreview`,
    params: {
      index: newIndex,
      media_id: mediaList.value[newIndex].id,
      query: query.value,
    },
  });
}
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
        router.push({
          name: `search`,
          params: {
            query: query,
          },
        });
      }
    "
  />
</template>
