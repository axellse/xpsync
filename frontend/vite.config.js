import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        error: resolve(__dirname, 'error.html'),
        transfer: resolve(__dirname, 'transfer.html'),
      }
    }
  },
  base: './'
})