<script setup lang="ts">
import { onMounted, ref } from "vue";
import LogoButton from "../Logo/LogoButton.vue";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { useRouter, type NavigationFailure, useRoute } from "vue-router";

const { postUserInfo } = useUserInfoStore();
const router = useRouter();
const route = useRoute();

const preferedTimezone = ref("");
const availableTimezones = ref<Array<string>>([]);
const encryptionKey = ref("");
const confirmEncryptionKey = ref("");
const errorMessage = ref("");
const warningMessage = ref("");

const isFormValid = ref(false);

const encryptionKeyRules: ((value: string) => true | string)[] = [
  (value: string) => {
    if (value.length < 8) {
      return "key should be atleast 8 character long";
    }
    return true;
  },
];

const confirmEncryptionKeyRules: ((value: string) => true | string)[] = [
  (value: string) => {
    if (value === encryptionKey.value) {
      return true;
    }
    return "does not match with the original key";
  },
];

onMounted(() => {
  const browserTimqZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  preferedTimezone.value = browserTimqZone;
  try {
    // @ts-ignore
    const allTimeZones = Intl.supportedValuesOf("timeZone");
    availableTimezones.value = allTimeZones;
  } catch (err) {
    availableTimezones.value = [browserTimqZone];
    warningMessage.value =
      "failed to load all timezone values. this will not allow you to select a different timezone.";
  }
});

async function returnToOriginalEndpoint() {
  const returnUriQuery = Array.isArray(route.query.return_uri)
    ? route.query.return_uri[0]
    : route.query.return_uri;
  let returnUri = "";
  if (returnUriQuery !== null) {
    returnUri = returnUriQuery;
  }
  let error: void | NavigationFailure | Error | undefined;
  try {
    error = await router.push(returnUri);
    if (!(error instanceof Error)) {
      return;
    }
  } catch (err) {
    if (err instanceof Error) {
      error = err;
    } else {
      error = new Error("unexpected error while returning to original page.");
    }
  }
  console.error(error);
  router.push({
    name: "Home",
  });
  return;
}

function submitHandler() {
  if (isFormValid.value == false || isFormValid.value == null) {
    return;
  }
  postUserInfo(preferedTimezone.value, encryptionKey.value)
    .then((ok) => {
      if (ok) {
        returnToOriginalEndpoint();
        return;
      }
      errorMessage.value = "request failed.";
    })
    .catch((err) => {
      errorMessage.value = "request failed. \n" + err;
    });
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
        <LogoButton :redirect="false" />
      </v-row>
      <v-row justify="center">
        <v-col>
          <v-row justify="center">
            <v-card-title> Set a Encryption key </v-card-title>
          </v-row>
        </v-col>
      </v-row>
      <v-row align="stretch" justify="center">
        <v-col cols="12">
          <v-form
            @submit.prevent="submitHandler"
            v-model="isFormValid"
            lazy-validation
          >
            <v-row class="px-4 py-2">
              <v-autocomplete
                class="input-fields"
                v-model="preferedTimezone"
                label="timezone"
                required
                prepend-inner-icon="mdi-map-clock-outline"
                :items="availableTimezones"
              />
            </v-row>
            <v-row class="px-4 py-2">
              <v-text-field
                class="input-fields"
                v-model="encryptionKey"
                :rules="encryptionKeyRules"
                label="encryption key"
                type="password"
                required
                autocomplete="on"
                prepend-inner-icon="mdi-lock-outline"
              ></v-text-field>
            </v-row>
            <v-row class="px-4 py-2">
              <v-text-field
                class="input-fields"
                v-model="confirmEncryptionKey"
                :rules="confirmEncryptionKeyRules"
                label="confirm encryption key"
                type="password"
                required
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
            <v-row justify="center">
              <v-alert
                type="info"
                title="Please Note"
                text="the above encryption key will be used to encrypt your files so make sure to set it to secure password. make sure to remember it if you forget your encryption key your will not be able read any of your uploaded data."
                variant="tonal"
              ></v-alert>
            </v-row>
            <v-row justify="center" class="px-4 py-2">
              <v-btn type="submit" color="primary"> Confirm </v-btn>
            </v-row>
          </v-form>
        </v-col>
      </v-row>
    </v-col>
  </v-sheet>
</template>
