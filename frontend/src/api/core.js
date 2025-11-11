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
