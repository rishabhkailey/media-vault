<script setup lang="ts">
import { download } from "@/js/encryptedFileDownload";
import { useErrorsStore } from "@/piniaStore/errors";

const props = defineProps<{
  media: Media;
}>();

const { appendError } = useErrorsStore();
function downloadMedia(media: Media) {
  download(media.url, media.name).catch((err) => {
    let errorMessage = "";
    if (typeof err === "string") {
      errorMessage = err;
    }
    if (err instanceof Error) {
      errorMessage = err.message + " " + err.stack;
    }
    appendError(`Download failed ${media.name}`, errorMessage, -1);
  });
}
</script>

<template>
  <v-card
    prepend-icon="mdi-alert"
    title="file type not supported"
    :subtitle="props.media.name"
    class="pa-15"
  >
    <template #text>
      This file type is not currently supported. You can try downloading the
      file and opening it in a different application
    </template>
    <v-card-actions>
      <v-btn
        prepend-icon="mdi-download"
        class="bg-primary mx-2"
        @click.stop="() => downloadMedia(props.media)"
        >Download</v-btn
      >
    </v-card-actions>
  </v-card>
</template>

<style scoped>
.container {
  max-width: 10em;
  max-height: 10em;
}
</style>
