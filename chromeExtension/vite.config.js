import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  build: {
    outDir: 'dist',
    lib: {
      formats: ['es'],
      entry: {
        contentScript: resolve(__dirname, 'contentScript.ts'),
        helper: resolve(__dirname, 'helper.ts'),
        index: resolve(__dirname, 'index.ts'),
        background: resolve(__dirname, 'src/service_worker.js')
      }
    },
    rollupOptions: {
      output: {
        format: 'es',
        entryFileNames: '[name].js'
      }
    }
  }
});