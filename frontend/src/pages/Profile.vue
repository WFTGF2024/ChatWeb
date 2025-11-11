<template>
  <div class="card col" style="max-width:640px;">
    <div class="row" style="justify-content:space-between;">
      <b>账户</b>
      <div class="row">
        <router-link class="btn ghost" to="/login" v-if="!user.isLogin">登录</router-link>
        <button class="btn danger" @click="logout" v-else>退出</button>
      </div>
    </div>
    <div v-if="user.isLogin">
      <pre style="white-space:pre-wrap;">{{ JSON.stringify(user.user, null, 2) }}</pre>
    </div>
    <small class="hint">免登录体验：不保存历史与个性化；登录后开启聊天记录保存、会员与订单。</small>
  </div>
</template>
<script setup>
import { useUserStore } from '../store/user'
const user = useUserStore()
function logout(){ user.logout() }
</script>
<style scoped>
/* 局部主题变量，仅作用于本组件 */
.card.col{
  --panel: #ffffff;
  --line: #e5e7eb;
  --bg: #f7f8fc;
  --muted: #6b7280;
  --primary: #2563eb;
  --danger: #e03131;
  --shadow: 0 10px 30px rgba(0,0,0,.06);

  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 16px;
  box-shadow: var(--shadow);
  padding: 16px;
  margin: 12px auto;
  gap: 12px;                 /* 让内部列有间距 */
}

/* 顶部行：标题 + 操作 */
.card.col > .row:first-child{
  align-items: center;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--line);
  margin-bottom: 6px;
}

/* 标题更清晰 */
.card.col > .row:first-child b{
  font-size: 18px;
  letter-spacing: .2px;
}

/* 通用行：间距 + 对齐 */
.row{
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;           /* 小屏不拥挤 */
}

/* 账号 JSON 展示区 */
pre{
  background: #fafafa;
  border: 1px dashed var(--line);
  border-radius: 12px;
  padding: 12px;
  margin: 10px 0 4px;
  max-height: 340px;
  overflow: auto;
  font-size: 13.5px;
  line-height: 1.55;
  color: #111827;
}

/* 提示信息 */
.hint{
  color: var(--muted);
  display: inline-block;
  margin-top: 2px;
}

/* 按钮：在本组件内定义不影响全局（scoped） */
.btn{
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  user-select: none;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
  text-decoration: none;     /* 兼容 router-link 的 a.btn */
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.btn:hover{
  transform: translateY(-1px);
  box-shadow: var(--shadow);
}
.btn:active{
  transform: translateY(0);
  box-shadow: none;
}

/* 幽灵按钮（登录） */
.btn.ghost{
  background: transparent;
  color: #374151;
}

/* 危险按钮（退出） */
.btn.danger{
  background: #fff5f5;
  color: var(--danger);
  border-color: #ffd6d6;
}
.btn.danger:hover{
  background: #ffecec;
}

/* 小屏优化 */
@media (max-width: 640px){
  .card.col{ padding: 14px; border-radius: 14px; }
  pre{ max-height: 300px; font-size: 13px; }
  .btn{ width: 100%; justify-content: center; }
}
</style>
