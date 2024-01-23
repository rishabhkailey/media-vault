<script setup lang="ts">
import { ref, onMounted, computed, reactive } from "vue";
import { chunkUpload } from "@/js/encryptFileUpload";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { useMediaStore } from "@/piniaStore/media";
import CompactFileUploadBar from "./CompactFileUploadBar.vue";
import FileUploadProgress from "./FileUploadProgress.vue";
const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);
const { usableEncryptionKey } = storeToRefs(useUserInfoStore());

const props = withDefaults(
  defineProps<{
    files: Array<File>;
    height: string;
    width: string;
  }>(),
  {},
);

const emits = defineEmits<{
  (e: "status:failed", value: boolean): void;
  (e: "status:completed", value: boolean): void;
  (e: "close"): void;
  (e: "collapse"): void;
}>();

const { appendMedia } = useMediaStore();
const collapsed = ref<boolean>(false);
const overallProgress = ref(0);

interface FileUploadStatus {
  progress: number;
  done: boolean;
  failed: boolean;
  errMessage: string;
  controller: AbortController;
}
const filesUploadStatus = computed<Array<FileUploadStatus>>(() => {
  let arr = new Array<FileUploadStatus>(props.files.length);
  for (let i = 0; i < arr.length; i++) {
    arr[i] = {
      progress: 0,
      done: false,
      failed: false,
      errMessage: "",
      controller: new AbortController(),
    };
  }
  return reactive<Array<FileUploadStatus>>(arr);
});

async function uploadFiles(files: Array<File>) {
  let totalSize = 0;
  let uploaded = 0;
  files.forEach((f) => totalSize + f.size);
  // forEach requres async function using which we can not upload synchronously
  for (let index = 0; index < files.length; index++) {
    try {
      if (filesUploadStatus.value[index].controller.signal.aborted) {
        throw new Error("canceled");
      }
      let file = files[index];
      let media = await chunkUpload(
        file,
        accessToken.value,
        usableEncryptionKey.value,
        filesUploadStatus.value[index].controller,
        (progress: number) => {
          filesUploadStatus.value[index].progress = progress;
          uploaded += (progress * file.size) / 100;
          overallProgress.value = (uploaded * 100) / totalSize;
        },
      );
      appendMedia([media]);
      filesUploadStatus.value[index].done = true;
      filesUploadStatus.value[index].failed = false;
    } catch (err) {
      filesUploadStatus.value[index].done = true;
      filesUploadStatus.value[index].failed = true;
      if (err instanceof Error) {
        filesUploadStatus.value[index].errMessage = err.toString();
      }
    }
  }
}

onMounted(() => {
  if (
    props.files.length > 0 &&
    props.files.length === filesUploadStatus.value.length
  ) {
    console.log("uploading...");
    uploadFiles(props.files);
  }
});
</script>

<template>
  <v-card
    :style="[
      'overflow-x: hidden; overflow-y: hidden',
      `height: ${props.height}`,
      `width: ${props.width}`,
    ]"
    v-if="!collapsed"
  >
    <v-col class="h-100 d-flex flex-column">
      <v-row class="flex-grow-0">
        <v-toolbar class="d-flex flex-column ma-0 pa-0" color="primary">
          <v-toolbar-title class="text-center ma-0 pa-0">
            Uploading
          </v-toolbar-title>
        </v-toolbar>
      </v-row>
      <v-row
        class="flex-grow-1"
        style="display: block; overflow-y: scroll; overflow-x: hidden"
      >
        <v-list lines="two">
          <FileUploadProgress
            v-for="(file, index) in props.files"
            :key="file.name"
            @cancel="
              () => {
                filesUploadStatus[index].controller.abort();
              }
            "
            :failed="filesUploadStatus[index].failed"
            :progress="filesUploadStatus[index].progress"
            :completed="filesUploadStatus[index].done"
            :name="file.name"
            :size="file.size"
          />
        </v-list>
      </v-row>
      <v-row class="flex-grow-0">
        <v-card-actions>
          <v-btn
            @click.stop="
              () => {
                filesUploadStatus.forEach((status) => {
                  status.controller.abort();
                });
                emits('close');
              }
            "
          >
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
        </v-card-actions>
      </v-row>
    </v-col>
  </v-card>
  <CompactFileUploadBar
    class="justify-end algin-end"
    v-else
    v-model="collapsed"
    :indeterminate="false"
    :progress="overallProgress"
  />
</template>
@/js/encryptFileUpload
