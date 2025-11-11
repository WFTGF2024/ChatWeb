import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// 临时抓隐藏报错（只为排查，之后可删）
app.config.errorHandler = (err)=>console.error('VueError:', err)
window.addEventListener('error', e=>console.error('WindowError:', e.error || e.message))
window.addEventListener('unhandledrejection', e=>console.error('Unhandled:', e.reason))

app.mount('#app')
