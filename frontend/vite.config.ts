import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    modulePreload: { polyfill: false }
  },
  server: {
    allowedHosts: ["splitzy.dak-shin.com"]
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@utils": path.resolve(__dirname, "./src/common"),
      "@comp": path.resolve(__dirname, "./src/components"),
      "@ui": path.resolve(__dirname, "./src/components/ui"),
      "@common": path.resolve(__dirname, "./src/components/common"),
      "@split-form": path.resolve(__dirname, "./src/components/split-form"),
    },
  }
})
