<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { daysShort, monthShort } from "@/js/date";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { storeToRefs } from "pinia";
import SelectWrapper from "@/components/SelectWrapper/SelectWrapper.vue";
import createJustifiedLayout from "justified-layout";
import { useErrorsStore } from "@/piniaStore/errors";
import MediaThumbnail from "./MediaThumbnail.vue";
// interface doesn't work https://github.com/vuejs/core/issues/4294
// const props = defineProps<DailyMedia>();
const props = defineProps<{
  day: number;
  date: number;
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  loadAllMediaUntil: (date: Date) => Promise<any>;
}>();

const emits = defineEmits<{
  (
    e: "thumbnailClick",
    mediaID: number,
    index: number,
    clickLocation: ThumbnailClickLocation | undefined,
  ): void;
}>();

function getThumbnailLocation(
  event: MouseEvent,
): ThumbnailClickLocation | undefined {
  let element: HTMLElement;
  if (event.target === null || !(event.target instanceof HTMLElement)) {
    return undefined;
  }
  element = event.target;
  try {
    let boundingRect = element.getBoundingClientRect();
    return {
      left: boundingRect.left,
      top: boundingRect.top,
      height: boundingRect.height,
      width: boundingRect.width,
      x: boundingRect.x,
      y: boundingRect.y,
    };
  } catch (err) {
    // ignore animation related error
    console.error(err);
    return undefined;
  }
}

const { appendError } = useErrorsStore();
const containerRef = ref<HTMLElement | undefined>(undefined);
let widowsCount = 0;
const justifiedLayout = ref<Array<WidthHeight>>([]);
watch(
  [() => props.indexMediaList, containerRef],
  ([newIndexMediaList, newContainerRef]) => {
    if (newContainerRef === undefined) {
      appendError(
        "media layout failed",
        "container reference is undefined unable to calculate media layout. ",
        -1,
      );
      return;
    }
    justifiedLayout.value = justifiedLayout.value.slice(
      0,
      justifiedLayout.value.length - widowsCount,
    );
    // let startIndex = justifiedLayout.value.length - widowsCount;
    let startIndex = justifiedLayout.value.length;
    let aspectRatios: Array<number> = [];
    for (let i = startIndex; i < newIndexMediaList.length; i++) {
      // for media with no thumbnail
      if (newIndexMediaList[i].media.thumbnail_aspect_ratio == 0) {
        newIndexMediaList[i].media.thumbnail_aspect_ratio = 1;
      }
      aspectRatios.push(newIndexMediaList[i].media.thumbnail_aspect_ratio);
    }
    const geometry = createJustifiedLayout(aspectRatios, {
      containerWidth: newContainerRef.getBoundingClientRect().width,
      targetRowHeight: 150,
      showWidows: true,
      boxSpacing: 10,
    });
    widowsCount = geometry.widowCount;
    const boxes = geometry.boxes.map((e) => ({
      width: e.width,
      height: e.height,
    }));
    justifiedLayout.value = [...justifiedLayout.value, ...boxes];
  },
);

const selectDayMediaLoading = ref(false);
const mediaSelectionStore = useMediaSelectionStore();
const { selectionMap, count: selectedMediaIDsCount } =
  storeToRefs(mediaSelectionStore);
const { updateSelection } = mediaSelectionStore;

function getSelection(index: number): boolean {
  return !!selectionMap.value?.get(index);
}

const dayMediaSelected = computed(() => {
  for (let { media } of props.indexMediaList) {
    if (!selectionMap.value.get(media.id)) {
      return false;
    }
  }
  return true;
});

async function selectDayMedia(value: boolean) {
  selectDayMediaLoading.value = true;
  await props.loadAllMediaUntil(
    // month is 0 indexed in js, so + 1
    new Date(`${props.year}-${props.month + 1}-${props.date}`),
  );
  props.indexMediaList.forEach(({ media }) => {
    updateSelection(media.id, value);
  });
  selectDayMediaLoading.value = false;
}
</script>

<template>
  <v-card class="bg-secondary-background w-100">
    <v-card-subtitle>
      <SelectWrapper
        :loading="selectDayMediaLoading"
        :absolute-position="false"
        :model-value="dayMediaSelected"
        @change="selectDayMedia"
        selectIconSize="small"
        :always-show-select-button="selectedMediaIDsCount > 0"
        :always-show-select-on-mobile="true"
        :show-select-button-on-hover="true"
        :select-on-content-click="selectedMediaIDsCount > 0"
      >
        <!-- todo unknown date -->
        <span
          style="
            font-size: 0.875rem;
            font-weight: 500;
            letter-spacing: 0.0175em;
          "
        >
          {{
            `${daysShort[props.day]}, ${monthShort[props.month]} ${
              props.date
            }, ${props.year}`
          }}
        </span>
      </SelectWrapper>
    </v-card-subtitle>
    <div
      ref="containerRef"
      style="padding: 10px; display: flex; flex-wrap: wrap; gap: 10px"
    >
      <div
        :key="`${index}`"
        v-for="({ width, height }, index) in justifiedLayout"
        :width="width"
        :height="height"
        class="bg-surface thumbnail-container"
        @click.stop="
          (e) => {
            emits(
              'thumbnailClick',
              props.indexMediaList[index].media.id,
              props.indexMediaList[index].index,
              getThumbnailLocation(e),
            );
          }
        "
      >
        <SelectWrapper
          :loading="false"
          :absolute-position="true"
          :model-value="getSelection(props.indexMediaList[index].media.id)"
          :always-show-select-button="selectedMediaIDsCount > 0"
          :always-show-select-on-mobile="true"
          :show-select-button-on-hover="true"
          :select-on-content-click="selectedMediaIDsCount > 0"
          @change="
            (value) => {
              updateSelection(props.indexMediaList[index].media.id, value);
            }
          "
          selectIconSize="large"
        >
          <MediaThumbnail
            :class="{
              'shrink-transition': true,
              shrink: getSelection(props.indexMediaList[index].media.id),
            }"
            :media="props.indexMediaList[index].media"
            :width="width"
            :height="height"
            transition="scale"
          />
        </SelectWrapper>
      </div>
    </div>
  </v-card>
</template>

<style scoped>
.thumbnail-container:hover {
  cursor: pointer;
  opacity: 90%;
  transition: all 0.15s ease-in-out;
}
.shrink-transition {
  transition: all 0.15s ease-in-out;
}
.shrink {
  transform: scale(0.9);
}
</style>
