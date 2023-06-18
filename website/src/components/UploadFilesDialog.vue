<script setup lang="ts">
import { ref, onMounted, computed, reactive } from "vue";
import { chunkUpload } from "@/utils/encryptFileUpload";
import { useAuthStore } from "@/piniaStore/auth";
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/piniaStore/userInfo";
import { useMediaStore } from "@/piniaStore/media";
const authStore = useAuthStore();
// todo store getter typing
// todo check
const { accessToken } = storeToRefs(authStore);
const { usableEncryptionKey } = storeToRefs(useUserInfoStore());

const { appendMedia } = useMediaStore();
const props = withDefaults(
  defineProps<{
    files: Array<File>;
  }>(),
  {}
);

interface FileUploadStatus {
  progress: number;
  done: boolean;
  failed: boolean;
  errMessage: string;
  controller: AbortController;
}
const filesUploadStatus = computed<Array<FileUploadStatus>>(() => {
  let arr = new Array<FileUploadStatus>(props.files.length);
  // arr.fill sets reference of same object to all elements and updating 1 update all
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

// todo move this to somewhere else on any update it reuploads the files again
onMounted(() => {
  if (
    props.files.length > 0 &&
    props.files.length === filesUploadStatus.value.length
  ) {
    console.log("uploading...");
    uploadFiles(props.files);
  }
});

async function uploadFiles(files: Array<File>) {
  // forEach requres async function using which we can not upload synchronously
  for (let index = 0; index < files.length; index++) {
    try {
      if (filesUploadStatus.value[index].controller.signal.aborted) {
        throw new Error("canceled");
      }
      let file = files[index];
      // not waiting for upload to finish
      let media = await chunkUpload(
        file,
        accessToken.value,
        usableEncryptionKey.value,
        filesUploadStatus.value[index].controller,
        (progress: number) => {
          filesUploadStatus.value[index].progress = progress;
        }
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
      console.log(err);
    }
  }
}
const collapsed = ref<boolean>(false);

const emits = defineEmits<{
  (e: "status:failed", value: boolean): void;
  (e: "status:completed", value: boolean): void;
  (e: "close"): void;
  (e: "collapse"): void;
}>();
</script>

<template>
  <v-card style="overflow-x: hidden; overflow-y: hidden" v-if="!collapsed">
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
          <v-list-item
            v-for="(file, index) in props.files"
            :key="file.name"
            :title="file.name"
            :subtitle="file.size"
          >
            <template v-slot:prepend>
              <!-- todo create a separate component with simple logic -->
              <v-avatar>
                <v-progress-circular
                  :size="70"
                  :width="7"
                  :color="filesUploadStatus[index].failed ? 'red' : 'primary'"
                  :model-value="
                    filesUploadStatus[index].failed
                      ? '100'
                      : filesUploadStatus[index].progress
                  "
                >
                  <!-- todo text size -->
                  {{
                    filesUploadStatus[index].failed
                      ? "!"
                      : Math.round(filesUploadStatus[index].progress) + "%"
                  }}
                </v-progress-circular>
              </v-avatar>
            </template>
            <template v-slot:append>
              <v-btn
                v-if="!filesUploadStatus[index].done"
                color="grey-lighten-1"
                icon="mdi-close"
                variant="text"
                @click.stop="
                  () => {
                    filesUploadStatus[index].controller.abort();
                  }
                "
              />
              <v-btn
                v-if="
                  filesUploadStatus[index].done &&
                  !filesUploadStatus[index].failed
                "
                color="grey-lighten-1"
                icon="mdi-check"
                variant="text"
              />
              <!-- rotate-right -->
              <v-tooltip
                location="top"
                text="retry"
                v-if="
                  filesUploadStatus[index].done &&
                  filesUploadStatus[index].failed
                "
              >
                <template v-slot:activator="{ props }">
                  <v-btn
                    color="grey-lighten-1"
                    icon="mdi-rotate-right"
                    variant="text"
                    v-bind="props"
                  />
                </template>
              </v-tooltip>
            </template>
          </v-list-item>
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
</template>
