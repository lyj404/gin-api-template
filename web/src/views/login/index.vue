<template>
  <div class="min-h-full flex-center bg-gray-100">
    <n-card class="w-400px" :bordered="false">
      <n-h2 class="text-center">管理员登录</n-h2>

      <n-form ref="formRef" :model="form" :rules="rules" size="large">
        <n-form-item path="email" label="邮箱">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" />
        </n-form-item>
        <n-form-item path="captcha" label="验证码">
          <n-input v-model:value="form.captcha" placeholder="请输入验证码" class="flex-1" />
          <img :src="captchaUrl" alt="验证码" class="ml-8 h-40px cursor-pointer" @click="refreshCaptcha" />
        </n-form-item>
      </n-form>

      <n-button type="primary" block :loading="loading" @click="handleLogin">登录</n-button>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NH2, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import type { FormRules, FormItemRule } from 'naive-ui'

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
