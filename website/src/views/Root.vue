<script setup lang="ts">
import AppBar from "@/components/AppBar/AppBar.vue";
import NavigationBar from "../components/NavigationBar.vue";
import { computed, ref, watch } from "vue";
import { useDisplay } from "vuetify";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { EncryptionKeyChannelClient } from "@/js/channels/encryptionKey";

const display = useDisplay();

const userInfoStore = useUserInfoStore();
const { usableEncryptionKey } = storeToRefs(userInfoStore);

const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);

const navigationBar = ref(!smallDisplay.value);

// todo is this here ok?
// or rename the component to app?
const encryptionKeyChannel = new EncryptionKeyChannelClient(
  usableEncryptionKey.value
);

watch(usableEncryptionKey, () => {
  encryptionKeyChannel.encryptionKey = usableEncryptionKey.value;
});
</script>

<template>
  <v-container class="pa-0 ma-0" fluid>
    <v-card>
      <v-layout>
        <AppBar
          :navigation-bar="navigationBar"
          @update:navigation-bar="
            (value) => {
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
