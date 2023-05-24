import { defineStore } from "pinia";
import { ref } from "vue";

export const useLoadingStore = defineStore("loading", () => {
  const loading = ref(false);
  const progress = ref(0);
  const indeterminate = ref(false);

  function setGlobalLoading(
    _loading: boolean,
    _indeterminate: boolean,
    _progress: number
  ) {
    loading.value = _loading;
    indeterminate.value = _indeterminate;
    progress.value = _progress;
  }

  function setProgress(_progress: number) {
    progress.value = _progress;
  }

  return {
    loading,
    progress,
    indeterminate,
    setGlobalLoading,
    setProgress,
  };
});
