<template>
  <div class="rg-page">
    <div class="rg-card">
      <div class="rg-head">
        <div>
          <h2 class="rg-title">创建账号</h2>
          <p class="rg-sub">填写基本信息与安全问题，完成后将跳转到登录页。</p>
        </div>
        <router-link to="/login" class="rg-btn rg-btn-ghost">去登录</router-link>
      </div>

      <form class="rg-form" @submit.prevent="submit">
        <!-- 基本信息 -->
        <div class="rg-section-title">基本信息</div>
        <div class="rg-grid">
          <div class="rg-field">
            <label class="rg-label">用户名 <span class="rg-req">*</span></label>
            <input
              class="rg-input"
              :class="{ 'rg-invalid': errors.username }"
              v-model.trim="f.username"
              placeholder="请输入用户名（3-20字符）"
              autocomplete="username"
            />
            <small v-if="errors.username" class="rg-err">{{ errors.username }}</small>
          </div>

          <div class="rg-field">
            <label class="rg-label">姓名</label>
            <input
              class="rg-input"
              v-model.trim="f.full_name"
              placeholder="姓名（可选）"
              autocomplete="name"
            />
          </div>

          <div class="rg-field">
            <label class="rg-label">邮箱</label>
            <input
              class="rg-input"
              :class="{ 'rg-invalid': errors.email }"
              v-model.trim="f.email"
              placeholder="name@example.com"
              autocomplete="email"
              inputmode="email"
            />
            <small v-if="errors.email" class="rg-err">{{ errors.email }}</small>
          </div>

          <div class="rg-field">
            <label class="rg-label">手机号</label>
            <input
              class="rg-input"
              :class="{ 'rg-invalid': errors.phone_number }"
              v-model.trim="f.phone_number"
              placeholder="仅数字或+区号"
              autocomplete="tel"
              inputmode="tel"
            />
            <small v-if="errors.phone_number" class="rg-err">{{ errors.phone_number }}</small>
          </div>

          <div class="rg-field rg-field-pw">
            <label class="rg-label">密码 <span class="rg-req">*</span></label>
            <div class="rg-pw-wrap">
              <input
                :type="showPwd ? 'text' : 'password'"
                class="rg-input"
                :class="{ 'rg-invalid': errors.password }"
                v-model="f.password"
                placeholder="至少8位，包含字母与数字"
                autocomplete="new-password"
              />
              <button type="button" class="rg-icon-btn" @click="showPwd = !showPwd" :aria-pressed="showPwd" title="显示/隐藏密码">
                {{ showPwd ? '🙈' : '👁️' }}
              </button>
            </div>
            <div class="rg-meter" aria-hidden="true">
              <div class="rg-meter-bar" :style="{ width: strength.percent + '%' }"></div>
            </div>
            <small class="rg-hint">强度：{{ strength.label }}</small>
            <small v-if="errors.password" class="rg-err">{{ errors.password }}</small>
          </div>
        </div>

        <!-- 安全问题 -->
        <div class="rg-section-title">安全问题</div>
        <div class="rg-grid">
          <div class="rg-field">
            <label class="rg-label">密保问题 1</label>
            <input class="rg-input" v-model.trim="f.security_question1" />
          </div>
          <div class="rg-field">
            <label class="rg-label">答案 1</label>
            <input class="rg-input" v-model.trim="f.security_answer1" />
          </div>

          <div class="rg-field">
            <label class="rg-label">密保问题 2</label>
            <input class="rg-input" v-model.trim="f.security_question2" />
          </div>
          <div class="rg-field">
            <label class="rg-label">答案 2</label>
            <input class="rg-input" v-model.trim="f.security_answer2" />
          </div>
        </div>

        <div class="rg-actions">
          <label class="rg-agree">
            <input type="checkbox" v-model="agree" />
            我已阅读并同意服务条款与隐私政策
          </label>

          <button class="rg-btn rg-btn-primary" type="submit" :disabled="submitting || !formValid">
            <span v-if="!submitting">注册</span>
            <span v-else>注册中…</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { register } from '../api/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const submitting = ref(false)
const showPwd = ref(false)
const agree = ref(true)

const f = reactive({
  username: '',
  password: '',
  full_name: '',
  email: '',
  phone_number: '',
  security_question1: '你的第一所学校？',
  security_answer1: '',
  security_question2: '你最喜欢的书？',
  security_answer2: ''
})

const errors = reactive({
  username: '',
  email: '',
  phone_number: '',
  password: ''
})

const strength = computed(() => {
  const p = f.password || ''
  let score = 0
  if (p.length >= 8) score += 1
  if (/[A-Z]/.test(p)) score += 1
  if (/[a-z]/.test(p)) score += 1
  if (/\d/.test(p)) score += 1
  if (/[^A-Za-z0-9]/.test(p)) score += 1
  const percent = Math.min(100, score * 20)
  const label = percent >= 80 ? '很强' : percent >= 60 ? '较强' : percent >= 40 ? '一般' : percent > 0 ? '较弱' : '无'
  return { percent, label }
})

