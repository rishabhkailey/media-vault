<script setup lang="ts">
import type { User } from "oidc-client-ts";
import { onMounted } from "vue";
import { useRouter, type LocationQueryRaw } from "vue-router";
import { handlePostLoginUsingUserManager } from "@/js/auth";
import { useAuthStore } from "@/piniaStore/auth";
import { getUserManager } from "@/js/auth";
import { errorScreenRoute } from "@/router/routesConstants";
const authStore = useAuthStore();
// not oidc state but state to persist some data after redirect, data will be stored in browser local storage
interface InternalState {
  internalRedirectPath: string;
  internalRedirectQuery: string;
  nonce: string;
}

const handlePostLogin = async () => {
  let router = useRouter();
  handlePostLoginUsingUserManager(getUserManager())
    .then((user: User) => {
      authStore.setUserAuthInfo(user);
      let internalState = user.state as InternalState;
      if (
        internalState.internalRedirectPath.length !== 0 ||
        internalState.internalRedirectQuery.length !== 0
      ) {
        // location.href = user.state?.internalRedirectUri;
        let query: LocationQueryRaw = {};
        let searchParams = new URLSearchParams(
          internalState.internalRedirectQuery,
        );
        searchParams.forEach((value, key) => {
          query[key] = value;
        });
        router.replace({
          path: internalState.internalRedirectPath,
          query: query,
        });
      }
    })
    .catch((err) => {
      console.error(err);
      router.push(errorScreenRoute("Sign in Failed", err));
    });
};

onMounted(() => {
  handlePostLogin();
});
</script>

<template>
  <v-card
    style="height: 100%"
    title="Sign In"
    subtitle="in progress"
    text="this is a placeholder component"
  />
</template>
