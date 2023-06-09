import { defineStore } from "pinia";
import { ref } from "vue";

export const useLoadingStore = defineStore("loading", () => {
  // global loading used in progress bar in app bar
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

  // initializing used to wait for service worker to register and other init things
  const initializing = ref(false);

  function setInitializing(_initializing: boolean) {
    initializing.value = _initializing;
  }

  return {
    // global
    loading,
    progress,
    indeterminate,
    setGlobalLoading,
    setProgress,
    // initializing
    initializing,
    setInitializing,
  };
});
