<script setup lang="ts">
const props = defineProps<{
  failed: boolean;
  progress: number;
  name: string;
  size: number;
  completed: boolean;
}>();
const emits = defineEmits<{
  (e: "cancel"): void;
}>();
</script>
<template>
  <v-list-item
    class="ma-0 pa-1"
    data-test-id="uploading-files-progress-list-item"
    :key="props.name"
    :title="props.name"
    :subtitle="props.size"
  >
    <template v-slot:prepend>
      <v-avatar>
        <v-progress-circular
          :size="70"
          :width="7"
          :color="props.failed ? 'red' : 'primary'"
          :model-value="props.failed ? '100' : props.progress"
        >
          <span
            style="font-size: 0.75em"
            data-test-id="uploading-files-progress-list-item-progress"
          >
            {{ props.failed ? "!" : Math.round(props.progress) + "%" }}</span
          >
        </v-progress-circular>
      </v-avatar>
    </template>
    <template v-slot:append>
      <v-btn
        v-if="!props.completed"
        color="grey-lighten-1"
        icon="mdi-close"
        variant="text"
        data-test-id="uploading-files-progress-list-item-cancel-button"
        @click.stop="() => emits('cancel')"
      />
      <v-icon
        v-if="props.completed && !props.failed"
        color="grey-lighten-1"
        icon="mdi-check"
        data-test-id="uploading-files-progress-list-item-completed-icon"
        variant="text"
      />
      <!-- rotate-right -->
      <v-tooltip
        location="top"
        text="retry"
        v-if="props.completed && props.failed"
      >
        <template v-slot:activator="{ props }">
          <v-btn
            color="grey-lighten-1"
            icon="mdi-rotate-right"
            variant="text"
            data-test-id="uploading-files-progress-list-item-retry-button"
            v-bind="props"
          />
        </template>
      </v-tooltip>
    </template>
  </v-list-item>
</template>
