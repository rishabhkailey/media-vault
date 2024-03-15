<script lang="ts" setup>
const props = defineProps<{
  title: string;
  message: string;
  cancelButtonText: string;
  cancelButtonColor: string;
  confirmButtonText: string;
  confirmButtonColor: string;
  confirmInProgress: boolean;
  modelValue: boolean;
  errorMessage: string;
  dataTestId: string; // used for e2e testing
}>();

const emits = defineEmits<{
  (e: "confirm"): void;
  (e: "cancel"): void;
  (e: "update:modelValue", value: boolean): void;
}>();
</script>
<template>
  <v-overlay
    :model-value="props.modelValue"
    @update:model-value="(value) => emits('update:modelValue', value)"
    class="align-center justify-center"
    :data-test-id="props.dataTestId"
  >
    <v-card :subtitle="props.title" style="min-width: 300px; max-width: 500px">
      <v-card-text>
        {{ props.message }}
        <v-alert type="error" v-if="props.errorMessage.length > 0">
          {{ props.errorMessage }}
        </v-alert>
      </v-card-text>

      <v-card-actions>
        <div class="d-flex flex-row justify-end flex-grow-1">
          <v-btn
            variant="text"
            @click="
              () => {
                emits('cancel');
              }
            "
            data-test-id="cancel-button"
            :color="props.cancelButtonText"
          >
            {{ props.cancelButtonText }}
          </v-btn>
          <v-btn
            :color="props.confirmButtonColor"
            variant="text"
            :loading="props.confirmInProgress"
            data-test-id="confirm-button"
            @click="() => emits('confirm')"
            >{{ props.confirmButtonText }}</v-btn
          >
        </div>
      </v-card-actions>
    </v-card>
  </v-overlay>
</template>
