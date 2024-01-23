import { createRouter, createWebHistory } from "vue-router";
import NotFoundPage from "@/views/NotFoundPage.vue";
import HomeView from "@/views/RootView.vue";
import ErrorScreenView from "@/views/ErrorScreenView.vue";
import PKCEVue from "@/views/PKCE.vue";
import TestView from "@/views/TestView.vue";
import SearchView from "@/views/SearchView.vue";
import HomePageVue from "@/views/HomePage.vue";
import AlbumsList from "@/components/Album/AlbumsList.vue";
import AlbumMediaGrid from "@/components/Album/AlbumMediaGrid.vue";
import AboutPage from "@/views/AboutPage.vue";
import InitialSetup from "@/views/InitialSetup.vue";
import EnterEncryptionKey from "@/views/EnterEncryptionKey.vue";
import UserMediaCarousel from "@/components/MediaCarousel/UserMediaCarousel.vue";
import SearchMediaCarousel from "@/components/MediaCarousel/SearchMediaCarousel.vue";
import AlbumMediaCarousel from "@/components/MediaCarousel/AlbumMediaCarousel.vue";
import { encryptionKeyGaurd, loginGaurd, serviceWrokerGaurd } from "./guards";
// todo pages without redirect from vue should be lazy loaded on external/server redirect
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition !== null) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  },
  routes: [
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: NotFoundPage,
    },
    {
      // tod rename it to root or app?
      path: "/",
      component: HomeView,
      beforeEnter: [loginGaurd, encryptionKeyGaurd, serviceWrokerGaurd],
      children: [
        {
          path: "",
          component: HomePageVue,
          name: "Home",
          children: [
            {
              path: "/media/:media_id/index/:index",
              name: "MediaPreview",
              component: UserMediaCarousel,
            },
          ],
        },
        {
          path: "search/:query",
          component: SearchView,
          name: "search",
          children: [
            {
              path: "/search/:query/media/:media_id/index/:index",
              name: "SearchMediaPreview",
              component: SearchMediaCarousel,
            },
          ],
        },
        {
          path: "/albums",
          component: AlbumsList,
          name: "Albums",
        },
        {
          path: "/album/:album_id",
          component: AlbumMediaGrid,
          name: "Album",
          children: [
            {
              path: "/album/:album_id/media/:media_id/index/:index",
              name: "AlbumMediaPreview",
              component: AlbumMediaCarousel,
            },
          ],
        },
      ],
    },
    {
      path: "/error",
      name: "errorscreen",
      component: ErrorScreenView,
    },
    {
      path: "/encryption-key",
      name: "encryptionKey",
      component: EnterEncryptionKey,
      beforeEnter: [loginGaurd],
    },
    {
      path: "/about",
      name: "about",
      component: AboutPage,
    },
    {
      path: "/initial-setup",
      name: "initialSetup",
      component: InitialSetup,
    },
    {
      path: "/pkce",
      name: "pkce",
      component: PKCEVue,
    },
    {
      path: "/test",
      name: "test",
      component: TestView,
    },
  ],
});

export default router;
