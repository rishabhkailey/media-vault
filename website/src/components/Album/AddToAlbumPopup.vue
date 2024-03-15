<script lang="ts" setup>
import { useAlbumStore } from "@/piniaStore/album";
import { storeToRefs } from "pinia";
import { ref } from "vue";

const props = defineProps<{
  title: string;
  message: string;
  cancelButtonText: string;
  cancelButtonColor: string;
  confirmButtonText: string;
  confirmButtonColor: string;
  confirmInProgress: boolean;
  modelValue: boolean;
  errorMessage: string;
}>();

const emits = defineEmits<{
  (e: "confirm", albumIDs: Array<number>): void;
  (e: "cancel"): void;
  (e: "update:modelValue", value: boolean): void;
}>();

const selectedAlbums = ref<Array<number>>([]);
const albumStore = useAlbumStore();
const { loadMoreAlbums } = albumStore;
const { albums } = storeToRefs(albumStore);
</script>
<template>
  <v-overlay
    :model-value="props.modelValue"
    @update:model-value="(value) => emits('update:modelValue', value)"
    class="align-center justify-center"
  >
    <v-card style="min-width: 300px; max-width: 500px">
      <v-card-title>
        {{ props.title }}
        <v-chip size="small">{{ selectedAlbums.length }} selected</v-chip>
      </v-card-title>
      <v-card-text>
        {{ props.message }}
        <v-alert type="error" v-if="props.errorMessage.length > 0">
          {{ props.errorMessage }}
        </v-alert>
      </v-card-text>
      <!-- albums -->
      <v-container class="pl-2" style="max-height: 50vh; overflow-y: scroll">
        <v-infinite-scroll
          :min-height="100"
          :min-width="100"
          :items="albums"
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
          <template #error> failed to load data from server </template>
          <template #default>
            <v-checkbox
              v-for="album in albums"
              :key="album.id"
              :value="album.id"
              v-model="selectedAlbums"
              :data-test-id="`appbar-add-to-album-checkbox-${album.id}`"
            >
              <template #label>
                <v-list-item
                  :title="album.name"
                  prepend-icon="mdi-image"
                  :value="album.name"
                  color="primary"
                />
              </template>
            </v-checkbox>
          </template>
        </v-infinite-scroll>
      </v-container>
      <v-card-actions>
        <div class="d-flex flex-row justify-end flex-grow-1">
          <v-btn
            variant="text"
            @click="
              () => {
                emits('cancel');
              }
            "
            :color="props.cancelButtonText"
          >
            {{ props.cancelButtonText }}
          </v-btn>
          <v-btn
            :color="props.confirmButtonColor"
            variant="text"
            :loading="props.confirmInProgress"
            @click="() => emits('confirm', selectedAlbums)"
            data-test-id="appbar-add-to-album-confirm-button"
          >
            {{ props.confirmButtonText }}
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-overlay>
</template>
