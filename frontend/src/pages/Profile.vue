<template>
  <div class="card col account-card">
    <!-- 顶部：标题 + 登录/退出 -->
    <div class="row head">
      <b>账户</b>
      <div class="row">
        <router-link class="btn ghost" to="/login" v-if="!user.isLogin">登录</router-link>
        <button class="btn danger" @click="logout" v-else>退出</button>
      </div>
    </div>

    <!-- 未登录提示 -->
    <div v-if="!user.isLogin" class="empty">
      当前为未登录模式，仅支持临时体验。请登录查看完整账户信息。
    </div>

    <!-- 已登录时展示结构化字段 -->
    <div v-else class="body">
      <div class="section-title">基本信息</div>
      <div class="grid">
        <div class="field">
          <span class="label">用户ID</span>
          <span class="value">{{ u.user_id ?? '-' }}</span>
        </div>
        <div class="field">
          <span class="label">账号（username）</span>
          <span class="value">{{ u.username ?? '-' }}</span>
        </div>
        <div class="field">
          <span class="label">姓名</span>
          <span class="value">{{ u.full_name ?? '-' }}</span>
        </div>
      </div>

      <div class="section-title">联系信息</div>
      <div class="grid">
        <div class="field">
          <span class="label">邮箱</span>
          <span class="value">{{ u.email ?? '-' }}</span>
        </div>
        <div class="field">
          <span class="label">手机号</span>
          <span class="value">{{ u.phone_number ?? '-' }}</span>
        </div>
      </div>

      <div class="section-title">系统信息</div>
      <div class="grid">
        <div class="field">
          <span class="label">创建时间</span>
          <span class="value">{{ formatDate(u.created_at) }}</span>
        </div>
        <div class="field">
          <span class="label">更新时间</span>
          <span class="value">{{ formatDate(u.updated_at) }}</span>
        </div>
      </div>

      <!-- 原始 JSON 也保留，方便调试 -->
      <details class="raw">
        <summary>查看原始 JSON</summary>
        <pre>{{ JSON.stringify(u, null, 2) }}</pre>
      </details>
    </div>

    <small class="hint">
      免登录体验：不保存历史与个性化；登录后开启聊天记录保存、会员与订单。
    </small>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from '../store/user'

const user = useUserStore()

// 把 store 里的 user 拿出来，防止没登录时报错
const u = computed(() => user.user || {})

function logout () {
  user.logout()
}

// 简单时间格式化，服务端给的是 ISO 字符串
function formatDate (v) {
  if (!v) return '-'
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return v   // 遇到奇怪格式就原样返回
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<style scoped>
.account-card{
  --panel: #ffffff;
  --line: #e5e7eb;
  --bg: #f7f8fc;
  --muted: #6b7280;
  --primary: #2563eb;
  --danger: #e03131;
  --shadow: 0 10px 30px rgba(0,0,0,.05);

  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 16px;
  box-shadow: var(--shadow);
  padding: 16px 16px 14px;
  margin: 12px auto;
  max-width: 640px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 顶部行 */
.head{
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--line);
}
.head b{
  font-size: 18px;
  letter-spacing: .1px;
}

/* 通用行 */
.row{
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

/* 未登录状态 */
.empty{
  background: #f3f4f6;
  border: 1px dashed rgba(203,213,225,.8);
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 13.5px;
  color: #374151;
}

/* 内容主体 */
.body{
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 小标题 */
.section-title{
  font-weight: 600;
  font-size: 13px;
  color: #111827;
  margin-top: 4px;
}

/* 两列布局 */
.grid{
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 12px;
}

/* 字段块 */
.field{
  background: #f9fafb;
  border: 1px solid rgba(229, 231, 235, .6);
  border-radius: 10px;
  padding: 7px 10px 6px;
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-height: 58px;
}
.label{
  font-size: 12px;
  color: var(--muted);
}
.value{
  font-size: 14px;
  color: #111827;
  word-break: break-all;
}

/* 原始 JSON 区域 */
.raw{
  margin-top: 6px;
}
.raw summary{
  cursor: pointer;
  font-size: 12.5px;
  color: #2563eb;
  user-select: none;
}
.raw pre{
  background: #f3f4f6;
  border: 1px dashed var(--line);
  border-radius: 12px;
  padding: 10px 12px;
  margin-top: 8px;
  max-height: 280px;
  overflow: auto;
  font-size: 13px;
  line-height: 1.55;
  color: #111827;
}

/* 提示信息 */
.hint{
  color: var(--muted);
  display: inline-block;
  margin-top: 2px;
  font-size: 12.5px;
}

/* 按钮 */
.btn{
  padding: 7px 11px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  user-select: none;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.btn:hover{
  transform: translateY(-1px);
  box-shadow: var(--shadow);
}
.btn.ghost{
  background: transparent;
  color: #374151;
}
.btn.danger{
  background: #fff5f5;
  color: var(--danger);
  border-color: #ffd6d6;
}
.btn.danger:hover{
  background: #ffecec;
}

/* 响应式：小屏改一列 */
@media (max-width: 540px){
  .grid{
    grid-template-columns: 1fr;
  }
  .btn{
    width: 100%;
    justify-content: center;
  }
  .account-card{
    padding: 14px 12px;
  }
}
</style>
