<script setup lang="ts">
import LazyLoading from "@/components/LazyLoading/LazyLoading.vue";
import { useAlbumStore } from "@/piniaStore/album";
import CreateAlbumOverlay from "./CreateAlbumOverlay.vue";
import { computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { useLoadingStore } from "@/piniaStore/loading";
import AlbumTitleThumbnail from "./AlbumTitleThumbnail.vue";
import { useRouter } from "vue-router";
import KebabMenuWrapper from "../KebabMenuWrapper/KebabMenuWrapper.vue";
import ConfirmationPopupVue from "../ConfirmationPopup.vue";
import { useDisplay } from "vuetify";

// cols="6"
// sm="4"
// md="3"
// lg="2"
// xl="2"
// xxl="1"
const display = useDisplay();
const width = computed<number>(() => {
  switch (display.name.value) {
    case "xs":
      return display.width.value / 2;
    case "sm":
      return display.width.value / 3;
    case "md":
      return display.width.value / 4;
    case "lg":
      return display.width.value / 6;
    case "xl":
      return display.width.value / 6;
    case "xxl":
      return display.width.value / 12;
    default:
      return display.width.value / 2;
  }
});

const router = useRouter();

const albumStore = useAlbumStore();
const { albums, allAlbumsLoaded } = storeToRefs(albumStore);
const { loadMoreAlbums, deleteAlbum } = albumStore;

const { initializing } = storeToRefs(useLoadingStore());

const createAlbumOverlay = ref(false);

const deleteConfirmationOverlay = ref(false);
const toDeleteAlbumID = ref(0);
const deleteInProgress = ref(false);
function onDeleteButtonClick(albumID: number) {
  deleteConfirmationOverlay.value = true;
  toDeleteAlbumID.value = albumID;
}

const deleteErrorMessage = ref("");
function onDeleteConfirm(albumID: number) {
  deleteInProgress.value = true;
  deleteErrorMessage.value = "";
  deleteAlbum(albumID)
    .then(() => {
      deleteInProgress.value = false;
      deleteConfirmationOverlay.value = false;
      toDeleteAlbumID.value = 0;
      deleteErrorMessage.value = "";
    })
    .catch((err) => {
      deleteErrorMessage.value = "something went wrong, " + err;
      deleteInProgress.value = false;
    });
}
</script>
<template>
  <div v-if="initializing">Loading...</div>
  <v-col v-else>
    <v-row>
      <v-toolbar :collapse="false" title="Albums" color="surface">
        <v-btn
          prepend-icon="mdi-plus"
          @click.stop="
            () => {
              createAlbumOverlay = true;
            }
          "
          >Create</v-btn
        >
        <CreateAlbumOverlay v-model="createAlbumOverlay" />
      </v-toolbar>
    </v-row>
    <v-row>
      <v-divider class="border-opacity-25"></v-divider>
    </v-row>
    <v-row>
      <div class="d-flex flex-row flex-wrap">
        <v-col
          :key="`${index}+${album.id}`"
          v-for="(album, index) in albums"
          class="d-flex child-flex pa-2"
        >
          <KebabMenuWrapper
            class="w-100"
            :show-select-button-on-hover="true"
            :select-on-content-click="false"
            :always-show-select-button="false"
            selectIconSize="large"
          >
            <AlbumTitleThumbnail
              :padding="0"
              :aspect-ratio="1"
              :width="width"
              class="w-100 h-100"
              :album="album"
              @click="
                () => {
                  router.push({
                    name: 'Album',
                    params: {
                      album_id: album.id,
                    },
                  });
                }
              "
            />
            <template #options>
              <v-list>
                <v-list-item>
                  <v-btn @click.stop="() => onDeleteButtonClick(album.id)">
                    Delete Album
                  </v-btn>
                </v-list-item>
              </v-list>
            </template>
          </KebabMenuWrapper>
        </v-col>
      </div>
      <ConfirmationPopupVue
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
            toDeleteAlbumID = 0;
            deleteErrorMessage = '';
          }
        "
        @confirm="() => onDeleteConfirm(toDeleteAlbumID)"
      />
      <LazyLoading
        v-if="!allAlbumsLoaded"
        :on-threshold-reach="() => loadMoreAlbums()"
        :threshold="0.1"
        :min-height="100"
        :min-width="100"
        :root-margin="10"
      ></LazyLoading>
    </v-row>
  </v-col>
</template>
