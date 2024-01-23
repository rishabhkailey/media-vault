<script setup lang="ts">
import { userManager } from "@/js/auth";
import { signinUsingUserManager } from "@/js/auth";
import AppLogoButton from "@/components/Logo/AppLogoButton.vue";

const logIn = () => {
  signinUsingUserManager(userManager, true);
};
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
          <v-btn class="bg-primary mx-2" @click.stop="logIn">
            <v-icon icon="mdi-login" />
            Sign In
          </v-btn>
        </v-row>
      </v-col>
    </v-row>
  </v-app-bar>
</template>
