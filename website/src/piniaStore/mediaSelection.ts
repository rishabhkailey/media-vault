import { defineStore } from "pinia";
import { computed, ref, type ComputedRef } from "vue";

export const useMediaSelectionStore = defineStore("mediaSelection", () => {
  // const selectedMediaIDs = ref<Set<number>>(new Set());
  const selectionMap = ref<Map<number, boolean>>(new Map());

  function updateSelection(id: number, value: boolean) {
    selectionMap.value.set(id, value);
  }

  // using action as a getter, we can add arguments in normal getters
  function isSelected(id: number): ComputedRef<boolean> {
    return computed(() => !!selectionMap.value.get(id));
  }

  return { selectionMap, isSelected, updateSelection };
});
