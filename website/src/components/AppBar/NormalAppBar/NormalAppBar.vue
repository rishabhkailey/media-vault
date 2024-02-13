<script setup lang="ts">
import { computed, ref } from "vue";
import FileUploadDialog from "@/components/FileUpload/FileUploadDialog.vue";
import FloatingWindow from "@/components/Modals/FloatingWindow.vue";
import { useAuthStore } from "@/piniaStore/auth";
import { userManager } from "@/js/auth";
import { signinUsingUserManager } from "@/js/auth";
import { useDisplay } from "vuetify";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import MobileAppBar from "./Size/MobileAppBar.vue";
import DesktopAppBar from "./Size/DesktopAppBar.vue";
import { logOut } from "@/components/AccountManagement/utils";
import { searchRoute } from "@/router/routesConstants";
import { UPLOAD_WINDOW_Z_INDEX } from "@/js/constants/z-index";

const router = useRouter();
const display = useDisplay();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value,
);

const props = defineProps<{
  sidebarOpen: boolean;
}>();

const emits = defineEmits<{
  (e: "update:sidebarOpen", value: boolean): void;
  (e: "selectFilesForUpload", selectedFiles: Array<File>): void;
}>();

const search = ref("");
const loading = ref(false);

const authStore = useAuthStore();
const { authenticated, userName, email } = storeToRefs(authStore);

const selectedFiles = ref<Array<File>>([]);
const FileUploadDialogModel = ref(false);

const logIn = () => {
  signinUsingUserManager(userManager, false);
};
function onLogOutClick() {
  loading.value = true;
  logOut().finally(() => {
    loading.value = false;
  });
}
const uploadFiles = (files: Array<File>) => {
  selectedFiles.value = files;
  FileUploadDialogModel.value = true;
  console.log(files);
};
console.log(display.mobile.value);

const searchSubmit = (query: string) => {
  if (search.value.trim().length === 0) {
    return;
  }
  router.push(searchRoute(query));
};
</script>

<template>
  <v-row v-if="smallDisplay">
    <MobileAppBar
      :navigation-bar="props.sidebarOpen"
      v-model:search-query="search"
      :authenticated="authenticated"
      :email="email"
      :user-name="userName"
      :user-auth-loading="loading"
      @search-submit="(query) => searchSubmit(query)"
      @select-files-for-upload="(files) => uploadFiles(files)"
      @update:navigation-bar="(value) => emits('update:sidebarOpen', value)"
      @login="logIn"
      @logout="onLogOutClick"
    />
  </v-row>
  <v-row v-else>
    <DesktopAppBar
      :navigation-bar="true"
      v-model:search-query="search"
      :authenticated="authenticated"
      :email="email"
      :user-name="userName"
      :user-auth-loading="loading"
      @search-submit="(query) => searchSubmit(query)"
      @select-files-for-upload="(files) => uploadFiles(files)"
      @login="logIn"
      @logout="onLogOutClick"
    />
  </v-row>
  <!-- todo move this to somewhere else? -->
  <!-- may be to selectButtom component? -->
  <FloatingWindow
    :z-index="UPLOAD_WINDOW_Z_INDEX"
    v-model="FileUploadDialogModel"
    :bottom="10"
    :right="10"
  >
    <FileUploadDialog
      :model-value="true"
      :files="selectedFiles"
      @close="() => (FileUploadDialogModel = false)"
      height="40vh"
      width="30vh"
    />
  </FloatingWindow>
</template>
