<script setup lang="ts">
import MediaCarousel from "@/components/MediaCarousel/MediaCarousel.vue";
import { storeToRefs } from "pinia";
import { onBeforeMount, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAlbumMediaStore } from "@/piniaStore/albumMedia";
import {
  albumMediaPreviewRoute,
  albumRoute,
  errorScreenRoute,
} from "@/router/routesConstants";
import { getQueryParamNumberValue } from "@/js/utils";
import { getSingleMediaById } from "@/js/api/media";
import { useErrorsStore } from "@/piniaStore/errors";

const router = useRouter();
const route = useRoute();
const { appendError } = useErrorsStore();

// params
const index = ref(0);
const mediaID = ref(0);
const albumID = ref(0);

function initParams() {
  // media_id
  let mediaIdParam = getQueryParamNumberValue(route.params, "media_id");
  if (mediaIdParam === undefined) {
    throw new Error("invalid media id param");
  }
  mediaID.value = Number(mediaIdParam);

  // album_id
  let albumIdParam = getQueryParamNumberValue(route.params, "album_id");
  if (albumIdParam === undefined) {
    throw new Error("invalid album id param.");
  }
  albumID.value = Number(albumIdParam);

  // media index
  let indexParam = getQueryParamNumberValue(route.params, "index");
  if (indexParam === undefined) {
    throw new Error("invalid index param.");
  }
  index.value = indexParam;
}

let allMediaLoaded = ref(true);
let mediaList = ref<Array<Media>>([]);
let loadMoreMedia: LoadMoreMedia;

function initMediaPreviewRefsAndStore() {
  const albumMediaStore = useAlbumMediaStore();
  if (albumMediaStore.albumID !== albumID.value) {
    albumMediaStore.setAlbumID(albumID.value);
  }
  if (
    albumMediaStore.mediaList.findIndex((m) => m.id === mediaID.value) !==
    index.value
  ) {
    initSingleMediaPreviewRefsAndStore();
    return;
  }

  ({ allMediaLoaded, mediaList } = storeToRefs(albumMediaStore));
  ({ loadMoreMedia } = albumMediaStore);
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
    albumMediaPreviewRoute(
      newIndex,
      mediaList.value[newIndex].id,
      albumID.value,
    ),
  );
}

onBeforeMount(() => {
  try {
    initParams();
    initMediaPreviewRefsAndStore();
  } catch (err) {
    console.debug(err);
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
        router.push(albumRoute(albumID));
      }
    "
  />
</template>
