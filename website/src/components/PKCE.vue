<script setup lang="ts">
import { UserManager, WebStorageStateStore, OidcClient } from "oidc-client-ts";
import { onMounted } from "vue";
import { useRoute } from "vue-router";
import { v4 } from "uuid";

const oidcClient = new OidcClient({
  client_id: "spa-test",
  authority: "http://localhost:8080",
  redirect_uri: window.location.origin + "/pkce",
  metadataUrl:
    "http://localhost:8080/v1/spa-test/.well-known/openid-configuration",
});

const signinUsingOidcClient = () => {
  let nonce = v4();
  let state = v4();
  localStorage.setItem("nonce", nonce);

  oidcClient
    .createSigninRequest({
      nonce: nonce,
      state: state,
    })
    .then((request) => {
      console.log(request);
      localStorage.setItem("state", request.state.id);
      location.href = request.url;
    })
    .catch((err) => {
      console.log(err);
    });
};

const handlePostLoginUsingOidcClient = () => {
  const route = useRoute();
  if (
    typeof route.query.code !== "string" &&
    typeof route.query.state !== "string"
  ) {
    console.log("code or state query param missing");
    return;
  }
  if (localStorage.getItem("loginInProgress") !== "true") {
    // todo ?
    console.log("hmmmm");
  }
  if (localStorage.getItem("state") !== route.query.state) {
    console.log(localStorage.getItem("state"), route.query.state);
    console.log("state mismatch");
    return;
  }
  console.log(location.href);
  oidcClient
    .processSigninResponse(location.href)
    .then((response) => {
      console.log(response);
    })
    .catch((err) => {
      console.log(err);
    });
};

var userManager = new UserManager({
  userStore: new WebStorageStateStore(),
  authority: "http://localhost:8080",
  metadataUrl:
    "http://localhost:8080/v1/spa-test/.well-known/openid-configuration",
  client_id: "spa-test",
  redirect_uri: window.location.origin + "/pkce",
  response_type: "code",
  scope: "openid profile",
  post_logout_redirect_uri: window.location.origin,
  // silent_redirect_uri: window.location.origin + "/static/silent-renew.html",
  accessTokenExpiringNotificationTimeInSeconds: 10,
  automaticSilentRenew: true,
  filterProtocolClaims: true,
  loadUserInfo: true,
});

const signinUsingUserManager = () => {
  userManager
    .signinRedirect()
    .then((response) => {
      console.log(response);
    })
    .catch((err) => {
      console.log(err);
    });
};

const handlePostLoginUsingUserManager = () => {
  userManager
    .signinRedirectCallback()
    .then((response) => {
      console.log(response);
    })
    .catch((err) => {
      console.log(err);
    });
};
onMounted(() => {
  // handlePostLoginUsingOidcClient();
  handlePostLoginUsingUserManager();
});
</script>

<template>
  <div>
    <v-btn @click.stop="signinUsingUserManager"> login </v-btn>
  </div>
</template>
