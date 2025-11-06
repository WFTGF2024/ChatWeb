<template>
  <div class="chat-wrap">
    <ChatHeader :role="role">
      <div class="row">
        <!-- 多会话：选择 / 新建 / 删除 -->
        <select v-model="chatId" class="select" style="margin-right:8px" v-if="chatList.length">
          <option v-for="c in chatList" :key="c.id" :value="c.id">{{ c.title }}</option>
        </select>
        <button class="btn ghost" @click="newChat" title="新建会话">＋新建</button>
        <button class="btn ghost danger" :disabled="!chatId" @click="removeChat" title="删除当前会话">🗑 删除</button>

        <!-- ① 角色风格 → 绑定音色 -->
        <select v-model="selectedVoicePreset" class="select" style="margin-left:8px; max-width:160px;">
          <option v-for="p in voicePresets" :key="p.id" :value="p.id">{{ p.label }}</option>
        </select>

        <button class="btn ghost" @click="toggleVoice" title="是否自动播放TTS">
          {{ settings.voiceEnabled ? '🔊 自动播放开' : '🔇 自动播放关' }}
        </button>
        <button class="btn ghost" @click="exportChat">导出</button>
        <button class="btn" :disabled="!canSave" @click="save">保存</button>
      </div>
    </ChatHeader>

    <LoginGate v-if="!isLogin" />

    <div class="chat-list">
      <MessageBubble
        v-for="(m,i) in messages"
        :key="m.ts ?? i"
        :who="m.role==='user' ? 'user' : 'ai'"
        :avatar="m.role==='user' ? '👤' : '🤖'"
      >
        <template #default>
          <div v-if="m.role==='assistant'" v-html="toHtml(m.content)"></div>
          <div v-else>{{ m.content }}</div>
        </template>
        <template #extra>
          <template v-if="m.role==='assistant'">
            <span v-if="m.audioUrl" style="margin-left:8px; opacity:.8;">WAV已生成</span>
            <button
              v-if="m.audioUrl && !isPlaying(m)"
              class="btn ghost"
              style="margin-left:8px"
              @click="play(m)"
            >▶ 播放</button>
            <button
              v-if="m.audioUrl && isPlaying(m)"
              class="btn ghost"
              style="margin-left:8px"
              @click="stop()"
            >■ 停止</button>
            <button
              v-if="m.audioUrl"
              class="btn ghost"
              style="margin-left:6px"
              @click="downloadFromUrl(m.audioUrl, `tts_${m.ts||Date.now()}.wav`)"
            >⬇️ 下载</button>
          </template>
          <a v-if="m.audioUrl" :href="m.audioUrl" target="_blank" style="margin-left:6px;">打开</a>
        </template>
      </MessageBubble>
    </div>

    <DeepQuestionChips :items="deepQuestions" @pick="useQuestion" />

    <div class="chat-input">
      <textarea
        v-model="text"
        class="input"
        rows="3"
        placeholder="说点什么……"
        @keyup.enter.exact.prevent="send()"
      ></textarea>
      <div class="row" style="justify-content:space-between; gap:8px; margin-top:6px;">
        <AudioRecorder @done="useASR" />
        <label class="row" style="gap:6px; align-items:center;">
          <input type="checkbox" v-model="autoSendASR" /> 语音识别后自动发送
        </label>
        <div class="row" style="gap:6px;">
          <button class="btn primary" :disabled="pending" @click="send">发送</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import { useChatStore } from '../store/chat'
import { useUserStore } from '../store/user'
import { buildSystemPrompt } from '../utils/prompts'
import { chatStream, chatOnce } from '../api/llm'
import { asrFull } from '../api/asr'
import { synthesizeTTS } from '../api/tts'   // 走 /tts 代理
import AudioRecorder from '../components/AudioRecorder.vue'
import MessageBubble from '../components/MessageBubble.vue'
import DeepQuestionChips from '../components/DeepQuestionChips.vue'
import LoginGate from '../components/LoginGate.vue'
import ChatHeader from '../components/ChatHeader.vue'
import MarkdownIt from 'markdown-it'
import {
  listChats, createChat as createSession, deleteChat as deleteSession,
  loadChat as loadSession, saveChat as saveSession, renameChat
} from '../utils/chatCache'

const chat = useChatStore()
const user = useUserStore()

const text = ref('')
const autoSendASR = ref(true)

