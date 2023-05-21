<script setup lang="ts">
import type { UserManager, User } from "oidc-client-ts";
import { onMounted, inject } from "vue";
import { useRouter, type LocationQueryRaw } from "vue-router";
import { userManagerKey } from "@/symbols/injectionSymbols";
import {
  handlePostLoginUsingUserManager,
  signinUsingUserManager,
} from "@/utils/auth";
import { useAuthStore } from "@/piniaStore/auth";
const authStore = useAuthStore();
// todo make this a component with slot for ui and callback methods as props, onSuccess, onError, getUserinfo: bool, onGetUserInfoSuccess, onGetUserInfoError
// button type signIn or singOut?

// const oidcClient = new OidcClient({
//   client_id: "spa-test",
//   authority: "http://localhost:8080",
//   redirect_uri: window.location.origin + "/pkce",
//   metadataUrl:
//     "http://localhost:8080/v1/spa-test/.well-known/openid-configuration",
// });

const userManager: UserManager | undefined = inject(userManagerKey);

// not oidc state but state to persist some data after redirect, data will be stored in browser local storage
interface InternalState {
  internalRedirectPath: string;
  internalRedirectQuery: string;
  nonce: string;
}

const login = () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    return;
  }
  signinUsingUserManager(userManager);
};

// todo remove code and state from the url
const handlePostLogin = async () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    return;
  }
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
          internalState.internalRedirectQuery
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
      console.log(err);
      // todo redirect to error page
    });
};

const logout = () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    return;
  }
  userManager
    ?.revokeTokens(["access_token", "refresh_token"])
    .then(() => {
      console.log("token revoked");
    })
    .catch((err) => {
      console.log(err);
    });

  // requires end session endpoint
  // userManager?.signoutPopup();
};
onMounted(() => {
  // handlePostLoginUsingOidcClient();
  handlePostLogin();
});
</script>

<template>
  <div>
    <v-btn @click.stop="login"> login </v-btn>
    <v-btn @click.stop="logout"> logout </v-btn>
  </div>
</template>
