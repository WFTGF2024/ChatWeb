
<template>
  <div class="col" style="gap:16px; max-width: 980px; margin: 0 auto;">
    <div class="row" style="justify-content:space-between; align-items:center;">
      <b>Web Search</b>
      <router-link class="btn ghost" to="/chat">返回聊天</router-link>
    </div>

    <div class="card col" style="gap:12px;">
      <div class="row" style="gap:8px; align-items:center;">
        <input class="input" v-model="q" placeholder="输入关键词，按回车搜索" @keyup.enter="go" style="flex:1" />
        <select class="select" v-model="mode" style="max-width:160px;">
          <option value="hybrid">Hybrid</option>
          <option value="semantic">Semantic</option>
          <option value="lexical">Lexical</option>
        </select>
        <input class="input" v-model.number="top_k" type="number" min="1" max="20" style="width:120px" />
        <button class="btn" @click="go" :disabled="!q || loading">搜索</button>
      </div>
      <small class="hint">该功能仅在此菜单页提供，聊天页面已移除联网搜索。</small>
    </div>

    <div class="card col" v-if="results.length">
      <div v-for="(r,idx) in results" :key="r.page_id" class="result-row">
        <div class="row" style="justify-content:space-between; align-items:center;">
          <div class="col">
            <a class="link" :href="r.url" target="_blank">{{ r.title || r.url }}</a>
            <small class="hint" style="margin-top:4px;">score={{ r.score?.toFixed?.(3) }} | {{ r.url }}</small>
          </div>
          <button class="btn ghost" @click="preview(r.page_id)">预览</button>
        </div>
        <p class="desc">{{ r.snippet }}</p>
      </div>
    </div>

    <div class="card col" v-if="previewPage">
      <div class="row" style="justify-content:space-between; align-items:center;">
        <b>页面预览</b>
        <button class="btn ghost" @click="previewPage=null">关闭</button>
      </div>
      <h3>{{ previewPage?.title }}</h3>
      <pre class="pre">{{ previewPage?.content }}</pre>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { webSearch, getPage } from '../api/websearch'

const q = ref('')
const mode = ref('hybrid')
const top_k = ref(6)
const results = ref([])
const previewPage = ref(null)
const loading = ref(false)

async function go(){
  if(!q.value) return
  loading.value = true
  try{
    const data = await webSearch({ q: q.value, top_k: top_k.value, mode: mode.value, alpha: 0.6 })
    results.value = data?.results || []
  } finally {
    loading.value = false
  }
}
async function preview(page_id){
  previewPage.value = await getPage(page_id)
}
</script>
<style scoped>
.result-row{ padding: 10px 0; border-bottom: 1px dashed rgba(0,0,0,.08); }
.result-row:last-child{ border-bottom: none; }
.desc{ color: var(--m-muted); margin: 8px 0 0; white-space: pre-wrap; }
.pre{ white-space: pre-wrap; background: rgba(0,0,0,.02); padding: 12px; border-radius: 8px; }
.link{ color: #1b6ac9; text-decoration: none; }
.link:hover{ text-decoration: underline; }
</style>
