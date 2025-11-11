
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
  --bg: #f7f8fc;
  --panel: #ffffff;
  --line: #e5e7eb;
  --muted: #6b7280;
  --primary: #2563eb;
  --radius: 14px;
  --shadow: 0 10px 30px rgba(0,0,0,.06);

  max-width: 1000px;
  margin: 0 auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: var(--bg);
}

.grid{
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  align-items: start;
}

.card{
  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
  transition: box-shadow .2s ease, transform .2s ease;
}
.card:hover{
  transform: translateY(-1px);
  box-shadow: 0 14px 36px rgba(0,0,0,.08);
}
.card h3{
  margin: 0 0 8px;
  font-size: 16px;
}

.hint{ color: var(--muted); margin: 0 0 8px; }
.sub{
  color: #111827;
  background: #fafafa;
  border: 1px dashed var(--line);
  border-radius: 10px;
  padding: 8px 10px;
  margin-top: 10px;
}

.row{
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap; /* 小屏自动换行，避免挤压 */
}

/* 按钮 */
.btn{
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  user-select: none;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
}
.btn:hover{ transform: translateY(-1px); box-shadow: var(--shadow); }
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
.btn.primary:hover{ filter: brightness(0.95); }

/* 输入类控件 */
.input, .select, textarea.input{
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  outline: none;
  font: inherit;
  line-height: 1.45;
  transition: border .15s ease, box-shadow .15s ease;
}
.input:focus, .select:focus, textarea.input:focus{
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(37,99,235,.15);
}
textarea.input{ resize: vertical; min-height: 120px; }
.select{ width: auto; min-width: 140px; }

/* 开关/标签 */
.chip{
  display: inline-flex;
  gap: 8px;
  align-items: center;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
}
.chip input{ accent-color: var(--primary); }

/* 状态色 */
.ok{ color: #10b981; font-weight: 600; }
.err{ color: #ef4444; font-weight: 600; }

/* 文件选择 */
input[type="file"]{
  display: block;
  width: 100%;
  padding: 8px 0;
}

/* 音频控件适配卡片宽度（scoped 深度选择器） */
:deep(audio){
  width: 100%;
  border-radius: 10px;
}

/* 响应式 */
@media (max-width: 960px){
  .grid{ grid-template-columns: 1fr; }
}
@media (max-width: 520px){
  .wrap{ padding: 10px 8px; }
  .btn{ width: 100%; justify-content: center; }
}
</style>
