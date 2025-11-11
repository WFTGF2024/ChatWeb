
<template>
  <div class="wrap">
    <h2>ASR / TTS 实验台</h2>
    <div class="grid">
      <div class="card">
        <h3>语音识别（ASR）</h3>
        <p class="hint">上传一段语音，调用 /asr 或 /asr/stream 识别。</p>
        <input type="file" accept="audio/*" @change="chooseASR">
        <div class="row" style="gap:8px;margin-top:8px;">
          <button class="btn" :disabled="!asrFile || loading" @click="runASR">识别</button>
          <button class="btn" :disabled="!asrFile || loading" @click="runASRStream">流式识别</button>
        </div>
        <p v-if="asrText" class="sub">识别结果：{{ asrText }}</p>
      </div>
      <div class="card">
        <h3>语音合成（TTS）</h3>
        <textarea class="input" rows="5" v-model="ttsText" placeholder="输入要朗读的文本"></textarea>
        <div class="row" style="gap:8px; align-items:center;">
          <label>风格</label>
          <select v-model="style" class="select">
            <option value="style1">面试官</option>
            <option value="style2">生活聊天</option>
            <option value="style3">苏格拉底</option>
          </select>
          <label class="chip">
            <input type="checkbox" v-model="usePrompt"> 参考音频
          </label>
          <input type="file" v-if="usePrompt" accept="audio/*" @change="choosePrompt">
          <button class="btn primary" :disabled="!ttsText || loading" @click="runTTS">合成并播放</button>
        </div>
        <audio v-if="audioUrl" :src="audioUrl" controls style="margin-top:8px;width:100%"></audio>
        <p class="sub">健康检查：<span :class="health?.status==='ok'?'ok':'err'">{{ health?.status || '-' }}</span></p>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { asr, asrStream } from '../api/asr'
import { synthesize, ttsHealth } from '../api/tts'
const loading = ref(false)
const asrFile = ref(null)
const asrText = ref('')
const ttsText = ref('你好，这是TTS合成测试。')
const style = ref('style2')
const usePrompt = ref(false)
const promptFile = ref(null)
const audioUrl = ref('')
const health = ref(null)
function chooseASR(e){ asrFile.value = e.target.files[0] || null }
function choosePrompt(e){ promptFile.value = e.target.files[0] || null }
async function runASR(){
  if(!asrFile.value) return
  loading.value = true
  try{
    const data = await asr(asrFile.value)
    asrText.value = data?.text || ''
  }finally{ loading.value = false }
}
async function runASRStream(){
  if(!asrFile.value) return
  loading.value = true
  try{
    let last = ''
    await asrStream(asrFile.value, (partial)=>{ last = partial; asrText.value = partial })
    asrText.value = last
  }finally{ loading.value = false }
}
async function runTTS(){
  loading.value = true
  try{
    const blob = await synthesize({ text: ttsText.value, style: style.value, prompt_audio: promptFile.value || undefined })
    audioUrl.value = URL.createObjectURL(blob)
    const a = new Audio(audioUrl.value); a.play()
  }finally{ loading.value = false }
}
onMounted(async ()=>{ try{ health.value = await ttsHealth() }catch{} })
</script>
<style scoped>
.wrap{ max-width: 980px; margin:0 auto; display:flex; flex-direction:column; gap:14px; }
.grid{ display:grid; grid-template-columns: 1fr 1fr; gap:14px; }
.card{ background:#fff; border-radius:14px; padding:14px; box-shadow: var(--shadow); }
.hint{ color:#6b7280; }
.row{ display:flex; }
.btn{ padding:8px 12px; border:none; border-radius:10px; background:#f3f4f6; cursor:pointer; }
.btn.primary{ background:#2563eb; color:white; }
.input, .select{ width:100%; padding:8px; border-radius:8px; border:1px solid #e5e7eb; }
.select{ width:auto; }
.chip{ display:flex; gap:8px; align-items:center; padding:6px 10px; border-radius:999px; border:1px solid #e5e7eb; }
.ok{ color:#10b981; } .err{ color:#ef4444; }
</style>
