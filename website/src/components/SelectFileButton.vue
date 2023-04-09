<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  label: string;
  prependIcon: string;
}>();

const emit = defineEmits<{
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
  <v-dialog
    v-model="selectFileDialog"
    width="auto"
    min-width="500px"
    scrollable
  >
    <template v-slot:activator>
      <v-btn
        class="bg-primary mx-2"
        @click.stop="() => (selectFileDialog = true)"
      >
        <v-icon :icon="props.prependIcon" />
        {{ props.label }}
      </v-btn>
    </template>
    <v-card>
      <v-card-title> Select Files </v-card-title>
      <v-card-text>
        <v-row>
          <v-file-input
            ref="selectFileElement"
            placeholder="select files"
            label="select files"
            v-model="selectedFiles"
            counter
            multiple
            clearable
          >
            <template v-slot:selection="{ fileNames }">
              <template v-for="(fileName, index) in fileNames" :key="fileName">
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
        </v-row>
      </v-card-text>
      <v-card-actions>
        <!-- todo preview button with option to remove some files -->
        <v-row class="d-flex flex-row justify-space-around">
          <!-- move red to theme under delete names -->
          <v-btn color="red" @click.stop="cancel">Cancel</v-btn>
          <v-btn
            color="primary"
            @click.stop="
              () => {
                emit('select', selectedFiles);
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
  </v-dialog>
</template>
