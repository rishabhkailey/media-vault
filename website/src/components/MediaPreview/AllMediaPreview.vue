<script setup lang="ts">
import MediaPreview from "./MediaPreview.vue";
import { useMediaStore } from "@/piniaStore/media";
import { storeToRefs } from "pinia";
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";

const router = useRouter();
const route = useRoute();

// params
const index = ref(0);
const mediaID = ref(0);

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

// function initRefsAndStores() {
//   // switch for route: media or album or search
//   switch (route.name) {
//     case "MediaPreview": {
//       initMediaPreviewRefsAndStore();
//       return;
//     }
//     case "SearchMediaPreview": {
//       initSearchMediaPreviewRefsAndStore();
//       return;
//     }
//     case "AlbumMediaPreview": {
//       initAlbumMediaPreviewRefsAndStore();
//       return;
//     }
//     default: {
//       return;
//     }
//   }
// }
function initMediaPreviewRefsAndStore() {
  const mediaStore = useMediaStore();
  if (
    mediaStore.mediaList.findIndex((m) => m.id === mediaID.value) !==
    index.value
  ) {
    initSingleMediaPreviewRefsAndStore();
    return;
  }

  ({ allMediaLoaded, mediaList } = storeToRefs(mediaStore));
  ({ loadMoreMedia } = mediaStore);
}

function initSingleMediaPreviewRefsAndStore() {
  allMediaLoaded.value = true;
  mediaList.value = [];
  loadMoreMedia = () => new Promise<boolean>((resolve) => resolve(true));
}
// function initAlbumMediaPreviewRefsAndStore() {
//   const albumMediaStore = useAlbumMediaStore();
//   // albumMediaStore.setAlbumID()
//   ({ allMediaLoaded, mediaList } = storeToRefs(albumMediaStore));
//   ({ loadMoreMedia } = albumMediaStore);
// }
// function initSearchMediaPreviewRefsAndStore() {}

function updateIndex(newIndex: number) {
  console.log(newIndex);
  index.value = newIndex;
  router.push({
    name: `MediaPreview`,
    params: {
      index: newIndex,
      media_id: mediaList.value[newIndex].id,
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
          name: `Home`,
        });
      }
    "
  />
</template>
