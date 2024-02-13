<script setup lang="ts">
const props = defineProps<{
  loading: boolean;
  userName: string;
  email: string;
  authenticated: boolean;
}>();

const emits = defineEmits<{
  (e: "logout"): void;
  (e: "login"): void;
}>();
</script>
<template>
  <v-btn
    v-if="!props.authenticated"
    class="bg-primary mx-2"
    @click.stop="() => emits('login')"
    data-test-id="signin-button"
  >
    <v-icon icon="mdi-login" />
    Sign In
  </v-btn>
  <v-menu v-else>
    <template v-slot:activator="{ props: menuProps }">
      <v-btn
        v-bind="menuProps"
        :loading="props.loading"
        color="primary"
        rounded="pill"
        icon="mdi-account"
      />
    </template>
    <v-card prepend-icon="mdi-account">
      <template v-slot:title>
        {{ userName }}
      </template>
      <template v-slot:subtitle>
        {{ email }}
      </template>
      <template v-slot:actions>
        <div class="d-flex justify-center flex-grow-1">
          <v-btn
            class="bg-primary mx-2"
            :loading="props.loading"
            @click.stop="() => emits('logout')"
          >
            <v-icon icon="mdi-logout" />
            Sign Out
          </v-btn>
        </div>
      </template>
    </v-card>
  </v-menu>
</template>
