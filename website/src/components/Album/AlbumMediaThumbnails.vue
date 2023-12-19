<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAlbumStore } from "@/piniaStore/album";
import { storeToRefs } from "pinia";
import { useAlbumMediaStore } from "@/piniaStore/albumMedia";
import LazyMediaThumbnailsPreviewVue from "../MediaThumbnailPreview/LazyMediaThumbnailsPreview.vue";
import ConfirmationPopup from "../ConfirmationPopup.vue";
import { base64UrlEncode } from "@/js/utls";

const albumStore = useAlbumStore();
const { getAlbumByID, deleteAlbum } = albumStore;

const albumMediaStore = useAlbumMediaStore();
const { loadMoreMedia, setAlbumID } = albumMediaStore;
const { mediaList, allMediaLoaded } = storeToRefs(albumMediaStore);

// todo move this to store?
async function loadAllMediaOfDate(date: Date): Promise<boolean> {
  let lastMediaDate = mediaList.value[mediaList.value.length - 1].date;
  console.log(date, lastMediaDate);
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

// in case of error it will return empty string
const getAlbumIdFromRoute = () => {
  let albumID = Array.isArray(route.params.album_id)
    ? route.params.album_id[0]
    : route.params.album_id;
  if (albumID.length == 0 || isNaN(Number(albumID))) {
    errorMessage.value = "invalid or empty album id";
    return "";
  }
  return albumID;
};

const route = useRoute();
const errorMessage = ref("");
let albumID = ref(getAlbumIdFromRoute());
console.log("album id in script", albumID.value);

watch(
  () => route.params.album_id,
  () => {
    let newAlbumId = getAlbumIdFromRoute();
    if (newAlbumId.length === 0 || newAlbumId === albumID.value) {
      return;
    }
    console.log("album id updated", newAlbumId);
    albumID.value = newAlbumId;
    loadAlbum();
  }
);

const album = ref<Album>({
  id: 0,
  name: "something went wrong",
  created_at: new Date(),
  media_count: 0,
  thumbnail_url: "",
  updated_at: new Date(),
});
const loading = ref(false);
const router = useRouter();

const deleteConfirmationOverlay = ref(false);
const deleteInProgress = ref(false);
const deleteErrorMessage = ref("");
function onDeleteConfirm() {
  deleteInProgress.value = true;
  deleteErrorMessage.value = "";
  deleteAlbum(Number(albumID))
    .then(() => {
      deleteInProgress.value = false;
      deleteConfirmationOverlay.value = false;
      deleteErrorMessage.value = "";
      router.push({
        name: "Albums",
      });
    })
    .catch((err) => {
      deleteErrorMessage.value = "something went wrong, " + err;
      deleteInProgress.value = false;
    });
}

const loadAlbum = () => {
  loading.value = true;
  setAlbumID(Number(albumID.value));
  getAlbumByID(Number(albumID.value))
    .then((_album) => {
      album.value = _album;
      loading.value = false;
    })
    .catch((err) => {
      errorMessage.value = "something went wrong. " + err;
    });
};

onMounted(() => {
  console.log("album id in mounted", albumID.value);
  if (albumID.value.length === 0) {
    return;
  }
  loadAlbum();
});
</script>
<template>
  <h1 v-if="loading">Loading...</h1>
  <v-col v-else>
    <v-row>
      <v-toolbar :collapse="false" :title="album.name ?? '!'" color="surface">
        <v-btn
          prepend-icon="mdi-delete"
          @click.stop="
            () => {
              deleteConfirmationOverlay = true;
            }
          "
          >Delete Album</v-btn
        >
        <ConfirmationPopup
          title="Delete album?"
          message="Deleting an album is permanent. Photos and videos that were in a
            deleted album will not be deleted."
          cancel-button-text="keep"
          cancel-button-color=""
          confirm-button-text="Delete"
          confirm-button-color="red"
          :confirm-in-progress="deleteInProgress"
          v-model:model-value="deleteConfirmationOverlay"
          :error-message="deleteErrorMessage"
          @cancel="
            () => {
              deleteInProgress = false;
              deleteConfirmationOverlay = false;
            }
          "
          @confirm="() => onDeleteConfirm()"
        />
      </v-toolbar>
    </v-row>
    <v-row>
      <v-divider class="border-opacity-25"></v-divider>
    </v-row>
    <v-row>
      <LazyMediaThumbnailsPreviewVue
        class="flex-grow-1"
        :media-list="mediaList"
        :all-media-loaded="allMediaLoaded"
        :load-more-media="() => loadMoreMedia()"
        :load-all-media-of-date="loadAllMediaOfDate"
        :media-date-getter="(media: Media) => media.uploaded_at"
        @thumbnail-click="
          (clickedMediaID, clickedIndex, thumbnailClickLocation) => {
            router.push({
              name: `AlbumMediaPreview`,
              params: {
                index: clickedIndex,
                media_id: clickedMediaID,
                album: albumID,
              },
              hash: `#${base64UrlEncode(thumbnailClickLocation)}`,
            });
          }
        "
      />
    </v-row>
  </v-col>
  <Teleport to="body">
    <div style="position: absolute; top: 0px; left: 0px; z-index: 9999999">
      <RouterView />
    </div>
  </Teleport>
</template>
