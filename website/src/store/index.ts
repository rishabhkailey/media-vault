import { createStore } from "vuex";
import { authModule } from "./modules/auth";
import { mediaModule } from "./modules/media";
import {
  MODULE_NAME as searchModuleName,
  searchModule,
} from "./modules/search";
// Create a new store instance.
const store = createStore({
  modules: {
    authModule: authModule,
    mediaModule: mediaModule,
    [searchModuleName]: searchModule,
  },
});

export default store;
