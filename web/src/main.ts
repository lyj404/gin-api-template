import { createApp } from 'vue'
import { addCollection } from '@iconify/vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import 'virtual:uno.css'
import '@styles/main.css'

// Use dynamic import to prevent tree-shaking
const { default: iconsData } = await import('@iconify-json/material-symbols/icons.json')
addCollection(iconsData)

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
