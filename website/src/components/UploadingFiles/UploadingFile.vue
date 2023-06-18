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
  <v-list-item :key="props.name" :title="props.name" :subtitle="props.size">
    <template v-slot:prepend>
      <!-- todo create a separate component with simple logic -->
      <v-avatar>
        <v-progress-circular
          :size="70"
          :width="7"
          :color="props.failed ? 'red' : 'primary'"
          :model-value="props.failed ? '100' : props.progress"
        >
          <!-- todo text size -->
          {{ props.failed ? "!" : Math.round(props.progress) + "%" }}
        </v-progress-circular>
      </v-avatar>
    </template>
    <template v-slot:append>
      <v-btn
        v-if="!props.completed"
        color="grey-lighten-1"
        icon="mdi-close"
        variant="text"
        @click.stop="() => emits('cancel')"
      />
      <v-btn
        v-if="props.completed && !props.failed"
        color="grey-lighten-1"
        icon="mdi-check"
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
            v-bind="props"
          />
        </template>
      </v-tooltip>
    </template>
  </v-list-item>
</template>
