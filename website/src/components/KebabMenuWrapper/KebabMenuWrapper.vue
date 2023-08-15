<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  selectIconSize: string | number;
  showSelectButtonOnHover: boolean;
  alwaysShowSelectButton: boolean;
  selectOnContentClick: boolean;
}>();

const clicked = ref(false);
const clickHandler = () => {
  clicked.value = true;
};

function onClickOutside() {
  console.log("clicked outside");
  clicked.value = false;
}

const contentWrapper = ref<HTMLElement | undefined>(undefined);
// const selected = ref(false);
</script>
<template>
  <div>
    <v-hover>
      <template v-slot:default="{ isHovering, props: hoverProps }">
        <div v-bind="hoverProps" class="d-flex child-flex pa-2">
          <!-- <v-scale-transition> -->
          <div
            v-if="props.alwaysShowSelectButton || isHovering || clicked"
            class="check-button-absolute"
          >
            <v-card class="pa-0 ma-0" v-click-outside="onClickOutside">
              <v-menu>
                <template v-slot:activator="{ props }">
                  <v-icon
                    icon="mdi-dots-vertical"
                    class="mr-2"
                    @click.stop="clickHandler"
                    color="grey"
                    v-bind="props"
                    :size="props.selectIconSize"
                  />
                </template>
                <slot name="options"></slot>
              </v-menu>
            </v-card>
          </div>
          <!-- </v-scale-transition> -->
          <div
            :class="{
              'pointer-cursor': props.selectOnContentClick,
              'w-100': true,
            }"
            ref="contentWrapper"
          >
            <!-- @click.stop.capture="props.selectOnContentClick ? contentClickHandler : null" -->
            <slot></slot>
          </div>
        </div>
      </template>
    </v-hover>
  </div>
</template>

<style scoped>
.pointer-cursor {
  cursor: pointer;
}
.check-button-absolute {
  position: absolute;
  right: 0;
  z-index: 1;
}
</style>
