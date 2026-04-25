import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte(), wails("./bindings")],
  resolve: {
    alias: {
      $lib: "/src/lib",
      $components: "/src/components",
      $views: "/src/views",
      $bindings: "/bindings",
    },
  },
});
