import { httpCore } from './http'
import { setCookie, eraseCookie } from '../utils/cookies'
import { useUserStore } from '../store/user'
export async function register(payload){
  const { data } = await httpCore.post('/api/auth/register', payload)
  return data
}

export async function login(payload) {
  const { data } = await httpCore.post('/api/auth/login', payload)
  const token = data?.token || data?.data?.token
  if (!token) throw new Error('Missing token in login response')
  localStorage.setItem('token', token)
  // also set cookies for compatibility
  try{ eraseCookie('jwt'); eraseCookie('token'); }catch{}
  setCookie('jwt', token)
  setCookie('token', token)
  setCookie('Authorization', 'Bearer ' + token)
  // 更新用户状态
  const user = useUserStore()
  user.setAuth(token, data?.user || data?.data?.user)
  return data
}


export async function me(){
  const { data } = await httpCore.get('/api/auth/me')
  return data
}

export async function updateUser(user_id, payload){
  const { data } = await httpCore.put(`/api/users/${user_id}`, payload)
  return data
}