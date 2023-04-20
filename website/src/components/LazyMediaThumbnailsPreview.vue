<script setup lang="ts">
import axios from "axios";
import { computed, inject, ref, type Ref, watch, provide } from "vue";
import { useStore } from "vuex";
import {
  initializingKey,
  loadMoreMediaKey,
  allMediaLoadedKey,
  mediaListKey,
} from "@/symbols/injectionSymbols";
import MonthlyThumbnailPreview from "./MonthlyThumbnailPreview.vue";
import { getMonthlyMediaIndex } from "@/utils/date";

const props = defineProps<{
  mediaList: Array<Media>;
  allMediaLoaded: boolean;
  loadMoreMedia: () => Promise<any>;
}>();

provide(
  mediaListKey,
  computed(() => props.mediaList)
);
provide(
  allMediaLoadedKey,
  computed(() => props.allMediaLoaded)
);
provide(loadMoreMediaKey, props.loadMoreMedia);

const store = useStore();
const initializing: Ref<boolean> | undefined = inject(initializingKey);
const accessToken = computed<string>(() => store.getters.accessToken);

const monthlyMediaList = computed<Array<MonthlyMedia>>(() =>
  getMonthlyMediaIndex(props.mediaList)
);
if (initializing === undefined) {
  throw new Error("undefined initializing");
}
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
    console.log(response);
  }
});
const lazyApiLoadObserverTarget = ref<HTMLElement | undefined>(undefined);
let lazyApiLoadTimedOut = false;
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
            props.loadMoreMedia().then(() => {
              if (
                !props.allMediaLoaded &&
                lazyApiLoadObserverTarget.value !== undefined
              ) {
                observer.observe(lazyApiLoadObserverTarget.value);
              }
            });
            // timeout of 100ms second to prevent any bug from overloading the browser with api calls
            lazyApiLoadTimedOut = true;
            setTimeout(() => {
              lazyApiLoadTimedOut = false;
            }, 100);
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
  if (newValue === undefined || !(newValue instanceof HTMLElement)) {
    console.warn("lazyApiLoadObserverTarget undefined");
    return;
  }
  observer.observe(newValue);
});
</script>

<template>
  <div class="pa-0 ma-0">
    <div v-if="initializing">loading...</div>
    <div v-else class="d-flex flex-column align-stretch">
      <div v-for="(monthlyMedia, index) in monthlyMediaList" :key="index">
        <MonthlyThumbnailPreview
          :month="monthlyMedia.month"
          :year="monthlyMedia.year"
          :index-media-list="monthlyMedia.media"
          :index-offset="0"
        />
        <v-divider />
      </div>

      <div
        ref="lazyApiLoadObserverTarget"
        style="min-height: 100px"
        v-if="!allMediaLoaded"
      ></div>
    </div>
  </div>
</template>
