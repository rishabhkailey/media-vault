<script setup lang="ts">
import { computed, ref } from "vue";
import { useStore } from "vuex";
// import { useRouter, useRoute } from "vue-router";
import {
  SET_LOGGED_IN_USERINFO_ACTION,
  LOGOUT_ACTION,
} from "@/store/actions-type";
import LoadingAnimation from "./LoadingAnimation.vue";
// import axios from "axios";

// todo types
// const router = useRouter();
// const route = useRoute();

const store = useStore();
const loading = ref(false);
const error = ref(false);
const authenticated = computed(() => store.getters.authenticated);
const userName = computed(() => store.getters.userName);
const email = computed(() => store.getters.email);

const logOut = () => {
  loading.value = true;
  store
    .dispatch(LOGOUT_ACTION)
    .then((message) => {
      console.log(message, store.state);
      loading.value = false;
    })
    .catch((message) => {
      console.log(message);
      loading.value = false;
      error.value = true;
    });
};
const logIn = () => {
  // vue router will handle it and request will not go to backend
  // router.push({
  //   path: "/v1/login",
  //   query: {
  //     returnUri: currentUri,
  //   },
  //   force: true,
  // });
  let currentUri = window.location.href;
  // base/server URL from env/config?
  let baseURL = window.location.origin;
  try {
    window.location.href = new URL(
      `?returnUri=${currentUri}`,
      new URL("/v1/login", baseURL)
    ).toString();
  } catch (_) {
    error.value = true;
  }

  // router.go()
};

loading.value = true;
store
  .dispatch(SET_LOGGED_IN_USERINFO_ACTION)
  .then((message) => {
    console.log(message);
    loading.value = false;
  })
  .catch((message) => {
    console.log(message);
    loading.value = false;
    error.value = true;
  });
</script>

<template>
  <div>
    <ol>
      <li>{{ loading }}</li>
      <li>{{ error }}</li>
      <li>{{ authenticated }}</li>
      <li>{{ userName }}</li>
      <li>{{ email }}</li>
    </ol>
    <LoadingAnimation v-if="loading" />
    <v-btn v-if="authenticated" @click.stop="logOut"> Log Out </v-btn>
    <v-btn v-else @click.stop="logIn"> Log In </v-btn>
  </div>
</template>
