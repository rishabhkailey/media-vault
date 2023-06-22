<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useDisplay } from "vuetify/lib/framework.mjs";

const props = defineProps<{
  modelValue: boolean;
  absolutePosition: boolean;
  selectIconSize: string | number;
  loading: boolean;
  showSelectButtonOnHover: boolean;
  alwaysShowSelectButton: boolean;
  alwaysShowSelectOnMobile: boolean;
  selectOnContentClick: boolean;
}>();

const { mobile: mobileDevice } = useDisplay();
const allwaysShowSelectButton = computed<boolean>(() => {
  return (
    props.alwaysShowSelectButton ||
    (mobileDevice.value && props.alwaysShowSelectOnMobile)
  );
});

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "change", value: boolean): void;
}>();

const clickHandler = () => {
  // selected.value = !selected.value;
  emit("update:modelValue", !props.modelValue);
  emit("change", !props.modelValue);
};

const contentClickHandler = (e: MouseEvent) => {
  e.stopPropagation();
  clickHandler();
};

const contentWrapper = ref<HTMLElement | undefined>(undefined);

watch(
  () => props.selectOnContentClick,
  (newValue, oldValue) => {
    if (contentWrapper.value === undefined) {
      return;
    }
    if (newValue) {
      contentWrapper.value.addEventListener("click", contentClickHandler, true);
      return;
    }
    contentWrapper.value.removeEventListener(
      "click",
      contentClickHandler,
      true
    );
  }
);
// const selected = ref(false);
</script>
<template>
  <v-container class="pa-0 ma-0">
    <v-hover>
      <template v-slot:default="{ isHovering, props: hoverProps }">
        <div v-bind="hoverProps" class="d-flex child-flex pa-2">
          <v-scale-transition>
            <div
              v-if="
                allwaysShowSelectButton ||
                isHovering ||
                props.modelValue ||
                props.loading
              "
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
                @click.stop="clickHandler"
                :color="props.modelValue ? 'primary' : 'grey  '"
                :size="props.selectIconSize"
              />
            </div>
          </v-scale-transition>
          <div
            :class="{
              'pointer-cursor': props.selectOnContentClick,
            }"
            ref="contentWrapper"
          >
            <!-- @click.stop.capture="props.selectOnContentClick ? contentClickHandler : null" -->
            <slot></slot>
          </div>
        </div>
      </template>
    </v-hover>
  </v-container>
</template>

<style scoped>
.pointer-cursor {
  cursor: pointer;
}
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
