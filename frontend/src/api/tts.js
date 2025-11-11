import { httpTTS } from './http'

export async function synthesize({ text, style = 'style2', prompt_audio = null, temperature=null, top_p=null }){
  const fd = new FormData()
  fd.append('text', text)
  if(style) fd.append('style', style)
  if(temperature != null) fd.append('temperature', String(temperature))
  if(top_p != null) fd.append('top_p', String(top_p))
  if(prompt_audio) fd.append('prompt_audio', prompt_audio, prompt_audio.name || 'prompt.wav')
  const resp = await fetch(`${httpTTS.defaults.baseURL}/synthesize`, { method:'POST', body: fd })
  if(!resp.ok) throw new Error('TTS synthesize failed')
  const blob = await resp.blob()
  return blob
}

export async function ttsHealth(){
  const { data } = await httpTTS.get('/health')
  return data
}