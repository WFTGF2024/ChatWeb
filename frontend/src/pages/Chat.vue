<template>
  <div class="chat-wrap">
    <ChatHeader :role="role" />

    <div class="row" style="gap:8px; align-items:center; margin:8px 0;">
      <label>会话：</label>
      <select v-model="chatId" class="select" v-if="chatList.length">
        <option v-for="c in chatList" :key="c.id" :value="c.id">{{ c.title || ('会话 ' + c.id) }}</option>
      </select>
      <button class="btn ghost" @click="newChat">＋新建</button>
      <button class="btn ghost danger" :disabled="!chatId" @click="removeChat">🗑 删除</button>
    </div>

    <div class="row" style="gap:8px;align-items:center;">
      <label class="chip"><input type="checkbox" v-model="chat.settings.voiceEnabled"> 语音播报</label>
      <small class="hint">开启后将自动分段合成并连续播放。</small>
      <a v-if="mergedAudioUrl" :href="mergedAudioUrl" download="chat_reply.wav" class="btn">下载整段音频</a>
    </div>

    <div class="panel messages" ref="msgBox">
      <div v-for="(m,idx) in messages" :key="idx" class="msg" :class="m.role === 'assistant' ? 'ai' : 'user'">
        <div class="bubble">
          <pre style="white-space:pre-wrap">{{ m.content }}</pre>
          <div class="row" style="gap:6px; margin-top:6px;" v-if="m.role === 'assistant'">
            <button class="chip" @click="speak(m.content)">🔊 朗读</button>
          </div>
        </div>
      </div>
      <div v-if="streamingText" class="msg ai">
        <div class="bubble">{{ streamingText }}</div>
      </div>
    </div>

    <div class="panel input-row">
      <textarea v-model="input" class="input" rows="3" placeholder="输入内容..."></textarea>
      <div class="row" style="gap:8px; align-items:center;">
        <AudioRecorder @done="onAudioDone" />
        <button class="btn primary" :disabled="sending || !input.trim()" @click="send">发送</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick, computed } from 'vue'
import ChatHeader from '../components/ChatHeader.vue'
import AudioRecorder from '../components/AudioRecorder.vue'
import { listChats, upsertChatMeta, removeChat as cacheRemove, loadChat, saveChat, getWithCache } from '../utils/chatCache'
import { createSession, listSessions, removeSession as apiRemove, appendMessage, listMessages } from '../api/core'
import { chatStream } from '../api/llm'
import { asr as asrOnce } from '../api/asr'
import { synthesize } from '../api/tts'

// --- TTS segmented playback & merge ---
import { mergeWaveBlobs } from '../utils/audio'
const speaking = ref(false)
const segments = ref([])  // { text, blob, url }
const mergedAudioUrl = ref('')

function extractParagraphs(text){
  // remove code blocks and metadata-like lines
  let t = text.replace(/```[\s\S]*?```/g, '').replace(/^>.*$/gm, '').replace(/^\s*\*\*.*\*\*\s*$/gm, '')
  // heuristically keep Chinese/English body, split by blank lines or long lines
  const paras = t.split(/\n{2,}/).map(s=>s.trim()).filter(s=>s && s.length > 2)
  // further split very long paragraphs
  const final = []
  for(const p of paras){
    if(p.length <= 220){ final.push(p); continue }
    // split by punctuation
    const parts = p.split(/(?<=[。！？.!?])/)
    let buf = ''
    for(const part of parts){
      if((buf + part).length > 220){ final.push(buf.trim()); buf = part }
      else buf += part
    }
    if(buf.trim()) final.push(buf.trim())
  }
  return final.slice(0, 24)  // cap to 24 segments to avoid abuse
}

async function ttsPlaySegments(text, style){
  speaking.value = true
  segments.value = []
  mergedAudioUrl.value = ''
  const paras = extractParagraphs(text)
  const blobs = []
  for(const p of paras){
    try{
      const blob = await synthesize({ text: p, style })
      const url = URL.createObjectURL(blob)
      segments.value.push({ text: p, blob, url })
      blobs.append? blobs.append(blob) : blobs.push(blob)
      // autoplay this segment now
      await new Promise((resolve)=>{
        const audio = new Audio(url)
        audio.onended = resolve
        audio.onerror = resolve
        audio.play().catch(()=>resolve())
      })
    }catch(e){
      console.warn('TTS segment failed', e)
    }
  }
  // merge for download
  try{
    const merged = await mergeWaveBlobs(blobs)
    mergedAudioUrl.value = URL.createObjectURL(merged)
  }catch(e){ console.warn('merge failed', e) }
  speaking.value = false
}

import { useChatStore } from '../store/chat'
import { useUserStore } from '../store/user'

const role = ref({ name: '默认', avatar: '🧠' })
const chat = useChatStore()
const user = useUserStore()
const chatList = ref([])
const chatId = ref(null)
const messages = ref([])
const input = ref('')
const sending = ref(false)
const streamingText = ref('')

// 添加登录状态计算属性
const isLogin = computed(() => user.isLogin)

