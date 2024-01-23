<script setup lang="ts">
import { ref } from "vue";
import AppLogoButton from "@/components/Logo/AppLogoButton.vue";
import { useUserInfoStore } from "@/piniaStore/userInfo";

const emits = defineEmits<{
  (e: "success", encryptionKey: string): void;
}>();

const { validateEncryptionKey } = useUserInfoStore();

const encryptionKey = ref("");
const warningMessage = ref("");
const errorMessage = ref("");

const isFormValid = ref(false);

function submitHandler() {
  if (isFormValid.value == false || isFormValid.value == null) {
    return;
  }
  if (validateEncryptionKey(encryptionKey.value)) {
    emits("success", encryptionKey.value);
    return;
  }
  errorMessage.value = "wrong encryption key";
}
</script>
<template>
  <v-sheet
    elevation="2"
    max-height="75vh"
    min-height="50vh"
    max-width="50vw"
    min-width="40vw"
    class="d-flex ma-0 pt-0"
  >
    <v-col>
      <v-row justify="center" align="center" class="pa-5">
        <AppLogoButton :redirect="false" />
      </v-row>
      <v-row align="stretch" justify="center">
        <v-col cols="12">
          <v-form
            @submit.prevent="submitHandler"
            v-model="isFormValid"
            lazy-validation
          >
            <v-row class="px-4 py-2">
              <v-text-field
                class="input-fields"
                v-model="encryptionKey"
                label="encryption key"
                type="password"
                required
                :rules="[
                  (value: string) => {
                    if (value.length > 0) {
                      return true;
                    }
                    return 'required';
                  },
                ]"
                autocomplete="on"
                prepend-inner-icon="mdi-lock-outline"
              ></v-text-field>
            </v-row>
            <v-row>
              <v-alert
                v-if="warningMessage.length > 0"
                type="warning"
                :text="warningMessage"
              />
              <v-alert
                v-if="errorMessage.length > 0"
                type="error"
                :text="errorMessage"
              />
            </v-row>
            <v-row justify="center" class="px-4 py-2">
              <v-btn
                type="submit"
                color="primary"
                prepend-icon="mdi-lock-open-outline"
              >
                Unlock
              </v-btn>
            </v-row>
          </v-form>
        </v-col>
      </v-row>
    </v-col>
  </v-sheet>
</template>
