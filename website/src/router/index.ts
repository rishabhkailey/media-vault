import { createRouter, createWebHistory } from "vue-router";
import NotFoundPage from "@/views/NotFoundPage.vue";
import HomeView from "@/views/RootView.vue";
import ErrorScreenView from "@/views/ErrorScreenView.vue";
import PKCEVue from "@/views/PKCE.vue";
import TestView from "@/views/TestView.vue";
import SearchView from "@/views/SearchView.vue";
import HomePageVue from "@/views/HomePage.vue";
import AlbumsGrid from "@/components/Album/AlbumsGrid.vue";
import AboutPage from "@/views/AboutPage.vue";
import InitialSetup from "@/views/InitialSetup.vue";
import EnterEncryptionKey from "@/views/EnterEncryptionKey.vue";
import UserMediaCarousel from "@/components/MediaCarousel/UserMediaCarousel.vue";
import SearchMediaCarousel from "@/components/MediaCarousel/SearchMediaCarousel.vue";
import AlbumMediaCarousel from "@/components/MediaCarousel/AlbumMediaCarousel.vue";
import {
  configurationGaurd,
  encryptionKeyGaurd,
  loginGaurd,
  serviceWrokerGaurd,
} from "./guards";
import {
  ABOUT_ROUTE_NAME,
  ALBUMS_ROUTE_NAME,
  ALBUM_MEDIA_PREVIEW_ROUTE_NAME,
  ALBUM_ROUTE_NAME,
  ENTER_ENCRYPTION_KEY_ROUTE_NAME,
  ERROR_SCREEN_ROUTE_NAME,
  HOME_ROUTE_NAME,
  INITIAL_SETUP_ROUTE_NAME,
  MEDIA_PREVIEW_ROUTE_NAME,
  NOT_FOUND_ROUTE_NAME,
  PKCE_ROUTE_NAME,
  SEARCH_MEDIA_PREVIEW_ROUTE_NAME,
  SEARCH_ROUTE_NAME,
} from "./routesConstants";
import AlbumMedia from "@/views/AlbumMedia.vue";

// todo pages without redirect from vue should be lazy loaded on external/server redirect
export const router = createRouter({
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
      name: NOT_FOUND_ROUTE_NAME,
      component: NotFoundPage,
    },
    {
      // tod rename it to root or app?
      path: "/",
      component: HomeView,
      beforeEnter: [
        configurationGaurd,
        loginGaurd,
        encryptionKeyGaurd,
        serviceWrokerGaurd,
      ],
      children: [
        {
          path: "",
          component: HomePageVue,
          name: HOME_ROUTE_NAME,
          children: [
            {
              path: "/media/:media_id/index/:index",
              name: MEDIA_PREVIEW_ROUTE_NAME,
              component: UserMediaCarousel,
            },
          ],
        },
        {
          path: "search/:query",
          component: SearchView,
          name: SEARCH_ROUTE_NAME,
          children: [
            {
              path: "/search/:query/media/:media_id/index/:index",
              name: SEARCH_MEDIA_PREVIEW_ROUTE_NAME,
              component: SearchMediaCarousel,
            },
          ],
        },
        {
          path: "/albums",
          component: AlbumsGrid,
          name: ALBUMS_ROUTE_NAME,
        },
        {
          path: "/album/:album_id",
          component: AlbumMedia,
          name: ALBUM_ROUTE_NAME,
          children: [
            {
              path: "/album/:album_id/media/:media_id/index/:index",
              name: ALBUM_MEDIA_PREVIEW_ROUTE_NAME,
              component: AlbumMediaCarousel,
            },
          ],
        },
      ],
    },
    {
      path: "/error",
      name: ERROR_SCREEN_ROUTE_NAME,
      component: ErrorScreenView,
    },
    {
      path: "/encryption-key",
      name: ENTER_ENCRYPTION_KEY_ROUTE_NAME,
      component: EnterEncryptionKey,
      beforeEnter: [configurationGaurd, loginGaurd],
    },
    {
      path: "/about",
      name: ABOUT_ROUTE_NAME,
      beforeEnter: [configurationGaurd],
      component: AboutPage,
    },
    {
      path: "/initial-setup",
      name: INITIAL_SETUP_ROUTE_NAME,
      beforeEnter: [configurationGaurd],
      component: InitialSetup,
    },
    {
      path: "/pkce",
      name: PKCE_ROUTE_NAME,
      beforeEnter: [configurationGaurd],
      component: PKCEVue,
    },
    {
      path: "/test",
      name: "test",
      component: TestView,
    },
  ],
});
