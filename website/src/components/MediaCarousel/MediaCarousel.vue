<script setup lang="ts">
import MediaViewer from "@/components/MediaViewer/MediaViewer.vue";
import MediaCarouselHeader from "@/components/MediaCarousel/MediaCarouselHeader.vue";
import { computed, onBeforeUnmount, onUpdated } from "vue";
import { onMounted } from "vue";
import { ref } from "vue";
import { MEDIA_CAROUSEL_HEADER_Z_INDEX } from "@/js/constants/z-index";
import { useErrorsStore } from "@/piniaStore/errors";

const props = withDefaults(
  defineProps<{
    index: number;
    mediaList: Array<Media>;
    loadMoreMedia: LoadMoreMedia;
    allMediaLoaded: boolean;
    loading?: boolean;
  }>(),
  {
    loading: false,
  },
);

const emits = defineEmits<{
  (e: "close"): void;
  (e: "update:index", value: number): void;
}>();

// required for single media init, because momentarily media list is empty
const media = computed<undefined | Media>(() => {
  if (props.index > props.mediaList.length - 1) {
    return undefined;
  }
  return props.mediaList[props.index];
});

const lastMouseMoveTime = ref<number>(new Date().getTime());
const showArrows = ref<boolean>(true);
const animationName = ref<string>("");
const { appendError } = useErrorsStore();

function loadMoreMediaIfRequired() {
  if (props.index > props.mediaList.length - 2 && !props.allMediaLoaded) {
    props.loadMoreMedia().catch((err) => {
      appendError(
        "Failed to load more media please refresh the page if facing any issues",
        `error - ${err}`,
        10,
      );
    });
  }
}

function left() {
  if (props.index === 0) {
    return;
  }
  animationName.value = "left-click-animation";
  // index.value--;
  emits("update:index", props.index - 1);
}

function right() {
  console.log("right");
  if (props.index === props.mediaList.length - 1) {
    return;
  }
  animationName.value = "right-click-animation";
  // index.value++;
  console.log("next");
  emits("update:index", props.index + 1);
}

let showArrowEvaluationTimer: undefined | NodeJS.Timer;
function updateLastMouseMoveTime() {
  lastMouseMoveTime.value = new Date().getTime();
}

onMounted(() => {
  document.addEventListener("mousemove", updateLastMouseMoveTime, false);
  document.addEventListener("mousedown", updateLastMouseMoveTime, false);

  showArrowEvaluationTimer = setInterval(() => {
    const secondsSinceMouseMove =
      (new Date().getTime() - lastMouseMoveTime.value) / 1000;
    // if less than 5
    // console.log("difference", secondsSinceMouseMove);
    showArrows.value = secondsSinceMouseMove < 5;
  }, 500);
});

onBeforeUnmount(() => {
  document.removeEventListener("mousemove", updateLastMouseMoveTime, false);
  document.removeEventListener("mousedown", updateLastMouseMoveTime, false);
  clearInterval(showArrowEvaluationTimer);
});

const touchStartX = ref<undefined | number>(undefined);
const touchEndX = ref<undefined | number>(undefined);
function touchStart(event: TouchEvent) {
  touchStartX.value = event.targetTouches.item(event.targetTouches.length - 1)
    ?.clientX;
}

function touchMove(event: TouchEvent) {
  touchEndX.value = event.targetTouches.item(event.targetTouches.length - 1)
    ?.clientX;
}

function touchEnd() {
  if (touchStartX.value === undefined || touchEndX.value === undefined) {
    return;
  }
  const xDiff = touchEndX.value - touchStartX.value;
  // if diff is less than 10% of the window width
  if (Math.abs(xDiff) < window.innerWidth * 0.1) {
    return;
  }
  if (xDiff < 0) {
    right();
    touchStartX.value = undefined;
    touchEndX.value = undefined;
    return;
  }
  left();
  touchStartX.value = undefined;
  touchEndX.value = undefined;
}

onMounted(() => {
  loadMoreMediaIfRequired();
});

onUpdated(() => {
  loadMoreMediaIfRequired();
});
</script>

<template>
  <!-- header -->
  <MediaCarouselHeader
    class="pt-1 pr-4 d-flex justify-end align-center media-carousel-header"
    @close="() => emits('close')"
  />
  <!-- content -->
  <div
    id="contianer"
    style="width: 100vw; height: 100vh"
    class="d-flex flex-row bg-secondary-background"
    @touchstart="touchStart"
    @touchmove="touchMove"
    @touchend="touchEnd"
  >
    <div
      v-if="media !== undefined"
      :style="`width: 100vw; height: 100vh; transform: translateX(${
        touchStartX == undefined || touchEndX == undefined
          ? 0
          : touchEndX - touchStartX
      }px)`"
      :class="`d-flex justify-center align-center ${animationName}`"
      :key="media.id"
    >
      <MediaViewer
        style="
          width: 100%;
          height: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
        :id="`media_window_${media.id}`"
        :media="media"
      >
      </MediaViewer>
    </div>
  </div>

  <!-- actions -->
  <!-- LEFT button -->
  <div
    v-if="index > 0 && showArrows"
    style="position: absolute; left: 0vw; top: 45vh; margin-left: 10px"
  >
    <v-btn icon="mdi-menu-left" @click.stop="left" />
  </div>

  <!-- RIGHT button -->
  <div
    v-if="index < mediaList.length - 1 && showArrows"
    style="position: absolute; right: 0vw; top: 45vh; margin-right: 10px"
  >
    <v-btn icon="mdi-menu-right" @click.stop="right" :loading="props.loading" />
  </div>
</template>

<style scoped>
.right-click-animation {
  animation-duration: 0.3s;
  animation-name: right-click;
  animation-timing-function: ease-in;
}

@keyframes right-click {
  0% {
    transform: translateX(50%);
  }
  100% {
    transform: translateX(0%);
  }
}

.left-click-animation {
  animation-duration: 0.3s;
  animation-name: left-click;
  animation-timing-function: ease-in;
}

@keyframes left-click {
  0% {
    transform: translateX(-50%);
  }
  100% {
    transform: translateX(0%);
  }
}

.media-carousel-header {
  box-shadow: 0px 5px 15px rgb(50, 50, 50, 0.3);
  position: absolute;
  width: 100vw;
  top: 0;
  z-index: v-bind(MEDIA_CAROUSEL_HEADER_Z_INDEX);
  background-color: rgba(0, 0, 0, 0.4);
  height: 4em;
}
</style>
