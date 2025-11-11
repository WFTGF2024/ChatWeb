<template>
  <div class="audio-recorder row" style="gap:8px; align-items:center;">
    <button class="btn" :class="{ danger: recording }" @click="toggle" :disabled="busy">
      <span v-if="!recording">🎙 开始</span>
      <span v-else>⏹ 停止 {{ seconds }}s</span>
    </button>
    <small class="hint" v-if="hint">{{ hint }}</small>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount, watchEffect } from 'vue'

/**
 * Props（可选）：
 * - deviceId: 指定麦克风设备 ID（不传则用默认）
 * - maxSeconds: 最长录音秒数，达到后自动停止（默认 300s）
 */
const props = defineProps({
  deviceId: { type: String, default: '' },
  maxSeconds: { type: Number, default: 300 }
})

const emit = defineEmits(['start','done','error'])

const recording = ref(false)
const busy = ref(false)
const seconds = ref(0)
const hint = ref('')

let stream = null
let mediaRecorder = null
let chunks = []
let timer = null

function startTimer(){
  clearInterval(timer)
  seconds.value = 0
  timer = setInterval(() => {
    seconds.value++
    if (props.maxSeconds > 0 && seconds.value >= props.maxSeconds) {
      // 达到最长录音时长，自动停止
      stop()
    }
  }, 1000)
}

/** 根据浏览器支持情况选择最合适的 mimeType */
function pickMimeType() {
  const candidates = [
    'audio/webm;codecs=opus',
    'audio/webm',
    'audio/mp4', // Safari 可能支持（但 MediaRecorder 在 Safari 的兼容性较弱）
    'audio/ogg;codecs=opus',
  ]
  for (const t of candidates) {
    try {
      if (window.MediaRecorder && MediaRecorder.isTypeSupported && MediaRecorder.isTypeSupported(t)) {
        return t
      }
    } catch {}
  }
  // 兜底
  return 'audio/webm'
}

/** 统一映射错误信息，给用户更有用的提示 */
function mapErrorMessage(e){
  const map = {
    NotAllowedError: '浏览器已阻止麦克风，请在地址栏站点权限里改为“允许”。',
    NotFoundError: '未检测到可用麦克风，请检查设备连接/驱动。',
    NotReadableError: '设备被占用或驱动异常，请关闭占用麦克风的程序后重试。',
    SecurityError: '当前非安全上下文（请使用 localhost 或 HTTPS 访问）。',
    OverconstrainedError: '设备不满足约束（deviceId/采样参数不被支持）。',
    AbortError: '录音被系统中断，请重试。'
  }
  // 某些浏览器只给 NotAllowedError/NotReadableError，不给 message；我们拼一个更完整的
  return map[e?.name] || `获取麦克风失败：${e?.name || 'Unknown'}${e?.message ? ' - ' + e.message : ''}`
}

/** 更稳妥的麦克风获取流程 */
async function ensureMic() {
  if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
    throw new Error('当前浏览器不支持录音（缺少 getUserMedia）。')
  }

  // 可选：先探测权限（有些浏览器不支持 permissions API）
  try {
    const st = await navigator.permissions.query({ name: 'microphone' })
    // st.state: 'granted' | 'denied' | 'prompt'
    // 仅用于日志/提示，不强依赖
    console.log('mic permission:', st.state)
  } catch {}

  // 录音约束：回声消除 & 降噪；可选设备 ID
  const constraints = {
    audio: {
      echoCancellation: true,
      noiseSuppression: true,
      deviceId: props.deviceId ? { exact: props.deviceId } : undefined,
    }
  }
  return await navigator.mediaDevices.getUserMedia(constraints)
}

/** 开始录音 */
async function start(){
  if (recording.value || busy.value) return
  busy.value = true
  hint.value = ''

  try{
    // 1) 权限 & 流
    stream = await ensureMic()

    // 2) MediaRecorder 支持判断
    if (!window.MediaRecorder) {
      throw new Error('当前浏览器不支持 MediaRecorder。')
    }
    const mt = pickMimeType()
    chunks = []
    mediaRecorder = new MediaRecorder(stream, { mimeType: mt })

    // 3) 事件
    mediaRecorder.ondataavailable = (e) => {
      if (e.data && e.data.size) chunks.push(e.data)
    }
    mediaRecorder.onstop = () => {
      try { clearInterval(timer) } catch {}
      timer = null
      recording.value = false

      const type = mediaRecorder && mediaRecorder.mimeType ? mediaRecorder.mimeType : 'audio/webm'
      const blob = new Blob(chunks, { type })
      const ext = type.includes('mp4') ? 'mp4' : (type.includes('ogg') ? 'ogg' : 'webm')
      const file = new File([blob], `recording_${Date.now()}.${ext}`, { type })
      emit('done', file)

      // 释放设备
      try { stream?.getTracks?.().forEach(t => t.stop()) } catch {}
      stream = null
      mediaRecorder = null
    }
    mediaRecorder.onerror = (e) => {
      const err = e?.error || e
      hint.value = mapErrorMessage(err)
      emit('error', err)
      stop()
    }

    // 4) 开始
    mediaRecorder.start(200) // timeslice 200ms，边录边吐
    startTimer()
    recording.value = true
    emit('start')
  }catch(e){
    hint.value = mapErrorMessage(e)
    emit('error', e)
    // 出错也要释放已占用的设备
    try { stream?.getTracks?.().forEach(t => t.stop()) } catch {}
    stream = null
    mediaRecorder = null
  }finally{
    busy.value = false
  }
}

/** 停止录音 */
function stop(){
  try {
    if (mediaRecorder && mediaRecorder.state !== 'inactive') {
      mediaRecorder.stop()
    }
  } catch {}
  try {
    stream?.getTracks?.().forEach(t => t.stop())
  } catch {}
}

/** 切换按钮 */
function toggle(){
  if (recording.value) stop()
  else start()
}

onBeforeUnmount(()=>{
  try { stop() } catch {}
  try { clearInterval(timer) } catch {}
  stream = null
  mediaRecorder = null
  chunks = []
})

// 若外部传入 deviceId 变化，提示用户重新开始以切换设备
watchEffect(() => {
  if (props.deviceId && recording.value) {
    hint.value = '已指定新的麦克风设备，停止后再次开始生效。'
  }
})
</script>
