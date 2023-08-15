<script setup lang="ts">
const props = defineProps<{
  src: string;
  lowResSrc: string;
}>();
</script>

<template>
  <div style="width: 100vw">
    <v-img
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
    </v-img>
  </div>
</template>
