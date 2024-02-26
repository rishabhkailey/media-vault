import { defineStore } from "pinia";
import { computed, ref, type ComputedRef } from "vue";

export const useMediaSelectionStore = defineStore("mediaSelection", () => {
  // used for checking if the media is selected or not
  const selectedMediaMap = ref<Map<number, boolean>>(new Map());

  const count = computed(() => selectedMediaMap.value.size);

  function updateSelection(id: number, value: boolean) {
    if (value) {
      selectedMediaMap.value.set(id, value);
    } else {
      selectedMediaMap.value.delete(id);
    }
  }

  // using action as a getter, we can add arguments in normal getters
  function isSelected(id: number): ComputedRef<boolean> {
    return computed(() => !!selectedMediaMap.value.get(id));
  }

  function reset() {
    selectedMediaMap.value = new Map();
  }

  return {
    selectedMediaMap,
    count,
    isSelected,
    updateSelection,
    reset,
  };
});