// ① 角色风格预设：一个风格绑定一个参考音色
const voicePresets = [
  { id: 'neutral', label: '通用助手', ttsStyle: 'style1', emoWeight: 0.65 },
  { id: 'interview', label: '面试官访谈', ttsStyle: 'style2', emoWeight: 0.55 },
  { id: 'story', label: '故事/角色', ttsStyle: 'style3', emoWeight: 0.8 }
]
const selectedVoicePreset = ref('neutral')

const role = computed(() => chat.currentRole)
const messages = computed(() => chat.messages)
const deepQuestions = computed(() => chat.deepQuestions)
const settings = chat.settings
const isLogin = computed(() => user.isLogin)
const pending = computed(() => chat.pending)
const canSave = computed(() => chat.messages.length > 0)

/* Markdown 渲染 */
const md = new MarkdownIt({ html: false, linkify: true, breaks: true })
const toHtml = (t) => md.render(t || '')

/* 简单播放器 */
const currentAudio = ref(null)
const currentUrl = ref('')
function isPlaying (m) {
  return !!currentAudio.value && currentUrl.value === m.audioUrl && !currentAudio.value.paused
}
function play (m) {
  try {
    if (!m?.audioUrl) return
    stop()
    currentUrl.value = m.audioUrl
    currentAudio.value = new Audio(m.audioUrl)
    currentAudio.value.onended = () => { currentAudio.value = null; currentUrl.value = '' }
    currentAudio.value.play().catch(() => {})
  } catch {}
}
function stop () {
  try {
    if (currentAudio.value) {
      currentAudio.value.pause()
      currentAudio.value.currentTime = 0
    }
  } catch {}
  currentAudio.value = null
  currentUrl.value = ''
}
function downloadFromUrl (url, filename) {
  const a = document.createElement('a')
  a.href = url
  a.download = filename || `tts_${Date.now()}.wav`
  document.body.appendChild(a)
  a.click()
  a.remove()
}

/* ====== 角色风格应用 & 保存到会话 ====== */
function applyVoicePreset (id) {
  const p = voicePresets.find(x => x.id === id)
  if (p) {
    chat.settings.ttsStyle = p.ttsStyle
    chat.settings.emoWeight = p.emoWeight
  }
}

watch(selectedVoicePreset, (v) => {
  applyVoicePreset(v)
  if (chatId.value) {
    saveSession(chatId.value, {
      messages: chat.messages,
      meta: { voicePreset: v }
    })
  }
})

