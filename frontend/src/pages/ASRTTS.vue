
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
.wrap{
  /* 局部主题变量 */
  --bg: #f5f5f6;
  --panel: #ffffff;
  --line: #e5e7eb;
  --muted: #6b7280;
  --heading: #111827;
  --primary: #2563eb;
  --radius: 14px;
  --shadow: 0 10px 30px rgba(15, 23, 42, 0.05);

  max-width: 1080px;
  margin: 0 auto;
  padding: 16px 14px 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  background: var(--bg);
  min-height: 100%;
}

.wrap > h2{
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--heading);
}

.grid{
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  align-items: stretch;
}

.card{
  background: var(--panel);
  border: 1px solid rgba(229, 231, 235, .8);
  border-radius: var(--radius);
  padding: 16px 16px 14px;
  box-shadow: var(--shadow);
  transition: box-shadow .2s ease, transform .2s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 100%;
}
.card:hover{
  transform: translateY(-1px);
  box-shadow: 0 14px 36px rgba(15, 23, 42, 0.08);
}

.card h3{
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--heading);
}

.hint{
  color: var(--muted);
  margin: 0;
  line-height: 1.45;
  font-size: 13px;
}

.row{
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

/* 按钮 */
.btn{
  padding: 8px 14px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  user-select: none;
  font-size: 13px;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}
.btn:hover{
  transform: translateY(-1px);
  box-shadow: var(--shadow);
}
.btn:disabled{
  opacity: .6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}
.btn.primary{
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}
.btn.primary:hover{ filter: brightness(0.96); }

/* 输入控件 */
.input, .select, textarea.input{
  width: 100%;
  padding: 9px 11px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  outline: none;
  font: inherit;
  line-height: 1.45;
  transition: border .15s ease, box-shadow .15s ease;
  font-size: 13px;
}
.input:focus, .select:focus, textarea.input:focus{
  border-color: rgba(37, 99, 235, 1);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, .12);
}
textarea.input{
  resize: vertical;
  min-height: 110px;
}
.select{
  width: auto;
  min-width: 150px;
}

/* 小标签/开关 */
.chip{
  display: inline-flex;
  gap: 7px;
  align-items: center;
  padding: 5px 10px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
  font-size: 12px;
  white-space: nowrap;
}
.chip input{ accent-color: var(--primary); }

/* 文件选择 */
input[type="file"]{
  display: block;
  width: 100%;
  font-size: 12px;
  color: #374151;
}

/* 识别结果/状态块 */
.sub{
  color: #111827;
  background: #fafafa;
  border: 1px dashed rgba(229, 231, 235, .9);
  border-radius: 10px;
  padding: 7px 10px;
  margin-top: 4px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}

.ok{ color: #10b981; font-weight: 600; }
.err{ color: #ef4444; font-weight: 600; }

/* audio 适配卡片宽度 */
:deep(audio){
  width: 100%;
  border-radius: 10px;
  outline: none;
}

/* 响应式 */
@media (max-width: 1024px){
  .grid{
    grid-template-columns: 1fr;
  }
}
@media (max-width: 520px){
  .wrap{ padding: 12px 10px 20px; }
  .row{ flex-direction: row; }
  .btn{
    flex: 1 1 auto;
  }
  .select{
    min-width: 130px;
  }
}
</style>
