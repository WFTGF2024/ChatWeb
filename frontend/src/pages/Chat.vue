<template>
  <div class="chat-page">
    <!-- 顶部工具栏（吸顶） -->
    <header class="toolbar">
      <ChatHeader :role="role" />

      <div class="toolbar-row">
        <div class="session-ops">
          <label class="label">会话</label>
          <select v-model="chatId" class="select" v-if="chatList.length">
            <option v-for="c in chatList" :key="c.id" :value="c.id">
              {{ c.title || ('会话 ' + c.id) }}
            </option>
          </select>
          <button class="btn ghost" @click="newChat">＋新建</button>
          <button class="btn ghost" :disabled="!chatId" @click="renameChat">✏️ 重命名</button>
          <button class="btn ghost danger" :disabled="!chatId" @click="removeChat">🗑 删除</button>
        </div>

        <div class="tts-ops">
          <label class="chip">
            <input type="checkbox" v-model="chat.settings.voiceEnabled" />
            语音播报
          </label>
          <small class="hint">开启后自动分段合成并连续播放</small>
          <a
            v-if="mergedAudioUrl"
            :href="mergedAudioUrl"
            download="chat_reply.wav"
            class="btn"
            title="下载整段音频"
          >下载音频</a>
        </div>
      </div>
    </header>

    <!-- 消息区（填满剩余空间，可滚动） -->
    <main class="messages" ref="msgBox">
      <div
        v-for="(m, idx) in messages"
        :key="idx"
        class="msg"
        :class="m.role === 'assistant' ? 'ai' : 'user'"
      >
        <div class="avatar" :aria-label="m.role">{{ m.role === 'assistant' ? '🤖' : '🧑' }}</div>
        <div class="bubble">
          <pre class="content">{{ m.content }}</pre>
          <div class="inline-actions" v-if="m.role === 'assistant'">
            <button class="chip" @click="speak(m.content)">🔊 朗读</button>
          </div>
        </div>
      </div>

      <div v-if="streamingText" class="msg ai">
        <div class="avatar">🤖</div>
        <div class="bubble typing">
          <span class="typing-dot"></span><span class="typing-dot"></span><span class="typing-dot"></span>
        <div class="content live">{{ streamingText }}</div>
        </div>
      </div>
    </main>

    <!-- 输入区（吸底） -->
    <footer class="composer">
      <textarea
        v-model="input"
        class="input"
        rows="3"
        placeholder="输入内容，Shift+Enter 换行，Enter 发送"
        @keydown.enter.exact.prevent="send"
        @keydown.shift.enter.stop
      ></textarea>

      <div class="composer-actions">
        <AudioRecorder @done="onAudioDone" />
        <button class="btn primary" :disabled="sending || !input.trim()" @click="send">发送</button>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick, computed } from 'vue'
import ChatHeader from '../components/ChatHeader.vue'
import AudioRecorder from '../components/AudioRecorder.vue'
import {
  listChats,
  upsertChatMeta,
  removeChat as cacheRemove,
  loadChat,
  saveChat,
  getWithCache
} from '../utils/chatCache'
import {
  createSession,
  listSessions,
  removeSession as apiRemove,
  appendMessage,
  listMessages,
  updateSession
} from '../api/core'
import { chatStream } from '../api/llm'
import { asr as asrOnce } from '../api/asr'
import { synthesize } from '../api/tts'

// --- TTS segmented playback & merge ---
import { mergeWaveBlobs } from '../utils/audio'
const speaking = ref(false)
const segments = ref([])  // { text, blob, url }
const mergedAudioUrl = ref('')

