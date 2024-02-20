<script setup lang="ts">
import AppBar from "@/components/AppBar/AppBar.vue";
import SideBar from "../components/SideBar/SideBar.vue";
import { computed, ref, watch } from "vue";
import { useDisplay } from "vuetify";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { storeToRefs } from "pinia";
import { EncryptionKeyChannelClient } from "@/js/channels/encryptionKey";
import PopUpError from "@/components/Error/PopUpError.vue";
import { onMounted } from "vue";
import { getUserManager } from "@/js/auth";
import { useAuthStore } from "@/piniaStore/auth";

const display = useDisplay();

const userInfoStore = useUserInfoStore();
const { usableEncryptionKey } = storeToRefs(userInfoStore);
const { setUserAuthInfo } = useAuthStore();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value,
);

const displaySidebar = ref(!smallDisplay.value);

const encryptionKeyChannel = new EncryptionKeyChannelClient(
  usableEncryptionKey.value,
);

watch(usableEncryptionKey, () => {
  encryptionKeyChannel.encryptionKey = usableEncryptionKey.value;
});

onMounted(() => {
  getUserManager().events.addUserLoaded((user) => {
    setUserAuthInfo(user);
  });
});
</script>

<template>
  <v-container class="pa-0 ma-0" fluid>
    <v-card>
      <v-layout>
        <AppBar
          :side-bar="displaySidebar"
          @update:side-bar="
            (value) => {
              displaySidebar = value;
            }
          "
        />
        <SideBar v-model="displaySidebar" />
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
          <PopUpError :max-messge-length="100" />
        </v-main>
      </v-layout>
    </v-card>
  </v-container>
</template>
