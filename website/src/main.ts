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

import VueVideoPlayer from "@videojs-player/vue";

import "./assets/main.css";

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

app.use(router);
app.use(vuetify);
app.use(VueVideoPlayer);
app.use(VueAxios, axios);

app.mount("#app");