function validate() {
  errors.username = ''
  errors.email = ''
  errors.phone_number = ''
  errors.password = ''

  if (!f.username || f.username.length < 3 || f.username.length > 20) {
    errors.username = '用户名长度需 3-20 位'
  }
  if (f.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(f.email)) {
    errors.email = '邮箱格式不正确'
  }
  if (f.phone_number && !/^\+?\d{5,20}$/.test(f.phone_number)) {
    errors.phone_number = '手机号格式不正确'
  }
  if (!f.password || f.password.length < 8 || !/[A-Za-z]/.test(f.password) || !/\d/.test(f.password)) {
    errors.password = '密码至少 8 位，且需包含字母与数字'
  }
  return Object.values(errors).every(v => !v)
}

const formValid = computed(() => validate() && agree.value)

async function submit() {
  if (!formValid.value) return
  submitting.value = true
  try {
    await register(f)
    router.push('/login')
  } catch (e) {
    const msg = e?.response?.data?.message || e.message || '未知错误'
    alert('注册失败：' + msg)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
/* ---- 基础与隔离 ---- */
.rg-page *,
.rg-page *::before,
.rg-page *::after {
  box-sizing: border-box;        /* 避免宽高计算导致的溢出重叠 */
}
.rg-page {
  --bg: #f6f7fb;
  --panel: #fff;
  --line: #e5e7eb;
  --muted: #6b7280;
  --primary: #428bff;
  --primary-weak: #e6efff;
  --danger: #e03131;
  --shadow: 0 8px 28px rgba(0,0,0,.06);

  min-height: calc(100vh - 32px);
  width: 100%;
  padding: 24px 12px;
  background: var(--bg);

  display: grid;
  place-items: start center;     /* 居中卡片，不会造成重叠 */
  overflow: auto;                 /* 超出时滚动，防止内容覆盖 */
}

.rg-card {
  width: 100%;
  max-width: 860px;
  background: var(--panel);
  border: 1px solid var(--line);
  border-radius: 16px;
  box-shadow: var(--shadow);
  padding: 18px 18px 20px;
  position: relative;            /* 形成定位上下文，避免子级绝对定位溢出 */
}

/* ---- 头部 ---- */
.rg-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  flex-wrap: wrap;                /* 小屏不重叠 */
}

.rg-title {
  margin: 0;
  font-size: 20px;
}
.rg-sub {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.5;
}

/* ---- 表单布局 ---- */
.rg-form {
  display: flex;
  flex-direction: column;
  gap: 16px;                      /* 更大间距，避免拥挤重叠 */
  width: 100%;
}

.rg-section-title {
  margin-top: 6px;
  font-weight: 600;
  font-size: 14px;
  color: #111827;
  padding-left: 8px;
  border-left: 3px solid var(--primary);
}

.rg-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  width: 100%;
}

.rg-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;                   /* 防止网格项内容挤压溢出 */
}

.rg-label {
  font-size: 13px;
  color: #111827;
  line-height: 1.3;
}

.rg-req { color: var(--danger); }

/* ---- 输入控件 ---- */
.rg-input {
  width: 100%;
  min-height: 40px;
  line-height: 1.4;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #fff;
  outline: none;
  transition: border .15s ease, box-shadow .15s ease;
}
.rg-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-weak);
}
.rg-invalid {
  border-color: #ffb2b2 !important;
  background: #fff8f8 !important;
}

.rg-err {
  color: var(--danger);
  font-size: 12px;
  line-height: 1.3;
}
.rg-hint {
  color: var(--muted);
  font-size: 12px;
}

/* ---- 密码字段 ---- */
.rg-field-pw .rg-pw-wrap {
  position: relative;
  display: grid;
}
.rg-icon-btn {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  border: 0;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 4px;
}
.rg-meter {
  margin-top: 6px;
  height: 6px;
  background: #f1f5f9;
  border-radius: 6px;
  overflow: hidden;
}
.rg-meter-bar {
  height: 100%;
  width: 0%;
  background: linear-gradient(90deg, #ff6b6b, #ffd166, #06d6a0);
  transition: width .25s ease;
}

/* ---- 底部动作 ---- */
.rg-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding-top: 10px;
  border-top: 1px dashed var(--line);
  margin-top: 2px;
  flex-wrap: wrap;                /* 小屏换行，避免覆盖 */
}
.rg-agree {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  color: #374151;
  font-size: 13px;
}

/* ---- 按钮 ---- */
.rg-btn {
  padding: 9px 16px;
  border-radius: 10px;
  border: 1px solid var(--line);
  background: #fff;
  cursor: pointer;
  transition: transform .15s ease, box-shadow .15s ease, background .15s ease, color .15s ease, border-color .15s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;            /* 避免文字换行挤压重叠 */
}
.rg-btn:hover { transform: translateY(-1px); box-shadow: var(--shadow); }
.rg-btn-ghost { background: transparent; }
.rg-btn-primary {
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}
.rg-btn-primary[disabled] {
  opacity: .6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* ---- 小屏单列 ---- */
@media (max-width: 768px) {
  .rg-grid { grid-template-columns: 1fr; }
  .rg-head { align-items: flex-start; }
  .rg-actions { flex-direction: column; align-items: stretch; }
  .rg-btn, .rg-btn-primary { width: 100%; justify-content: center; }
}
</style>
