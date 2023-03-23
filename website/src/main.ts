import { createApp } from "vue";

import App from "./App.vue";
import "@/assets/main.css";

// Vuetify;
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
import { aliases, mdi } from "vuetify/iconsets/mdi";
import axios from "axios";
import VueAxios from "vue-axios";
import router from "./router";
import store from "./store";

import VueVideoPlayer from "@videojs-player/vue";
import { UserManager, WebStorageStateStore } from "oidc-client-ts";
import "./assets/main.css";
import { userManagerKey } from "./symbols/injectionSymbols";

const app = createApp(App);

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi,
    },
  },
});

const userManager = new UserManager({
  userStore: new WebStorageStateStore(),
  authority: "http://localhost:8080",
  metadataUrl:
    "http://localhost:8080/v1/spa-test/.well-known/openid-configuration",
  client_id: "spa-test",
  redirect_uri: window.location.origin + "/pkce",
  response_type: "code",
  scope: "openid profile email user",
  post_logout_redirect_uri: window.location.origin,
  // silent_redirect_uri: window.location.origin + "/static/silent-renew.html",
  accessTokenExpiringNotificationTimeInSeconds: 10,
  automaticSilentRenew: true,
  // if true it removes the nonce
  filterProtocolClaims: false,
  loadUserInfo: true,
});

app.use(router);
app.use(store);
app.use(vuetify);
app.use(VueVideoPlayer);
app.use(VueAxios, axios);
app.provide(userManagerKey, userManager)
app.mount("#app");
