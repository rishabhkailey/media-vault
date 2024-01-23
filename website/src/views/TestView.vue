<script setup lang="ts">
import PopUpError from "@/components/Error/PopUpError.vue";
import { useErrorsStore } from "@/piniaStore/errors";
const { appendError } = useErrorsStore();
</script>

<template>
  <div style="display: block; height: 100vh; overflow-y: scroll">
    <PopUpError :max-messge-length="100" />

    <v-btn
      @click="
        () =>
          appendError(
            'test error',
            'a realy loooong error message. booring looong erroorrr messaggeeeeeebooring looong erroorrr messaggeeeeeebooring looong erroorrr messaggeeeeeebooring looong erroorrr messaggeeeeeebooring looong erroorrr messaggeeeeeebooring looong erroorrr messaggeeeeee.',
            -1,
          )
      "
      >pop up error test</v-btn
    >
  </div>
</template>
