<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean;
  indeterminate: boolean;
  progress?: number;
}>();
const emits = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
</script>
<template>
  <v-container v-if="props.modelValue">
    <v-btn
      flat
      rounded="pill"
      @click.stop="
        () => {
          emits('update:modelValue', false);
        }
      "
    >
      <v-avatar>
        <v-progress-circular
          :size="70"
          :width="7"
          color="primary"
          :indeterminate="props.indeterminate"
          :model-value="props.progress"
        >
          <template v-slot:default>
            <v-icon color="grey">mdi-arrow-expand</v-icon>
          </template>
        </v-progress-circular>
      </v-avatar>
    </v-btn>
  </v-container>
</template>
