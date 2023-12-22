<script setup lang="ts">
import AppBar from "@/components/AppBar/AppBar.vue";
import NavigationBar from "../components/NavigationBar.vue";
import { computed, inject, onMounted, ref, onBeforeMount, watch } from "vue";
import { userManagerKey } from "@/symbols/injectionSymbols";
import type { UserManager } from "oidc-client-ts";
// import decryptWorker from "@/worker/dist/bundle.js?url";
import decryptWorker from "@/worker/decrypt?url";
// todo if not authenticated redirect to some different page
// maybe /about
import { useDisplay } from "vuetify";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/piniaStore/auth";
import { useLoadingStore } from "@/piniaStore/loading";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { EncryptionKeyChannelClient } from "@/js/channels/encryptionKey";
import axios from "axios";
import { updateOrRegisterServiceWorker } from "@/js/serviceWorker/registeration";

const authStore = useAuthStore();
const { authenticated } = storeToRefs(authStore);
const display = useDisplay();

const userInfoStore = useUserInfoStore();
const { loadUserInfo } = userInfoStore;
const {
  initRequired: userInfoInitRequired,
  encryptionKeyValidated,
  usableEncryptionKey,
} = storeToRefs(userInfoStore);

const { setInitializing } = useLoadingStore();
// setInitializing(true);

const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const router = useRouter();
// const initializingRef = ref(true);
const navigationBar = ref(!smallDisplay.value);
// provide(initializingKey, initializingRef);
// const userManager: UserManager | undefined = inject(userManagerKey);

// const userInit = () => {
//   return new Promise<boolean>((resolve, reject) => {
//     if (userManager === undefined) {
//       reject(new Error("undefined userManager"));
//       return;
//     }
//     userManager
//       .getUser()
//       .then((user) => {
//         console.log(user);
//         if (user?.expired) {
//           resolve(false);
//           return;
//         }
//         if (user === null) {
//           resolve(false);
//           return;
//         }
//         authStore.setUserInfo(user);
//         resolve(true);
//       })
//       .catch((err) => {
//         reject(err);
//         return;
//       });
//   });
// };

// const updateOrRegisterServiceWorker = () => {
//   return new Promise<ServiceWorker>((resolve, reject) => {
//     navigator.serviceWorker
//       .getRegistration(decryptWorker)
//       .then((registration) => {
//         if (registration === undefined) {
//           registerServiceWorker()
//             .then((serviceWorker) => {
//               resolve(serviceWorker);
//               return;
//             })
//             .catch((err) => {
//               reject(err);
//               return;
//             });
//         } else {
//           updateServiceWorker(registration)
//             .then((serviceWorker) => {
//               resolve(serviceWorker);
//               return;
//             })
//             .catch((err) => {
//               reject(err);
//               return;
//             });
//         }
//       });
//   });
// };

// function updateServiceWorker(registration: ServiceWorkerRegistration) {
//   // if hard reload or service worker url changed
//   if (
//     navigator.serviceWorker.controller === null ||
//     !navigator.serviceWorker.controller.scriptURL.endsWith(decryptWorker)
//   ) {
//     // https://github.com/rishabhkailey/media-service/issues/2
//     // https://stackoverflow.com/a/66816077
//     return registration.unregister().then(() => registerServiceWorker());
//   }
//   return new Promise<ServiceWorker>((resolve, reject) => {
//     registration
//       .update()
//       .then(() => {
//         console.log("updated");
//         if (registration.active === null) {
//           reject(new Error("got null service worker after update"));
//           return;
//         }
//         resolve(registration.active);
//         return;
//       })
//       .catch((err) => {
//         reject(err);
//         return;
//       });
//   });
// }

// // https://github.com/jimmywarting/StreamSaver.js/blob/master/mitm.html#L39
// function registerServiceWorker(): Promise<ServiceWorker> {
//   return new Promise<ServiceWorker>((resolve, reject) => {
//     if ("serviceWorker" in navigator) {
//       // unregister the existing service worker
//       navigator.serviceWorker
//         .register(decryptWorker, {
//           scope: "./",
//           type: "module",
//         })
//         .then((swReg) => {
//           if (swReg.active !== null) {
//             console.debug("Service Worker registsered");
//             resolve(swReg.active);
//             return;
//           }
//           const swRegTmp = swReg.installing || swReg.waiting;
//           if (swRegTmp === null) {
//             reject(new Error("got null service worker registration"));
//             return;
//           }
//           let callback: () => void;
//           console.debug("waiting for Service Worder to registser");
//           swRegTmp.addEventListener(
//             "statechange",
//             (callback = () => {
//               if (swRegTmp.state === "activated") {
//                 console.debug("Service Worder registed and active");
//                 swRegTmp.removeEventListener("statechange", callback);
//                 resolve(swRegTmp);
//               }
//             })
//           );
//         })
//         .catch((err) => {
//           reject(err);
//           return;
//         });
//     }
//   });
// }
async function init(): Promise<boolean> {
  return new Promise<boolean>(() => true);
  // try {
  //   await updateOrRegisterServiceWorker();
  //   // await userInit();
  // } catch (error) {
  //   console.log(error);
  //   router.push({
  //     name: "errorscreen",
  //     query: {
  //       title: "init failed",
  //       message: "error message - " + error,
  //     },
  //   });
  // }
  // if (!authenticated.value) {
  //   // route to about page
  //   router.push({
  //     name: "about",
  //   });
  // }
  // try {
  //   await loadUserInfo();
  // } catch (err: any) {
  //   if (axios.isAxiosError(err) && err.response?.status === 401) {
  //     router.push({
  //       name: "errorscreen",
  //       query: {
  //         title: "Access Denied",
  //         message: `server denied the access. try clearing the browser cookies if the problem presists. \nerror from server: ${err}`,
  //       },
  //     });
  //   }
  // }

  // if (userInfoInitRequired.value) {
  //   router.push({
  //     name: "onboarding",
  //   });
  //   return false;
  // }
  // if (encryptionKeyValidated.value == false) {
  //   // encryptionKeyOverlay.value = true;
  //   router.push({
  //     name: "encryptionKey",
  //     query: {
  //       return_uri: window.location.pathname + window.location.hash,
  //     },
  //   });
  //   return false;
  // }
  // setInitializing(false);
  // return true;
}

// todo is this here ok?
// or rename the component to app?
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
            <RouterView v-slot="{ Component }">
              <KeepAlive>
                <component :is="Component" />
              </KeepAlive>
            </RouterView>
          </div>
        </v-main>
      </v-layout>
    </v-card>
  </v-container>
</template>

//
https://stackoverflow.com/questions/48859119/why-my-service-worker-is-always-waiting-to-activate
// https://web.dev/service-worker-lifecycle/ @/js/channels/encryptionKey
