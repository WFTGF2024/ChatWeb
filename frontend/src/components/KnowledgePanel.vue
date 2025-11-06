<template>
  <div class="col">
    <!-- 新增提示 -->
    <div
      class="alert"
      style="background:#fff4e6; border:1px solid #ffdda7; padding:6px 10px; border-radius:6px; margin-bottom:8px; color:#c06a00;"
    >
      该功能正在开发，敬请期待
    </div>

    <div class="row" style="justify-content:space-between;">
      <b>知识面板</b>
      <span class="faint">用于检索/注入上下文</span>
    </div>
    <div class="card col">
      <div class="row">
        <input class="input" v-model="url" placeholder="粘贴网页 URL，入库并切块" style="flex:1;" />
        <button class="btn" @click="ingest">入库</button>
      </div>
      <small class="hint">抓取 → 清洗 → 切块 → 向量化 → Qdrant 检索</small>
      <small class="hint" v-if="ingested">已入库：{{ ingested.title }}</small>
    </div>

    <div class="card col">
      <div class="row" style="gap:8px;">
        <input class="input" v-model="q" placeholder="检索问句" />
        <button class="btn" @click="search">搜索</button>
      </div>
      <div v-if="results.length" class="col" style="gap:4px; margin-top:6px;">
        <div
          v-for="r in results"
          :key="r.id || r.page_id"
          class="row"
          style="justify-content:space-between; gap:8px;"
        >
          <div style="max-width:240px;">
            <b>{{ r.title || '命中文档' }}</b>
            <div class="faint" style="font-size:12px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">
              {{ r.snippet || r.text }}
            </div>
          </div>
          <button class="btn ghost" @click="addToContext(r)">注入</button>
        </div>
      </div>
    </div>

    <div class="card col">
      <label>已经注入的片段：</label>
      <textarea class="input" rows="5" v-model="contextText"></textarea>
      <button class="btn" style="margin-top:6px;" @click="inject">应用到当前对话</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ingestUrl, webSearch } from '../api/websearch'
import { useChatStore } from '../store/chat'

const chat = useChatStore()
const url = ref('')
const ingested = ref(null)
const q = ref('')
const results = ref([])
const contextText = ref('')
const ctxPieces = ref([])

async function ingest () {
  if (!url.value) return
  ingested.value = await ingestUrl(url.value)
}
async function search () {
  if (!q.value) return
  const data = await webSearch({ q: q.value, top_k: 5, mode: 'hybrid', alpha: 0.6 })
  results.value = data.results || []
}
function addToContext (r) {
  ctxPieces.value.push(r)
  contextText.value = ctxPieces.value.map(x => x.text || x.snippet || '').join('\n')
}
function inject () {
  chat.setKbContext(contextText.value)
  ctxPieces.value = []
}
</script>
