<script setup lang="ts">
import axios from "axios";
import { computed, inject, onMounted, ref, type Ref, watch } from "vue";
import { useStore } from "vuex";
import decryptWorker from "@/worker/decrypt?url";
import { initializingKey } from "@/symbols/injectionSymbols";
const store = useStore();
interface media {
  name: string;
  type: string;
  date: Date;
  size: number;
  thumbnail: boolean;
  url: string;
  thumbnail_url: string;
}

const initializing: Ref<boolean> | undefined = inject(initializingKey);
if (initializing === undefined) {
  throw new Error("undefined initializing");
}
const accessToken = computed<string>(() => store.getters.accessToken);
const mediaList = ref<Array<media>>([]);
console.log(accessToken.value);
watch(initializing, async (newValue, oldValue) => {
  console.log("initializing changed to ", newValue);
  if (newValue === oldValue) {
    return;
  }
  if (!newValue) {
    let response = await axios.post(
      "/v1/refreshSession",
      {},
      {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      }
    );
    console.log(response, undefined);
    axios
      .get<Array<media>>("/v1/mediaList", {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      })
      .then((response) => {
        console.log(response);
        if (response.status == 200) {
          mediaList.value = response.data;
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }
});
</script>

<template>
  <div v-if="initializing">loading...</div>
  <div v-else>
    <div :key="`${index}+${media.name}`" v-for="(media, index) in mediaList">
      <img v-if="media.thumbnail" :src="media.thumbnail_url" />
    </div>
  </div>
</template>
