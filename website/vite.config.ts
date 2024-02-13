import { fileURLToPath, URL } from "node:url";
// import mkcert from "vite-plugin-mkcert";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "/v1": "http://localhost:8090/",
      "/realms": "http://localhost:8081/",
    },
    headers: {
      "Service-Worker-Allowed": "/",
    },
    watch: {
      usePolling: true,
    },
  },
});
