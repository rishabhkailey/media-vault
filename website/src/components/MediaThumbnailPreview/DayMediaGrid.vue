<script setup lang="ts">
import { computed, ref } from "vue";
import { daysShort, monthShort } from "@/js/date";
import SelectWrapper from "@/components/SelectWrapper/SelectWrapper.vue";
import createJustifiedLayout from "justified-layout";
import MediaThumbnail from "./MediaThumbnail.vue";

const props = defineProps<{
  day: number;
  date: number;
  month: number;
  year: number;
  mediaList: Array<Media>;
  loadAllMediaUntil: (date: Date) => Promise<any>;
  selectedMediaMap: Map<number, boolean>;
}>();

const emits = defineEmits<{
  (e: "thumbnailClick", mediaID: number): void;
  (e: "selectMedia", mediaID: number, value: boolean): void;
  (e: "selectAllMedia", value: boolean): void;
}>();

const containerRef = ref<HTMLElement | undefined>(undefined);
const justifiedLayout = computed<Array<WidthHeight>>(() => {
  if (containerRef.value === undefined) {
    return [];
  }
  let mediaList = props.mediaList;
  let aspectRatios = mediaList.map<number>((media) => {
    if (media.thumbnail_aspect_ratio == 0) {
      return 1;
    }
    return media.thumbnail_aspect_ratio;
  });
  const geometry = createJustifiedLayout(aspectRatios, {
    containerWidth: containerRef.value.getBoundingClientRect().width,
    targetRowHeight: 150,
    showWidows: true,
    boxSpacing: 10,
  });
  return geometry.boxes.map((e) => ({
    width: e.width,
    height: e.height,
  }));
});

function getSelection(index: number): boolean {
  return !!props.selectedMediaMap?.get(index);
}
const dayMediaSelected = computed(() => {
  for (let media of props.mediaList) {
    if (!props.selectedMediaMap.get(media.id)) {
      return false;
    }
  }
  return true;
});
</script>

<template>
  <v-card class="bg-secondary-background w-100">
    <v-card-subtitle>
      <SelectWrapper
        :absolute-position="false"
        :model-value="dayMediaSelected"
        @change="(value) => emits('selectAllMedia', value)"
        selectIconSize="small"
        :always-show-select-button="selectedMediaMap.size > 0"
        :always-show-select-on-mobile="true"
        :show-select-button-on-hover="true"
        :select-on-content-click="selectedMediaMap.size > 0"
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
        v-for="({ width, height }, index) in justifiedLayout"
        :key="`${props.mediaList[index].id}_${index}`"
        :width="width"
        :height="height"
        class="bg-surface thumbnail-container"
        @click.stop="
          (e) => {
            emits('thumbnailClick', props.mediaList[index].id);
          }
        "
      >
        <SelectWrapper
          :loading="false"
          :absolute-position="true"
          :model-value="getSelection(props.mediaList[index].id)"
          :always-show-select-button="selectedMediaMap.size > 0"
          :always-show-select-on-mobile="true"
          :show-select-button-on-hover="true"
          :select-on-content-click="selectedMediaMap.size > 0"
          @change="
            (value) => {
              emits('selectMedia', props.mediaList[index].id, value);
            }
          "
          selectIconSize="large"
        >
          <MediaThumbnail
            :class="{
              'shrink-transition': true,
              shrink: getSelection(props.mediaList[index].id),
            }"
            :media="props.mediaList[index]"
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
