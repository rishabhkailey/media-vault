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
      "/v1": "http://127.0.0.1:8090/",
    },
    headers: {
      "Service-Worker-Allowed": "/",
    },
  },
});
