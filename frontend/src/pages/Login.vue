<template>
  <div class="login-wrap">
    <div class="card">
      <h2>登录</h2>

      <p v-if="err" class="err">{{ err }}</p>

      <div class="form">
        <label>用户名</label>
        <input v-model="username" class="input" placeholder="用户名" />
        <label>密码</label>
        <input v-model="password" class="input" type="password" placeholder="密码" />

        <button class="btn primary" :disabled="loading" @click="submit">
          {{ loading ? '登录中…' : '登录' }}
        </button>
      </div>

      <div class="hint">
        没有账号？
        <router-link :to="`/register?redirect=${encodeURIComponent(redirectTarget)}`">去注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../store/user'
import { login, me } from '../api/core'   // ← 继续沿用你已有的 /api/core 接口

const router = useRouter()
const route  = useRoute()
const user   = useUserStore()

const username = ref('')
const password = ref('')
const loading  = ref(false)
const err      = ref('')

// 登录成功后回跳的目标：优先用 redirect，其次默认回到 /chat
const redirectTarget = computed(() => {
  const r = route.query.redirect
  try {
    if (typeof r === 'string' && r.length) return decodeURIComponent(r)
  } catch {}
  return '/chat'
})

async function submit () {
  err.value = ''
  loading.value = true
  try {
    // 1) 调后端登录
    const res = await login({ username: username.value, password: password.value })
    // 兼容不同返回结构：{token} 或 {data:{token}}
    const token = res?.token || res?.data?.token
    if (!token) throw new Error('登录成功但未返回 token')

    // 2) 落地 token（供 http.js 拦截器读取）
    try { localStorage.setItem('token', token) } catch {}

    // 3) 更新本地用户态（若你有 user.setAuth）
    if (typeof user?.setAuth === 'function') {
      try {
        const info = await me()
        user.setAuth(token, info)
      } catch {
        user.setAuth(token, null)
      }
    }

    // 4) 回跳来源路径或默认 /chat
    router.replace(redirectTarget.value)
  } catch (e) {
    err.value = e?.response?.data?.message || e?.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-wrap{ min-height: 60vh; display: grid; place-items: center; padding: 24px; }
.card{ width: min(420px, 92vw); background: #fff; padding: 20px 24px; border-radius: 12px; box-shadow: 0 6px 24px rgba(0,0,0,.08); }
.form{ display: grid; gap: 10px; margin-top: 12px; }
.input{ border:1px solid #e5e7eb; border-radius:8px; padding:10px 12px; outline:none; }
.input:focus{ border-color:#2563eb; box-shadow: 0 0 0 3px rgba(37,99,235,.15); }
.btn{ padding: 10px 14px; border-radius:8px; border:none; cursor:pointer; }
.btn.primary{ background:#2563eb; color:#fff; }
.err{ color:#ef4444; margin: 6px 0 0 0; }
.hint{ margin-top: 10px; color:#666; }
</style>
