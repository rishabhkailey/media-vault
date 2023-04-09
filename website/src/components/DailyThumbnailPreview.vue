<script setup lang="ts">
import MediaThumbnail from "./MediaThumbnail.vue";
import MediaPreview from "./MediaPreview.vue";
import { ref } from "vue";
import { daysShort, monthShort } from "@/utils/date";
// interface doesn't work https://github.com/vuejs/core/issues/4294
// const props = defineProps<DailyMedia>();
const props = defineProps<{
  day: number;
  date: number;
  month: number;
  year: number;
  indexMediaList: Array<IndexMedia>;
}>();

const selectedMediaIndex = ref<number | undefined>(undefined);
const prviewOverlay = ref<boolean>(false);
</script>

<template>
  <v-card
    :subtitle="`${daysShort[props.day]}, ${monthShort[props.month]} ${
      props.date
    }, ${props.year}`"
    class="bg-secondary-background"
  >
    <div>
      <div class="d-flex flex-row flex-wrap">
        <div
          :key="`${index}+${media.name}`"
          v-for="{ media, index } in props.indexMediaList"
          class="d-flex child-flex pa-2"
          style="max-height: 300px"
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
        </div>
      </div>
      <!-- todo move this to media thumbnail? -->
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
