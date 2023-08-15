<script setup lang="ts">
import { computed, inject, ref } from "vue";
import SelectFileButton from "@/components/SelectFileButton.vue";
import UploadFilesDialog from "@/components/UploadingFiles/UploadingFilesDialog.vue";
import SearchInputField from "@/components/SearchInputField.vue";
import FloatingWindow from "@/components/FloatingWindow/FloatingWindow.vue";
import LogoButton from "@/components/Logo/LogoButton.vue";
import { useAuthStore } from "@/piniaStore/auth";
import { userManagerKey } from "@/symbols/injectionSymbols";
import { signinUsingUserManager } from "@/js/auth";
import type { UserManager } from "oidc-client-ts";
import axios from "axios";
import { useDisplay } from "vuetify/lib/framework.mjs";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import ProfileButton from "./ProfileButton.vue";
import MobileAppBar from "./Size/MobileAppBar.vue";
import DesktopAppBar from "./Size/DesktopAppBar.vue";

const router = useRouter();
const display = useDisplay();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const props = defineProps<{
  navigationBar: boolean;
}>();

const emits = defineEmits<{
  (e: "update:navigationBar", value: boolean): void;
}>();

const search = ref("");
const loading = ref(false);
// todo error and error message pop up
const error = ref(false);
const authStore = useAuthStore();
const { authenticated, userName, email } = storeToRefs(authStore);

const selectedFiles = ref<Array<File>>([]);
const uploadFilesDialogModel = ref(false);
const userManager: UserManager | undefined = inject(userManagerKey);

const logOut = async () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    return;
  }
  loading.value = true;
  try {
    await userManager.revokeTokens(["access_token", "refresh_token"]);
    authStore.reset();
    await userManager.removeUser();
    let response = await axios.post("/v1/terminateSession");
    if (response.status !== 200) {
      console.log("terminate session failed");
    }
  } catch (err) {
    console.log("logout failed ", err);
  } finally {
    loading.value = false;
  }
};
const logIn = () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    error.value = true;
    return;
  }
  signinUsingUserManager(userManager, false);
};
const uploadFiles = (files: Array<File>) => {
  selectedFiles.value = files;
  uploadFilesDialogModel.value = true;
  console.log(files);
};
console.log(display.mobile.value);

const searchSubmit = (query: string) => {
  if (search.value.trim().length === 0) {
    return;
  }
  router.push({
    name: `search`,
    params: {
      query: query,
    },
  });
};
// todo different commonent for mobile app bar instead of if else
</script>

<template>
  <v-row v-if="smallDisplay">
    <MobileAppBar
      :navigation-bar="props.navigationBar"
      :search-query="search"
      :authenticated="authenticated"
      :email="email"
      :user-name="userName"
      :user-auth-loading="loading"
      @search-submit="(query) => searchSubmit(query)"
      @upload-files="(files) => uploadFiles(files)"
      @update:navigation-bar="(value) => emits('update:navigationBar', value)"
      @login="logIn"
      @logout="logOut"
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
      @upload-files="(files) => uploadFiles(files)"
      @login="logIn"
      @logout="logOut"
    />
  </v-row>
  <!-- todo move this to somewhere else? -->
  <!-- may be to selectButtom component? -->
  <FloatingWindow v-model="uploadFilesDialogModel" :bottom="10" :right="10">
    <UploadFilesDialog
      :model-value="true"
      :files="selectedFiles"
      @close="() => (uploadFilesDialogModel = false)"
      height="40vh"
      width="30vh"
    />
  </FloatingWindow>
</template>
@/js/auth