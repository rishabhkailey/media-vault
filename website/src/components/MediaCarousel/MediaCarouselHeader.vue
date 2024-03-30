<script setup lang="ts">
import type { StyleValue } from "vue";

const props = withDefaults(
  defineProps<{
    style?: StyleValue;
    class: string;
  }>(),
  {
    class: "",
    style: "",
  },
);

const emits = defineEmits<{
  (e: "openInNewWindow"): void;
  (e: "close"): void;
  (e: "download"): void;
}>();
</script>
<template>
  <!--  -->
  <div :class="props.class" :style="props.style">
    <v-tooltip
      text="Opens in new tab or downloads if unsupported"
      location="bottom"
    >
      <template v-slot:activator="{ props }">
        <v-btn
          icon="mdi-open-in-new"
          @click.stop="() => emits('openInNewWindow')"
          style="background: none; border: none; box-shadow: none"
          v-bind="props"
        />
      </template>
    </v-tooltip>

    <v-tooltip text="Download" location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn
          icon="mdi-download"
          @click.stop="() => emits('download')"
          style="background: none; border: none; box-shadow: none"
          v-bind="props"
        />
      </template>
    </v-tooltip>

    <v-tooltip text="Close" location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn
          icon="mdi-close"
          @click.stop="() => emits('close')"
          style="background: none; border: none; box-shadow: none"
          data-test-id="media-carousel-close-button"
          v-bind="props"
        />
      </template>
    </v-tooltip>
  </div>
</template>
