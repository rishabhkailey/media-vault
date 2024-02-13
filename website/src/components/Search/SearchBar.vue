<script setup lang="ts">
import { ref, watch } from "vue";

const props = withDefaults(
  defineProps<{
    collapsed: boolean;
    modelValue: string;
  }>(),
  { collapsed: false, modelValue: "" },
);
const emits = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "submit", value: string): void;
}>();
const searchInputRules: Array<any> = [];
const searchDialog = ref(false);
const searchElement = ref<HTMLElement | undefined>(undefined);

watch(searchElement, (newValue) => {
  if (newValue === undefined) {
    return;
  }
  // https://github.com/vuetifyjs/vuetify/issues/10659#issuecomment-594329553
  setTimeout(() => {
    console.log("focused ", searchElement.value?.focus);
    searchElement.value?.focus();
  }, 100);
});
</script>

<template>
  <!-- normal search bar -->
  <v-form
    v-if="!props.collapsed"
    class="d-flex flex-grow-1"
    @submit.prevent="
      (e) => {
        console.log('form submit');
        emits('submit', props.modelValue);
      }
    "
  >
    <v-text-field
      :clearable="true"
      clear-icon="mdi-close"
      :model-value="props.modelValue"
      @update:model-value="
        (value) => {
          emits('update:modelValue', value);
        }
      "
      :rules="searchInputRules"
      label="search"
      :hide-details="true"
    >
      <template #append-inner>
        <v-icon
          icon="mdi-magnify"
          @click="
            () => {
              console.log('icon submit');
              emits('submit', props.modelValue);
            }
          "
        />
      </template>
    </v-text-field>
  </v-form>

  <!-- mobile search button -->
  <v-dialog v-else v-model="searchDialog" location="top">
    <template v-slot:activator="{ props }">
      <v-btn color="primary" v-bind="props" icon="mdi-magnify"> </v-btn>
    </template>
    <v-container
      class="ma-0 pa-0 justify-center align-center d-flex w-100"
      style="max-width: 100vw"
    >
      <v-col cols="12" xs="12" sm="10" md="6" lg="6" xl="6" xxl="6">
        <v-card>
          <v-card-text>
            <v-form
              class="d-flex flex-grow-1"
              @submit.prevent="
                (e) => {
                  console.log('form submit');
                  emits('submit', props.modelValue);
                  searchDialog = false;
                }
              "
            >
              <v-text-field
                :clearable="true"
                clear-icon="mdi-close"
                :model-value="props.modelValue"
                @update:model-value="
                  (value) => {
                    emits('update:modelValue', value);
                  }
                "
                :rules="searchInputRules"
                label="search"
                ref="searchElement"
                focused
                hide-details
              >
                <template #append-inner>
                  <v-icon
                    icon="mdi-magnify"
                    @click="
                      (e) => {
                        console.log('collapsed icon submit');
                        emits('submit', props.modelValue);
                        searchDialog = false;
                      }
                    "
                  />
                </template>
              </v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" block @click="searchDialog = false">
              Close Dialog
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-container>
  </v-dialog>
</template>
