<script setup lang="ts"></script>

<template>
  <v-navigation-drawer permanent :rounded="false">
    <v-list nav>
      <v-list-item prepend-icon="mdi-email" title="Inbox" value="inbox"></v-list-item>
      <v-list-item
        prepend-icon="mdi-account-supervisor-circle"
        title="Supervisors"
        value="supervisors"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-clock-start"
        title="Clock-in"
        value="clockin"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>
