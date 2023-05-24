<script setup lang="ts">
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { useMediaStore } from "@/piniaStore/media";
import { useLoadingStore } from "@/piniaStore/loading";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/piniaStore/auth";
import { ref } from "vue-demi";

const mediaSelectionStore = useMediaSelectionStore();
const { reset: resetMediaSelection, updateSelection } = mediaSelectionStore;
const { count: selectedMediaCount, selectedMediaIDs } =
  storeToRefs(mediaSelectionStore);

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

const { deleteMultipleMedia } = useMediaStore();
const { setGlobalLoading, setProgress } = useLoadingStore();

const deleteConfirmationPopUp = ref(false);
// const deleteButton = ref<HTMLElement | undefined>(undefined);

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
</script>
<template>
  <v-row class="d-flex align-center ml-2 justify-start mx-2">
    <!-- start -->
    <v-col>
      <v-row class="d-flex align-center">
        <v-btn icon="mdi-close" @click.stop="resetMediaSelection" />
        <div class="text-subtitle-1">{{ selectedMediaCount }} selected</div>
      </v-row>
    </v-col>

    <!-- end -->
    <v-col class="d-flex flex-row justify-end">
      <v-row class="d-flex align-center justify-end">
        <!-- delete button -->
        <transition name="bounce">
          <v-tooltip
            v-if="!deleteConfirmationPopUp"
            text="delete"
            location="bottom"
          >
            <template v-slot:activator="{ props }">
              <v-btn
                :icon="
                  deleteConfirmationPopUp
                    ? 'mdi-trash-can'
                    : 'mdi-trash-can-outline'
                "
                @click.stop="() => (deleteConfirmationPopUp = true)"
                color="white"
                v-bind="props"
              />
            </template>
          </v-tooltip>
        </transition>
        <!-- delete confirmation -->
        <transition name="bounce">
          <div class="d-flex flex-row" v-if="deleteConfirmationPopUp">
            <v-tooltip text="confirm delete" location="bottom">
              <template v-slot:activator="{ props }">
                <v-btn
                  icon="mdi-check"
                  v-bind="props"
                  @click.stop="deleteSelectedMedia"
                />
              </template>
            </v-tooltip>
            <v-tooltip text="cancel delete" location="bottom">
              <template v-slot:activator="{ props }">
                <v-btn
                  icon="mdi-close"
                  @click.stop="() => (deleteConfirmationPopUp = false)"
                  v-bind="props"
                />
              </template>
            </v-tooltip>
          </div>
        </transition>
      </v-row>
    </v-col>
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
