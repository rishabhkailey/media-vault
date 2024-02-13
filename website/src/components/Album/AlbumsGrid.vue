<script setup lang="ts">
import { useAlbumStore } from "@/piniaStore/album";
import CreateAlbumModal from "./CreateAlbumModal.vue";
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { useLoadingStore } from "@/piniaStore/loading";
import AlbumCard from "./AlbumCard.vue";
import { useRouter } from "vue-router";
import KebabMenuWrapper from "../KebabMenuWrapper/KebabMenuWrapper.vue";
import ConfirmationPopupVue from "@/components/Modals/ConfirmationModal.vue";
import { albumRoute } from "@/router/routesConstants";

const router = useRouter();
const albumStore = useAlbumStore();
const { albums } = storeToRefs(albumStore);
const { loadMoreAlbums, deleteAlbum } = albumStore;
const showCreateAlbumModal = ref(false);
const deleteConfirmationOverlay = ref(false);
const toDeleteAlbumID = ref(0);
const deleteInProgress = ref(false);
const deleteErrorMessage = ref("");

function onDeleteButtonClick(albumID: number) {
  deleteConfirmationOverlay.value = true;
  toDeleteAlbumID.value = albumID;
}

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
  <v-col>
    <v-row>
      <v-toolbar :collapse="false" title="Albums" color="surface">
        <v-btn
          prepend-icon="mdi-plus"
          @click.stop="
            () => {
              showCreateAlbumModal = true;
            }
          "
          >Create</v-btn
        >
        <CreateAlbumModal v-model="showCreateAlbumModal" />
      </v-toolbar>
    </v-row>
    <v-row>
      <v-divider class="border-opacity-25"></v-divider>
    </v-row>
    <v-row>
      <!-- <SizeWrapper v-slot="{ width }"> -->
      <!-- cols = 12 (default will only work for xs) -->
      <v-infinite-scroll
        :items="albums"
        style="width: 100%; height: 100%"
        @load="
          ({ done }) => {
            loadMoreAlbums()
              .then((status) => {
                done(status);
              })
              .catch((_) => {
                done('error');
              });
          }
        "
      >
        <v-col
          :xxl="2"
          :xl="2"
          :lg="2"
          :md="3"
          :sm="6"
          :xs="12"
          :cols="12"
          :key="`${index}+${album.id}`"
          v-for="(album, index) in albums"
        >
          <KebabMenuWrapper
            :show-select-button-on-hover="true"
            :select-on-content-click="false"
            :always-show-select-button="false"
            :always-show-select-on-mobile="true"
            selectIconSize="large"
          >
            <AlbumCard
              :padding="0"
              :aspect-ratio="1"
              :album="album"
              @click="
                () => {
                  router.push(albumRoute(album.id));
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
      </v-infinite-scroll>
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
    </v-row>
  </v-col>
</template>
