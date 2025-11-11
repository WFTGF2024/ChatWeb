// src/api/membership.js
// 会员 & 订单 API —— 走 Core 后端
import { httpCore } from './http'

/* ==================== 订单 ==================== */

/**
 * 创建订单
 * @param {Object} payload
 * @param {number} payload.user_id
 * @param {number} payload.duration_months  购买时长（月）
 * @param {number} payload.amount           金额（元）
 * @param {string} [payload.payment_method='other'] 支付方式
 */
export async function createOrder({ user_id, duration_months, amount, payment_method = 'other' }) {
  const res = await httpCore.post('/api/membership/orders', {
    user_id,
    duration_months,
    amount,
    payment_method
  })
  return res.data
}

/** 按用户列出全部订单 */
export async function listOrdersByUser(user_id) {
  const res = await httpCore.get(`/api/membership/orders/${user_id}`)
  return res.data
}

/** 获取用户最新一笔订单 */
export async function getLatestOrder(user_id) {
  const res = await httpCore.get(`/api/membership/orders/${user_id}/latest`)
  return res.data
}

/** 获取用户最近 N 笔订单（默认 5） */
export async function listRecentOrders(user_id, n = 5) {
  const res = await httpCore.get(`/api/membership/orders/${user_id}/recent`, {
    params: { n }
  })
  return res.data
}


/* ==================== 会员 ==================== */

/** 列出所有会员记录（通常仅管理员使用） */
export async function listMembership() {
  const res = await httpCore.get('/api/membership')
  return res.data
}

/**
 * 获取某用户当前会员信息
 * 兼容两种后端风格：
 *  - GET /api/membership?user_id=xxx
 *  - GET /api/membership/:user_id
 * 若返回 404，视为“未开通”，函数返回 null（而不是抛错）
 */
export async function getMembershipByUser(user_id) {
  try {
    const res = await httpCore.get(`/api/membership/${user_id}`)
    return res.data
  } catch (e) {
    if (e?.response?.status === 404) return null // 没开通过会员 -> 返回 null
    throw e
  }
}

/**
 * 新建会员记录（是否允许由前端直接开通取决于后端策略）
 * @param {{user_id:number, start_date?:string, expire_date?:string, status?:string}} payload
 */
export async function createMembership(payload) {
  const res = await httpCore.post('/api/membership', payload)
  return res.data
}

/** 更新会员记录（需要 membership_id） */
export async function updateMembership(membership_id, payload) {
  const res = await httpCore.put(`/api/membership/${membership_id}`, payload)
  return res.data
}


/* ==================== 组合 helper（联调/验收友好） ==================== */

/**
 * 一次性下单并刷新视图所需数据
 * - openImmediately: 为 true 时，模拟“支付成功即开通会员”（需后端允许或提供测试接口）
 * - membershipPayload: 开通会员时的补充字段，如 { start_date, expire_date, status:'active' }
 */
export async function purchaseAndRefresh({
  user_id,
  duration_months,
  amount,
  payment_method = 'other',
  openImmediately = false,
  membershipPayload = {}
}) {
  const order = await createOrder({ user_id, duration_months, amount, payment_method })

  let membership = null
  if (openImmediately) {
    // 这里假设后端开放了直接创建/开通的接口用于联调
    membership = await createMembership({ user_id, status: 'active', ...membershipPayload })
  }

  const [latest, orders] = await Promise.all([
    getLatestOrder(user_id).catch(() => null),
    listOrdersByUser(user_id).catch(() => []),
  ])

  return { order, membership, latest, orders }
}

/* ==================== 认证（可选） ==================== */

/** /api/auth/me */
export async function me() {
  const res = await httpCore.get('/api/auth/me')
  return res.data
}

// 兼容旧命名（如果页面里用到了 getMembership）
export const getMembership = getMembershipByUser
