import { createStore } from "vuex";
import { authModule } from "./modules/auth";
import { mediaModule } from "./modules/media";
// Create a new store instance.
const store = createStore({
  modules: {
    authModule: authModule,
    mediaModule: mediaModule,
  },
});

export default store;
