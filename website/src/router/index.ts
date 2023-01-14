import { createRouter, createWebHistory } from "vue-router";
import NotFound from "../views/NotFoundView.vue";
import ErrorScreenView from "../views/ErrorScreenView.vue";
import TestVideoScreen from "../views/TestVideoScreen.vue";
import FileUploadView from "../views/FileUploadView.vue";

// todo pages without redirect from vue should be lazy loaded on external/server redirect
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:pathMatch(.*)*",
      name: "home",
      component: NotFound,
    },
    {
      path: "/error",
      name: "errorscreen",
      component: ErrorScreenView,
    },
    {
      path: "/testVideo",
      name: "testVideoScreen",
      component: TestVideoScreen,
    },
    {
      path: "/upload",
      name: "FileUploadView",
      component: FileUploadView,
    },
  ],
});

export default router;
