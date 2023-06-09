<script setup lang="ts">
import MediaThumbnail from "./MediaThumbnail.vue";
import MediaPreview from "@/components/MediaPreview.vue";
import { computed, ref } from "vue";
import { daysShort, monthShort } from "@/utils/date";
import { useMediaSelectionStore } from "@/piniaStore/mediaSelection";
import { storeToRefs } from "pinia";
import SelectWrapper from "@/components/SelectWrapper/SelectWrapper.vue";
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
  <v-card class="bg-secondary-background">
    <v-card-subtitle>
      <SelectWrapper
        :loading="selectDayMediaLoading"
        :absolute-position="false"
        :model-value="dayMediaSelected"
        @change="selectDayMedia"
        selectIconSize="small"
        :always-show-select-button="selectedMediaIDsCount > 0"
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
    <div>
      <div class="d-flex flex-row flex-wrap">
        <div
          :key="`${index}+${media.name}`"
          v-for="{ media, index } in props.indexMediaList"
          class="d-flex child-flex pa-2"
        >
          <SelectWrapper
            :loading="false"
            :absolute-position="true"
            :model-value="getSelection(media.id)"
            :always-show-select-button="selectedMediaIDsCount > 0"
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
              :aspect-ratio="1"
              :height="175"
              :width="175"
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
        </div>
      </div>
      <!-- todo move this to media thumbnail? -->
      <v-overlay
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
      </v-overlay>
    </div>
  </v-card>
</template>
