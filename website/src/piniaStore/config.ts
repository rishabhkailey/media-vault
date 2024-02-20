import { getConfig, type Config } from "@/js/api/config";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useConfigStore = defineStore("config", () => {
  const config = ref<Config | undefined>(undefined);
  const loaded = ref<boolean>(false);
  const loadFailed = ref<boolean>(false);

  function loadConfig(): Promise<Config> {
    return new Promise<Config>((resolve, reject) => {
      getConfig()
        .then((c) => {
          loaded.value = true;
          loadFailed.value = false;
          config.value = c;
          resolve(c);
        })
        .catch((err) => {
          loaded.value = false;
          loadFailed.value = true;
          reject(err);
        });
    });
  }
  return {
    config,
    loaded,
    loadFailed,
    loadConfig,
  };
});
