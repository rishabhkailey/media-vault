<script setup lang="ts">
import SelectFileButton from "@/components/SelectFileButton.vue";
import SearchInputField from "@/components/SearchInputField.vue";
import ProfileButton from "@/components/AppBar/NormalAppBar/ProfileButton.vue";
import LogoButton from "@/components/Logo/LogoButton.vue";

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
          <LogoButton :redirect="true" redirect-component-name="Home" />
        </div>
      </v-toolbar-title>
    </v-col>
    <!-- mid -->
    <v-col cols="4">
      <SearchInputField
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
          <SelectFileButton
            :icon-only="false"
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
