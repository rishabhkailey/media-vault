<script setup lang="ts">
import MediaViewer from "@/components/MediaViewer/MediaViewer.vue";
import { download } from "@/js/encryptedFileDownload";
import { useErrorsStore } from "@/piniaStore/errors";
import { onMounted, onUpdated, ref } from "vue";
import { computed } from "vue";
import { MEDIA_CAROUSEL_HEADER_Z_INDEX } from "@/js/constants/z-index";

const props = withDefaults(
  defineProps<{
    index: number;
    mediaList: Array<Media>;
    loadMoreMedia: LoadMoreMedia;
    allMediaLoaded: boolean;
    routeName: string;
    animationOriginSelector: string;
    loading?: boolean;
  }>(),
  {
    loading: false,
  },
);

const emits = defineEmits<{
  (e: "close"): void;
  (e: "next"): void;
  (e: "previous"): void;
  (e: "update:index", value: number): void;
}>();

const media = computed<Media>(() => {
  return props.mediaList[props.index];
});
const rootContainer = ref<HTMLElement | null>(null);

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

function scaleWindowToThumbnailSizeWithoutCrop(
  window: DOMRect,
  thumbnail: DOMRect,
) {
  let initialHeight = thumbnail.height;
  let initialWidth = thumbnail.width;
  let initialX = thumbnail.x;
  let initialY = thumbnail.y;
  // aspect ration of the element should be equal to display aspect ratio
  // we don't want cut/crop the image
  if (window.width / window.height > media.value.thumbnail_aspect_ratio) {
    // actual display is horizontaly longer than the photo
    // bottle neck will be the height, we will be increasing the height
    initialWidth = initialHeight * (window.width / window.height);
    initialX = initialX - Math.abs(thumbnail.width - initialWidth) / 2;
  } else {
    initialHeight = initialWidth / (window.width / window.height);
    initialY = initialY - Math.abs(thumbnail.height - initialHeight) / 2;
  }
  return {
    x: initialX,
    y: initialY,
    width: initialWidth,
    height: initialHeight,
  };
}
function startImageOpenAnimation() {
  try {
    let thumbnailElement = document.querySelector(
      props.animationOriginSelector,
    );
    let mediaWindowElement = document.getElementById(
      `media_window_${media.value.id}`,
    );
    if (thumbnailElement === null || mediaWindowElement === null) {
      return;
    }
    // right now the thumbnails are always rendered, we are rendering media on top of thumbnails
    // if we change this behavior this might stop working, in that case we can use route hash to get thumbnailRect
    let thumbnailRect = thumbnailElement.getBoundingClientRect();
    let mediaWindowRect = mediaWindowElement.getBoundingClientRect();
    let {
      x: initialX,
      y: initialY,
      height: initialHeight,
      width: initialWidth,
    } = scaleWindowToThumbnailSizeWithoutCrop(mediaWindowRect, thumbnailRect);
    let startFrame: Keyframe = {
      transform: `translate3D(${
        initialX -
        mediaWindowRect.x -
        (mediaWindowRect.width - initialWidth) / 2
      }px, ${
        initialY -
        mediaWindowRect.y -
        (mediaWindowRect.height - initialHeight) / 2
      }px, 0) scale3d(${initialWidth / mediaWindowRect.width}, ${
        initialHeight / mediaWindowRect.height
      }, 1)`,
    };
    let endFrame: Keyframe = {
      transform: `translate3D(0px, 0px, 0px) scale3d(1, 1, 1)`,
    };
    return mediaWindowElement.animate([startFrame, endFrame], {
      duration: 150,
      easing: "ease-in",
      fill: "both",
    });
  } catch (err) {
    // ignore animation related error
    console.error(err);
    return;
  }
}
function startBackgroundOpenAnimation() {
  if (rootContainer.value === null) {
    return;
  }
  return rootContainer.value.animate(
    [
      {
        backgroundColor: "rgba(0, 0, 0, 0)",
      },
      {
        backgroundColor: "rgba(0, 0, 0, 1)",
      },
    ],
    {
      duration: 150,
      easing: "ease-in",
      fill: "both",
    },
  );
}
function startImageCloseAnimation() {
  try {
    let thumbnailElement = document.querySelector(
      props.animationOriginSelector,
    );
    let mediaWindowElement = document.getElementById(
      `media_window_${media.value.id}`,
    );
    if (thumbnailElement === null || mediaWindowElement === null) {
      return;
    }
    // right now the thumbnails are always rendered, we are rendering media on top of thumbnails
    // if we change this behavior this might stop working, in that case we can use route hash to get thumbnailRect
    let thumbnailRect = thumbnailElement.getBoundingClientRect();
    let mediaWindowRect = mediaWindowElement.getBoundingClientRect();
    let {
      x: initialX,
      y: initialY,
      height: initialHeight,
      width: initialWidth,
    } = scaleWindowToThumbnailSizeWithoutCrop(mediaWindowRect, thumbnailRect);

    let startFrame: Keyframe = {
      transform: `translate3D(0px, 0px, 0px) scale3d(1, 1, 1)`,
    };
    let endFrame: Keyframe = {
      transform: `translate3D(${
        initialX -
        mediaWindowRect.x -
        (mediaWindowRect.width - initialWidth) / 2
      }px, ${
        initialY -
        mediaWindowRect.y -
        (mediaWindowRect.height - initialHeight) / 2
      }px, 0) scale3d(${initialWidth / mediaWindowRect.width}, ${
        initialHeight / mediaWindowRect.height
      }, 1)`,
    };
    return mediaWindowElement.animate([startFrame, endFrame], {
      duration: 150,
      fill: "both",
      easing: "ease-in",
    });
  } catch (err) {
    // ignore animation related error
    console.error(err);
    return;
  }
}
function startBackgroundCloseAnimation() {
  if (rootContainer.value === null) {
    return;
  }
  return rootContainer.value.animate(
    [
      {
        backgroundColor: "rgba(0, 0, 0, 1)",
      },
      {
        backgroundColor: "rgba(0, 0, 0, 0)",
      },
    ],
    {
      duration: 150,
      easing: "ease-in",
      fill: "both",
    },
  );
}

