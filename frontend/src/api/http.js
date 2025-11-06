import axios from 'axios'

// ====== Base URLs ======
const LLM_BASE  = import.meta.env.VITE_LLM_BASE  || 'http://127.0.0.1:7207'
const ASR_BASE  = import.meta.env.VITE_ASR_BASE  || 'http://127.0.0.1:7205'
const TTS_BASE  = import.meta.env.VITE_TTS_BASE  || 'http://127.0.0.1:7206'
const CORE_BASE = import.meta.env.VITE_CORE_BASE || 'http://127.0.0.1:7210'

// ====== Axios Instances ======
// Core：你的主后端（/api/*、/web/*）
export const httpCore = axios.create({ baseURL: CORE_BASE, timeout: 20000 })
// 其他独立服务
export const httpLLM  = axios.create({ baseURL: LLM_BASE,  timeout: 60000 })
export const httpASR  = axios.create({ baseURL: ASR_BASE,  timeout: 60000 })
export const httpTTS  = axios.create({ baseURL: TTS_BASE,  timeout: 60000 })
// Web Search 可单独实例（也可以直接用 httpCore）
export const httpWEB  = axios.create({ baseURL: CORE_BASE, timeout: 30000 })

// 默认 JSON 头（ASR/TTS 场景按需在各自 API 里传 multipart/form-data）
httpCore.defaults.headers.post['Content-Type'] = 'application/json'
httpCore.defaults.headers.put['Content-Type']  = 'application/json'
httpWEB .defaults.headers.post['Content-Type'] = 'application/json'
httpWEB .defaults.headers.put['Content-Type']  = 'application/json'

// ====== Token 注入 ======
function getToken () {
  try { return localStorage.getItem('token') || '' } catch { return '' }
}

;[httpCore, httpWEB, httpLLM, httpASR, httpTTS].forEach(inst => {
  inst.interceptors.request.use(cfg => {
    const t = getToken()
    if (t) cfg.headers.Authorization = `Bearer ${t}`
    return cfg
  })
})

// ====== 统一 401 处理（Core/Web 生效；LLM/ASR/TTS 不干预）======
function onUnauthorized () {
  try { localStorage.removeItem('token') } catch {}
  // 已在登录页就不重定向，避免死循环
  const here = location.pathname + location.search + location.hash
  if (!location.pathname.startsWith('/login')) {
    const back = encodeURIComponent(here || '/chat')
    location.href = `/login?redirect=${back}`
  }
}

;[httpCore, httpWEB].forEach(inst => {
  inst.interceptors.response.use(
    (res) => res,
    (err) => {
      const resp = err?.response
      // 只拦 Core/Web 的 401；其余状态直接抛出
      if (resp && resp.status === 401) {
        onUnauthorized()
      }
      return Promise.reject(err)
    }
  )
})

// ====== 便捷导出：环境端点 ======
export const endpoints = { LLM_BASE, ASR_BASE, TTS_BASE, CORE_BASE }

// ====== （可选）Web Search 便捷 API ======
// 也可以单独建文件 src/api/websearch.js 引入 httpWEB 使用
export const webAPI = {
  search: (payload) => httpCore.post('/web/search', payload),                 // 返回 axios Response
  page:   (id)      => httpCore.get(`/web/page/${id}`),
  list:   (q='', limit=20, offset=0) =>
             httpCore.get('/web/items', { params: { q, limit, offset } }),
  create: (body)    => httpCore.post('/web/items', body),
  update: (id, b)   => httpCore.put(`/web/items/${id}`, b),
  remove: (id)      => httpCore.delete(`/web/items/${id}`)
}
