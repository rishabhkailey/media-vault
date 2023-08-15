<script setup lang="ts">
import AppBar from "@/components/AppBar/AppBar.vue";
import NavigationBar from "../components/NavigationBar.vue";
import EncryptionKeyInput from "@/components/EncryptionKeyInput.vue";
import { computed, inject, onMounted, ref, onBeforeMount, watch } from "vue";
import { userManagerKey } from "@/symbols/injectionSymbols";
import type { UserManager } from "oidc-client-ts";
import decryptWorker from "@/worker/dist/bundle.js?url";
// todo if not authenticated redirect to some different page
// maybe /about
import { useDisplay } from "vuetify";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/piniaStore/auth";
import { useLoadingStore } from "@/piniaStore/loading";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { EncryptionKeyChannelClient } from "@/js/channels/encryptionKey";

const authStore = useAuthStore();
const { authenticated } = storeToRefs(authStore);
const display = useDisplay();

const userInfoStore = useUserInfoStore();
const { updateUserInfo } = userInfoStore;
const {
  initRequired: userInfoInitRequired,
  encryptionKeyValidated,
  usableEncryptionKey,
} = storeToRefs(userInfoStore);

const { setInitializing } = useLoadingStore();
setInitializing(true);

const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const route = useRoute();
const router = useRouter();
// const initializingRef = ref(true);
const navigationBar = ref(!smallDisplay.value);
const encryptionKeyOverlay = ref(false);
// provide(initializingKey, initializingRef);
const userManager: UserManager | undefined = inject(userManagerKey);

const userInit = () => {
  return new Promise<boolean>((resolve, reject) => {
    if (userManager === undefined) {
      reject(new Error("undefined userManager"));
      return;
    }
    userManager
      .getUser()
      .then((user) => {
        if (user?.expired) {
          resolve(false);
          return;
        }
        if (user === null) {
          resolve(false);
          return;
        }
        authStore.setUserInfo(user);
        resolve(true);
      })
      .catch((err) => {
        reject(err);
        return;
      });
  });
};

const updateOrRegisterServiceWorker = () => {
  return new Promise<ServiceWorker>((resolve, reject) => {
    navigator.serviceWorker.getRegistration().then((registration) => {
      if (registration === undefined) {
        registerServiceWorker()
          .then((serviceWorker) => {
            resolve(serviceWorker);
            return;
          })
          .catch((err) => {
            reject(err);
            return;
          });
      } else {
        updateServiceWorker()
          .then((serviceWorker) => {
            resolve(serviceWorker);
            return;
          })
          .catch((err) => {
            reject(err);
            return;
          });
      }
    });
  });
};

const updateServiceWorker = () => {
  // todo instead delete and install again?
  // because if we are changing the path or file name of worker then this fails
  return new Promise<ServiceWorker>((resolve, reject) => {
    navigator.serviceWorker.getRegistration().then((registration) => {
      if (registration === undefined) {
        reject(true);
        return;
      }
      registration
        .update()
        .then(() => {
          console.log("updated");
          if (registration.active === null) {
            reject(new Error("got null service worker after update"));
            return;
          }
          resolve(registration.active);
          return;
        })
        .catch((err) => {
          reject(err);
          return;
        });
    });
  });
};

// https://github.com/jimmywarting/StreamSaver.js/blob/master/mitm.html#L39
function registerServiceWorker(): Promise<ServiceWorker> {
  return new Promise<ServiceWorker>((resolve, reject) => {
    if ("serviceWorker" in navigator) {
      // unregister the existing service worker
      navigator.serviceWorker
        .register(decryptWorker, {
          scope: "./",
          type: "module",
        })
        .then((swReg) => {
          if (swReg.active !== null) {
            resolve(swReg.active);
            return;
          }
          const swRegTmp = swReg.installing || swReg.waiting;
          if (swRegTmp === null) {
            reject(new Error("got null service worker registration"));
            return;
          }
          let callback: () => void;
          swRegTmp.addEventListener(
            "statechange",
            (callback = () => {
              if (swRegTmp.state === "activated") {
                console.log("registed and activated");
                swRegTmp.removeEventListener("statechange", callback);
                resolve(swRegTmp);
              }
            })
          );
        })
        .catch((err) => {
          reject(err);
          return;
        });
    }
  });
}
async function init(): Promise<boolean> {
  try {
    // todo try catch
    await updateOrRegisterServiceWorker();
    await userInit();
  } catch (error) {
    console.log(error);
    router.push({
      name: "errorscreen",
      query: {
        title: "init failed",
        message: "error message - " + error,
      },
    });
  }
  if (!authenticated.value) {
    // route to about page
    router.push({
      name: "about",
    });
  }
  await updateUserInfo();
  if (userInfoInitRequired.value) {
    router.push({
      name: "onboarding",
    });
    return false;
  }
  if (encryptionKeyValidated.value == false) {
    encryptionKeyOverlay.value = true;
    // router.push({
    //   name: "onboarding",
    // });
    return false;
  }
  setInitializing(false);
  return true;
}

function onValidationEncryptionKey() {
  encryptionKeyOverlay.value = false;
  setInitializing(false);
  navigator.serviceWorker.ready.then((registration) => {
    if (registration.active === null) {
      // todo handle errors
      throw new Error("service worker not active");
    }
    registration?.active?.postMessage({
      encryptionKey: usableEncryptionKey.value,
    });
  });
}

const encryptionKeyChannel = new EncryptionKeyChannelClient(
  usableEncryptionKey.value
);

watch(usableEncryptionKey, () => {
  encryptionKeyChannel.encryptionKey = usableEncryptionKey.value;
});

onMounted(() => {
  init()
    .then(() => {})
    .catch((err) => {
      console.log(err);
      // setInitializing(false);
    });
});

onBeforeMount(() => {});
const test = (value: boolean) => {
  console.log("called", value);
};
</script>

<template>
  <v-container class="pa-0 ma-0" fluid>
    <v-overlay
      v-model="encryptionKeyOverlay"
      class="justify-center align-center"
    >
      <EncryptionKeyInput @success="onValidationEncryptionKey" />
    </v-overlay>
    <v-card>
      <v-layout>
        <AppBar
          :navigation-bar="navigationBar"
          @update:navigation-bar="
            (value) => {
              test(value);
              navigationBar = value;
            }
          "
        />
        <NavigationBar v-model="navigationBar" />
        <v-main
          style="height: 100vh; overflow-y: hidden"
          class="d-flex flex-column align-stretch"
        >
          <div class="flex-grow-1" style="overflow-y: scroll">
            <!-- <HomePage /> -->
            <RouterView :key="route.fullPath" />
          </div>
        </v-main>
      </v-layout>
    </v-card>
  </v-container>
</template>

//
https://stackoverflow.com/questions/48859119/why-my-service-worker-is-always-waiting-to-activate
// https://web.dev/service-worker-lifecycle/
@/js/channels/encryptionKey