/* ====== TTS 清洗 ====== */
function cleanTextForTTS (raw) {
  if (!raw) return ''
  let t = String(raw)
  // 去掉深度问题段
  t = t.replace(/\[DEEP_QUESTIONS[\s\S]*$/i, '')
  // 去掉 markdown 代码块
  t = t.replace(/```[\s\S]*?```/g, '')
  // 去掉行级标记
  t = t.replace(/^[-*+#>\s]+/gm, '')
  // 合并空白
  t = t.replace(/\s+/g, ' ')
  return t.trim()
}

async function doTTS (text, msgIndex) {
  text = cleanTextForTTS(text)
  if (!text) return
  try {
    const res = await synthesizeTTS({
      text,
      style: chat.settings.ttsStyle,
      emoWeight: chat.settings.emoWeight,
      format: 'wav'
    })
    let url = null
    if (res instanceof Blob) url = URL.createObjectURL(res)
    else if (res?.url) url = res.url
    else if (res?.blob) url = URL.createObjectURL(res.blob)
    if (!url) return

    if (chat.messages[msgIndex]) {
      chat.messages[msgIndex].audioUrl = url
      if (settings.voiceEnabled) play(chat.messages[msgIndex])
    }
  } catch (e) {
    console.warn('TTS 失败：', e)
    chat.addMessage({ role: 'assistant', content: `【系统】TTS失败：${e.message}`, ts: Date.now() })
  }
}

async function converse (userText) {
  chat.pending = true
  try {
    const system = buildSystemPrompt({
      role: role.value,
      memorySummary: chat.memorySummary,
      userPrefs: {}
    })
    const sysWithKB = chat.kbContext ? system + `\n\n【外部上下文，供参考】\n` + chat.kbContext : system

    let msgs = []
    // ③ 登录后，同一会话第二次开始带上下文
    if (user.isLogin) {
      const userMsgCount = chat.messages.filter(m => m.role === 'user').length
      if (userMsgCount >= 1) {
        msgs = [{ role: 'system', content: sysWithKB }]
        for (const m of chat.messages) {
          if (m.role === 'user' || m.role === 'assistant') {
            msgs.push({ role: m.role, content: m.content })
          }
        }
        msgs.push({ role: 'user', content: userText })
      } else {
        // 第一次提问：不带历史
        msgs = [{ role: 'system', content: sysWithKB }, { role: 'user', content: userText }]
      }
    } else {
      msgs = [{ role: 'system', content: sysWithKB }, { role: 'user', content: userText }]
    }

    let full = ''
    let aiIndex = -1

    if (settings.stream) {
      await chatStream({
        messages: msgs,
        onDelta: (delta) => {
          if (!full) {
            chat.addMessage({ role: 'assistant', content: delta, ts: Date.now() })
            full = delta
            aiIndex = chat.messages.length - 1
          } else {
            full += delta
            chat.messages[aiIndex].content = full
          }
        },
        onDone: async () => {
          if (aiIndex >= 0) await doTTS(full, aiIndex)
          save()
        }
      })
    } else {
      const content = await chatOnce(msgs)
      chat.addMessage({ role: 'assistant', content, ts: Date.now() })
      full = content
      const idx = chat.messages.length - 1
      await doTTS(full, idx)
      save()
    }

    // 解析深度问题
    const qs = parseDeepQuestions(full)
    chat.setDeepQuestions(qs)
  } catch (e) {
    chat.addMessage({ role: 'assistant', content: '【系统】对话失败：' + e.message, ts: Date.now() })
  } finally {
    chat.pending = false
  }
}

/* 语音识别 → 填到输入框 */
async function useASR (blob) {
  try {
    const txt = await asrFull(blob)
    text.value = txt
    if (autoSendASR.value && txt && txt.trim()) await send()
  } catch (e) {
    console.warn('ASR失败', e)
  }
}

/* 深度问题按钮注入 */
function useQuestion (q) {
  text.value = q
}

function parseDeepQuestions (content) {
  if (!content) return []
  const m = content.match(/\[DEEP_QUESTIONS\]([\s\S]*?)\[END\]/)
  if (!m) return []
  return m[1].split('\n').map(s => s.trim()).filter(Boolean)
}

/* ============ 多会话本地缓存 ============ */
const chatList = ref(listChats())
const chatId = ref(chatList.value[0]?.id || '')

onMounted(() => {
  if (!chatId.value) {
    const c = createSession(role.value?.name ? `${role.value.name} 的会话` : '新会话')
    chatList.value = listChats()
    chatId.value = c.id
  }
  const initData = loadSession(chatId.value) || { messages: [], meta: {} }
  chat.messages.splice(0, chat.messages.length, ...initData.messages)
  selectedVoicePreset.value = initData.meta?.voicePreset || 'neutral'
  applyVoicePreset(selectedVoicePreset.value)
})

watch(chatId, (id) => {
  if (!id) return
  const data = loadSession(id) || { messages: [], meta: {} }
  chat.messages.splice(0, chat.messages.length, ...data.messages)
  selectedVoicePreset.value = data.meta?.voicePreset || 'neutral'
  applyVoicePreset(selectedVoicePreset.value)
})

function newChat () {
  const c = createSession(role.value?.name ? `${role.value.name} 的会话` : '新会话')
  chatList.value = listChats()
  chatId.value = c.id
  chat.clear()
  // 重置成默认音色
  selectedVoicePreset.value = 'neutral'
  applyVoicePreset('neutral')
}
function removeChat () {
  if (!chatId.value) return
  deleteSession(chatId.value)
  chatList.value = listChats()
  chatId.value = chatList.value[0]?.id || ''
  const data = chatId.value ? loadSession(chatId.value) : { messages: [], meta: {} }
  chat.messages.splice(0, chat.messages.length, ...(data.messages || []))
  selectedVoicePreset.value = data.meta?.voicePreset || 'neutral'
  applyVoicePreset(selectedVoicePreset.value)
}
function save () {
  if (chatId.value) {
    saveSession(chatId.value, {
      messages: chat.messages,
      meta: { voicePreset: selectedVoicePreset.value }
    })
  }
  console.info('已保存到本地：', chatId.value)
}

function exportChat () {
  const payload = { role: role.value, messages: chat.messages, ts: Date.now() }
  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `chat-${role.value.id}-${Date.now()}.json`
  a.click()
  URL.revokeObjectURL(url)
}

function toggleVoice () { chat.settings.voiceEnabled = !chat.settings.voiceEnabled }

async function send () {
  const v = text.value.trim()
  if (!v) return
  chat.addMessage({ role: 'user', content: v, ts: Date.now() })
  text.value = ''
  await converse(v)
}
</script>

<style scoped>
.chat-wrap {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.chat-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 200px;
}
.chat-input {
  margin-top: 10px;
}
</style>