const msgBox = ref(null)
function scrollToBottom(){ nextTick(()=>{ if(msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight }) }

async function hydrateSessions(){
  // cache first
  chatList.value = listChats()
  // then pull from DB (10s TTL)
  if(isLogin.value) {
    const server = await getWithCache('cache:sessions', 10_000, async ()=>{
      try{
        const data = await listSessions()
        return data
      }catch{ return [] }
    })
    if(Array.isArray(server) && server.length){
      chatList.value = server.map(x => ({ id: x.id, title: x.title }))
      // update cache metas
      server.forEach(x => upsertChatMeta({ id: x.id, title: x.title }))
    }
  }
  if(!chatId.value && chatList.value.length) chatId.value = chatList.value[0].id
}

watch(chatId, async (id)=>{
  if(!id) return
  // load from cache first
  const s = loadChat(id)
  messages.value = s.messages || []
  // then load from DB if empty and user is logged in
  if(messages.value.length === 0 && isLogin.value){
    try{
      const data = await listMessages(id)
      if(Array.isArray(data)) messages.value = data
      saveChat(id, { messages: messages.value })
    }catch{}
  }
  scrollToBottom()
})

async function newChat(){
  try{
    let data
    if(isLogin.value) {
      data = await createSession('新会话')
    } else {
      data = { id: Date.now(), title: '临时会话' }
    }
    const id = data.id || data.session_id || Date.now()
    upsertChatMeta({ id, title: data.title || '临时会话' })
    chatList.value = listChats()
    chatId.value = id
    messages.value = []
    saveChat(id, { messages: [] })
  }catch(e){
    alert('创建会话失败：' + e?.message)
  }
}

async function removeChat(){
  if(!chatId.value) return
  const id = chatId.value
  if(isLogin.value) {
    try{ await apiRemove(id) }catch{ /* ignore */ }
  }
  cacheRemove(id)
  chatList.value = listChats()
  chatId.value = chatList.value[0]?.id || null
  messages.value = chatId.value ? (loadChat(chatId.value).messages) : []
}

async function send(){
  if(!chatId.value) await newChat()
  const id = chatId.value
  const userText = input.value.trim()
  if(!userText) return
  input.value = ''
  const userMsg = { role:'user', content: userText }
  messages.value.push(userMsg)
  saveChat(id, { messages: messages.value })

  // also persist to DB if logged in
  if(isLogin.value){ try{ await appendMessage(id, 'user', userText) }catch{} }

  // stream assistant
  sending.value = true
  streamingText.value = ''
  scrollToBottom()
  try{
    const payload = isLogin.value ? messages.value : [{ role: 'user', content: userText }]
    await chatStream(payload, (delta)=>{
      streamingText.value += delta
      scrollToBottom()
    })
    const aiMsg = { role:'assistant', content: streamingText.value }
    messages.value.push(aiMsg)
    streamingText.value = ''
    saveChat(id, { messages: messages.value })
    if(chat.settings.voiceEnabled){ try{ await ttsPlaySegments(aiMsg.content, chat.settings.ttsStyle || 'style2') }catch(e){ console.warn(e) } }(id, { messages: messages.value })
    if(isLogin.value){ try{ await appendMessage(id, 'assistant', aiMsg.content) }catch{} }
  }catch(e){
    alert('LLM 生成失败：' + e?.message)
  }finally{
    sending.value = false
    scrollToBottom()
  }
}

async function onAudioDone(file){
  try{
    const data = await asrOnce(file)
    const text = data?.text || ''
    if(text) input.value = (input.value ? (input.value + ' ') : '') + text
  }catch(e){
    alert('ASR 失败：' + e?.message)
  }
}

async function speak(text){
  try{
    const blob = await synthesize({ text, style: 'style2' })
    const url = URL.createObjectURL(blob)
    const audio = new Audio(url)
    audio.play()
  }catch(e){
    alert('TTS 失败：' + e?.message)
  }
}

onMounted(hydrateSessions)
</script>

<style scoped>
.chat-wrap{ max-width: 960px; margin: 0 auto; display:flex; flex-direction:column; gap:12px; }
.panel{ background: white; border-radius: 12px; padding: 12px; box-shadow: var(--shadow); }
.messages{ height: 52vh; overflow: auto; display:flex; flex-direction:column; gap:12px; }
.msg{ display:flex; }
.msg.user{ justify-content: flex-end; }
.msg.ai{ justify-content: flex-start; }
.bubble{ background: #f7f7f7; padding: 10px 12px; border-radius: 12px; max-width: 80%; white-space: pre-wrap; }
.input-row{ display:flex; flex-direction:column; gap:8px; }
.input{ width:100%; padding:8px; border-radius:8px; border:1px solid #ddd; }
.row{ display:flex; }
.select{ padding:6px 8px; border-radius:8px; border:1px solid #ddd; }
.btn{ padding:6px 12px; border-radius:10px; border:1px solid #ddd; background:#fff; cursor:pointer; }
.btn.ghost{ background: transparent; }
.btn.primary{ background: #4a8; color:white; border-color:#4a8; }
.btn.danger{ border-color:#b55; color:#b55; }
.chip{ padding:4px 8px; border-radius: 999px; border:1px solid #ddd; cursor:pointer; }
</style>
