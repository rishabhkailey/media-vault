<script setup lang="ts">
import AppBar from "../components/AppBar.vue";
import NavigationBar from "../components/NavigationBar.vue";
import HomePage from "../components/HomePage.vue";
// todo if not authenticated redirect to some different page
// maybe /about
</script>

<template>
  <v-container class="pa-0 ma-0" fluid>
    <v-card>
      <v-layout>
        <AppBar />
        <NavigationBar />
        <v-main style="height: 100vh">
          <HomePage />
        </v-main>
      </v-layout>
    </v-card>
  </v-container>
</template>