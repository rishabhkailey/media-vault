<script setup lang="ts">
import axios from "axios";
import { computed, inject, ref, type Ref, watch } from "vue";
import { useStore } from "vuex";
import { initializingKey } from "@/symbols/injectionSymbols";
import MediaThumbnail from "./MediaThumbnail.vue";
import MediaPreview from "./MediaPreview.vue";
import { LOAD_MORE_MEDIA_ACTION } from "@/store/modules/media";

const store = useStore();
const initializing: Ref<boolean> | undefined = inject(initializingKey);
const accessToken = computed<string>(() => store.getters.accessToken);
const mediaList = computed<Array<Media>>(() => store.getters.mediaList);
const allMediaLoaded = computed<boolean>(() => store.getters.allMediaLoaded);

const selectedMediaIndex = ref<number | undefined>(undefined);
const prviewOverlay = ref<boolean>(false);

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
            store.dispatch(LOAD_MORE_MEDIA_ACTION).then(() => {
              if (
                !allMediaLoaded.value &&
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
  <v-container>
    <div v-if="initializing">loading...</div>
    <div v-else>
      <v-row>
        <v-col
          :key="`${index}+${media.name}`"
          v-for="(media, index) in mediaList"
          class="d-flex child-flex"
          xs="12"
          sm="6"
          md="4"
          lg="3"
          xl="2"
        >
          <MediaThumbnail
            :media="media"
            @click="
              () => {
                prviewOverlay = true;
                selectedMediaIndex = index;
              }
            "
          />
        </v-col>
        <div
          ref="lazyApiLoadObserverTarget"
          style="min-height: 100px"
          v-if="!allMediaLoaded"
        ></div>
      </v-row>
      <v-overlay
        v-model="prviewOverlay"
        :close-on-content-click="false"
        :close-delay="0"
        :open-delay="0"
        class="d-flex justify-center align-center"
        :z-index="2000"
      >
        <MediaPreview
          v-if="selectedMediaIndex !== undefined"
          :media-list="mediaList"
          :index="selectedMediaIndex"
          :allMediaLoaded="allMediaLoaded"
          :loadMoreMedia="() => store.dispatch(LOAD_MORE_MEDIA_ACTION)"
          @close="
            () => {
              prviewOverlay = false;
            }
          "
        />
      </v-overlay>
    </div>
  </v-container>
</template>
