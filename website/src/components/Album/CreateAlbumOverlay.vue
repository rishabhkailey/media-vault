<script setup lang="ts">
import { ref } from "vue";
import { useAlbumStore } from "@/piniaStore/album";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/piniaStore/auth";
const props = defineProps<{
  modelValue: boolean;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

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
  createAlbum(accessToken.value, albumName.value, "")
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
      >
        <v-text-field
          v-model="albumName"
          :rules="albumNameRules"
          label="Album name"
        ></v-text-field>
        <v-btn
          :loading="albumCreationInProgress"
          type="submit"
          block
          class="mt-2"
          text="Create"
        ></v-btn>
      </v-form>
      <v-alert type="error" v-if="errorMessage.length > 0">
        {{ errorMessage }}
      </v-alert>
    </v-card>
  </v-overlay>
</template>
