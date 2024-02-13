<script setup lang="ts">
import ErrorMessage from "./ErrorMessage.vue";
import { storeToRefs } from "pinia";
import { useErrorsStore } from "@/piniaStore/errors";
import FloatingWindow from "@/components/Modals/FloatingWindow.vue";
import { POPUP_ERROR_WINDOW_Z_INDEX } from "@/js/constants/z-index";

const props = withDefaults(
  defineProps<{
    maxMessgeLength: number;
  }>(),
  {
    maxMessgeLength: 100,
  },
);

const errorsStore = useErrorsStore();
const { popUpErrors } = storeToRefs(errorsStore);
const { removeError } = errorsStore;
</script>
<template>
  <FloatingWindow
    :model-value="popUpErrors.length > 0"
    style="max-height: 300px; max-width: 600px"
    :right="10"
    :bottom="10"
    :z-index="POPUP_ERROR_WINDOW_Z_INDEX"
  >
    <div style="max-height: 300px; overflow-y: auto">
      <ErrorMessage
        v-for="{ title, message, id } in popUpErrors"
        :message="message"
        :title="title"
        :expanded-by-default="false"
        :max-messge-length="props.maxMessgeLength"
        :key="id"
        @close="() => removeError(id)"
      />
    </div>
  </FloatingWindow>
</template>
