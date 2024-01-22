import { defineStore } from "pinia";
import { ref } from "vue";

interface popUpError {
  id: number;
  message: string;
  title: string;
  timeoutSeconds: number;
}

export const useErrorsStore = defineStore("errors", () => {
  const popUpErrors = ref<Array<popUpError>>([]);
  let id = 0;

  // -1 for never timeout
  function appendError(
    title: string,
    message: string,
    timeoutSeconds: number,
  ): number {
    const errId = id++;
    popUpErrors.value.push({
      title: title,
      message: message,
      timeoutSeconds: timeoutSeconds,
      id: errId,
    });
    if (timeoutSeconds === -1) {
      return errId;
    }
    setTimeout(() => {
      removeError(errId);
    }, timeoutSeconds * 1000);
    return errId;
  }

  function removeError(id: number): boolean {
    const index = popUpErrors.value.findIndex((err) => err.id === id);
    if (index !== -1) {
      popUpErrors.value.splice(index, 1);
      return true;
    }
    return false;
  }

  return {
    popUpErrors,
    appendError,
    removeError,
  };
});
