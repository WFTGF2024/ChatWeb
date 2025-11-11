
// src/utils/audio.js
// Merge multiple WAV blobs into one WAV using WebAudio
export async function mergeWaveBlobs(blobs){
  if(!blobs?.length) throw new Error('no blobs')
  const ctx = new (window.OfflineAudioContext || window.webkitOfflineAudioContext)(1, 44100*60*10, 44100) // up to 10 minutes mono
  const buffers = []
  const audioCtx = new (window.AudioContext || window.webkitAudioContext)()
  for(const b of blobs){
    const arr = await b.arrayBuffer()
    const buf = await audioCtx.decodeAudioData(arr.slice(0))  // decode in real ctx
    buffers.push(buf)
  }
  // concatenate
  let totalLength = buffers.reduce((s,b)=> s + b.length, 0)
  const out = ctx.createBuffer(1, totalLength, 44100)
  let offset = 0
  for(const b of buffers){
    out.getChannelData(1-1).set(b.getChannelData(1-1), offset)
    offset += b.length
  }
  // encode WAV
  return encodeWAV(out)
}
function encodeWAV(buffer){
  const numOfChan = buffer.numberOfChannels
  const length = buffer.length * numOfChan * 2 + 44
  const out = new ArrayBuffer(length)
  const view = new DataView(out)
  // RIFF header
  writeString(view, 0, 'RIFF')
  view.setUint32(4, 36 + buffer.length * numOfChan * 2, true)
  writeString(view, 8, 'WAVE')
  // fmt chunk
  writeString(view, 12, 'fmt ')
  view.setUint32(16, 16, true)
  view.setUint16(20, 1, true) // PCM
  view.setUint16(22, numOfChan, true)
  view.setUint32(24, buffer.sampleRate, true)
  view.setUint32(28, buffer.sampleRate * numOfChan * 2, true)
  view.setUint16(32, numOfChan * 2, true)
  view.setUint16(34, 16, true)
  // data
  writeString(view, 36, 'data')
  view.setUint32(40, buffer.length * numOfChan * 2, true)
  // interleave & write
  let offset = 44
  const channels = []
  for (let i = 0; i < numOfChan; i++) channels.push(buffer.getChannelData(i))
  const interleaved = interleave(channels)
  for (let i = 0; i < interleaved.length; i++, offset += 2) {
    const s = Math.max(-1, Math.min(1, interleaved[i]))
    view.setInt16(offset, s < 0 ? s * 0x8000 : s * 0x7FFF, true)
  }
  return new Blob([view], { type: 'audio/wav' })
}
function interleave(channels){
  const length = channels[0].length
  const result = new Float32Array(length * channels.length)
  let index = 0
  for (let i = 0; i < length; i++) {
    for (let c = 0; c < channels.length; c++) result[index++] = channels[c][i]
  }
  return result
}
function writeString(view, offset, string){
  for (let i=0; i<string.length; i++) view.setUint8(offset + i, string.charCodeAt(i))
}
