<script setup lang="ts">
import { getQueryParamStringValue } from "@/js/utils";
import { homeRoute } from "@/router/routesConstants";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const subTitle = ref("");
const message = ref("");
const returnUri = ref("");

onMounted(() => {
  subTitle.value =
    getQueryParamStringValue(route.query, "title") ?? "unknown error";

  message.value = getQueryParamStringValue(route.query, "message") ?? "";

  returnUri.value = getQueryParamStringValue(route.query, "return_uri") ?? "";
});
</script>

<template>
  <v-card
    class="h-100 w-100"
    min-width="400px"
    min-height="300px"
    title="Error processing request"
    prepend-icon="mdi-alert"
    :subtitle="subTitle"
    :text="message"
  >
    <v-card-actions>
      <v-btn prepend-icon="mdi-home" :to="homeRoute()">go home</v-btn>
      <v-btn v-if="returnUri.length > 0" prepend-icon="" :to="returnUri"
        >Return to Last Page</v-btn
      >
    </v-card-actions>
  </v-card>
</template>
