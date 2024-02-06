<script setup lang="ts">
import { reactive } from "vue";
import { storeToRefs } from "pinia";
import { useErrorsStore } from "@/piniaStore/errors";
import FloatingWindow from "@/components/Modals/FloatingWindow.vue";

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
const expandedErrorMessageMap = reactive<Map<number, boolean>>(new Map());
function getTrucatedMessage(id: number, message: string): string {
  if (
    message.length <= props.maxMessgeLength ||
    expandedErrorMessageMap.get(id) === true
  ) {
    return message;
  }
  return message.substring(0, props.maxMessgeLength) + "...";
}
</script>
<template>
  <FloatingWindow
    :model-value="popUpErrors.length > 0"
    style="max-height: 300px; max-width: 600px"
    :right="10"
    :bottom="10"
  >
    <div style="height: 100%; overflow-y: auto">
      <v-alert
        v-for="{ title, message, id } in popUpErrors"
        :key="id"
        :elevation="3"
        type="error"
        :title="title"
        :text="getTrucatedMessage(id, message)"
        variant="elevated"
        class="ma-1"
      >
        <template #append>
          <div class="d-flex flex-column">
            <v-btn
              icon="mdi-close"
              @click="() => removeError(id)"
              size="x-small"
              variant="text"
            />
            <v-btn
              v-if="
                message.length > props.maxMessgeLength &&
                expandedErrorMessageMap.get(id) !== true
              "
              icon="mdi-arrow-expand"
              @click="() => expandedErrorMessageMap.set(id, true)"
              size="x-small"
              variant="text"
            />
            <v-btn
              v-if="
                message.length > props.maxMessgeLength &&
                expandedErrorMessageMap.get(id) === true
              "
              icon="mdi-arrow-collapse"
              @click="() => expandedErrorMessageMap.set(id, false)"
              size="x-small"
              variant="text"
            />
          </div>
        </template>
      </v-alert>
    </div>
  </FloatingWindow>
</template>
