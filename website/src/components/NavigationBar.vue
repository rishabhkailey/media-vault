<script setup lang="ts">
import { computed } from "vue";
import { useDisplay } from "vuetify";
const display = useDisplay();

const props = defineProps<{
  modelValue: boolean;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value
);
</script>

<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    @update:model-value="
      (newValue) => {
        emit('update:modelValue', newValue);
      }
    "
    :permanent="!smallDisplay"
    :temporary="smallDisplay"
    :rounded="false"
    elevation="2"
  >
    <v-list nav>
      <v-list-item
        prepend-icon="mdi-email"
        title="Inbox"
        value="inbox"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-account-supervisor-circle"
        title="Supervisors"
        value="supervisors"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-clock-start"
        title="Clock-in"
        value="clockin"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>
