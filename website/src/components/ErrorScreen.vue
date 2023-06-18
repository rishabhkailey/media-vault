<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();

const subTitle = ref("");
const message = ref("");

onMounted(() => {
  let titleParam = route.query?.title;
  if (titleParam && typeof titleParam === "string") {
    subTitle.value = titleParam;
  } else {
    subTitle.value = "Received invalid error title from server";
  }

  let messageParam = route.query?.message;
  if (messageParam && typeof messageParam === "string") {
    message.value = messageParam;
  } else {
    message.value = "Received invalid error message from server";
  }
  console.log(route.query);
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
      <v-btn
        prepend-icon="mdi-home"
        :to="{
          name: 'Home',
        }"
        >go home</v-btn
      >
    </v-card-actions>
  </v-card>
</template>
