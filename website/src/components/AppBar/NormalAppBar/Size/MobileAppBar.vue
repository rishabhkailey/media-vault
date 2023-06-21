<script setup lang="ts">
import SelectFileButton from "@/components/SelectFileButton.vue";
import SearchInputField from "@/components/SearchInputField.vue";
import ProfileButton from "@/components/AppBar/NormalAppBar/ProfileButton.vue";

const props = defineProps<{
  navigationBar: boolean;
  searchQuery: string;
  authenticated: boolean;
  userAuthLoading: boolean;
  userName: string;
  email: string;
}>();

const emits = defineEmits<{
  (e: "update:navigationBar", value: boolean): void;
  (e: "update:searchQuery", value: string): void;
  (e: "login"): void;
  (e: "logout"): void;
  (e: "searchSubmit", query: string): void;
  (e: "uploadFiles", selectedFiles: Array<File>): void;
}>();
</script>

<template>
  <v-row class="d-flex align-center ml-2 mr-3">
    <!-- start -->
    <v-col
      cols="2"
      @click.stop="
        () => {
          emits('update:navigationBar', !props.navigationBar);
        }
      "
    >
      <v-btn icon="mdi-menu"> </v-btn>
    </v-col>
    <!-- end -->
    <v-col cols="10">
      <v-row class="d-flex flex-row flex-nowrap justify-end align-center">
        <div>
          <SearchInputField
            :model-value="props.searchQuery"
            @update:model-value="(value) => emits('update:searchQuery', value)"
            :collapsed="true"
            @submit="(query) => emits('searchSubmit', query)"
          />
        </div>
        <div>
          <SelectFileButton
            :icon-only="true"
            label="upload"
            icon="mdi-upload"
            @select="(files) => emits('uploadFiles', files)"
          />
        </div>
        <ProfileButton
          :authenticated="authenticated"
          :loading="props.userAuthLoading"
          :user-name="userName"
          :email="email"
          @logout="() => emits('logout')"
          @login="emits('login')"
        />
      </v-row>
    </v-col>
  </v-row>
</template>
