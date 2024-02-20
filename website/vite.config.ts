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
      "/ui": {
        target: "http://localhost:8082",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/ui/, ""),
      },
      "/v1": "http://localhost:8090/",
      "/realms": "http://localhost:8081/",
      "/accounts": {
        target: "http://localhost:8081",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/accounts/, ""),
      },
    },
    headers: {
      "Service-Worker-Allowed": "/",
    },
    watch: {
      usePolling: true,
    },
  },
});
