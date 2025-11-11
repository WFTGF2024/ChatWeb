
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
:root{
  /* 保底变量，外层没有也能用 */
  --surface: #ffffff;
  --line: rgba(226, 232, 240, .9);
  --shadow-sm: 0 4px 18px rgba(15, 23, 42, .08);
  --muted: #6b7280;
}

.role-picker{
  position: relative;
}

/* 触发按钮 */
.chip{
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #f3f4f6;
  border: 1px solid rgba(148, 163, 184, .2);
  padding: 6px 12px 6px 6px;
  border-radius: 999px;
  cursor: pointer;
  transition: background .12s ease, box-shadow .12s ease, border .12s ease;
  font-size: 13.5px;
  color: #111827;
}
.chip:hover{
  background: #e5edff;
  border-color: rgba(37, 99, 235, .35);
  box-shadow: 0 6px 20px rgba(15, 23, 42, .08);
}
.chip svg{
  opacity: .6;
}

/* 头像小圆 */
.avatar{
  width: 30px;
  height: 30px;
  border-radius: 999px;
  background: #fff;
  border: 1px solid rgba(148, 163, 184, .3);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  line-height: 1;
}

/* 下拉菜单 */
.menu{
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  min-width: 240px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 14px;
  box-shadow: var(--shadow-sm);
  padding: 6px;
  z-index: 100;
  backdrop-filter: saturate(150%);
}

/* 每一项 */
.menu-item{
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 6px 6px 6px 4px;
  border-radius: 10px;
  cursor: pointer;
  transition: background .12s ease;
}
.menu-item:hover{
  background: #eff6ff;
}
.menu-item .col{
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.menu-item b{
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
}
.menu-item .hint{
  font-size: 11.5px;
  color: var(--muted);
}

/* 小屏下避免菜单太宽 */
@media (max-width: 520px){
  .menu{
    min-width: 210px;
  }
  .chip{
    max-width: 180px;
  }
}
</style>
