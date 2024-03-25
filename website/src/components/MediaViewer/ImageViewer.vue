<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  src: string;
  lowResSrc: string;
}>();
const reloadCount = ref(0); // used for reload
async function reload() {
  reloadCount.value += 1;
}
</script>

<template>
  <v-img
    :key="reloadCount"
    style="max-width: 100vw; max-height: 100vh"
    :src="props.src"
    :lazy-src="props.lowResSrc"
    transition="fade-transition"
  >
    <template v-slot:placeholder>
      <div class="d-flex align-center justify-center fill-height">
        <v-progress-circular
          color="grey-lighten-4"
          indeterminate
        ></v-progress-circular>
      </div>
    </template>

    <template v-slot:error>
      <div class="d-flex flex-column align-center justify-center fill-height">
        <v-btn icon="mdi-refresh" color="grey-lighten-4" @click="reload" />
        <div
          style="
            color: white;
            background-color: rgba(0, 0, 0, 0.2);
            padding: 0.5em;
            border-radius: 0.5em;
          "
        >
          something went wrong
        </div>
      </div>
    </template>
  </v-img>
</template>
