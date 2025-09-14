import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/sets': 'http://localhost:8080',
      '/collections': 'http://localhost:8080',
      '/sightings': 'http://localhost:8080'
    }
  }
})
