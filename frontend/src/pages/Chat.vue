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
        <div class="avatar" :aria-label="m.role">
          {{ m.role === 'assistant' ? '🤖' : '🧑' }}
        </div>
        <div class="bubble">
          <div class="content markdown-body" v-html="renderMarkdown(m.content)"></div>

          <div class="inline-actions" v-if="m.role === 'assistant'">
            <button class="chip" @click="speak(m.content)" :disabled="speaking">
              🔊 朗读
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- 苏格拉底深度追问（显示在输入框上方，不在消息里显示） -->
    <section v-if="isSocrates && deepQuestions.length" class="dq-panel">
      <div class="dq-title">你可以继续问：</div>
      <div class="dq-list">
        <button
          v-for="(q, i) in deepQuestions"
          :key="i"
          class="dq-item"
          @click="fillQuestion(q)"
        >
          {{ q }}
        </button>
      </div>
    </section>

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
        <button class="btn primary" :disabled="sending || !input.trim()" @click="send">
          发送
        </button>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick, computed } from 'vue'
import MarkdownIt from 'markdown-it'
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
import { buildSystemPrompt } from '../utils/prompts'
import { mergeWaveBlobs } from '../utils/audio'
import { useChatStore } from '../store/chat'
import { useUserStore } from '../store/user'

/** ===== Markdown 渲染 ===== */
const md = new MarkdownIt({ linkify: true, breaks: true })
function renderMarkdown (text) { return md.render(text || '') }

/** ===== Store & 状态 ===== */
const chat = useChatStore()
const user = useUserStore()
const role = computed(() => chat.currentRole || { name: '默认', avatar: '🧠' })
const isLogin = computed(() => user.isLogin)
const isSocrates = computed(() => chat?.currentRole?.id === 'socrates')

const chatList = ref([])
const chatId = ref(null)
const messages = ref([])
const input = ref('')
const sending = ref(false)
const streamingText = ref('')

/** ===== TTS segmented playback & merge ===== */
const speaking = ref(false)
const segments = ref([])  // { text, blob, url }
const mergedAudioUrl = ref('')

/** 深挖问题（仅苏格拉底角色使用） */
const deepQuestions = ref([])

