<script setup lang="ts">
import { computed } from "vue";
import { ref } from "vue";

const props = withDefaults(
  defineProps<{
    maxMessgeLength?: number;
    title: string;
    message: string;
    expandedByDefault?: boolean;
  }>(),
  {
    maxMessgeLength: 100,
    expandedByDefault: false,
  },
);

const emits = defineEmits<{
  (e: "close"): void;
}>();

const expanded = ref<boolean>(props.expandedByDefault);
const errorMessage = computed(() => {
  if (props.message.length <= props.maxMessgeLength || expanded.value) {
    return props.message;
  }
  return props.message.substring(0, props.maxMessgeLength) + "...";
});
</script>
<template>
  <v-alert
    :elevation="3"
    type="error"
    :title="title"
    :text="errorMessage"
    variant="elevated"
    class="ma-1"
  >
    <template #append>
      <div class="d-flex flex-column">
        <v-btn
          icon="mdi-close"
          @click="() => emits('close')"
          size="x-small"
          variant="text"
        />
        <v-btn
          v-if="message.length > props.maxMessgeLength && !expanded"
          icon="mdi-arrow-expand"
          @click="() => (expanded = true)"
          size="x-small"
          variant="text"
        />
        <v-btn
          v-if="message.length > props.maxMessgeLength && expanded"
          icon="mdi-arrow-collapse"
          @click="() => (expanded = false)"
          size="x-small"
          variant="text"
        />
      </div>
    </template>
  </v-alert>
</template>
