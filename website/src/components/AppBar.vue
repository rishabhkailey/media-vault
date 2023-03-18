<script setup lang="ts">
import SelectFileButton from "@/components/SelectFileButton.vue";
import UploadFilesDialog from "@/components/UploadFilesDialog.vue";
import { computed, inject, ref } from "vue";
import { useStore } from "vuex";
import { RESET_USERINFO_ACTION } from "@/store/actions-type";
import { userManagerKey } from "@/symbols/injectionSymbols";
import type { UserManager } from "oidc-client-ts";
import { signinUsingUserManager } from "@/utils/auth";

const search = ref("");
const searchInputRules: Array<any> = [];

const store = useStore();
const loading = ref(false);
// todo error and error message pop up
const error = ref(false);
const authenticated = computed(() => store.getters.authenticated);
const userName = computed(() => store.getters.userName);
const email = computed(() => store.getters.email);

const selectedFiles = ref<Array<File>>([]);
const uploadFilesDialogModel = ref(false);
const userManager: UserManager | undefined = inject(userManagerKey);

const logOut = () => {
  loading.value = true;
  if (userManager === undefined) {
    console.error("userManager not defined");
    return;
  }
  userManager
    ?.revokeTokens(["access_token", "refresh_token"])
    .then(() => {
      store.dispatch(RESET_USERINFO_ACTION).catch(() => {
        error.value = true;
      });
      console.log("token revoked");
    })
    .catch((err) => {
      console.log(err);
    });
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
</script>

<template>
  <v-app-bar :rounded="false">
    <v-row class="d-flex align-center ml-2">
      <!-- start -->
      <v-col class="d-flex flex-row justify-start">
        <v-toolbar-title>TODO</v-toolbar-title>
      </v-col>
      <!-- mid -->
      <v-col class="d-flex flex-row justify-center">
        <v-form class="d-flex flex-grow-1">
          <v-text-field
            :clearable="true"
            clear-icon="mdi-close"
            append-inner-icon="mdi-magnify"
            v-model="search"
            :rules="searchInputRules"
            label="search"
            :hide-details="true"
          ></v-text-field>
        </v-form>
      </v-col>
      <!-- end -->
      <v-col>
        <v-row class="d-flex flex-row justify-end align-center mr-2">
          <div>
            <SelectFileButton
              label="upload"
              prepend-icon="mdi-upload"
              @select="uploadFiles"
            />
          </div>
          <div v-if="authenticated">
            <v-btn
              :loading="loading"
              color="primary"
              class="mx-2"
              rounded="pill"
            >
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
  </v-app-bar>
  <UploadFilesDialog
    :height="400"
    :width="300"
    v-if="uploadFilesDialogModel"
    v-model="uploadFilesDialogModel"
    :files="selectedFiles"
  />
</template>
