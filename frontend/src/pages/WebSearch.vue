<template>
  <div class="wrap">
    <h2>Web Search</h2>

    <div class="card">
      <div class="row" style="gap:10px; align-items:center;">
        <label class="chip"><input type="checkbox" v-model="persistToDB"> 存入数据库</label>
        <input class="input flex1" v-model="q" placeholder="输入关键词，如：golang、LLM……"/>
        <select v-model="mode" class="select">
          <option value="semantic">Semantic</option>
          <option value="keyword">Keyword</option>
        </select>
        <input class="input w80" v-model.number="topK" type="number" min="1" placeholder="TopK"/>
        <button class="btn primary" :disabled="loading" @click="doSearch">
          {{ loading ? '搜索中…' : '搜索' }}
        </button>
      </div>

      <div class="row" style="gap:10px; margin-top:8px;">
        <input class="input flex1" v-model="urlsText" placeholder="可选：输入要抓取的URL，支持逗号/空格分隔，回车后会入库"/>
        <label class="row" style="gap:6px; align-items:center;">
          <input type="checkbox" v-model="fetchOn" />
          抓取外网并入库
        </label>
      </div>
    </div>

    <!-- 结果区 -->
    <div class="card" v-if="err">
      <div class="err">{{ err }}</div>
    </div>

    <div class="card" v-if="!loading && results.length===0 && !err">
      <div class="empty">
        <div>没有结果。</div>
        <div class="sub">建议：先在“URL”输入框里放一个网址（如 <code>https://go.dev/</code>），勾选“抓取外网并入库”，再搜索关键字。</div>
      </div>
    </div>

    <div class="list" v-if="results.length">
      <div class="item" v-for="r in results" :key="r.id || r.url">
        <div class="title">
          <a :href="r.url" target="_blank" rel="noreferrer">{{ r.title || r.url }}</a>
        </div>
        <div class="meta">
          <code>{{ r.url }}</code>
          <span v-if="r.score" class="score">score: {{ r.score }}</span>
        </div>
        <div class="snippet" v-if="r.snippet">{{ r.snippet }}</div>
        <div class="ops">
          <button class="btn ghost" @click="openPage(r)">查看全文</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { webAPI } from '../api/http'

const q = ref('')
const persistToDB = ref(false)
const urlsText = ref('')
const topK = ref(6)
const mode = ref('semantic')   // 目前后端未用到，只做保留
const fetchOn = ref(true)

const loading = ref(false)
const err = ref('')
const results = ref([])

function parseUrls(text){
  return (text || '')
    .split(/[\s,;]+/)
    .map(s => s.trim())
    .filter(s => /^https?:\/\//i.test(s))
}

async function doSearch(){
  // cache-first
  const cacheKey = `websearch:${mode.value}:${topK.value}:${q.value}`
  const cached = localStorage.getItem(cacheKey)
  if(cached){ try{ results.value = JSON.parse(cached); return }catch{} }

  err.value = ''
  results.value = []
  loading.value = true
  try{
    const payload = {
      q: q.value.trim(),
      urls: parseUrls(urlsText.value),
      fetch: !!fetchOn.value,
      top_k: Number(topK.value) || 6
    }
    // 关键：取 data.results（兼容 items/raw）
    const { data } = await webAPI.search(payload)
    results.value = data?.results || data?.items || data || []
  }catch(e){
    err.value = e?.response?.data?.error || e?.message || '搜索失败'
  }finally{
    loading.value = false
  }
}

async function openPage(r){
  try {
    const id = r.id
    if (!id) {
      // 没有 id 也允许直接打开 url
      if (r.url) window.open(r.url, '_blank')
      return
    }
    const { data } = await webAPI.page(id)
    const blob = new Blob([data.content || ''], { type:'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    window.open(url, '_blank')
    setTimeout(()=>URL.revokeObjectURL(url), 30000)
  } catch(e) {
    console.warn('get page failed', e)
  }
}
</script>

<style scoped>
.wrap { padding: 10px 12px; }
.card { background:#fff; border-radius:12px; padding:14px; box-shadow:0 4px 16px rgba(0,0,0,.05); margin:12px 0; }
.row { display:flex; }
.flex1 { flex:1; }
.input { flex:1; border:1px solid #e5e7eb; border-radius:8px; padding:10px 12px; outline:none; }
.input:focus{ border-color:#2563eb; box-shadow:0 0 0 3px rgba(37,99,235,.14); }
.select { border:1px solid #e5e7eb; border-radius:8px; padding:10px 12px; }
.w80 { width:80px; }
.btn { padding:10px 14px; border-radius:8px; border:none; cursor:pointer; }
.btn.primary { background:#2563eb; color:#fff; }
.btn.ghost { background:#f3f4f6; color:#111827; }
.tip { margin-top:8px; color:#6b7280; }
.err { color:#ef4444; }
.list { display:flex; flex-direction:column; gap:10px; }
.item { background:#fff; border-radius:12px; padding:12px 14px; box-shadow:0 2px 10px rgba(0,0,0,.04); }
.title a { font-weight:600; color:#111827; text-decoration:none; }
.title a:hover { text-decoration:underline; }
.meta { color:#6b7280; margin-top:4px; display:flex; gap:8px; align-items:center; }
.score { background:#eef2ff; color:#3730a3; padding:2px 6px; border-radius:6px; font-size:12px; }
.snippet { margin-top:6px; color:#374151; line-height:1.6; }
.ops { margin-top:8px; }
.empty { color:#6b7280; }
.sub { font-size:13px; margin-top:4px; }
.chip{ padding:4px 8px; border-radius:999px; border:1px solid #ddd; display:flex; gap:6px; align-items:center; }
</style>
