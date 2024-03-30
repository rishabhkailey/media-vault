<script setup lang="ts">
import MediaCarousel from "@/components/MediaCarousel/MediaCarousel.vue";
import OpenCloseAnimation from "@/components/Animations/OpenCloseAnimation.vue";
import { getQueryParamNumberValue } from "@/js/utils";
import { useMediaStore } from "@/piniaStore/media";
import {
  errorScreenRoute,
  homeRoute,
  mediaPreviewRoute,
} from "@/router/routesConstants";
import { storeToRefs } from "pinia";
import { onBeforeMount, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getSingleMediaById } from "@/js/api/media";
import { useErrorsStore } from "@/piniaStore/errors";

const router = useRouter();
const route = useRoute();

const { appendError } = useErrorsStore();

// params
const index = ref(-1);
const mediaID = ref(-1);

function initParams() {
  // media_id
  let mediaIdParam = getQueryParamNumberValue(route.params, "media_id");
  if (mediaIdParam === undefined) {
    throw new Error(`invalid media id param`);
  }
  mediaID.value = Number(mediaIdParam);

  // media index
  let indexParam = getQueryParamNumberValue(route.params, "index");
  if (indexParam === undefined) {
    throw new Error(`invalid index param`);
  }
  index.value = indexParam;
}

let allMediaLoaded = ref(true);
let mediaList = ref<Array<Media>>([]);
let loadMoreMedia: LoadMoreMedia;

const mediaStore = useMediaStore();
function initMediaPreviewRefsAndStore() {
  if (
    mediaStore.mediaList.findIndex((m) => m.id === mediaID.value) !==
    index.value
  ) {
    initSingleMediaPreviewRefsAndStore();
    return;
  }

  ({ allMediaLoaded, mediaList } = storeToRefs(mediaStore));
  ({ loadMoreMedia } = mediaStore);
}

const loading = ref<boolean>(false);
function initSingleMediaPreviewRefsAndStore() {
  loading.value = true;
  getSingleMediaById(mediaID.value)
    .then((media) => {
      allMediaLoaded.value = true;
      mediaList.value = [media];
      index.value = 0;
      loadMoreMedia = () =>
        new Promise<LoadMoreMediaStatus>((resolve) => resolve("empty"));
    })
    .catch((err) => {
      appendError(
        "failed to get media info from server",
        `error message - ${err}`,
        -1,
      );
    })
    .finally(() => {
      loading.value = false;
    });
}

function updateIndex(newIndex: number) {
  index.value = newIndex;
  router.push(mediaPreviewRoute(newIndex, mediaList.value[newIndex].id));
}

onBeforeMount(() => {
  try {
    initParams();
    initMediaPreviewRefsAndStore();
  } catch (err) {
    console.error(err);
    router.push(
      errorScreenRoute(
        "UserMediaCarousel component intialization failed",
        `error message = "${err}"`,
      ),
    );
  }
});
</script>
<template>
  <OpenCloseAnimation
    :source-element-selector="`#thumbnail_${mediaList[index]?.id}`"
    :target-element-selector="`#user_media_preview_container`"
    :_key="`${mediaList[index]?.id}`"
    :time-ms="150"
    class=""
    style="width: 100vw; height: 100vh"
  >
    <template #default="{ closeAnimation }">
      <MediaCarousel
        :loading="loading"
        :index="index"
        :media-list="mediaList"
        :load-more-media="loadMoreMedia"
        :all-media-loaded="allMediaLoaded"
        :animation-origin-selector="`#thumbnail_${mediaList[index]?.id}`"
        @update:index="updateIndex"
        @close="
          async () => {
            try {
              await closeAnimation();
            } catch {}
            router.push(homeRoute());
          }
        "
      />
    </template>
  </OpenCloseAnimation>
</template>
