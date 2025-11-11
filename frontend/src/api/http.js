// src/api/http.js
import axios from 'axios'
import { getCookie } from '../utils/cookies'

// ---------------- 基础实例（只用 JWT 头，不走凭证） ----------------
export const httpCore = axios.create({
  baseURL: import.meta.env.VITE_CORE_BASE, // dev 建议 /core（Vite 代理），prod 也是 /core（Nginx 反代）
  timeout: 20000,
  withCredentials: false, // 只用 Authorization 头
})

export const httpLLM = axios.create({
  baseURL: import.meta.env.VITE_LLM_BASE,
  timeout: 120000,
  withCredentials: false,
})

export const httpASR = axios.create({
  baseURL: import.meta.env.VITE_ASR_BASE,
  timeout: 120000,
  withCredentials: false,
})

export const httpTTS = axios.create({
  baseURL: import.meta.env.VITE_TTS_BASE,
  timeout: 120000,
  withCredentials: false,
})

// ---------------- 统一加 Authorization 头（Core） ----------------
httpCore.interceptors.request.use((cfg) => {
  let token =
    localStorage.getItem('token') ||
    getCookie('jwt') ||
    getCookie('token') ||
    ''

  if (token) {
    token = token.startsWith('Bearer ') ? token.slice(7) : token
    cfg.headers.Authorization = `Bearer ${token}`
  }
  return cfg
})

// ---------------- 401 统一处理（Core） ----------------
let redirecting401 = false
httpCore.interceptors.response.use(
  (r) => r,
  async (err) => {
    if (err?.response?.status === 401 && !redirecting401) {
      redirecting401 = true
      try {
        const { default: router } = await import('@/router')
        const from = router.currentRoute.value.fullPath
        localStorage.removeItem('token')
        router.replace({ path: '/login', query: { from } })
      } finally {
        setTimeout(() => (redirecting401 = false), 300)
      }
    }
    return Promise.reject(err)
  }
)

// ---------------- Web 搜索相关 API（供 WebSearch.vue 使用） ----------------
export const webAPI = {
  // 搜索：支持 { q, mode, top_k } 或 { urls, top_k }
  search: (body) => httpCore.post('/web/search', body),

  // 列表：已抓取的条目
  items: () => httpCore.get('/web/items'),

  // 新建/抓取页面（可选）
  create: (payload) => httpCore.post('/web/items', payload),

  // 查看页面全文（WebSearch.vue 的 openPage 用到）
  page: (id) => httpCore.get(`/web/page/${id}`),

  // 更新/删除（可选）
  update: (id, payload) => httpCore.put(`/web/items/${id}`, payload),
  remove: (id) => httpCore.delete(`/web/items/${id}`),

  // 触发入库/分块（可选）
  ingest: () => httpCore.post('/web/ingest'),
  chunk: () => httpCore.post('/web/chunk'),
}

// ---------------- 工具：逐行读取 NDJSON ----------------
export async function readNDJsonStream(response, onChunk) {
  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buf = ''
  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    const parts = buf.split(/\r?\n/)
    buf = parts.pop() || ''
    for (const line of parts) {
      if (!line.trim()) continue
      try { onChunk(JSON.parse(line)) } catch {}
    }
  }
  if (buf.trim()) { try { onChunk(JSON.parse(buf.trim())) } catch {} }
}
