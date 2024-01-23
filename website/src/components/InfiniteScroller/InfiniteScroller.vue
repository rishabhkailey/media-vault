<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from "vue";

const props = withDefaults(
  defineProps<{
    // todo on error timeout
    minHeight: number;
    minWidth: number;
    rootMargin: number;
    threshold: number;
    onThresholdReach: () => Promise<any>;
  }>(),
  {
    minHeight: 100,
    minWidth: 100,
    rootMargin: 10,
    threshold: 0.1,
    loading: false,
  },
);
const emits = defineEmits<{
  (e: "thresholdReached"): void;
}>();

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
            props.onThresholdReach().then(() => {
              // timeout of 100ms second to prevent any bug from overloading the browser with api calls
              lazyApiLoadTimedOut = true;
              setTimeout(() => {
                lazyApiLoadTimedOut = false;
                if (lazyApiLoadObserverTarget.value !== undefined) {
                  console.log("observing again");
                  observer.observe(lazyApiLoadObserverTarget.value);
                }
              }, 100);
            });
          }
      }
    });
  },
  {
    root: null,
    rootMargin: `${props.rootMargin}px`,
    threshold: props.threshold,
  },
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

onBeforeUnmount(() => {
  console.log("observer unmounting");
  if (lazyApiLoadObserverTarget.value !== undefined) {
    try {
      observer.unobserve(lazyApiLoadObserverTarget.value);
      lazyApiLoadObserverTarget.value = undefined;
    } catch (err) {
      /* ignore error */
      console.error("unbserve failed", err);
    }
  }
});
</script>
<template>
  <div
    ref="lazyApiLoadObserverTarget"
    :style="{
      'min-height': `${props.minHeight}px`,
      'min-width': `${props.minWidth}px`,
    }"
  >
    <slot></slot>
  </div>
</template>
