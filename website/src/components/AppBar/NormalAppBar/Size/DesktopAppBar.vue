<script setup lang="ts">
import FilePickerButton from "@/components/FileUpload/FilePickerButton.vue";
import SearchBar from "@/components/Search/SearchBar.vue";
import AccountButton from "@/components/AppBar/NormalAppBar/AccountButton.vue";
import AppLogoButton from "@/components/Logo/AppLogoButton.vue";

const props = defineProps<{
  searchQuery: string;
  authenticated: boolean;
  userAuthLoading: boolean;
  userName: string;
  email: string;
}>();

const emits = defineEmits<{
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
      col="4"
      class="d-flex flex-row justify-start align-stretch pa-0 ma-0"
    >
      <v-toolbar-title>
        <div>
          <AppLogoButton :redirect="true" redirect-component-name="Home" />
        </div>
      </v-toolbar-title>
    </v-col>
    <!-- mid -->
    <v-col cols="4">
      <SearchBar
        :model-value="props.searchQuery"
        @update:model-value="(value) => emits('update:searchQuery', value)"
        :collapsed="false"
        @submit="(query) => emits('searchSubmit', query)"
      />
    </v-col>
    <!-- end -->
    <v-col cols="4">
      <v-row class="d-flex flex-row flex-nowrap justify-end align-center">
        <div>
          <FilePickerButton
            :icon-only="false"
            label="upload"
            icon="mdi-upload"
            @select="(files) => emits('uploadFiles', files)"
          />
        </div>
        <AccountButton
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
