<script setup lang="ts">
import AppBar from "@/components/AboutPage/AppBar.vue";
import MainPage from "@/components/AboutPage/MainPage.vue";
</script>
<template>
  <v-layout>
    <AppBar />
    <v-main
      style="height: 100vh; overflow-y: hidden"
      class="d-flex flex-column align-stretch"
    >
      <div class="flex-grow-1" style="overflow-y: scroll">
        <MainPage />
      </div>
    </v-main>
  </v-layout>
</template>
