<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  label: string;
  icon: string;
  iconOnly: boolean;
}>();

const emits = defineEmits<{
  (e: "select", value: Array<File>): void;
}>();

const selectFileElement = ref<HTMLInputElement | undefined>(undefined);
const selectFileDialog = ref<boolean>(false);
const selectedFiles = ref<Array<File>>([]);

const cancel = () => {
  selectFileDialog.value = false;
  selectedFiles.value = [];
};
</script>

<template>
  <!-- todo common component for select file and search dialog? -->
  <v-form
    class="d-flex flex-grow-1 justify-center align-center"
    @submit.prevent="(e) => emits('select', selectedFiles)"
  >
    <v-dialog v-model="selectFileDialog" scrollable>
      <template v-slot:activator>
        <v-btn
          v-if="props.iconOnly"
          :icon="props.icon"
          color="primary"
          data-test-id="upload-file-button"
          @click.stop="() => (selectFileDialog = true)"
        />
        <v-btn
          v-else
          class="bg-primary mx-2"
          :prepend-icon="props.icon"
          :text="props.label"
          data-test-id="upload-file-button"
          @click.stop="() => (selectFileDialog = true)"
        />
      </template>

      <v-container
        class="ma-0 pa-0 justify-center align-center d-flex w-100"
        style="max-width: 100vw"
      >
        <v-col cols="12" xs="12" sm="10" md="6" lg="6" xl="6" xxl="6">
          <v-card>
            <v-card-title> Select Files </v-card-title>
            <v-card-text>
              <v-file-input
                ref="selectFileElement"
                placeholder="select files"
                label="select files"
                v-model="selectedFiles"
                data-test-id="upload-file-input"
                counter
                multiple
                clearable
              >
                <template v-slot:selection="{ fileNames }">
                  <template
                    v-for="(fileName, index) in fileNames"
                    :key="fileName"
                  >
                    <v-chip
                      v-if="index < 2"
                      color="deep-purple-accent-4"
                      label
                      size="small"
                      class="me-2"
                    >
                      {{ fileName }}
                    </v-chip>
                    <span
                      v-else-if="index === 2"
                      class="text-overline text-grey-darken-3 mx-2"
                    >
                      +{{ selectedFiles.length - 2 }} File(s)
                    </span>
                  </template>
                </template>
              </v-file-input>
            </v-card-text>
            <v-card-actions>
              <!-- feature -->
              <!-- todo preview button with option to remove some files -->
              <v-row class="d-flex flex-row justify-space-around">
                <!-- move red to theme under delete names -->
                <v-btn color="red" @click.stop="cancel">Cancel</v-btn>
                <v-btn
                  color="primary"
                  data-test-id="upload-files-selection-confirm-button"
                  @click.stop="
                    () => {
                      emits('select', selectedFiles);
                      selectFileDialog = false;
                      selectedFiles = [];
                    }
                  "
                  :disabled="selectedFiles.length === 0"
                  >{{ props.label }}</v-btn
                >
              </v-row>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-container>
    </v-dialog>
  </v-form>
</template>
