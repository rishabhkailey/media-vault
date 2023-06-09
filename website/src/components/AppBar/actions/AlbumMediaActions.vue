<script setup lang="ts">
import { ref } from "vue";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { useMediaStore } from "@/piniaStore/media";
import { useLoadingStore } from "@/piniaStore/loading";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/piniaStore/auth";
import ConfirmationPopup from "@/components/ConfirmationPopup.vue";
import { useRoute } from "vue-router";
import { useAlbumStore } from "@/piniaStore/album";

const mediaSelectionStore = useMediaSelectionStore();
const { reset: resetMediaSelection, updateSelection } = mediaSelectionStore;
const { selectedMediaIDs } = storeToRefs(mediaSelectionStore);

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

const { deleteMultipleMedia } = useMediaStore();
const { setGlobalLoading, setProgress } = useLoadingStore();

const deleteConfirmationPopUp = ref(false);

// batch deleting
// if we have a lot of media selected deleting 1 by 1 cause rerender for every delete, and it causes high cpu usage in browser
async function deleteSelectedMedia() {
  // we don't want this to be reactive
  let mediaIDs = [...selectedMediaIDs.value];
  let failedIDs = new Array<number>();
  let count = mediaIDs.length;
  resetMediaSelection();
  setGlobalLoading(true, false, 0);
  const batchSize = 30;
  for (let index = 0; index < count; index += batchSize) {
    let end = Math.min(index + batchSize, mediaIDs.length);
    let mediaIDsToDelete = mediaIDs.slice(index, end);
    try {
      let failedMediaIDs = await deleteMultipleMedia(
        accessToken.value,
        mediaIDsToDelete
      );
      failedIDs.push(...failedMediaIDs);
      setProgress((100 * (index + batchSize)) / count);
    } catch (err) {
      failedIDs.push(...mediaIDsToDelete);
      // todo user feedback component for errors
      console.log(err);
    }
  }
  failedIDs.forEach((id) => updateSelection(id, true));
  setGlobalLoading(false, false, 0);
}

const route = useRoute();
const removeSelectedMediaError = ref("");
const removeSelectedMediaInProgress = ref(false);
const removedMediaFromAlbumPopUp = ref(false);
const albumID = Array.isArray(route.params.album_id)
  ? route.params.album_id[0]
  : route.params.album_id;
const { removeMediaFromAlbum } = useAlbumStore();

function removeSelectedMedia() {
  removeSelectedMediaInProgress.value = true;
  removeMediaFromAlbum(accessToken.value, Number(albumID), [
    ...selectedMediaIDs.value,
  ])
    .then(() => {
      removeSelectedMediaInProgress.value = false;
      removedMediaFromAlbumPopUp.value = false;
      resetMediaSelection();
      // todo remove media from album media list
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
        <ConfirmationPopup
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
        <ConfirmationPopup
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
