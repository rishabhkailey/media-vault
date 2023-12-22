import { createApp } from "vue";

import App from "./App.vue";
import "@/assets/main.css";

// Vuetify;
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";
import { createVuetify, type ThemeDefinition } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
import { aliases, mdi } from "vuetify/iconsets/mdi";
import axios from "axios";
import VueAxios from "vue-axios";

import VueVideoPlayer from "@videojs-player/vue";
import { UserManager, WebStorageStateStore } from "oidc-client-ts";
import "./assets/main.css";
import { userManagerKey } from "./symbols/injectionSymbols";
import { createPinia } from "pinia";

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);

const darkTheme: ThemeDefinition = {
  dark: true,
  colors: {
    primary: "#29B6F6", // indigo
    secondary: "#ffca28", // yellow
    accent: "#00bcd4", // cyan
    success: "#4caf50", // green
    info: "#2196f3", // light blue
    warning: "#ffc107", // amber
    error: "#f44336", // red
    background: "#121212", // dark grey
    "secondary-background": "#000000", // black
    surface: "#191919", // grey
    onPrimary: "#ffffff",
    onSecondary: "#000000",
    onAccent: "#000000",
    onSuccess: "#ffffff",
    onInfo: "#ffffff",
    onWarning: "#000000",
    onError: "#ffffff",
    onBackground: "#ffffff",
    onSurface: "#ffffff",
  },
};

const lightTheme: ThemeDefinition = {
  dark: false,
  colors: {
    primary: "#2C3E50",
    secondary: "#E67E22",
    accent: "#1ABC9C",
    success: "#27AE60",
    info: "#3498DB",
    warning: "#F1C40F",
    error: "#E74C3C",
    background: "#F5F5F5",
    "secondary-background": "#F5F5F5",
    surface: "#FFFFFF",
    onPrimary: "#FFFFFF",
    onSecondary: "#FFFFFF",
    onAccent: "#FFFFFF",
    onSuccess: "#FFFFFF",
    onInfo: "#FFFFFF",
    onWarning: "#000000",
    onError: "#FFFFFF",
    onBackground: "#000000",
    onSurface: "#000000",
  },
};

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
  theme: {
    defaultTheme: "darkTheme",
    themes: { darkTheme, lightTheme },
  },
});

const authServiceUrl = import.meta.env.VITE_AUTH_SERVICE_URL;
const authServiceClientID = import.meta.env.VITE_AUTH_SERVICE_CLIENT_ID;
if (authServiceUrl === undefined || authServiceClientID === undefined) {
  throw new Error("VITE_AUTH_SERVICE_URL, VITE_AUTH_SERVICE_CLIENT_ID not set");
}
const userManager = new UserManager({
  authority: import.meta.env.VITE_AUTH_SERVICE_URL,
  metadataUrl: import.meta.env.VITE_AUTH_SERVICE_DISCOVERY_ENDPOINT,
  client_id: import.meta.env.VITE_AUTH_SERVICE_CLIENT_ID,
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
  stateStore: new WebStorageStateStore({ store: window.localStorage }),
  userStore: new WebStorageStateStore({ store: window.localStorage }),
});

app.use(vuetify);
app.use(VueVideoPlayer);
app.use(VueAxios, axios);
app.provide(userManagerKey, userManager);

// router depends on pinia and userMananger
import router from "./router";
app.use(router);

app.mount("#app");