const { appendError } = useErrorsStore();
function downloadMedia(media: Media) {
  download(media.url, media.name).catch((err) => {
    let errorMessage = "";
    if (typeof err === "string") {
      errorMessage = err;
    }
    if (err instanceof Error) {
      errorMessage = err.message + " " + err.stack;
    }
    appendError(`Download failed ${media.name}`, errorMessage, -1);
  });
}

onMounted(() => {
  startBackgroundOpenAnimation();
  startImageOpenAnimation();
  loadMoreMediaIfRequired();
});

onUpdated(() => {
  loadMoreMediaIfRequired();
});

async function close() {
  let closeAnimation = startImageCloseAnimation();
  let backgroundCloseAnimation = startBackgroundCloseAnimation();
  if (closeAnimation !== undefined) {
    try {
      await closeAnimation.finished;
    } catch (err) {
      // ignore animation related error
      console.error(err);
    }
  }
  if (backgroundCloseAnimation !== undefined) {
    try {
      await backgroundCloseAnimation.finished;
    } catch (err) {
      // ignore animation related error
      console.error(err);
    }
  }
  emits("close");
}
</script>

<template>
  <div
    class="d-flex flex-column justify-center align-stretch ma-0 pa-0"
    style="width: 100vw; height: 100vh"
    ref="rootContainer"
  >
    <!-- header -->
    <div
      class="pt-1 pr-4 d-flex justify-end align-center media-carousel-header"
    >
      <v-btn
        icon="mdi-download"
        @click.stop="() => downloadMedia(media)"
        style="background: none; border: none; box-shadow: none"
      >
        <v-icon color="white">mdi-download</v-icon>
      </v-btn>
      <v-btn
        icon="mdi-close"
        @click.stop="close"
        style="background: none; border: none; box-shadow: none"
        data-test-id="media-carousel-close-button"
      >
        <v-icon color="white">mdi-close</v-icon>
      </v-btn>
    </div>

    <!-- v-row -->
    <div
      class="ma-0 pa-0 d-flex flex-row flex-grow-1 d-flex flex-column justify-center align-center"
      style="height: 100%; width: 100%"
    >
      <v-progress-circular
        v-if="props.loading"
        indeterminate
        :size="80"
        :width="5"
      ></v-progress-circular>
      <v-window
        v-else
        style="height: 100%; width: 100%"
        :model-value="index"
        @update:model-value="
          (value) => {
            emits('update:index', value);
          }
        "
        :continuous="false"
        :show-arrows="true"
        touch
        direction="horizontal"
      >
        <v-window-item
          style="height: 100%; width: 100%"
          v-for="(media, index) in mediaList"
          :key="index"
        >
          <template v-slot:default>
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
          </template>
        </v-window-item>
      </v-window>
    </div>
  </div>
</template>

<style scoped>
.media-carousel-header {
  box-shadow: 0px 5px 15px rgb(50, 50, 50, 0.3);
  position: absolute;
  width: 100vw;
  top: 0;
  z-index: v-bind(MEDIA_CAROUSEL_HEADER_Z_INDEX);
  background-color: rgba(0, 0, 0, 0.4);
}
</style>
