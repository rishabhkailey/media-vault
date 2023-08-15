<script setup lang="ts">
import MediaPreview from "./MediaPreview.vue";
import { storeToRefs } from "pinia";
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAlbumMediaStore } from "@/piniaStore/albumMedia";

const router = useRouter();
const route = useRoute();

// params
const index = ref(0);
const mediaID = ref(0);
const albumID = ref(0);

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

  // album_id
  let albumIdParam = Array.isArray(route.params.album_id)
    ? route.params.album_id[0]
    : route.params.album_id;
  if (Number.isNaN(mediaIdParam)) {
    router.replace({
      name: "errorscreen",
      query: {
        title: "Invalid Album ID",
        message: `got album id "${albumIdParam}", expected a number.`,
      },
    });
    return;
  }
  albumID.value = Number(albumIdParam);

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

function initSingleMediaPreviewRefsAndStore() {
  allMediaLoaded.value = true;
  mediaList.value = [];
  loadMoreMedia = () => new Promise<boolean>((resolve) => resolve(true));
}

function updateIndex(newIndex: number) {
  console.log(newIndex);
  index.value = newIndex;
  router.push({
    name: `AlbumMediaPreview`,
    params: {
      index: newIndex,
      media_id: mediaList.value[newIndex].id,
      album: albumID.value,
    },
  });
}
</script>
<template>
  <MediaPreview
    :index="index"
    @update:index="updateIndex"
    :media-list="mediaList"
    :load-more-media="loadMoreMedia"
    :all-media-loaded="allMediaLoaded"
    route-name="MediaPreview"
    @close="
      () => {
        router.push({
          name: `Album`,
          params: {
            album_id: albumID,
          },
        });
      }
    "
  />
</template>
