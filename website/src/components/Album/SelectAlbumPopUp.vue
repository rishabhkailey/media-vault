<script lang="ts" setup>
import { useAlbumStore } from "@/piniaStore/album";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
import { ref } from "vue";
import LazyLoading from "../LazyLoading/LazyLoading.vue";

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
const { albums, allAlbumsLoaded } = storeToRefs(albumStore);
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
        <v-checkbox
          v-for="album in albums"
          :key="album.id"
          :value="album.id"
          v-model="selectedAlbums"
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
        <v-row class="justify-center">
          <LazyLoading
            v-if="!allAlbumsLoaded"
            :on-threshold-reach="loadMoreAlbums"
            :threshold="0.1"
            :min-height="100"
            :min-width="100"
            :root-margin="10"
          >
            <v-progress-circular indeterminate></v-progress-circular>
          </LazyLoading>
        </v-row>
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
          >
            {{ props.confirmButtonText }}
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-overlay>
</template>
