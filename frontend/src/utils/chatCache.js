// 简单 localStorage 多会话缓存
const LIST_KEY = 'chat:list:v1'
const SESSION_KEY = (id) => `chat:session:${id}`

const readJSON = (k, d = null) => {
  try {
    const v = localStorage.getItem(k)
    return v ? JSON.parse(v) : d
  } catch {
    return d
  }
}
const writeJSON = (k, v) => localStorage.setItem(k, JSON.stringify(v))

export function listChats () {
  return readJSON(LIST_KEY, [])
}

export function createChat (title = '新会话') {
  const id = String(Date.now()) + Math.random().toString(16).slice(2, 8)
  const list = listChats()
  const item = { id, title, updatedAt: Date.now() }
  writeJSON(LIST_KEY, [item, ...list])
  return item
}

export function renameChat (id, title) {
  const list = listChats().map(x => x.id === id ? { ...x, title, updatedAt: Date.now() } : x)
  writeJSON(LIST_KEY, list)
}

export function deleteChat (id) {
  const list = listChats().filter(x => x.id !== id)
  writeJSON(LIST_KEY, list)
  localStorage.removeItem(SESSION_KEY(id))
}

// 新版：不仅能拿到 messages，还能拿到 meta（比如你选的角色风格）
export function loadChat (id) {
  const s = readJSON(SESSION_KEY(id))
  if (!s) return { messages: [], meta: {} }

  // 兼容老数据：老的就是一个数组
  if (Array.isArray(s)) {
    return { messages: s, meta: {} }
  }

  const messages = Array.isArray(s.messages) ? s.messages : []
  const meta = s.meta || {}
  return { messages, meta }
}

// 新版：允许存 { messages, meta }，也兼容老的只传数组的方式
export function saveChat (id, payload) {
  const messages = Array.isArray(payload) ? payload : (payload.messages || [])
  const meta = Array.isArray(payload) ? {} : (payload.meta || {})
  writeJSON(SESSION_KEY(id), { ts: Date.now(), messages, meta })
  const list = listChats().map(x => x.id === id ? { ...x, updatedAt: Date.now() } : x)
  writeJSON(LIST_KEY, list)
}
