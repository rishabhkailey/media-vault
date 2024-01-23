<script setup lang="ts">
import type { User } from "oidc-client-ts";
import { onMounted } from "vue";
import { useRouter, type LocationQueryRaw } from "vue-router";
import { handlePostLoginUsingUserManager } from "@/js/auth";
import { useAuthStore } from "@/piniaStore/auth";
import { userManager } from "@/js/auth";
const authStore = useAuthStore();
// not oidc state but state to persist some data after redirect, data will be stored in browser local storage
interface InternalState {
  internalRedirectPath: string;
  internalRedirectQuery: string;
  nonce: string;
}

const handlePostLogin = async () => {
  let router = useRouter();
  handlePostLoginUsingUserManager(userManager)
    .then((user: User) => {
      console.log(user);
      if (user.profile.email === undefined) {
        throw new Error("email missing from the response");
      }
      authStore.setUserInfo(user);
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
      router.push({
        name: "errorscreen",
        query: {
          title: "Sign in Failed",
          message: err,
        },
      });
      console.log(err);
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
