<script setup lang="ts">
import { HOME_ROUTE_NAME } from "@/router/routesConstants";
import { useRouter } from "vue-router";

const props = withDefaults(
  defineProps<{
    redirect?: boolean;
    redirectComponentName?: string;
  }>(),
  {
    redirect: false,
    redirectComponentName: HOME_ROUTE_NAME,
  },
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
    <span
      style="
        font-size: 20px;
        text-transform: capitalize;
        text-shadow: 7px 7px 2px #0f0f0f;
      "
    >
      <span>Media </span><span>Vault</span>
    </span>
  </v-btn>
</template>
