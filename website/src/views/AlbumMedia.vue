<script setup lang="ts">
import { useRoute } from "vue-router";
import { generateQueryParamsKey } from "@/js/utils";
import AlbumMediaGrid from "@/components/Album/AlbumMediaGrid.vue";
const route = useRoute();
</script>

<template>
  <AlbumMediaGrid :key="generateQueryParamsKey(route.params, 'album_id')" />
</template>
