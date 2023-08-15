<script setup lang="ts">
import { ref } from "vue";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { useMediaStore } from "@/piniaStore/media";
import { useLoadingStore } from "@/piniaStore/loading";
import { storeToRefs } from "pinia";
import ConfirmationPopup from "@/components/ConfirmationPopup.vue";
import SelectAlbumPopUp from "@/components/Album/SelectAlbumPopUp.vue";
import { useAlbumStore } from "@/piniaStore/album";
import { useAlbumMediaStore } from "@/piniaStore/albumMedia";

const mediaSelectionStore = useMediaSelectionStore();
const { reset: resetMediaSelection, updateSelection } = mediaSelectionStore;
const { selectedMediaIDs } = storeToRefs(mediaSelectionStore);

const { deleteMultipleMedia } = useMediaStore();
const { setGlobalLoading, setProgress } = useLoadingStore();

const { removeMediaByIDsFromLocalState } = useAlbumMediaStore();
const deleteConfirmationPopUp = ref(false);
// const deleteButton = ref<HTMLElement | undefined>(undefined);

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
      let failedMediaIDs = await deleteMultipleMedia(mediaIDsToDelete);
      failedIDs.push(...failedMediaIDs);
      setProgress((100 * (index + batchSize)) / count);
    } catch (err) {
      failedIDs.push(...mediaIDsToDelete);
      // todo user feedback component for errors
      console.log(err);
    }
  }
  failedIDs.forEach((id) => updateSelection(id, true));
  removeMediaByIDsFromLocalState(
    mediaIDs.filter((id) => !failedIDs.includes(id))
  );
  setGlobalLoading(false, false, 0);
}

const { addMediaToAlbum } = useAlbumStore();
const { reset: resetAlbumMedia } = useAlbumMediaStore();
const addToAlbumConfirmationPopUp = ref(false);
const addToAlbumInProgress = ref(false);
const addToAlbumErrorMessage = ref("");
async function addToAlbumsConfirm(
  albumIDs: Array<number>
): Promise<Array<number>> {
  let mediaIDs = [...selectedMediaIDs.value];
  let failedAlbumIDs = [];
  addToAlbumInProgress.value = true;
  for (const albumID of albumIDs) {
    try {
      await addMediaToAlbum(albumID, mediaIDs);
    } catch (err) {
      failedAlbumIDs.push(albumID);
    }
  }
  if (failedAlbumIDs.length !== 0) {
    addToAlbumErrorMessage.value = "something went wrong";
  } else {
    addToAlbumConfirmationPopUp.value = false;
    resetMediaSelection();
  }
  addToAlbumInProgress.value = false;
  resetAlbumMedia();
  return failedAlbumIDs;
}
</script>
<template>
  <v-row class="d-flex align-center justify-end">
    <!-- add media in album button -->
    <v-tooltip text="Add to album" location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn
          icon="mdi-plus"
          @click.stop="() => (addToAlbumConfirmationPopUp = true)"
          color="white"
          v-bind="props"
        />
      </template>
    </v-tooltip>
    <!-- add media in album confirmation -->
    <SelectAlbumPopUp
      title="Add to"
      cancel-button-color=""
      cancel-button-text="cancel"
      message=""
      v-model:model-value="addToAlbumConfirmationPopUp"
      error-message=""
      confirm-button-text="Confirm"
      confirm-button-color="primary"
      :confirm-in-progress="addToAlbumInProgress"
      @cancel="() => (addToAlbumConfirmationPopUp = false)"
      @confirm="addToAlbumsConfirm"
    />

    <!-- delete button -->
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
    <!-- delete confirmation -->
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
  </v-row>
</template>

<style scoped>
/* slide animation 1 */
.delete-confirmation-enter-active,
.delete-confirmation-leave-active {
  transition: width 0.3s ease, opacity 0.3s ease;
  width: 100px;
  opacity: 100%;
}

.delete-confirmation-enter-from,
.delete-confirmation-leave-to {
  width: 0px;
  opacity: 0;
}

/* slide animation 2 */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.3s ease-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(20px);
  width: 0;
  opacity: 0;
}

/* bounce transitio */
.bounce-enter-active {
  animation: bounce-in 0.3s;
}
.bounce-leave-active {
  animation: bounce-in 0.3s reverse;
}

@keyframes bounce-in {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}
</style>
