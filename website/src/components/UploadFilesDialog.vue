<script setup lang="ts">
import { ref, defineProps, withDefaults, defineEmits } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    width?: number;
    height?: number;
    margin?: number;
    files: Array<File>;
  }>(),
  {
    width: 250,
    height: 250,
    margin: 25,
  }
);

// const attachedTo = document.createElement("div");
// attachedTo.style.position = "absolute";
// attachedTo.style.bottom = "0";
// attachedTo.style.right = "0";

const attachedTo = ref<HTMLDivElement | undefined>(undefined);
const collapsed = ref<boolean>(false);

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "status:failed", value: boolean): void;
  (e: "status:completed", value: boolean): void;
}>();
</script>

<template>
  <div
    :style="`position: absolute; bottom: ${
      props.height + props.margin
    }px; right: ${props.width + props.margin}px`"
    ref="attachedTo"
  ></div>
  <!-- https://vuejs.org/guide/components/v-model.html#component-v-model -->
  <v-menu
    :model-value="props.modelValue"
    @update:model-value="
      () => {
        /* ignore all model updates we will manually update when to close/open menu */
      }
    "
    open-delay="0"
    close-delay="0"
    :width="props.width"
    :height="props.height"
    location-strategy="static"
    :scrollable="false"
    :location="`${collapsed ? 'top left' : 'bottom left'}`"
    :attach="attachedTo"
    :close-on-content-click="false"
    :close-on-back="false"
    style="position: absolute; bottom: 0; right: 0"
  >
    <v-card style="overflow-x: hidden; overflow-y: hidden" v-if="!collapsed">
      <v-toolbar
        class="d-flex flex-column ma-0 pa-0"
        :style="`max-height: ${(props.height * 1.5) / 10}px;`"
        color="primary"
      >
        <v-toolbar-title class="text-center ma-0 pa-0">
          Uploading
        </v-toolbar-title>
      </v-toolbar>

      <v-card-item
        class="flex-grow-1 ma-0 pa-0"
        :style="`display: block; overflow-y: scroll; overflow-x: hidden; max-height: ${
          (props.height * 7) / 10
        }px;`"
      >
        <v-container class="ma-0 pa-0">
          <v-list lines="two">
            <v-list-item
              v-for="file in files"
              :key="file.name"
              :title="file.name"
              :subtitle="file.size"
            >
              <template v-slot:prepend>
                <v-avatar>
                  <v-progress-circular
                    :size="70"
                    :width="7"
                    color="primary"
                    indeterminate
                  >
                    <!-- todo remove this and add logic for progress -->
                    <template v-slot:default>
                      <v-icon color="grey">mdi-folder</v-icon>
                    </template>
                  </v-progress-circular>
                </v-avatar>
              </template>
              <template v-slot:append>
                <v-btn
                  color="grey-lighten-1"
                  icon="mdi-close"
                  variant="text"
                ></v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-container>
      </v-card-item>
      <v-bottom-navigation
        :style="`max-height: ${(props.height * 1.5) / 10}px`"
      >
        <v-btn>
          <v-icon>mdi-close</v-icon>
          cancel
        </v-btn>
        <v-btn
          @click.stop="
            () => {
              collapsed = true;
            }
          "
        >
          <v-icon>mdi-arrow-collapse </v-icon>
          collapse
        </v-btn>
      </v-bottom-navigation>
    </v-card>
    <v-container
      class="ma-0 pa-0 d-flex justify-end align-end flex-grow-1"
      v-else
    >
      <v-btn
        flat
        rounded="pill"
        @click.stop="
          () => {
            collapsed = false;
          }
        "
      >
        <v-avatar>
          <v-progress-circular
            :size="70"
            :width="7"
            color="primary"
            indeterminate
          >
            <!-- todo remove this and add logic for progress -->
            <template v-slot:default>
              <v-icon color="grey">mdi-arrow-expand</v-icon>
            </template>
          </v-progress-circular>
        </v-avatar>
      </v-btn>
    </v-container>
  </v-menu>
</template>
