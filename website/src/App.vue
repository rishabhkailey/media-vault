<script setup lang="ts">
import { RouterView } from "vue-router";
</script>

<template>
  <v-app>
    <v-main>
      <v-container class="pa-0 ma-0" fluid>
        <RouterView />
      </v-container>
    </v-main>
  </v-app>
</template>
