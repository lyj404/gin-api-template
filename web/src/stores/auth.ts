import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi } from '@/api'
import type { LoginRequest, LoginResponse } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('accessToken') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  const captchaUrl = ref<string>('')

  const setToken = (access: string, refresh: string) => {
    token.value = access
    refreshToken.value = refresh
    localStorage.setItem('accessToken', access)
    localStorage.setItem('refreshToken', refresh)
  }

  const clearToken = () => {
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  }

  const login = async (data: LoginRequest): Promise<LoginResponse> => {
    const res = await loginApi(data)
    const result = res.data as any
    if (result?.data) {
      setToken(result.data.accessToken, result.data.refreshToken)
    }
    return result.data
  }

  const refreshCaptcha = async (): Promise<string> => {
    const res = await fetch('/api/captcha', { credentials: 'include' })
    const url = URL.createObjectURL(await res.blob())
    captchaUrl.value = url
    return url
  }

  const logout = () => {
    clearToken()
    window.location.href = '/login'
  }

  return { token, refreshToken, captchaUrl, login, logout, clearToken, refreshCaptcha }
})
