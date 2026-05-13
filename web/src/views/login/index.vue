<template>
  <div class="login-container">
    <n-card class="login-card" :bordered="false" shadow>
      <template #header>
        <div class="login-title">管理员登录</div>
      </template>

      <n-form ref="formRef" :model="form" :rules="rules" size="large">
        <n-form-item path="email" label="邮箱">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" maxlength="50" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" maxlength="20" />
        </n-form-item>
        <n-form-item path="captcha" label="验证码">
          <div class="captcha-wrapper">
            <n-input v-model:value="form.captcha" placeholder="请输入验证码" maxlength="4" />
            <img :src="captchaUrl" alt="验证码" class="captcha-img cursor-pointer" @click="refreshCaptcha" />
          </div>
        </n-form-item>
      </n-form>

      <n-button type="primary" block :loading="loading" @click="handleLogin">登录</n-button>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
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
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  background-size: cover;
  background-position: center;
}

.login-card {
  width: 480px;
  background: rgba(255, 255, 255, 0.98);
}

.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: 500;
  color: #333;
}

.captcha-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.captcha-wrapper :deep(.n-input) {
  flex: 1;
}

.captcha-img {
  width: 100px;
  height: 40px;
  border-radius: 4px;
  flex-shrink: 0;
}
</style>
