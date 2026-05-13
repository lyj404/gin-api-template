<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
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
            class="custom-input"
          >
            <template #prefix>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
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
            class="custom-input"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
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
              class="custom-input"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <svg class="input-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
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
          <span class="btn-text">登 录</span>
        </n-button>
      </n-form>

      <div class="login-footer">
        <span>© {{ new Date().getFullYear() }} Admin System. All rights reserved.</span>
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
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f5f7fa;
  padding: 20px;
}

.login-card {
  width: 420px;
  max-width: 100%;
  padding: 44px 40px 32px;
  background: #ffffff;
  border: 1px solid #eaecef;
  border-radius: 12px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04), 0 8px 24px rgba(15, 23, 42, 0.06);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #2563eb;
  border-radius: 14px;
  color: #fff;
}

.logo-icon svg {
  width: 28px;
  height: 28px;
}

.login-title {
  margin: 0 0 6px;
  font-size: 22px;
  font-weight: 600;
  color: #0f172a;
  letter-spacing: 0.3px;
}

.login-subtitle {
  margin: 0;
  font-size: 13px;
  color: #64748b;
  font-weight: 400;
}

.login-form {
  margin-bottom: 8px;
}

.custom-input :deep(.n-input__input-el) {
  height: 44px;
  font-size: 14px;
  color: #0f172a;
}

.custom-input :deep(.n-input__input-el::placeholder) {
  color: #94a3b8;
}

.custom-input :deep(.n-input) {
  background: #f8fafc !important;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.custom-input :deep(.n-input .n-input__border),
.custom-input :deep(.n-input .n-input__state-border) {
  border: 1px solid #e2e8f0 !important;
  border-radius: 8px !important;
}

.custom-input :deep(.n-input:hover .n-input__state-border) {
  border-color: #cbd5e1 !important;
}

.custom-input :deep(.n-input--focus) {
  background: #ffffff !important;
}

.custom-input :deep(.n-input--focus .n-input__state-border) {
  border-color: #2563eb !important;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12) !important;
}

.custom-input :deep(.n-input__prefix) {
  margin-right: 8px;
  color: #94a3b8;
}

.input-icon {
  width: 16px;
  height: 16px;
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
  height: 44px;
  border-radius: 8px;
  cursor: pointer;
  flex-shrink: 0;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  transition: border-color 0.2s ease;
}

.captcha-img:hover {
  border-color: #2563eb;
}

.login-btn {
  height: 46px !important;
  margin-top: 8px;
  border-radius: 8px !important;
  font-size: 15px !important;
  font-weight: 500 !important;
  letter-spacing: 3px;
  background: #2563eb !important;
  border: none !important;
  transition: background-color 0.2s ease !important;
}

.login-btn:hover {
  background: #1d4ed8 !important;
}

.login-btn:active {
  background: #1e40af !important;
}

.btn-text {
  font-weight: 500;
}

.login-footer {
  margin-top: 28px;
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
}

@media (max-width: 480px) {
  .login-card {
    width: 100%;
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
