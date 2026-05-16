<template>
  <div class="login-page">
    <div class="login-bg" />
    <div class="login-bg-grid" />
    <div class="login-container">
      <div class="login-card">
        <div class="login-card-inner">
          <div class="login-header">
            <div class="logo-icon">
              <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M16 3L3 9.5L16 16L29 9.5L16 3Z" stroke="currentColor" stroke-width="2.2" stroke-linejoin="round"/>
                <path d="M3 16.5L16 23L29 16.5" stroke="currentColor" stroke-width="2.2" stroke-linejoin="round"/>
                <path d="M3 23L16 29.5L29 23" stroke="currentColor" stroke-width="2.2" stroke-linejoin="round"/>
                <path d="M16 16V29.5" stroke="currentColor" stroke-width="2.2" stroke-linejoin="round"/>
              </svg>
            </div>
            <h1 class="login-title">欢迎回来</h1>
            <p class="login-subtitle">请登录您的管理员账户</p>
          </div>

          <n-form ref="formRef" :model="form" :rules="rules" size="large" class="login-form">
            <n-form-item path="email">
              <n-input
                v-model:value="form.email"
                placeholder="邮箱地址"
                maxlength="50"
                class="login-input"
              >
                <template #prefix>
                  <svg class="input-icon" viewBox="0 0 24 24" fill="none">
                    <path d="M4 4H20C21.1 4 22 4.9 22 6V18C22 19.1 21.1 20 20 20H4C2.9 20 2 19.1 2 18V6C2 4.9 2.9 4 4 4Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M22 6L12 13L2 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="password">
              <n-input
                v-model:value="form.password"
                type="password"
                placeholder="密码"
                show-password-on="click"
                maxlength="20"
                class="login-input"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <svg class="input-icon" viewBox="0 0 24 24" fill="none">
                    <rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M7 11V7C7 5.67 7.53 4.42 8.47 3.47C9.42 2.53 10.67 2 12 2C13.33 2 14.58 2.53 15.53 3.47C16.47 4.42 17 5.67 17 7V11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="captcha">
              <div class="captcha-wrapper">
                <n-input
                  v-model:value="form.captcha"
                  placeholder="验证码"
                  maxlength="4"
                  class="login-input"
                  @keyup.enter="handleLogin"
                >
                  <template #prefix>
                    <svg class="input-icon" viewBox="0 0 24 24" fill="none">
                      <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M21 12C21 16.97 16.97 21 12 21C7.03 21 3 16.97 3 12C3 7.03 7.03 3 12 3C16.97 3 21 7.03 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </template>
                </n-input>
                <img
                  :src="captchaUrl"
                  alt="验证码"
                  class="captcha-img"
                  title="点击刷新"
                  @click="refreshCaptcha"
                />
              </div>
            </n-form-item>

            <n-button
              type="primary"
              block
              size="large"
              :loading="loading"
              class="login-btn"
              @click="handleLogin"
            >
              登 录
            </n-button>
          </n-form>

          <div class="login-footer">
            &copy; {{ new Date().getFullYear() }} Admin System
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import type { FormRules } from 'naive-ui'

const router = useRouter()
const message = useMessage()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)

const form = ref({
  email: '',
  password: '',
  captcha: ''
})

const rules: FormRules = {
  email: { required: true, message: '请输入邮箱', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  captcha: { required: true, message: '请输入验证码', trigger: 'blur' }
}

const captchaUrl = ref('')

const refreshCaptcha = async () => {
  await auth.refreshCaptcha()
  captchaUrl.value = auth.captchaUrl
}

const handleLogin = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true
    await auth.login(form.value)
    message.success('登录成功')
    router.push('/')
  } catch (err: any) {
    if (err?.response?.data?.message) {
      message.error(err.response.data.message)
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--color-bg, #f7f4f0);
}

.login-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 80% 60% at 50% -10%, rgba(194, 112, 74, 0.08) 0%, transparent 60%),
    radial-gradient(ellipse 60% 50% at 80% 90%, rgba(194, 112, 74, 0.06) 0%, transparent 50%),
    radial-gradient(ellipse 50% 40% at 20% 80%, rgba(217, 119, 6, 0.04) 0%, transparent 50%);
  pointer-events: none;
}