/** ===== 滚动到底部 ===== */
const msgBox = ref(null)
function scrollToBottom () {
  nextTick(() => { if (msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight })
}

/** ===== 会话装载 ===== */
async function hydrateSessions () {
  // 本地优先
  chatList.value = listChats()
  // 已登录则拉取服务端（10s TTL）
  if (isLogin.value) {
    const server = await getWithCache('cache:sessions', 10_000, async () => {
      try { return await listSessions() } catch { return [] }
    })
    if (Array.isArray(server) && server.length) {
      chatList.value = server.map(x => ({ id: x.id, title: x.title }))
      server.forEach(x => upsertChatMeta({ id: x.id, title: x.title }))
    }
  }
  if (!chatId.value && chatList.value.length) chatId.value = chatList.value[0].id
}

watch(chatId, async (id) => {
  if (!id) return
  const s = loadChat(id)
  messages.value = s.messages || []
  if (messages.value.length === 0 && isLogin.value) {
    try {
      const data = await listMessages(id)
      if (Array.isArray(data)) messages.value = data
      saveChat(id, { messages: messages.value })
    } catch {}
  }
  scrollToBottom()
})

/** ===== 新建/重命名/删除会话 ===== */
async function newChat () {
  try {
    let data
    if (isLogin.value) data = await createSession('新会话')
    else data = { id: Date.now(), title: '临时会话' }

    const id = data.id || data.session_id || Date.now()
    upsertChatMeta({ id, title: data.title || '临时会话' })
    chatList.value = listChats()
    chatId.value = id
    messages.value = []
    saveChat(id, { messages: [] })
  } catch (e) {
    alert('创建会话失败：' + e?.message)
  }
}

async function renameChat () {
  if (!chatId.value) return
  const id = chatId.value
  const old = chatList.value.find(x => x.id === id)?.title || `会话 ${id}`
  let title = window.prompt('输入新的会话名称（1-40 字）', old)
  if (title == null) return
  title = title.trim().slice(0, 40)
  if (!title || title === old) return

  const snapshot = [...chatList.value]
  try {
    upsertChatMeta({ id, title })
    chatList.value = listChats()
    if (isLogin.value) await updateSession(id, title)
  } catch (e) {
    chatList.value = snapshot
    upsertChatMeta({ id, title: old })
    alert('重命名失败：' + (e?.response?.data?.error || e.message || '未知错误'))
  }
}

async function removeChat () {
  if (!chatId.value) return
  const id = chatId.value
  if (isLogin.value) { try { await apiRemove(id) } catch {} }
  cacheRemove(id)
  chatList.value = listChats()
  chatId.value = chatList.value[0]?.id || null
  messages.value = chatId.value ? (loadChat(chatId.value).messages) : []
}

/** ===== 发送（流式 + System + 深挖抽取 + 自动TTS） ===== */
async function send () {
  if (!chatId.value) {
    await newChat()
  }
  const id = chatId.value
  const userText = input.value.trim()
  if (!userText) return

  // 开始新一轮前，清空上轮 DEEP_QUESTIONS
  deepQuestions.value = []

  // 用户发言先落本地
  input.value = ''
  const userMsg = { role: 'user', content: userText }
  messages.value.push(userMsg)
  saveChat(id, { messages: messages.value })

  if (chat.isLogin) {
    try { await appendMessage(id, 'user', userText) } catch (e) { console.warn('appendMessage failed:', e) }
  }

  sending.value = true
  streamingText.value = ''
  scrollToBottom()

  try {
    // 1) 角色 System Prompt
    const sysContent = buildSystemPrompt({
      role: chat.currentRole,
      memorySummary: chat.memorySummary || '',
      userPrefs: {}
    })
    const sysMsg = { role: 'system', content: sysContent }

    // 2) 构造 payload（不要把占位 assistant 混进去）
    const history = messages.value.filter(m => m.role !== 'system')
    const payload = chat.isLogin
      ? [sysMsg, ...history]
      : [sysMsg, { role: 'user', content: userText }]

    // 3) 现在再推一个占位 assistant，流式往里写
    const aiPlaceholder = { role: 'assistant', content: '' }
    messages.value.push(aiPlaceholder)
    scrollToBottom()

    // 4) 开始流式（苏格拉底模式：边流式边隐藏 [DEEP_QUESTIONS] 区块）
    await chatStream(payload, (delta) => {
      streamingText.value += delta
      const visible = isSocrates.value
        ? removeDeepQuestions(streamingText.value, /*hidePartial*/ true)
        : streamingText.value
      aiPlaceholder.content = visible
      scrollToBottom()
    })

    // 5) 流结束：抽取两条深挖（仅苏格拉底），并从显示内容中彻底移除该区块
    if (isSocrates.value) {
      const qs = extractDeepQuestions(streamingText.value)
      deepQuestions.value = qs.slice(0, 2)
      aiPlaceholder.content = removeDeepQuestions(streamingText.value, /*hidePartial*/ false)
    } else {
      aiPlaceholder.content = streamingText.value
    }

    // 6) 存一把（确保存的是“去掉深挖块”的可视内容）
    saveChat(id, { messages: messages.value })

    // 7) 自动 TTS：仅朗读正文，分段连续播放并合并提供下载
    if (chat.settings.voiceEnabled) {
      const main = extractMainText(aiPlaceholder.content)
      if (main) {
        await ttsPlaySegments(main, chat.settings.voiceStyle || 'style2')
      }
    }
  } catch (e) {
    alert('LLM 生成失败：' + (e?.message || e))
  } finally {
    sending.value = false
    scrollToBottom()
  }
}

/** ===== 深挖块抽取/移除 ===== */
function extractDeepQuestions (content) {
  const dq = []
  if (!content) return dq
  const start = content.indexOf('[DEEP_QUESTIONS]')
  const end = content.indexOf('[END]')
  if (start === -1 || end === -1 || end <= start) return dq

  const block = content.slice(start + '[DEEP_QUESTIONS]'.length, end).trim()
  block.split('\n').forEach(line => {
    const l = line
      .replace(/^\s*[-*]\s*/, '')   // 列表符号
      .replace(/^\s*\d+\.\s*/, '')  // 编号
      .trim()
    if (l) dq.push(l)
  })
  return dq
}

/**
 * 移除苏格拉底回复中的 [DEEP_QUESTIONS] 块：
 * - hidePartial=true：若只有起始标记，无结束标记，则从起始处开始都隐藏（用于流式过程）
 * - hidePartial=false：需要完整移除成品块（用于流式完成后）
 */
function removeDeepQuestions (content, hidePartial = true) {
  if (!content) return ''
  const start = content.indexOf('[DEEP_QUESTIONS]')
  if (start === -1) return content
  const end = content.indexOf('[END]', start)
  if (end === -1) {
    return hidePartial ? content.slice(0, start).trim() : content
  }
  const head = content.slice(0, start).trimEnd()
  const tail = content.slice(end + '[END]'.length).trimStart()
  return (head + (head && tail ? '\n\n' : '') + tail).trim()
}

/** ===== 朗读：过滤正文 + 分段合成并顺序播放 ===== */
function extractMainText (content) {
  // 现在消息里已无 [DEEP_QUESTIONS] 块，但依然做一次兜底处理
  if (!content) return ''
  const i = content.indexOf('[DEEP_QUESTIONS]')
  const main = (i === -1 ? content : content.slice(0, i)).trim()
  return main
    .replace(/```[\s\S]*?```/g, '')       // 去代码块
    .replace(/^>.*$/gm, '')               // 去引用
    .replace(/^\s*\*\*.*\*\*\s*$/gm, '')  // 去纯粗体行
    .trim()
}

function extractParagraphs (text) {
  let t = text
  const paras = t.split(/\n{2,}/).map(s => s.trim()).filter(s => s && s.length > 2)
  const final = []
  for (const p of paras) {
    if (p.length <= 220) { final.push(p); continue }
    const parts = p.split(/(?<=[。！？.!?])/)
    let buf = ''
    for (const part of parts) {
      if ((buf + part).length > 220) { final.push(buf.trim()); buf = part }
      else buf += part
    }
    if (buf.trim()) final.push(buf.trim())
  }
  return final.slice(0, 24) // 最多 24 段，避免过长
}

async function ttsPlaySegments (text, style) {
  speaking.value = true
  segments.value = []
  mergedAudioUrl.value = ''
  const paras = extractParagraphs(text)
  const blobs = []
  for (const p of paras) {
    try {
      const blob = await synthesize({ text: p, style })
      const url = URL.createObjectURL(blob)
      segments.value.push({ text: p, blob, url })
      blobs.push(blob)
      // 顺序播放
      await new Promise((resolve) => {
        const audio = new Audio(url)
        audio.onended = resolve
        audio.onerror = resolve
        audio.play().catch(() => resolve())
      })
    } catch (e) {
      console.warn('TTS segment failed', e)
    }
  }
  try {
    if (blobs.length) {
      const merged = await mergeWaveBlobs(blobs)
      mergedAudioUrl.value = URL.createObjectURL(merged)
    }
  } catch (e) { console.warn('merge failed', e) }
  speaking.value = false
}

/** 点击“朗读”按钮：对该条消息做同样的正文过滤 + 分段播报 */
async function speak (rawContent) {
  try {
    const main = extractMainText(rawContent || '')
    if (!main) return
    await ttsPlaySegments(main, chat.settings.voiceStyle || 'style2')
  } catch (e) {
    alert('TTS 失败：' + e?.message)
  }
}

/** 语音转文字 */
async function onAudioDone (file) {
  try {
    const data = await asrOnce(file)
    const text = data?.text || ''
    if (text) input.value = (input.value ? (input.value + ' ') : '') + text
  } catch (e) {
    alert('ASR 失败：' + e?.message)
  }
}

/** 点选深挖问题，回填输入框 */
function fillQuestion (q) { input.value = q }

/** 生命周期 */
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

.inline-actions{ display: flex; gap: 8px; margin-top: 6px; }

.typing{
  display: inline-block; position: relative;
}
.typing-dot{
  display:inline-block;width:6px;height:6px;border-radius:50%;
  background:#c0c4ce;margin-right:4px;animation:blink 1.2s infinite ease-in-out;
}
.typing-dot:nth-child(2){ animation-delay: .2s }
.typing-dot:nth-child(3){ animation-delay: .4s }
@keyframes blink { 0%, 80%, 100% { opacity:.2 } 40% { opacity: 1 } }

/* 苏格拉底深挖面板（紧贴输入框上方） */
.dq-panel{
  position: sticky;
  bottom: 0;
  z-index: 4;
  background: #fff;
  border-top: 1px solid var(--line);
  padding: .5rem .75rem;
  display: flex;
  flex-wrap: wrap;
  gap: .5rem;
  align-items: center;
}
.dq-title{ font-size: 12px; color: var(--muted); margin-right: .5rem; }
.dq-list{ display: flex; gap: .5rem; flex-wrap: wrap; }
.dq-item{
  border: none; background: #eef2ff; color: #312e81;
  border-radius: 999px; padding: 2px 10px; font-size: 12px; cursor: pointer;
}

/* 输入区（吸底） */
.composer{
  position: sticky; bottom: 0; z-index: 3;
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

/* Markdown 基本样式 */
.content.markdown-body h1,
.content.markdown-body h2,
.content.markdown-body h3 {
  margin: .4rem 0 .25rem;
  font-weight: 600;
}
.content.markdown-body pre {
  background: #f6f8fa;
  padding: .5rem .75rem;
  border-radius: 6px;
  overflow: auto;
}
.content.markdown-body code {
  background: #f6f8fa;
  padding: 0 .25rem;
  border-radius: 4px;
}
.content.markdown-body ul {
  padding-left: 1.2rem;
}
</style>
