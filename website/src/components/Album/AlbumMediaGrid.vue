<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAlbumStore } from "@/piniaStore/album";
import { storeToRefs } from "pinia";
import { useAlbumMediaStore } from "@/piniaStore/albumMedia";
import MediaGrid from "../MediaThumbnailPreview/MediaGrid.vue";
import ConfirmationModal from "@/components/Modals/ConfirmationModal.vue";
import { getQueryParamNumberValue } from "@/js/utils";
import { albumsRoute, albumMediaPreviewRoute } from "@/router/routesConstants";
import { MEDIA_PREVIEW_CONTAINER_Z_INDEX } from "@/js/constants/z-index";
import ErrorMessage from "../Error/ErrorMessage.vue";

const errorMessage = ref("");
const albumStore = useAlbumStore();
const { getAlbumByID, deleteAlbum } = albumStore;

const albumMediaStore = useAlbumMediaStore();
const { loadMoreMedia, setAlbumID, loadAllMediaUntil } = albumMediaStore;
const { mediaList, allMediaLoaded } = storeToRefs(albumMediaStore);

function getAlbumIdFromRoute(): number {
  let albumID = getQueryParamNumberValue(route.params, "album_id");
  if (albumID === undefined) {
    errorMessage.value = "invalid or empty album id";
    return -1;
  }
  return albumID;
}

const route = useRoute();
let albumID = ref<number>(-1);
console.log("album id in script", albumID.value);

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
      router.push(albumsRoute());
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
      errorMessage.value = "invalid or empty album id" + err;
    });
};

onMounted(() => {
  albumID.value = getAlbumIdFromRoute();
  console.log("album id in mounted", albumID.value);
  if (albumID.value === -1) {
    return;
  }
  loadAlbum();
});
</script>
<template>
  <ErrorMessage
    v-if="errorMessage.length > 0"
    :message="errorMessage"
    title="Something went wrong"
    :expanded-by-default="true"
    @close="() => (errorMessage = '')"
  />
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
        <ConfirmationModal
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
      <MediaGrid
        class="flex-grow-1"
        :media-list="mediaList"
        :all-media-loaded="allMediaLoaded"
        :load-more-media="() => loadMoreMedia()"
        :load-all-media-until="loadAllMediaUntil"
        :media-date-getter="(media: Media) => media.uploaded_at"
        @thumbnail-click="
          (clickedMediaID, clickedIndex, thumbnailClickLocation) => {
            router.push(
              albumMediaPreviewRoute(
                clickedIndex,
                clickedMediaID,
                albumID,
                thumbnailClickLocation,
              ),
            );
          }
        "
      />
    </v-row>
  </v-col>
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
