<script setup lang="ts">
import { useRouter } from "vue-router";

const props = withDefaults(
  defineProps<{
    redirect: boolean;
    redirectComponentName: string;
  }>(),
  {
    redirect: false,
    redirectComponentName: "Home",
  }
);

const emits = defineEmits<{
  (e: "click"): void;
}>();

function clickHandler() {
  if (props.redirect) {
    router.push({
      name: props.redirectComponentName,
    });
  }
  emits("click");
}

const router = useRouter();
</script>
<template>
  <v-btn @click.stop="clickHandler">
    <img src="/src/assets/logo.svg" style="height: 30px" /> placeholder
  </v-btn>
</template>
