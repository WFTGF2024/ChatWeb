import { httpCore } from './http'

// ---- Chat sessions & messages ----
export async function createSession(title){
  const { data } = await httpCore.post('/chat/sessions', { title })
  return data
}
export async function listSessions(){
  const { data } = await httpCore.get('/chat/sessions')
  return data
}
export async function removeSession(id){
  const { data } = await httpCore.delete(`/chat/sessions/${id}`)
  return data
}
export async function appendMessage(session_id, role, content){
  const body = { session_id, role, content }
  const { data } = await httpCore.post('/chat/messages', body)
  return data
}
export async function listMessages(session_id){
  const { data } = await httpCore.get(`/chat/sessions/${session_id}/messages`)
  return data
}
export async function streamCompletion(session_id, content, model=''){
  // Backend supports: POST /chat/sessions/:id/stream { content, model }
  const resp = await httpCore.post(`/chat/sessions/${session_id}/stream`, { content, model }, { responseType: 'stream' })
  return resp
}

export async function updateSession(id, title) {
  const t = (title || '').trim()
  if (!t) {
    // 可以直接抛个错误，让调用方去提示“标题不能为空”
    throw new Error('title required')
  }
  const { data } = await httpCore.put(`/chat/sessions/${id}`, { title: t })
  return data
}
