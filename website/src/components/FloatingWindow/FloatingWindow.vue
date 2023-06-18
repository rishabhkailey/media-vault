<script setup lang="ts">
import type { StyleValue } from "vue";
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    height: string;
    width: string;
    top?: number;
    bottom?: number;
    left?: number;
    right?: number;
    margin: number;
    modelValue: boolean;
  }>(),
  {
    height: "",
    width: "",
    margin: 10,
    modelValue: false,
  }
);

const positionStyle = computed<Array<StyleValue>>(() => {
  let style: Array<string> = [
    "position: absolute",
    `margin: ${props.margin}px`,
  ];
  if (props.height.length !== 0) {
    style.push(`height: ${props.height}`);
  }
  if (props.width.length !== 0) {
    style.push(`width: ${props.width}`);
  }
  if (props.top !== undefined) {
    style.push(`top: ${props.top}px`);
  }
  if (props.bottom !== undefined) {
    style.push(`bottom: ${props.bottom}px`);
  }
  if (props.left !== undefined) {
    style.push(`left: ${props.left}px`);
  }
  if (props.right !== undefined) {
    style.push(`right: ${props.right}px`);
  }
  if (props.top !== undefined) {
    style.push(`top: ${props.top}px`);
  }
  return style;
});
</script>
<template>
  <Teleport to="body">
    <div
      v-if="props.modelValue"
      class="floating-over-body"
      :style="positionStyle"
    >
      <slot></slot>
    </div>
  </Teleport>
</template>