function extractParagraphs(text){
  let t = text
    .replace(/```[\s\S]*?```/g, '') // 去掉代码块
    .replace(/^>.*$/gm, '')         // 去掉引用
    .replace(/^\s*\*\*.*\*\*\s*$/gm, '') // 去掉粗体标题行
  const paras = t.split(/\n{2,}/).map(s=>s.trim()).filter(s=>s && s.length > 2)
  const final = []
  for(const p of paras){
    if(p.length <= 220){ final.push(p); continue }
    const parts = p.split(/(?<=[。！？.!?])/)
    let buf = ''
    for(const part of parts){
      if((buf + part).length > 220){ final.push(buf.trim()); buf = part }
      else buf += part
    }
    if(buf.trim()) final.push(buf.trim())
  }
  return final.slice(0, 24)
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
      blobs.push(blob)
      // 顺序播放
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

const isLogin = computed(() => user.isLogin)

const msgBox = ref(null)
function scrollToBottom(){ nextTick(()=>{ if(msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight }) }

async function hydrateSessions(){
  // 本地优先
  chatList.value = listChats()
  // 已登录则拉取服务端（10s TTL）
  if(isLogin.value) {
    const server = await getWithCache('cache:sessions', 10_000, async ()=>{
      try{
        const data = await listSessions()
        return data
      }catch{ return [] }
    })
    if(Array.isArray(server) && server.length){
      chatList.value = server.map(x => ({ id: x.id, title: x.title }))
      server.forEach(x => upsertChatMeta({ id: x.id, title: x.title }))
    }
  }
  if(!chatId.value && chatList.value.length) chatId.value = chatList.value[0].id
}

watch(chatId, async (id)=>{
  if(!id) return
  const s = loadChat(id)
  messages.value = s.messages || []
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

async function renameChat() {
  if (!chatId.value) return

  const id = chatId.value
  const old = chatList.value.find(x => x.id === id)?.title || `会话 ${id}`

  // 1. 让用户输入
  let title = window.prompt('输入新的会话名称（1-40 字）', old)
  if (title == null) return

  // 2. 本地规整：去空格 + 截断
  title = title.trim().slice(0, 40)

  // 3. 空的就不改；跟原来一样也不改
  if (!title || title === old) return

  // 先备份一份，失败回滚
  const snapshot = [...chatList.value]

  try {
    // 4. 本地乐观更新
    upsertChatMeta({ id, title })
    chatList.value = listChats()

    // 5. 登录状态下同步到后端
    if (isLogin.value) {
      // 这里传“字符串”，不是对象
      await updateSession(id, title)
    }
  } catch (e) {
    // 6. 出错回滚
    chatList.value = snapshot
    upsertChatMeta({ id, title: old })
    alert('重命名失败：' + (e?.response?.data?.error || e.message || '未知错误'))
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

  if(isLogin.value){ try{ await appendMessage(id, 'user', userText) }catch{} }

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

    if(chat.settings.voiceEnabled){
      try{ await ttsPlaySegments(aiMsg.content, chat.settings.ttsStyle || 'style2') }catch(e){ console.warn(e) }
    }
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
/* 结构布局 */
.chat-page{
  --bg: #f6f7fb;
  --panel: #fff;
  --muted: #6b7280;
  --line: #e5e7eb;
  --shadow: 0 6px 24px rgba(0,0,0,0.06);
  --radius: 14px;

  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: var(--bg);
}

.toolbar{
  position: sticky;
  top: 0;
  z-index: 5;
  background: linear-gradient(180deg, #ffffff 80%, rgba(255,255,255,0.6));
  box-shadow: 0 1px 0 var(--line);
  padding: 8px 12px 10px;
}

.toolbar-row{
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: center;
  margin-top: 8px;
}

.session-ops{
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.tts-ops{
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.label{ color: var(--muted); font-size: 13px; }

.messages{
  flex: 1;
  min-height: 0; /* 允许子元素滚动 */
  overflow: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.msg{ display: grid; grid-template-columns: 36px 1fr; gap: 10px; align-items: flex-start; }
.msg.user{ grid-template-columns: 1fr 36px; }
.msg.user .avatar{ order: 2; }
.msg.user .bubble{ order: 1; justify-self: end; background: #eef7f3; }
.msg.ai .bubble{ background: #f7f7f9; }

.avatar{
  width: 36px; height: 36px; border-radius: 50%;
  display:flex; align-items:center; justify-content:center;
  background: #fff; border: 1px solid var(--line);
  box-shadow: var(--shadow);
  font-size: 18px;
}

.bubble{
  max-width: min(820px, 88%);
  border-radius: 16px;
  padding: 8px 14px;
  box-shadow: var(--shadow);
  border: 1px solid var(--line);

}


.content.live{ margin-top: 6px; }

.inline-actions{
  display: flex; gap: 8px; margin-top: 6px;
}

.typing{
  display: inline-block;
  position: relative;
}
.typing-dot{
  display:inline-block;width:6px;height:6px;border-radius:50%;
  background:#c0c4ce;margin-right:4px;animation:blink 1.2s infinite ease-in-out;
}
.typing-dot:nth-child(2){ animation-delay: .2s }
.typing-dot:nth-child(3){ animation-delay: .4s }
@keyframes blink { 0%, 80%, 100% { opacity:.2 } 40% { opacity: 1 } }

/* 输入区（吸底） */
.composer{
  position: sticky; bottom: 0; z-index: 4;
  background: var(--panel);
  border-top: 1px solid var(--line);
  box-shadow: 0 -6px 24px rgba(0,0,0,0.04);
  display: grid; gap: 10px;
  padding: 12px;
}

.input{
  width: 100%; padding: 10px 12px; border-radius: var(--radius);
  border: 1px solid var(--line); resize: vertical; min-height: 80px;
  font-size: 14px; line-height: 1.6; background: #fff;
}

.composer-actions{
  display: flex; align-items: center; justify-content: space-between; gap: 10px;
}

/* 控件风格 */
.select{ padding:6px 10px; border-radius:10px; border:1px solid var(--line); background:#fff; }
.btn{
  padding: 6px 12px; border-radius: 10px; border: 1px solid var(--line);
  background: #fff; cursor: pointer; transition: .15s ease;
}
.btn:hover{ transform: translateY(-1px); box-shadow: var(--shadow); }
.btn.ghost{ background: transparent; }
.btn.primary{ background: #428bff; color: #fff; border-color: #428bff; }
.btn.danger{ color: #c44545; border-color: #e4b7b7; }
.chip{
  padding: 4px 10px; border-radius: 999px; border:1px solid var(--line);
  display: inline-flex; align-items:center; gap:6px; background:#fff;
}
.hint{ color: var(--muted); }





</style>
