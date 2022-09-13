import { defineNuxtConfig } from 'nuxt'

export default defineNuxtConfig({
  buildModules: ['@nuxtjs/supabase', '@nuxthq/ui', '@nuxtjs/color-mode'],
  build: {
    transpile: ['charts.js', 'vue-chart-3'],
  },
  ssr: false,
  ui: {
    colors: {
      primary: 'green',
    },
  },
  server: {
    port: 4000,
  },
})
