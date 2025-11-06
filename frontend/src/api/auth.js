import { httpCore } from './http'

export async function login(payload) {
  const { data } = await httpCore.post('/api/auth/login', payload)
  // 后端可能返回 {token: "..."} 或 {data:{token:"..."}}
  const token = data?.token || data?.data?.token
  if (!token) throw new Error('Missing token in login response')
  localStorage.setItem('token', token)
  return data
}
