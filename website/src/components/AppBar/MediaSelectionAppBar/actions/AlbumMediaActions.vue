<script setup lang="ts">
import { ref } from "vue";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { useMediaStore } from "@/piniaStore/media";
import { useLoadingStore } from "@/piniaStore/loading";
import { storeToRefs } from "pinia";
import ConfirmationModal from "@/components/Modals/ConfirmationModal.vue";
import { useRoute, useRouter } from "vue-router";
import { useAlbumStore } from "@/piniaStore/album";
import { useErrorsStore } from "@/piniaStore/errors";
import { getQueryParamNumberValue } from "@/js/utils";
import { errorScreenRoute } from "@/router/routesConstants";

const mediaSelectionStore = useMediaSelectionStore();
const { reset: resetMediaSelection, updateSelection } = mediaSelectionStore;
const { selectedMediaMap } = storeToRefs(mediaSelectionStore);
const { deleteMultipleMedia } = useMediaStore();
const { setGlobalLoading } = useLoadingStore();
const { appendError } = useErrorsStore();
const deleteConfirmationPopUp = ref(false);
const router = useRouter();

function getAlbumIdFromRoute(): number {
  let albumID = getQueryParamNumberValue(route.params, "album_id");
  if (albumID === undefined) {
    router.push(
      errorScreenRoute(
        "Failed to render Album's Media page",
        "invalid or empty album id",
      ),
    );
    return 0;
  }
  return albumID;
}

// batch deleting
// if we have a lot of media selected deleting 1 by 1 cause rerender for every delete, and it causes high cpu usage in browser
async function deleteSelectedMedia() {
  // we don't want this to be reactive
  let mediaIDs = [...selectedMediaMap.value.keys()];
  resetMediaSelection();
  setGlobalLoading(true, false, 0);
  try {
    await deleteMultipleMedia(mediaIDs);
  } catch (err) {
    mediaIDs.forEach((id) => updateSelection(id, true));
    appendError(
      "deletion failed",
      `there was an issue in deletion if the files are still selected please try to delete again. error message - ${err}`,
      -1,
    );
  }
  setGlobalLoading(false, false, 0);
}

const route = useRoute();
const removeSelectedMediaError = ref("");
const removeSelectedMediaInProgress = ref(false);
const removedMediaFromAlbumPopUp = ref(false);
const albumID = getAlbumIdFromRoute();
const { removeMediaFromAlbum } = useAlbumStore();

function removeSelectedMedia() {
  removeSelectedMediaInProgress.value = true;
  removeMediaFromAlbum(Number(albumID), [...selectedMediaMap.value.keys()])
    .then(() => {
      removeSelectedMediaInProgress.value = false;
      removedMediaFromAlbumPopUp.value = false;
      resetMediaSelection();
    })
    .catch((err) => {
      removeSelectedMediaInProgress.value = false;
      removeSelectedMediaError.value = "something went wrong, " + err;
    });
}
</script>
<template>
  <v-row class="d-flex align-center justify-end">
    <!-- remove from album button -->
    <div>
      <!-- button -->
      <transition name="bounce">
        <v-tooltip text="remove from album" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn
              icon="mdi-image-remove-outline"
              @click.stop="() => (removedMediaFromAlbumPopUp = true)"
              color="white"
              v-bind="props"
            />
          </template>
        </v-tooltip>
      </transition>
      <!-- confirmation -->
      <transition name="bounce">
        <ConfirmationModal
          title="Remove selected?"
          cancel-button-color=""
          cancel-button-text="cancel"
          message="Remove selected items from Album"
          v-model:model-value="removedMediaFromAlbumPopUp"
          :error-message="removeSelectedMediaError"
          confirm-button-text="Remove"
          confirm-button-color="red"
          :confirm-in-progress="removeSelectedMediaInProgress"
          @cancel="() => (removedMediaFromAlbumPopUp = false)"
          @confirm="removeSelectedMedia"
        />
      </transition>
    </div>

    <!-- delete button -->
    <div>
      <!-- button -->
      <transition name="bounce">
        <v-tooltip text="delete" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn
              icon="mdi-trash-can-outline"
              @click.stop="() => (deleteConfirmationPopUp = true)"
              color="white"
              v-bind="props"
            />
          </template>
        </v-tooltip>
      </transition>
      <!-- confirmation -->
      <transition name="bounce">
        <ConfirmationModal
          title="Delete selected?"
          cancel-button-color=""
          cancel-button-text="cancel"
          message="Selected items will be deleted forever"
          v-model:model-value="deleteConfirmationPopUp"
          error-message=""
          confirm-button-text="Delete"
          confirm-button-color="red"
          :confirm-in-progress="false"
          @cancel="() => (deleteConfirmationPopUp = false)"
          @confirm="deleteSelectedMedia"
        />
      </transition>
    </div>
  </v-row>
</template>
