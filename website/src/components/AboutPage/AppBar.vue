<script setup lang="ts">
import AppLogoButton from "@/components/Logo/AppLogoButton.vue";
import AccountButton from "../AppBar/NormalAppBar/AccountButton.vue";
import { userManager } from "@/js/auth";
import { signinUsingUserManager } from "@/js/auth";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/piniaStore/auth";
import { ref } from "vue";
import { logOut } from "../AccountManagement/utils";

const authStore = useAuthStore();
const { authenticated, userName, email } = storeToRefs(authStore);

function logIn() {
  signinUsingUserManager(userManager, true);
}

const loading = ref(false);
function onLogOutClick() {
  loading.value = true;
  logOut().finally(() => {
    loading.value = false;
  });
}
</script>

<template>
  <v-app-bar
    :rounded="false"
    elevation="2"
    class="pa-0 ma-0 d-flex justify-center align-center"
    style="height: inherit"
  >
    <v-row class="d-flex align-center mx-2">
      <!-- start -->
      <v-col class="d-flex flex-row justify-start align-stretch pa-0 ma-0">
        <v-toolbar-title>
          <div>
            <AppLogoButton />
          </div>
          <!-- <v-list-item prepend-icon="mdi-home" title="Home" /> -->
        </v-toolbar-title>
      </v-col>
      <v-col>
        <v-row class="d-flex flex-row flex-nowrap justify-end align-center">
          <AccountButton
            :authenticated="authenticated"
            :loading="loading"
            :user-name="userName"
            :email="email"
            @logout="onLogOutClick"
            @login="logIn"
          />
        </v-row>
      </v-col>
    </v-row>
  </v-app-bar>
</template>
