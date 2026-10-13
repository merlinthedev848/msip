import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('sip.js')) {
              return 'vendor-sip'
            }
            if (id.includes('svelte')) {
              return 'vendor-svelte'
            }
            return 'vendor'
          }
        }
      }
    }
  },
  server: {
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        ws: true
      },
      '/ws': {
        target: 'ws://localhost:8081',
        ws: true
      }
    }
  }
})

