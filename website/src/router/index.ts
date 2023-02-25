import { createRouter, createWebHistory } from "vue-router";
import NotFound from "@/views/NotFoundView.vue";
import ErrorScreenView from "@/views/ErrorScreenView.vue";
import TestVideoScreen from "@/views/TestVideoScreen.vue";
import FileUploadView from "@/views/FileUploadView.vue";
import TestScreenViewVue from "@/views/TestScreenView.vue";
import TestImageUploadScreen from "@/views/TestImageUploadScreen.vue";
import TestVideoUploadScreen from "@/views/TestVideoUploadScreen.vue";
import WebWorkerModifyResponseView from "@/views/WebWorkerModifyResponseView.vue";
import ChunkedUploadFormView from "@/views/ChunkedUploadFormView.vue";
import EncryptedChunkedUploadFormView from "@/views/EncryptedChunkedUploadFormView.vue";
import TestEncryptedFileDownload from "@/views/TestEncryptedFileDownload.vue";
import EncryptedChunkedUploadUsingTsFormViewVue from "@/views/EncryptedChunkedUploadUsingTsFormView.vue";
import TestEncryptedFileDownloadUsingTs from "@/views/TestEncryptedFileDownloadUsingTs.vue";
import EncryptedVideoPlayView from "@/views/EncryptedVideoPlayView.vue"; 
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
    {
      path: "/testScreen",
      name: "testScreen",
      component: TestScreenViewVue,
    },
    {
      path: "/testImageUpload",
      name: "testImageUpload",
      component: TestImageUploadScreen,
    },
    {
      path: "/testVideoUpload",
      name: "testVideoUpload",
      component: TestVideoUploadScreen,
    },
    {
      path: "/decrypt",
      name: "decrypt",
      component: WebWorkerModifyResponseView,
    },
    {
      path: "/chunkUpload",
      name: "chunkUpload",
      component: ChunkedUploadFormView,
    },
    {
      path: "/encryptedChunkUpload",
      name: "encryptedChunkUpload",
      component: EncryptedChunkedUploadFormView,
    },
    {
      path: "/encryptedChunkUploadUsingTs",
      name: "encryptedChunkUploadUsingTs",
      component: EncryptedChunkedUploadUsingTsFormViewVue,
    },
    {
      path: "/encryptedDownload",
      name: "encryptedDownload",
      component: TestEncryptedFileDownload,
    },
    {
      path: "/encryptedDownloadUsingTs",
      name: "encryptedDownloadUsingTs",
      component: TestEncryptedFileDownloadUsingTs,
    },
    {
      path: "/encryptedVideoPlay",
      name: "encryptedVideoPlay",
      component: EncryptedVideoPlayView,
    }
  ],
});

export default router;
