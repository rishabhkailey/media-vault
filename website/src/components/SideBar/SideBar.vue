<script setup lang="ts">
import { computed } from "vue";
import { useDisplay } from "vuetify";
import AlbumsVerticalList from "@/components/Album/AlbumsVerticalList.vue";
import { homeRoute } from "@/router/routesConstants";
const display = useDisplay();

const props = defineProps<{
  modelValue: boolean;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();
const smallDisplay = computed(
  () => display.mobile.value || display.smAndDown.value,
);
</script>

<template>
  <v-navigation-drawer
    data-test-id="side-bar"
    :model-value="props.modelValue"
    @update:model-value="
      (newValue) => {
        emit('update:modelValue', newValue);
      }
    "
    :permanent="!smallDisplay"
    :temporary="smallDisplay"
    :rounded="false"
    elevation="2"
  >
    <v-list nav>
      <v-list-item
        prepend-icon="mdi-home"
        title="Home"
        value="Home"
        :to="homeRoute()"
        :exact="true"
        color="primary"
        data-test-id="sidebar-home-group"
      ></v-list-item>
      <AlbumsVerticalList />
    </v-list>
  </v-navigation-drawer>
</template>
