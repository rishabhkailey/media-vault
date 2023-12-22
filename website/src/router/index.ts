import { createRouter, createWebHistory } from "vue-router";
import NotFound from "@/views/NotFoundView.vue";
import HomeView from "@/views/Root.vue";
import ErrorScreenView from "@/views/ErrorScreenView.vue";
import PKCEVue from "@/views/PKCE.vue";
import TestView from "@/views/TestView.vue";
import SearchView from "@/views/SearchView.vue";
import HomePageVue from "@/components/HomePage.vue";
import AlbumsPageVue from "@/components/Album/AlbumsPage.vue";
import AlbumMediaThumbnailsVue from "@/components/Album/AlbumMediaThumbnails.vue";
import AboutPage from "@/views/AboutPage.vue";
import UserOnboarding from "@/views/UserOnboarding.vue";
import EnterEncryptionKey from "@/views/EnterEncryptionKey.vue";
import AllMediaPreviewVue from "@/components/MediaPreview/AllMediaPreview.vue";
import SearchMediaPreviewVue from "@/components/MediaPreview/SearchMediaPreview.vue";
import AlbumMediaPreviewVue from "@/components/MediaPreview/AlbumMediaPreview.vue";
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
      component: NotFound,
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
              component: AllMediaPreviewVue,
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
              component: SearchMediaPreviewVue,
            },
          ],
        },
        {
          path: "/albums",
          component: AlbumsPageVue,
          name: "Albums",
        },
        {
          path: "/album/:album_id",
          component: AlbumMediaThumbnailsVue,
          name: "Album",
          children: [
            {
              path: "/album/:album_id/media/:media_id/index/:index",
              name: "AlbumMediaPreview",
              component: AlbumMediaPreviewVue,
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
      path: "/onboarding",
      name: "onboarding",
      component: UserOnboarding,
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
