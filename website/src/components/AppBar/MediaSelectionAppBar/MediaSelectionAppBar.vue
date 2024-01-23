<script setup lang="ts">
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { storeToRefs } from "pinia";
import MediaActions from "./actions/MediaActions.vue";
import AlbumMediaActions from "./actions/AlbumMediaActions.vue";
import { useRoute } from "vue-router";

const route = useRoute();
console.log("route name", route.name);
const mediaSelectionStore = useMediaSelectionStore();
const { reset: resetMediaSelection } = mediaSelectionStore;
const { count: selectedMediaCount } = storeToRefs(mediaSelectionStore);

console.log(typeof AlbumMediaActions);
let Actions = MediaActions;
switch (route.name) {
  case "Album": {
    Actions = AlbumMediaActions;
    break;
  }
  default: {
    Actions = MediaActions;
    break;
  }
}
</script>
<template>
  <v-row class="d-flex align-center ml-2 justify-start mx-2">
    <!-- start -->
    <v-col>
      <v-row class="d-flex align-center">
        <v-btn icon="mdi-close" @click.stop="resetMediaSelection" />
        <div class="text-subtitle-1">{{ selectedMediaCount }} selected</div>
      </v-row>
    </v-col>

    <!-- end -->
    <v-col class="d-flex flex-row justify-end">
      <Actions />
    </v-col>
  </v-row>
</template>
