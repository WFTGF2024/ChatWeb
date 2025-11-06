
<template>
  <div class="row" style="justify-content:space-between;">
    <div class="row" style="gap:8px; align-items:center;">
      <div class="role-picker">
        <button class="chip" @click="open = !open">
          <span class="avatar">{{ current.avatar }}</span>
          <span>{{ current.name }}</span>
          <svg viewBox="0 0 24 24" width="16" height="16"><path d="M7 10l5 5 5-5z" fill="currentColor"/></svg>
        </button>
        <div v-if="open" class="menu">
          <div class="menu-item" v-for="r in roles" :key="r.id" @click="choose(r)">
            <span class="avatar">{{ r.avatar }}</span>
            <div class="col">
              <b>{{ r.name }}</b>
              <small class="hint">{{ r.id }}</small>
            </div>
          </div>
        </div>
      </div>
      <small class="hint">可切换角色，自动映射到对应音色</small>
    </div>
    <slot />
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import roles from '../data/roles'
import { useChatStore } from '../store/chat'

const props = defineProps({ role: Object })
const chat = useChatStore()
const open = ref(false)
const current = computed(()=> chat.currentRole || roles[0])

const ROLE_TTS_STYLE = {
  'interviewer-java': 'style1',
  'interviewer-ml': 'style1',
  'assistant-tech': 'style2',
  'helper': 'style2',
  'harry': 'style3'
}

function choose(r){
  chat.currentRole = r
  chat.settings.ttsStyle = ROLE_TTS_STYLE[r.id] || 'style1'
  open.value = false
}
</script>
<style scoped>
.role-picker{ position: relative; }
.chip{ display:flex; align-items:center; gap:8px; background: rgba(0,0,0,.04); padding:8px 12px; border-radius: 999px; border: none; cursor: pointer; }
.avatar{ width: 22px; height: 22px; display:inline-flex; align-items:center; justify-content:center; }
.menu{ position:absolute; top: 110%; left: 0; min-width: 260px; background: white; border-radius: 12px; box-shadow: var(--shadow); padding: 8px; z-index: 100; }
.menu-item{ display:flex; gap:10px; align-items:center; padding:8px; border-radius:10px; cursor:pointer; }
.menu-item:hover{ background: rgba(0,0,0,.05); }
</style>
