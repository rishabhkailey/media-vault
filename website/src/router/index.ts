import { createRouter, createWebHistory } from "vue-router";
import NotFound from "@/views/NotFoundView.vue";
import HomeView from "@/views/HomeView.vue";
import ErrorScreenView from "@/views/ErrorScreenView.vue";
import PKCEVue from "@/views/PKCE.vue";
import TestView from "@/views/TestView.vue";
import SearchView from "@/views/SearchView.vue";
import HomePageVue from "@/components/HomePage.vue";
import AlbumsPageVue from "@/components/Album/AlbumsPage.vue";
import AlbumMediaThumbnailsVue from "@/components/Album/AlbumMediaThumbnails.vue";
// todo pages without redirect from vue should be lazy loaded on external/server redirect
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: NotFound,
    },
    {
      path: "/",
      component: HomeView,
      children: [
        {
          path: "",
          component: HomePageVue,
          name: "Home",
        },
        {
          path: "search/:query",
          component: SearchView,
          name: "search",
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
        },
      ],
    },
    {
      path: "/error",
      name: "errorscreen",
      component: ErrorScreenView,
    },
    // {
    //   path: "/testVideo",
    //   name: "testVideoScreen",
    //   component: TestVideoScreen,
    // },
    // {
    //   path: "/upload",
    //   name: "FileUploadView",
    //   component: FileUploadView,
    // },
    // {
    //   path: "/testScreen",
    //   name: "testScreen",
    //   component: TestScreenViewVue,
    // },
    // {
    //   path: "/testImageUpload",
    //   name: "testImageUpload",
    //   component: TestImageUploadScreen,
    // },
    // {
    //   path: "/testVideoUpload",
    //   name: "testVideoUpload",
    //   component: TestVideoUploadScreen,
    // },
    // {
    //   path: "/decrypt",
    //   name: "decrypt",
    //   component: WebWorkerModifyResponseView,
    // },
    // {
    //   path: "/chunkUpload",
    //   name: "chunkUpload",
    //   component: ChunkedUploadFormView,
    // },
    // {
    //   path: "/encryptedChunkUpload",
    //   name: "encryptedChunkUpload",
    //   component: EncryptedChunkedUploadFormView,
    // },
    // {
    //   path: "/encryptedChunkUploadUsingTs",
    //   name: "encryptedChunkUploadUsingTs",
    //   component: EncryptedChunkedUploadUsingTsFormViewVue,
    // },
    // {
    //   path: "/encryptedDownload",
    //   name: "encryptedDownload",
    //   component: TestEncryptedFileDownload,
    // },
    // {
    //   path: "/encryptedDownloadUsingTs",
    //   name: "encryptedDownloadUsingTs",
    //   component: TestEncryptedFileDownloadUsingTs,
    // },
    // {
    //   path: "/encryptedVideoPlay",
    //   name: "encryptedVideoPlay",
    //   component: EncryptedVideoPlayView,
    // },
    // {
    //   path: "/videoThumbnailView",
    //   name: "videoThumbnailView",
    //   component: VideoThumbnailView,
    // },
    {
      path: "/pkce",
      name: "pkce",
      component: PKCEVue,
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
