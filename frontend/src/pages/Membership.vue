<template>
  <div class="col">
    <!-- 顶部栏 -->
    <div class="row topbar">
      <b>会员中心</b>
      <router-link class="btn ghost" to="/profile">账号</router-link>
    </div>

    <div class="card col">
      <!-- 未登录 -->
      <div class="row" v-if="!user.isLogin">
        <small class="hint">登录后可查看会员信息与订单。</small>
        <router-link to="/login" class="btn primary">去登录</router-link>
      </div>

      <!-- 已登录 -->
      <template v-else>
        <!-- 操作栏 -->
        <div class="row actions">
          <button class="btn" :disabled="loading" @click="refresh">
            {{ loading ? '刷新中…' : '刷新' }}
          </button>
          <button class="btn secondary" :disabled="buying" @click="buy">
            {{ buying ? '下单中…' : '购买月度会员（¥29）' }}
          </button>
          <small v-if="toast" class="hint">{{ toast }}</small>
        </div>

        <!-- 当前会员 -->
        <div class="card">
          <div class="card-head">
            <b>当前会员</b>
            <span v-if="normMembership" class="badge" :class="badgeClass">{{ statusText }}</span>
          </div>

          <div v-if="loading" class="skeleton-block"></div>

          <template v-else>
            <div v-if="normMembership" class="kv">
              <div class="kv-row"><span>用户ID</span><span>{{ normMembership.user_id ?? '—' }}</span></div>
              <div class="kv-row"><span>状态</span><span>{{ statusText }}</span></div>
              <div class="kv-row"><span>开始日期</span><span>{{ fmtDate(normMembership.start_date) }}</span></div>
              <div class="kv-row"><span>到期日期</span><span>{{ fmtDate(normMembership.expire_date) }}</span></div>
              <div class="kv-row"><span>剩余天数</span><span>{{ daysLeftText }}</span></div>
            </div>
            <div v-else class="empty">
              尚未开通会员。
            </div>
          </template>
        </div>

        <!-- 最近订单 -->
        <div class="card">
          <b>最近订单</b>
          <div v-if="loading" class="skeleton-block"></div>
          <template v-else>
            <table v-if="normLatestOrder" class="table">
              <thead>
                <tr>
                  <th>订单号</th>
                  <th>金额(元)</th>
                  <th>时长(月)</th>
                  <th>支付方式</th>
                  <th>状态</th>
                  <th>下单时间</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>{{ normLatestOrder.id ?? '—' }}</td>
                  <td>{{ normLatestOrder.amount ?? '—' }}</td>
                  <td>{{ normLatestOrder.duration ?? '—' }}</td>
                  <td>{{ normLatestOrder.payment ?? '—' }}</td>
                  <td>{{ normLatestOrder.status ?? '—' }}</td>
                  <td>{{ fmtDateTime(normLatestOrder.created_at) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty">暂无最近订单。</div>
          </template>
        </div>

        <!-- 历史订单 -->
        <div class="card">
          <b>历史订单</b>
          <div v-if="loading" class="skeleton-block" style="height:120px;"></div>
          <template v-else>
            <table v-if="normOrders.length" class="table">
              <thead>
                <tr>
                  <th>订单号</th>
                  <th>金额(元)</th>
                  <th>时长(月)</th>
                  <th>支付方式</th>
                  <th>状态</th>
                  <th>下单时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="o in normOrders" :key="o.id">
                  <td>{{ o.id }}</td>
                  <td>{{ o.amount ?? '—' }}</td>
                  <td>{{ o.duration ?? '—' }}</td>
                  <td>{{ o.payment ?? '—' }}</td>
                  <td>{{ o.status ?? '—' }}</td>
                  <td>{{ fmtDateTime(o.created_at) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty">暂无订单记录。</div>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../store/user'
import {
  getMembershipByUser,
  listOrdersByUser,
  getLatestOrder,
  createOrder,
  createMembership
} from '../api/membership'

const user = useUserStore()

// 原始数据
const membership = ref(null)
const latestOrder = ref(null)
const orders = ref([])

const loading = ref(false)
const buying  = ref(false)
const toast   = ref('')

/** ———————— 规范化：兼容不同字段命名 ———————— */
function normalizeMembership(m) {
  if (!m) return null
  return {
    membership_id: m.membership_id ?? m.MembershipID ?? m.id ?? m.ID,
    user_id:       m.user_id       ?? m.UserID,
    start_date:    m.start_date    ?? m.StartDate,
    expire_date:   m.expire_date   ?? m.ExpireDate,
    status:        m.status        ?? m.Status
  }
}
function normalizeOrder(o) {
  if (!o) return null
  return {
    id:         o.order_id        ?? o.OrderID ?? o.id ?? o.ID,
    user_id:    o.user_id         ?? o.UserID,
    amount:     o.amount          ?? o.Amount,
    duration:   o.duration_months ?? o.DurationMonths,
    payment:    o.payment_method  ?? o.PaymentMethod,
    status:     o.status          ?? o.Status,
    created_at: o.purchase_date   ?? o.PurchaseDate ?? o.created_at ?? o.CreatedAt
  }
}

const normMembership   = computed(() => normalizeMembership(membership.value))
const normLatestOrder  = computed(() => normalizeOrder(latestOrder.value))
const normOrders       = computed(() => orders.value.map(normalizeOrder))

/** ———————— 展示用派生状态 ———————— */
const statusText = computed(() => {
  if (!normMembership.value) return '未开通'
  const s = (normMembership.value.status || '').toLowerCase()
  if (s === 'active') return '已开通'
  if (s === 'expired') return '已过期'
  return normMembership.value.status || '—'
})
const badgeClass = computed(() => {
  const s = statusText.value
  return s === '已开通' ? 'ok' : s === '已过期' ? 'warn' : ''
})
const daysLeftText = computed(() => {
  if (!normMembership.value?.expire_date) return '—'
  const left = daysLeft(normMembership.value.expire_date)
  return (left >= 0 ? `${left} 天` : `已过期 ${Math.abs(left)} 天`)
})

/** ———————— 工具函数 ———————— */
function fmtDate(d) {
  if (!d) return '—'
  try {
    const t = new Date(d)
    if (isNaN(t.getTime())) return '—'
    return t.toISOString().slice(0, 10)
  } catch { return '—' }
}
function fmtDateTime(d) {
  if (!d) return '—'
  try {
    const t = new Date(d)
    if (isNaN(t.getTime())) return '—'
    const iso = t.toISOString()
    return iso.slice(0, 10) + ' ' + iso.slice(11, 19)
  } catch { return '—' }
}
function daysLeft(expire) {
  try {
    const end = new Date(expire).getTime()
    const now = Date.now()
    return Math.floor((end - now) / (24 * 3600 * 1000))
  } catch { return NaN }
}

/** ———————— 视图动作 ———————— */
async function refresh() {
  if (!user.isLogin) { toast.value = '请先登录'; return }
  const uid = user.user?.user_id
  if (!uid) { toast.value = '找不到用户ID'; return }

  loading.value = true
  toast.value = ''
  try {
    // 只使用现有 API：/api/membership/:user_id
    const [m, lo, os] = await Promise.all([
      getMembershipByUser(uid).catch(() => null), // 404 → null
      getLatestOrder(uid).catch(() => null),
      listOrdersByUser(uid).catch(() => []),
    ])
    membership.value = m
    latestOrder.value = lo
    orders.value = Array.isArray(os) ? os : []
  } catch (e) {
    toast.value = `刷新失败：${e?.response?.data?.message || e.message || '未知错误'}`
  } finally {
    loading.value = false
  }
}

async function buy() {
  if (!user.isLogin) { toast.value = '请先登录'; return }
  const uid = user.user?.user_id
  if (!uid) { toast.value = '找不到用户ID'; return }

  buying.value = true
  toast.value = ''
  try {
    // 1) 先下单（已有 /api/membership/orders）
    const order = await createOrder({
      user_id: uid,
      duration_months: 1,
      amount: 29,
      payment_method: 'wechat'
    })
    toast.value = '下单成功：订单号 ' + (order?.order_id || order?.OrderID || order?.id || '')

    // 2) 直接创建会员（已有 /api/membership），便于联调验收；生产可改为等待支付回调入库
    const now = new Date()
    const expires = new Date(now.getTime() + 30*24*3600*1000)
    await createMembership({
      user_id: uid,
      status: 'active',
      start_date: now.toISOString(),
      expire_date: expires.toISOString()
    })

    // 3) 刷新视图
    await refresh()
  } catch (e) {
    toast.value = `下单/开通失败：${e?.response?.data?.message || e.message || '未知错误'}`
  } finally {
    buying.value = false
  }
}

onMounted(() => { if (user.isLogin) refresh() })
</script>

<style scoped>
/* 基础主题变量（仅当前组件生效） */
.col {
  --bg: #f6f7fb;
  --panel: #ffffff;
  --line: #e5e7eb;
  --muted: #6b7280;
  --primary: #428bff;
  --primary-weak: #e6efff;
  --secondary: #12b886;
  --secondary-weak: #e8fbf5;
  --danger: #e03131;
  --shadow: 0 6px 24px rgba(0, 0, 0, 0.06);

  max-width: 980px;
  margin: 0 auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  background: var(--bg);
  min-height: calc(100vh - 32px);
}

/* 顶部行 */
.topbar {
  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 14px;
  padding: 10px 12px;
  box-shadow: var(--shadow);
  justify-content: space-between;
  align-items: center;
}
.topbar b { font-size: 18px; }

/* 卡片容器 */
.card {
  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 16px;
  padding: 14px;
  box-shadow: var(--shadow);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 行/列通用 */
.row { display: flex; gap: 10px; align-items: center; }
.col { display: flex; flex-direction: column; gap: 10px; }

/* 操作栏 */
.actions { gap: 16px; align-items: center; flex-wrap: wrap; }

/* 提示文字 */
.hint { color: var(--muted); }

/* 按钮样式 */
.btn {
  padding: 8px 14px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
  user-select: none;
}
.btn:hover { transform: translateY(-1px); box-shadow: var(--shadow); }
.btn:disabled { opacity: .6; cursor: not-allowed; transform: none; box-shadow: none; }

.btn.ghost { background: transparent; }

.btn.secondary {
  background: var(--secondary);
  color: #fff;
  border-color: var(--secondary);
}
.btn.secondary:hover { background: #0fa276; border-color: #0fa276; }

.btn.primary {
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}
.btn.primary:hover { background: #2f79ea; border-color: #2f79ea; }

/* 表格 */
.table {
  width: 100%;
  border-collapse: collapse;
  border: 1px solid var(--line);
  border-radius: 12px;
  overflow: hidden;
  font-size: 14px;
}
.table th, .table td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--line);
  text-align: left;
}
.table thead { background: #fafafa; }
.table tbody tr:hover { background: #fbfdff; }

/* 空状态 */
.empty {
  color: var(--muted);
  background: #fafafa;
  border: 1px dashed var(--line);
  border-radius: 10px;
  padding: 12px;
}

/* KV 展示 */
.kv { display: grid; grid-template-columns: 160px 1fr; row-gap: 8px; column-gap: 12px; }
.kv-row { display: contents; }
.kv-row > span:first-child { color: var(--muted); }
.kv-row > span:last-child { color: #111827; }

/* 头部徽章 */
.card-head { display: flex; align-items: center; justify-content: space-between; }
.badge {
  font-size: 12px;
  padding: 3px 8px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
}
.badge.ok {
  color: #0d7a50;
  background: var(--secondary-weak);
  border-color: #bff0df;
}
.badge.warn {
  color: #9a3412;
  background: #fff7ed;
  border-color: #fed7aa;
}

/* 骨架屏 */
.skeleton-block {
  height: 52px;
  border-radius: 10px;
  background: linear-gradient(90deg, #f2f3f6 25%, #e9ecf3 37%, #f2f3f6 63%);
  background-size: 400% 100%;
  animation: shimmer 1.2s infinite;
}
@keyframes shimmer { 0%{background-position: 100% 0;} 100%{background-position: -100% 0;} }

/* RouterLink 按钮 */
a.btn {
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}

/* 小屏优化 */
@media (max-width: 720px) {
  .col { padding: 12px; gap: 12px; }
  .btn { padding: 7px 12px; }
  .row { flex-wrap: wrap; }
  .kv { grid-template-columns: 120px 1fr; }
}
</style>
