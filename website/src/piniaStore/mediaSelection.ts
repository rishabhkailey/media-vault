import { defineStore } from "pinia";
import { computed, ref, type ComputedRef } from "vue";

export const useMediaSelectionStore = defineStore("mediaSelection", () => {
  // used for getting the selected media IDs and to check if any media is selected or not
  const selectedMediaIDs = ref<Set<number>>(new Set());
  // used for checking if the media is selected or not
  const selectionMap = ref<Map<number, boolean>>(new Map());

  const count = computed(() => selectedMediaIDs.value.size);

  function updateSelection(id: number, value: boolean) {
    if (value) {
      selectedMediaIDs.value.add(id);
    } else {
      selectedMediaIDs.value.delete(id);
    }
    selectionMap.value.set(id, value);
  }

  // using action as a getter, we can add arguments in normal getters
  function isSelected(id: number): ComputedRef<boolean> {
    return computed(() => !!selectionMap.value.get(id));
  }

  function reset() {
    selectedMediaIDs.value = new Set();
    selectionMap.value = new Map();
  }

  return {
    selectionMap,
    selectedMediaIDs,
    count,
    isSelected,
    updateSelection,
    reset,
  };
});
