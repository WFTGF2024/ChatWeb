import { httpASR, readNDJsonStream } from './http'

export async function asr(file){
  const fd = new FormData()
  fd.append('file', file, file.name || 'audio.wav')
  const { data } = await httpASR.post('/asr', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
  // {"text":"你好..."}
  return data
}

export async function asrStream(file, onPartial){
  const fd = new FormData()
  fd.append('file', file, file.name || 'audio.wav')
  const resp = await fetch(`${httpASR.defaults.baseURL}/asr/stream`, {
    method: 'POST',
    body: fd
  })
  if(!resp.ok) throw new Error('ASR stream failed')
  await readNDJsonStream(resp, (obj)=>{
    if(obj.partial_text) onPartial(obj.partial_text)
  })
}