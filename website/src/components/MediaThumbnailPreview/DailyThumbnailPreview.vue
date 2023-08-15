<script setup lang="ts">
import MediaThumbnail from "./MediaThumbnail.vue";
import MediaPreview from "@/components/MediaPreview.vue";
import { computed, ref, toRef, watch } from "vue";
import { daysShort, monthShort } from "@/js/date";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { storeToRefs } from "pinia";
import SelectWrapper from "@/components/SelectWrapper/SelectWrapper.vue";
import { useDisplay } from "vuetify";
import createJustifiedLayout from "justified-layout";
import { useRouter } from "vue-router";
// interface doesn't work https://github.com/vuejs/core/issues/4294
// const props = defineProps<DailyMedia>();
const props = defineProps<{
  day: number;
  date: number;
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
  loadAllMediaOfDate: (date: Date) => Promise<any>;
  // todo load all day media function
}>();
const router = useRouter();

// const display = useDisplay();
// const width = computed<number>(() => {
//   switch (display.name.value) {
//     case "xs":
//       return display.width.value / 2.2;
//     case "sm":
//       return display.width.value / 4.2;
//     case "md":
//       return display.width.value / 6.2;
//     case "lg":
//       return display.width.value / 6.2;
//     case "xl":
//       return display.width.value / 6.2;
//     case "xxl":
//       return display.width.value / 12;
//     default:
//       return display.width.value / 2.2;
//   }
// });
console.log(props.indexMediaList);
const containerRef = ref<HTMLElement | undefined>(undefined);
let widowsCount = 0;
// const justifiedLayout = computed<Array<WidthHeight>>(() => {
//   if (containerRef.value === undefined) {
//     return [];
//   }

//   return [];
// });
const justifiedLayout = ref<Array<WidthHeight>>([]);
watch(
  [() => props.indexMediaList, containerRef],
  ([newIndexMediaList, newContainerRef], [oldIndexMediaList]) => {
    console.log("in watch", newContainerRef, containerRef.value);
    if (newContainerRef === undefined) {
      console.error("containerRef undefined");
      return;
    }
    console.log("justifiedLayout", containerRef.value);
    justifiedLayout.value = justifiedLayout.value.slice(
      0,
      justifiedLayout.value.length - widowsCount
    );
    // let startIndex = justifiedLayout.value.length - widowsCount;
    let startIndex = justifiedLayout.value.length;
    console.log(
      justifiedLayout.value,
      justifiedLayout.value.length,
      widowsCount
    );
    let aspectRatios: Array<number> = [];
    console.log(newIndexMediaList);
    for (let i = startIndex; i < newIndexMediaList.length; i++) {
      console.log(i);
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
    console.log(
      "justifiedLayout",
      aspectRatios,
      justifiedLayout.value,
      widowsCount,
      newIndexMediaList.length
    );
  }
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
  await props.loadAllMediaOfDate(
    // month is 0 indexed in js, so + 1
    new Date(`${props.year}-${props.month + 1}-${props.date}`)
  );
  props.indexMediaList.forEach(({ media }) => {
    updateSelection(media.id, value);
  });
  selectDayMediaLoading.value = false;
}

const selectedMediaIndex = ref<number | undefined>(undefined);
const prviewOverlay = ref<boolean>(false);
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
        {{
          `${daysShort[props.day]}, ${monthShort[props.month]} ${props.date}, ${
            props.year
          }`
        }}
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
          () => {
            router.push({
              name: `MediaPreview`,
              params: {
                index: props.indexMediaList[index].index,
                media_id: props.indexMediaList[index].media.id,
              },
            });
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
          <v-img
            :src="props.indexMediaList[index].media.thumbnail_url"
            :width="width"
            :height="height"
            transition="scale"
            :class="{
              'shrink-transition': true,
              shrink: getSelection(props.indexMediaList[index].media.id),
            }"
            cover
          />
        </SelectWrapper>
      </div>
      <!-- <MediaThumbnail
          class="h-100 w-100"
          style="aspect-ratio: 1"
          :aspect-ratio="1"
          :padding="getSelection(media.id) ? 10 : 0"
          :media="media"
          @click="
            () => {
              prviewOverlay = true;
              selectedMediaIndex = index;
            }
          "
        /> -->
      <!-- <div class="d-flex flex-row flex-wrap">
        <v-col
          cols="6"
          sm="4"
          md="3"
          lg="2"
          xl="2"
          xxl="1"
          :key="`${index}+${media.name}`"
          v-for="{ media, index } in props.indexMediaList"
          class="d-flex child-flex pa-1"
        >
          <SelectWrapper
            class="h-100 w-100"
            style="aspect-ratio: 1"
            :loading="false"
            :absolute-position="true"
            :model-value="getSelection(media.id)"
            :always-show-select-button="selectedMediaIDsCount > 0"
            :always-show-select-on-mobile="true"
            :show-select-button-on-hover="true"
            :select-on-content-click="selectedMediaIDsCount > 0"
            @change="
              (value) => {
                updateSelection(media.id, value);
              }
            "
            selectIconSize="large"
          >
            <MediaThumbnail
              class="h-100 w-100"
              style="aspect-ratio: 1"
              :aspect-ratio="1"
              :padding="getSelection(media.id) ? 10 : 0"
              :media="media"
              @click="
                () => {
                  prviewOverlay = true;
                  selectedMediaIndex = index;
                }
              "
            />
          </SelectWrapper>
        </v-col>
      </div> -->
      <!-- todo move this to media thumbnail? -->
      <!-- <v-overlay
        v-model="prviewOverlay"
        transition="fade-transition"
        :close-on-content-click="false"
        :close-delay="0"
        :open-delay="0"
        class="d-flex justify-center align-center"
        :z-index="2000"
      >
        <MediaPreview
          v-if="selectedMediaIndex !== undefined"
          :index="selectedMediaIndex"
          @close="
            () => {
              prviewOverlay = false;
            }
          "
        />
      </v-overlay> -->
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
