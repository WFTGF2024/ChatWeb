<template>
  <div class="bubble" :class="whoClass">
    <div class="avatar">{{ avatar }}</div>
    <div class="content">
      <slot />
      <div class="extra">
        <slot name="extra" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  who: { type: String, default: 'ai' },     // 'user' | 'ai' | 其他
  avatar: { type: String, default: '🤖' }
})

const whoClass = computed(()=>{
  // 统一左/右侧风格：用户在右，助手在左，其他按助手处理
  return props.who === 'user' ? 'right' : 'left'
})
</script>

<style scoped>
.bubble{
  display:flex; gap:10px; align-items:flex-start; margin:6px 0;
}
.bubble.left { flex-direction: row; }
.bubble.right{ flex-direction: row-reverse; }

.avatar{
  width:32px; height:32px; border-radius:50%;
  display:flex; align-items:center; justify-content:center;
  background: rgba(0,0,0,.04);
}

.content{
  max-width: min(760px, 80vw);
  padding: 10px 12px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 1px 2px rgba(0,0,0,.05);
  color: #333;
  word-break: break-word;
}

.bubble.right .content{ background: #f1f7ff; }
.extra{ margin-top:6px; display:flex; gap:6px; align-items:center; opacity:.9; }
</style>
