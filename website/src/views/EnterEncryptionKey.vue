<script setup lang="ts">
import { useUserInfoStore } from "@/piniaStore/userInfo";
import EncryptionKeyInput from "../components/EncryptionKeyInput.vue";
import { storeToRefs } from "pinia";
import { useRoute, useRouter, type NavigationFailure } from "vue-router";
const route = useRoute();
const router = useRouter();
const userInfoStore = useUserInfoStore();
const { usableEncryptionKey } = storeToRefs(userInfoStore);

function onValidationEncryptionKey() {
  console.debug("encryption key validated");
  navigator.serviceWorker.ready.then((registration) => {
    if (registration.active === null) {
      // todo handle errors
      throw new Error("service worker not active");
    }
    registration?.active?.postMessage({
      encryptionKey: usableEncryptionKey.value,
    });
  });
  returnToOriginalEndpoint();
}

async function returnToOriginalEndpoint() {
  console.debug("returning to original validated");
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
  // todo return uri
  router.push({
    name: "Home",
  });
  return;
}
</script>

<template>
  <v-container class="w-100 h-100 d-flex justify-center align-center">
    <EncryptionKeyInput @success="onValidationEncryptionKey" />
  </v-container>
</template>
