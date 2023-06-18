<script setup lang="ts">
import { inject, ref } from "vue";
import { userManagerKey } from "@/symbols/injectionSymbols";
import { signinUsingUserManager } from "@/utils/auth";
import type { UserManager } from "oidc-client-ts";
import LogoButton from "@/components/Logo/LogoButton.vue";

// todo error and error message pop up
const error = ref(false);

const userManager: UserManager | undefined = inject(userManagerKey);

const logIn = () => {
  if (userManager === undefined) {
    console.error("userManager not defined");
    error.value = true;
    return;
  }
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
            <LogoButton />
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
