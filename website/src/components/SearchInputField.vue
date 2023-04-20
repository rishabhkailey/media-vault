<script setup lang="ts">
import { ref, watch } from "vue";
import type { SubmitEventPromise } from "vuetify/lib/framework.mjs";

const props = withDefaults(
  defineProps<{
    collapsed: boolean;
    modelValue: string;
  }>(),
  { collapsed: false, modelValue: "" }
);
const emits = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "submit", value: SubmitEventPromise): void;
}>();
const searchInputRules: Array<any> = [];
const searchDialog = ref(false);
const searchElement = ref<HTMLElement | undefined>(undefined);

watch(searchElement, (newValue) => {
  if (newValue === undefined) {
    return;
  }
  // https://github.com/vuetifyjs/vuetify/issues/10659#issuecomment-594329553
  setTimeout(() => {
    console.log("focused ", searchElement.value?.focus);
    searchElement.value?.focus();
  }, 100);
});
</script>

<template>
  <v-form
    class="d-flex flex-grow-1"
    @submit.prevent="(e) => emits('submit', e)"
  >
    <v-text-field
      v-if="!props.collapsed"
      :clearable="true"
      clear-icon="mdi-close"
      append-inner-icon="mdi-magnify"
      :model-value="props.modelValue"
      @update:model-value="
        (value) => {
          emits('update:modelValue', value);
        }
      "
      :rules="searchInputRules"
      label="search"
      :hide-details="true"
    />
    <v-dialog v-else v-model="searchDialog" location="top">
      <template v-slot:activator="{ props }">
        <v-btn color="primary" v-bind="props" icon="mdi-magnify"> </v-btn>
      </template>
      <v-card>
        <v-card-text>
          <v-text-field
            :clearable="true"
            clear-icon="mdi-close"
            append-inner-icon="mdi-magnify"
            :model-value="props.modelValue"
            @update:model-value="
              (value) => {
                emits('update:modelValue', value);
              }
            "
            :rules="searchInputRules"
            label="search"
            ref="searchElement"
            focused
            hide-details
          />
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" block @click="searchDialog = false"
            >Close Dialog</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
