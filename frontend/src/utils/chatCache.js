// localStorage-based cache with TTL + chat helpers
const LIST_KEY = 'chat:list:v2'
const SESSION_KEY = (id) => `chat:session:${id}:v2`

const readJSON = (k, d=null) => {
  try{ const v = localStorage.getItem(k); return v ? JSON.parse(v) : d }catch{ return d }
}
const writeJSON = (k, v) => localStorage.setItem(k, JSON.stringify(v))

export function listChats(){
  return readJSON(LIST_KEY, [])
}
export function upsertChatMeta({ id, title }){
  const now = Date.now()
  const list = listChats()
  const i = list.findIndex(x => x.id === id)
  if(i >= 0){
    list[i] = { ...list[i], title, updatedAt: now }
  }else{
    list.unshift({ id, title, updatedAt: now })
  }
  writeJSON(LIST_KEY, list)
  return list
}
export function removeChat(id){
  const list = listChats().filter(x => x.id !== id)
  writeJSON(LIST_KEY, list)
  localStorage.removeItem(SESSION_KEY(id))
  return list
}
export function loadChat(id){
  const s = readJSON(SESSION_KEY(id), null)
  if(!s) return { messages: [], meta: {}, ts: 0 }
  return s
}
export function saveChat(id, payload){
  const messages = Array.isArray(payload) ? payload : (payload.messages || [])
  const meta = Array.isArray(payload) ? {} : (payload.meta || {})
  const s = { messages, meta, ts: Date.now() }
  writeJSON(SESSION_KEY(id), s)
  return s
}

// Generic cache with TTL
export async function getWithCache(cacheKey, ttlMs, fetcher){
  const cached = readJSON(cacheKey, null)
  const now = Date.now()
  if(cached && (now - (cached.ts || 0) < ttlMs)){
    return cached.value
  }
  const value = await fetcher()
  writeJSON(cacheKey, { ts: now, value })
  return value
}
export function deleteChat(id) {
  const list = listChats().filter(x => x.id !== id)
  writeJSON(LIST_KEY, list)
  localStorage.removeItem(SESSION_KEY(id))
  // 同步到后端
  if (user.isLogin) {
    httpCore.delete(`/chat/sessions/${id}`)
      .catch(err => console.error('Failed to delete session:', err))
  }
}
