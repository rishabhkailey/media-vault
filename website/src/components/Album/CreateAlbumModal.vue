<script setup lang="ts">
import { ref } from "vue";
import { useAlbumStore } from "@/piniaStore/album";
const props = defineProps<{
  modelValue: boolean;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();

const albumStore = useAlbumStore();
const { createAlbum } = albumStore;

const albumName = ref("");
const albumCreationInProgress = ref(false);
const albumNameRules = [
  (value: string) => {
    if (value) return true;

    return "You must enter a Album name.";
  },
];
const isFormValid = ref(false);

const errorMessage = ref("");
function createAlbumSubmit() {
  if (isFormValid.value == false || isFormValid.value == null) {
    return;
  }
  albumCreationInProgress.value = true;
  createAlbum(albumName.value, "")
    .then(() => {
      emit("update:modelValue", false);
      albumCreationInProgress.value = false;
    })
    .catch(() => {
      errorMessage.value = "Album creation failed";
      albumCreationInProgress.value = false;
    });
}
</script>
<template>
  <v-overlay
    :model-value="props.modelValue"
    @update:model-value="
      (newValue) => {
        emit('update:modelValue', newValue);
      }
    "
    class="align-center justify-center"
  >
    <v-card title="New Album" min-width="500">
      <v-form
        validate-on="input"
        v-model="isFormValid"
        @submit.prevent="createAlbumSubmit"
        class="my-2"
        data-test-id="create-album-form"
      >
        <v-text-field
          v-model="albumName"
          :rules="albumNameRules"
          label="Album name"
          data-test-id="create-album-form-name-input"
        ></v-text-field>
        <v-btn
          :loading="albumCreationInProgress"
          type="submit"
          block
          class="mt-2"
          text="Create"
          data-test-id="create-album-form-submit-button"
        ></v-btn>
      </v-form>
      <v-alert type="error" v-if="errorMessage.length > 0">
        {{ errorMessage }}
      </v-alert>
    </v-card>
  </v-overlay>
</template>
