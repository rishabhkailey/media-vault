<script setup lang="ts">
// import MediaThumbnail from "@/components/MediaThumbnail.vue";
// import { ref } from "vue";

// const media: Media = {
//   thumbnail_url: "/v1/thumbnail/c72ce9d9-0763-4876-825f-a6b1791bfc9f",
//   thumbnail: true,
//   date: new Date(),
//   name: "test",
//   size: 100,
//   type: "image/jpeg",
//   url: "",
// };

const props = defineProps<{
  modelValue: boolean;
  absolutePosition: boolean;
  selectIconSize: string | number;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "change", value: boolean): void;
}>();

const clickHanlder = () => {
  // selected.value = !selected.value;
  emit("update:modelValue", !props.modelValue);
  emit("change", !props.modelValue);
};
// const selected = ref(false);
</script>
<template>
  <!-- using badge, hover 
  <v-hover>
    <template v-slot:default="{ isHovering, props }">
      <div v-bind="props" class="d-flex child-flex pa-2">
        <v-badge
          location="top start"
          color="transparent"
          icon="mdi-check-circle-outline"
          :model-value="isHovering"
        >
          <div class="d-flex child-flex pa-2">
            <MediaThumbnail :media="media" @click="() => {}" />
          </div>
        </v-badge>
      </div>
    </template>
  </v-hover>
  -->

  <!-- using badge (slot) and hover
  <v-hover>
    <template v-slot:default="{ isHovering, props }">
      <div v-bind="props" class="d-flex child-flex pa-2">
        <v-badge
          location="top start"
          color="transparent"
          :model-value="isHovering || selected"
        >
          <div class="d-flex child-flex pa-2">
            <MediaThumbnail :media="media" @click="() => {}" />
          </div>
          <template v-slot:badge>
            <v-icon
              icon="mdi-check-circle-outline"
              class="select-badge"
              @click="onSelect"
              :color="selected ? 'success' : 'white'"
              size="x-large"
            />
          </template>
        </v-badge>
      </div>
    </template>
  </v-hover>
  -->

  <!-- todo try custom method so we can set the size of the select button
  <v-hover>
    <template v-slot:default="{ isHovering, props }">
      <div v-bind="props" class="d-flex child-flex pa-2">
        <v-scale-transition>
          <div
            v-if="isHovering || selected"
            style="position: absolute; z-index: 1"
          >
            <v-icon
              icon="mdi-check-circle"
              class="select-badge"
              @click="onSelect"
              :color="selected ? 'success' : 'white'"
              size="large"
            />
          </div>
        </v-scale-transition>
        <div :class="`d-flex child-flex ${selected ? 'pa-2' : ''}`">
          <MediaThumbnail :media="media" @click="() => {}" />
        </div>
      </div>
    </template>
  </v-hover> -->
  <v-hover>
    <template v-slot:default="{ isHovering, props: hoverProps }">
      <div v-bind="hoverProps" class="d-flex child-flex pa-2">
        <v-scale-transition>
          <div
            v-if="isHovering || props.modelValue"
            :class="{
              'check-button-absolute': props.absolutePosition,
            }"
          >
            <v-icon
              v-if="props.loading"
              icon="mdi-loading"
              class="mr-2 loading"
              color="grey"
              :size="props.selectIconSize"
            />
            <v-icon
              v-else
              icon="mdi-check-circle"
              class="mr-2"
              @click="clickHanlder"
              :color="props.modelValue ? 'primary' : 'grey  '"
              :size="props.selectIconSize"
            />
          </div>
        </v-scale-transition>
        <div
          :class="{
            slide: isHovering || props.modelValue,
          }"
        >
          <slot></slot>
        </div>
      </div>
    </template>
  </v-hover>
</template>

<style scoped>
.check-button-absolute {
  position: absolute;
  z-index: 1;
}

.loading {
  animation: rotation 2s infinite linear;
}

@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(359deg);
  }
}

.check-button {
}
.slide {
  animation: slide 1s cubic cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}
@keyframes slide {
  from {
    transform: translate3d(0, 0, 0);
  }
  to {
    transform: translate3d(20, 0, 0);
  }
}
</style>
