<script setup lang="ts">
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { useLoadingStore } from "@/piniaStore/loading";
import NormalAppBar from "@/components/AppBar/NormalAppBar/NormalAppBar.vue";
import MediaSelectionAppBar from "@/components/AppBar/MediaSelectionAppBar/MediaSelectionAppBar.vue";
import { storeToRefs } from "pinia";

const props = defineProps<{
  navigationBar: boolean;
}>();

const emits = defineEmits<{
  (e: "update:navigationBar", value: boolean): void;
}>();

const mediaSelectionStore = useMediaSelectionStore();
const { count: selectedMediaCount } = storeToRefs(mediaSelectionStore);

const loadingStore = useLoadingStore();
const { loading, progress, indeterminate } = storeToRefs(loadingStore);
</script>
<template>
  <!-- <v-scale-transition> -->
  <v-app-bar
    :rounded="false"
    elevation="2"
    class="pa-0 ma-0 d-flex justify-center align-center"
    style="height: inherit"
  >
    <v-progress-linear
      color="primary"
      location="top"
      :absolute="true"
      :active="loading"
      :indeterminate="indeterminate"
      :model-value="progress"
      class="pa-0 ma-0"
    ></v-progress-linear>

    <MediaSelectionAppBar v-if="selectedMediaCount > 0" />
    <NormalAppBar
      v-else
      :navigation-bar="props.navigationBar"
      @update:navigation-bar="
        (value) => {
          emits('update:navigationBar', value);
        }
      "
    />
  </v-app-bar>
  <!-- </v-scale-transition> -->
</template>
