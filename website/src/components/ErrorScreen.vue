<script setup lang="ts">
import router from "@/router";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const subTitle = ref("");
const message = ref("");

onMounted(async () => {
  await router.isReady;
  let titleParam = router.currentRoute.value.query?.title;
  if (titleParam && typeof titleParam === "string") {
    subTitle.value = titleParam;
  } else {
    subTitle.value = "Received invalid error title from server";
  }

  let messageParam = router.currentRoute.value.query?.message;
  if (messageParam && typeof messageParam === "string") {
    message.value = messageParam;
  } else {
    message.value = "Received invalid error message from server";
  }
  console.log(router.currentRoute.value.query);
});
</script>

<template>
  <v-card
    min-width="400px"
    max-width="600px"
    min-height="300px"
    title="Error processing request"
    prepend-icon="mdi-alert"
    :subtitle="subTitle"
    :text="message"
  />
</template>