.login-bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(194, 112, 74, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(194, 112, 74, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
  pointer-events: none;
  -webkit-mask-image: radial-gradient(ellipse 70% 60% at 50% 50%, black 30%, transparent 70%);
  mask-image: radial-gradient(ellipse 70% 60% at 50% 50%, black 30%, transparent 70%);
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: 20px;
  position: relative;
  z-index: 1;
  animation: fadeInUp 0.6s ease both;
}

.login-card {
  background: var(--color-surface, #ffffff);
  border-radius: 16px;
  box-shadow: var(--shadow-lg, 0 8px 30px rgba(28, 25, 23, 0.08));
  border: 1px solid var(--color-border-light, #efe9e2);
  position: relative;
  overflow: hidden;
}

.login-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #c2704a, #d97706, #c2704a);
  opacity: 0.6;
}

.dark .login-card::before {
  background: linear-gradient(90deg, #d97a4a, #fbbf24, #d97a4a);
}

.login-card-inner {
  padding: 40px 36px 28px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 18px;
  background: linear-gradient(135deg, #c2704a, #d97706);
  border-radius: 14px;
  color: #fff;
  animation: fadeInUp 0.6s ease 0.05s both;
}

.logo-icon svg {
  width: 28px;
  height: 28px;
}

.login-title {
  margin: 0 0 6px;
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text, #1c1917);
  letter-spacing: -0.3px;
  animation: fadeInUp 0.6s ease 0.1s both;
}

.login-subtitle {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-secondary, #78716c);
  font-weight: 400;
  animation: fadeInUp 0.6s ease 0.15s both;
}

.login-form {
  animation: fadeInUp 0.6s ease 0.2s both;
}

.login-form :deep(.n-form-item) {
  margin-bottom: 20px;
}

.login-form :deep(.n-form-item:last-child) {
  margin-bottom: 0;
}

.login-input :deep(.n-input-wrapper) {
  padding-left: 0;
}

.login-input :deep(.n-input__input-el) {
  height: 46px;
  font-size: 14px;
  color: var(--color-text, #1c1917);
}

.login-input :deep(.n-input__input-el::placeholder) {
  color: var(--color-text-muted, #a8a29e);
}

.login-input :deep(.n-input) {
  background: var(--color-bg, #f7f4f0) !important;
  border-radius: 10px;
  transition: all 0.25s ease;
}

.login-input :deep(.n-input .n-input__border),
.login-input :deep(.n-input .n-input__state-border) {
  border: 1.5px solid var(--color-border-light, #efe9e2) !important;
  border-radius: 10px !important;
  transition: all 0.25s ease;
}

.login-input :deep(.n-input:hover .n-input__state-border) {
  border-color: var(--color-primary-soft, rgba(194, 112, 74, 0.3)) !important;
}

.login-input :deep(.n-input--focus) {
  background: var(--color-surface, #ffffff) !important;
}

.login-input :deep(.n-input--focus .n-input__state-border) {
  border-color: var(--color-primary, #c2704a) !important;
  box-shadow: 0 0 0 3px rgba(194, 112, 74, 0.1) !important;
}

.login-input :deep(.n-input__prefix) {
  margin-right: 10px;
  color: var(--color-text-muted, #a8a29e);
}

.dark .login-input :deep(.n-input) {
  background: var(--color-surface, #26211e) !important;
}

.dark .login-input :deep(.n-input--focus) {
  background: #2d2824 !important;
}

.input-icon {
  width: 17px;
  height: 17px;
}

.captcha-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.captcha-wrapper :deep(.n-input) {
  flex: 1;
}

.captcha-img {
  width: 110px;
  height: 46px;
  border-radius: 10px;
  cursor: pointer;
  flex-shrink: 0;
  border: 1.5px solid var(--color-border-light, #efe9e2);
  background: var(--color-surface-2, #faf8f5);
  transition: border-color 0.25s ease, box-shadow 0.25s ease;
  object-fit: cover;
}

.captcha-img:hover {
  border-color: var(--color-primary, #c2704a);
  box-shadow: 0 0 0 3px rgba(194, 112, 74, 0.1);
}

.login-btn {
  height: 48px !important;
  margin-top: 24px;
  border-radius: 10px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
  letter-spacing: 2px;
  transition: all 0.25s ease !important;
  box-shadow: 0 2px 8px rgba(194, 112, 74, 0.25) !important;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(194, 112, 74, 0.35) !important;
}

.login-btn:active {
  transform: translateY(0);
}

.login-footer {
  margin-top: 28px;
  text-align: center;
  font-size: 12px;
  color: var(--color-text-muted, #a8a29e);
  animation: fadeInUp 0.6s ease 0.3s both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 480px) {
  .login-card-inner {
    padding: 32px 24px 24px;
  }

  .login-title {
    font-size: 20px;
  }

  .captcha-img {
    width: 96px;
  }
}
</style>
