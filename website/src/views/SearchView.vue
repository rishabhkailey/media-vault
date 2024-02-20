<script setup lang="ts">
import { useRoute } from "vue-router";
import SearchMediaGrid from "../components/Search/SearchMediaGrid.vue";
import { generateQueryParamsKey } from "@/js/utils";
const route = useRoute();
</script>

<template>
  <SearchMediaGrid :key="generateQueryParamsKey(route.params, 'query')" />
</template>
