import { createRouter, createWebHistory } from 'vue-router'
import Chat from '../pages/Chat.vue'
import Login from '../pages/Login.vue'
import Register from '../pages/Register.vue'
import RoleLibrary from '../pages/RoleLibrary.vue'
import Membership from '../pages/Membership.vue'
import Profile from '../pages/Profile.vue'
import WebSearch from '../pages/WebSearch.vue'
import ASRTTS from '../pages/ASRTTS.vue'

const routes = [
  { path: '/', redirect: '/chat' },
  { path: '/chat', component: Chat },          // 游客可进
  { path: '/role', component: RoleLibrary },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/membership', component: Membership, meta: { requiresAuth: true } },
  { path: '/profile', component: Profile, meta: { requiresAuth: true } },
  { path: '/web-search', component: WebSearch },
  { path: '/asr-tts', component: ASRTTS },
]

const router = createRouter({ history: createWebHistory(), routes })

// 若怀疑这里导致白屏，先全部注释掉再试
router.beforeEach((to, _from, next) => {
  if (to.meta?.requiresAuth && !localStorage.getItem('token')) {
    return next({ path: '/login', query: { from: to.fullPath } })
  }
  next()
})

export default router
