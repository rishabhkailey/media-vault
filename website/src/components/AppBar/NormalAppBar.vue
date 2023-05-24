<script setup lang="ts">
import { computed, inject, ref } from "vue";
import SelectFileButton from "@/components/SelectFileButton.vue";
import UploadFilesDialog from "@/components/UploadFilesDialog.vue";
import SearchInputField from "@/components/SearchInputField.vue";
import { useAuthStore } from "@/piniaStore/auth";
import { userManagerKey } from "@/symbols/injectionSymbols";
import { signinUsingUserManager } from "@/utils/auth";
import type { UserManager } from "oidc-client-ts";
import axios from "axios";
import { useDisplay } from "vuetify/lib/framework.mjs";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";

const router = useRouter();
const display = useDisplay();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const props = defineProps<{
  navigationBar: boolean;
}>();

const emit = defineEmits<{
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
  signinUsingUserManager(userManager);
};
const uploadFiles = (files: Array<File>) => {
  selectedFiles.value = files;
  uploadFilesDialogModel.value = true;
  console.log(files);
};
console.log(display.mobile.value);

const searchSubmit = () => {
  if (search.value.trim().length === 0) {
    return;
  }
  router.push({
    name: `search`,
    params: {
      query: search.value,
    },
  });
};
</script>

<template>
  <v-row class="d-flex align-center mx-2">
    <!-- start -->
    <v-col
      v-if="smallDisplay"
      @click.stop="
        () => {
          emit('update:navigationBar', !props.navigationBar);
        }
      "
    >
      <v-btn icon="mdi-menu"> </v-btn>
    </v-col>
    <v-col v-else class="d-flex flex-row justify-start">
      <v-toolbar-title>TODO</v-toolbar-title>
    </v-col>
    <!-- mid -->
    <v-col class="d-flex flex-row justify-center">
      <SearchInputField
        v-if="!smallDisplay"
        v-model="search"
        :collapsed="false"
        @submit="searchSubmit"
      />
    </v-col>
    <!-- end -->
    <v-col>
      <v-row class="d-flex flex-row flex-nowrap justify-end align-center">
        <div>
          <SearchInputField
            v-if="smallDisplay"
            v-model="search"
            :collapsed="true"
            @submit="searchSubmit"
          />
        </div>
        <div>
          <SelectFileButton
            label="upload"
            prepend-icon="mdi-upload"
            @select="uploadFiles"
          />
        </div>
        <div v-if="authenticated">
          <v-btn :loading="loading" color="primary" class="mx-2" rounded="pill">
            <v-icon icon="mdi-account" color="primary" size="x-large" />
            <v-menu activator="parent">
              <v-card prepend-icon="mdi-account">
                <template v-slot:title>
                  {{ userName }}
                </template>
                <template v-slot:subtitle>
                  {{ email }}
                </template>
                <template v-slot:actions>
                  <div class="d-flex justify-center flex-grow-1">
                    <v-btn class="bg-primary mx-2" @click.stop="logOut">
                      <v-icon icon="mdi-logout" />
                      Sign Out
                    </v-btn>
                  </div>
                </template>
              </v-card>
            </v-menu>
          </v-btn>
        </div>
        <v-btn v-else class="bg-primary mx-2" @click.stop="logIn">
          <v-icon icon="mdi-login" />
          Sign In
        </v-btn>
      </v-row>
    </v-col>
  </v-row>
  <UploadFilesDialog
    :height="400"
    :width="300"
    v-if="uploadFilesDialogModel"
    v-model="uploadFilesDialogModel"
    :files="selectedFiles"
  />
</template>
