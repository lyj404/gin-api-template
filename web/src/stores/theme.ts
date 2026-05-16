import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref<boolean>(localStorage.getItem('theme-dark') === 'true')

  const applyTheme = () => {
    const html = document.documentElement
    if (isDark.value) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  const toggle = () => {
    const html = document.documentElement
    html.classList.add('theming')
    isDark.value = !isDark.value
    setTimeout(() => html.classList.remove('theming'), 400)
  }

  watch(isDark, (val) => {
    localStorage.setItem('theme-dark', String(val))
    applyTheme()
  }, { immediate: true })

  return { isDark, toggle, applyTheme }
})
