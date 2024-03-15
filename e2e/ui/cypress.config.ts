import { defineConfig } from 'cypress'

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:8181',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
  env: {
    // warning: default user/password, do not use following credentials these are visible to everyone
    "TEST_USER_NAME": "test-user",
    "TEST_USER_PASSWORD": "ZmuBYNurSS",
    "TEST_USER_ENCRYPTION_KEY": "ZmuBYNurSS"
  }
})

