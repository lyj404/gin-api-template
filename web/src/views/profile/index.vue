<template>
  <div class="page-padding">
    <n-h2>个人信息</n-h2>

    <n-card title="基本信息" class="mb-6">
      <n-form ref="profileFormRef" :model="profileForm" :rules="profileRules" label-placement="left" :label-width="100">
        <n-form-item label="姓名" path="name">
          <n-input v-model:value="profileForm.name" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="profileForm.email" />
        </n-form-item>
        <n-form-item label="创建时间">
          <n-input :value="profile.created_at" disabled />
        </n-form-item>
        <n-form-item label="更新时间">
          <n-input :value="profile.updated_at" disabled />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="profileLoading" @click="handleUpdateProfile">保存</n-button>
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="修改密码">
      <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-placement="left" :label-width="100">
        <n-form-item label="原密码" path="old_password">
          <n-input v-model:value="passwordForm.old_password" type="password" show-password-on="click" />
        </n-form-item>
        <n-form-item label="新密码" path="new_password">
          <n-input v-model:value="passwordForm.new_password" type="password" show-password-on="click" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="passwordLoading" @click="handleChangePassword">确认修改</n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  NCard, NH2, NForm, NFormItem, NInput, NButton, useMessage
} from 'naive-ui'
import { getProfile, updateProfile, changePassword } from '@/api'
import type { ProfileResponse, UpdateProfileRequest, ChangePasswordRequest } from '@/types'

const message = useMessage()

const profile = ref<ProfileResponse>({ id: '', name: '', email: '', created_at: '', updated_at: '' })
const profileLoading = ref(false)
const passwordLoading = ref(false)
const profileFormRef = ref<any>(null)
const passwordFormRef = ref<any>(null)

const profileForm = reactive<UpdateProfileRequest>({ name: '', email: '' })

const profileRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email' as const, message: '邮箱格式错误', trigger: 'blur' }]
}

const passwordForm = reactive<ChangePasswordRequest>({ old_password: '', new_password: '' })

const passwordRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }]
}

const loadProfile = async () => {
  try {
    const res = await getProfile()
    profile.value = res.data.data!
    profileForm.name = profile.value.name
    profileForm.email = profile.value.email
  } catch {
    message.error('获取个人信息失败')
  }
}

const handleUpdateProfile = async () => {
  try {
    await profileFormRef.value?.validate()
  } catch {
    return
  }
  profileLoading.value = true
  try {
    await updateProfile(profileForm)
    message.success('个人信息更新成功')
    await loadProfile()
  } catch (err: any) {
    message.error(err?.response?.data?.message || '更新失败')
  } finally {
    profileLoading.value = false
  }
}

const handleChangePassword = async () => {
  try {
    await passwordFormRef.value?.validate()
  } catch {
    return
  }
  passwordLoading.value = true
  try {
    await changePassword(passwordForm)
    message.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
  } catch (err: any) {
    message.error(err?.response?.data?.message || '修改失败')
  } finally {
    passwordLoading.value = false
  }
}

onMounted(loadProfile)
</script>
