import { httpLLM, readNDJsonStream } from './http'

export async function chat(messages, model){
  const body = { messages }
  if(model) body.model = model
  const { data } = await httpLLM.post('/api/chat', body)
  // returns { content, model, raw? }
  return data
}

export async function chatStream(messages, onDelta, model){
  const body = { messages }
  if(model) body.model = model
  const resp = await fetch(`${httpLLM.defaults.baseURL}/api/chat/stream`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  if(!resp.ok) throw new Error('LLM stream failed')
  await readNDJsonStream(resp, (obj)=>{
    if(obj.delta) onDelta(obj.delta)
  })
}