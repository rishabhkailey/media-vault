<script setup lang="ts">
import axios from "axios";
import { computed, inject, onMounted, ref, type Ref, watch } from "vue";
import { useStore } from "vuex";
import decryptWorker from "@/worker/decrypt?url";
import { initializingKey } from "@/symbols/injectionSymbols";
import MediaPreviewPartVue from "@/components/MediaPreviewPart.vue";
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
    // axios
    //   .get<Array<media>>("/v1/mediaList", {
    //     headers: {
    //       Authorization: `Bearer ${accessToken.value}`,
    //     },
    //   })
    //   .then((response) => {
    //     console.log(response);
    //     if (response.status == 200) {
    //       mediaList.value = response.data;
    //     }
    //   })
    //   .catch((err) => {
    //     console.log(err);
    //   });
  }
});
function loadMedia(page: number) {
  return new Promise<{
    n: Number;
    completed: boolean;
  }>((resolve, reject) => {
    axios
      .get<Array<media>>(`/v1/mediaList?page=${page}`, {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      })
      .then((response) => {
        console.log(response);
        if (response.status == 200) {
          mediaList.value = [...mediaList.value, ...response.data];
          resolve({
            n: response.data.length,
            completed: response.data.length === 0,
          });
          return;
        }
        reject(new Error("non 200 status code"));
        return;
      })
      .catch((err) => {
        console.log(err);
        reject(err);
      });
  });
}
const lazyApiLoadObserverTarget = ref<HTMLElement | undefined>(undefined);
let lazyApiLoadTimedOut = false;
let nextPageNumber = 0;
let allMediaLoaded = ref(false);
const observer = new IntersectionObserver(
  (entries, observer) => {
    console.log(entries, observer);
    entries.forEach((entry) => {
      if (
        !entry.isIntersecting ||
        lazyApiLoadObserverTarget.value === undefined
      ) {
        return;
      }
      switch (entry.target) {
        case lazyApiLoadObserverTarget.value:
          console.log("lazyApiLoadObserverTarget matched");
          if (!lazyApiLoadTimedOut) {
            observer.unobserve(lazyApiLoadObserverTarget.value);
            loadMedia(nextPageNumber).then(({ n, completed }) => {
              console.log(`loaded ${n} media`);
              if (completed) {
                console.log("all media loaded");
                allMediaLoaded.value = true;
                return;
              }
              nextPageNumber++;
              if (lazyApiLoadObserverTarget.value !== undefined) {
                observer.observe(lazyApiLoadObserverTarget.value);
              }
            });
            // timeout of 1 second to prevent any bug from overloading the browser with api calls
            lazyApiLoadTimedOut = true;
            setTimeout(() => {
              lazyApiLoadTimedOut = false;
            }, 1000);
          }
      }
    });
  },
  {
    root: null,
    rootMargin: "10px",
    threshold: 0.1,
  }
);
watch(lazyApiLoadObserverTarget, (newValue, oldvalue) => {
  if (oldvalue !== undefined) {
    observer.unobserve(oldvalue);
  }
  if (newValue === undefined) {
    console.warn("lazyApiLoadObserverTarget undefined");
    return;
  }
  observer.observe(newValue);
});
</script>

<template>
  <v-container>
    <div v-if="initializing">loading...</div>
    <div v-else>
      <div :key="`${index}+${media.name}`" v-for="(media, index) in mediaList">
        <img v-if="media.thumbnail" :src="media.thumbnail_url" />
      </div>
      <!-- <v-lazy :min-height="200" :options="{ threshold: 0.5 }">
        <v-container>
          <v-container style="min-height: 100px">
            <MediaPreviewPartVue
              :page="0"
              :perPage="50"
              :accessToken="accessToken"
            />
          </v-container>
          <v-container style="min-height: 100px">
            <MediaPreviewPartVue
              :page="1"
              :perPage="50"
              :accessToken="accessToken"
            />
          </v-container>
        </v-container>
      </v-lazy> -->
      <div
        ref="lazyApiLoadObserverTarget"
        style="min-height: 100px"
        v-if="!allMediaLoaded"
      ></div>
    </div>
  </v-container>
</template>
