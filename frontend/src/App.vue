<template>
  <div class="app-shell">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="row left">
        <span class="badge">{{ appName }}</span>
        <span class="faint">莫兰迪配色 · 统一语音交互平台</span>
      </div>
      <div class="row right">
        <router-link class="btn ghost" to="/chat">聊天</router-link>
        <!-- 这里把 角色库 去掉了 -->
        <router-link class="btn ghost" to="/membership">会员</router-link>
        <router-link class="btn ghost" to="/web-search">搜索</router-link>
                <router-link class="btn ghost" to="/profile">账户</router-link>
      </div>
    </header>

    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
const appName = import.meta.env.VITE_APP_NAME || '软件设计实践周项目'
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f6;
}
.header {
  height: 56px;
  padding: 0 16px;
  background: #fff;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.badge {
  background: #4b5563;
  color: #fff;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}
.faint {
  color: #888;
  font-size: 13px;
}
.btn {
  border: none;
  background: transparent;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 13px;
  text-decoration: none;
  color: #374151;
  transition: background .15s ease;
}
.btn.ghost.router-link-active,
.btn.ghost:hover {
  background: rgba(0,0,0,0.04);
}
.main {
  flex: 1;
  padding: 16px;
}
</style>
