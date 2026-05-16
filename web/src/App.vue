<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="themeStore.isDark ? darkTheme : undefined" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <router-view />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
    <button class="theme-toggle" :class="{ 'theme-toggle--dark': themeStore.isDark }" @click="themeStore.toggle">
      <span class="theme-toggle-icon" />
    </button>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, NNotificationProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import zhCN from 'naive-ui/es/locales/common/zhCN'
import dateZhCN from 'naive-ui/es/locales/date/zhCN'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const isDark = themeStore.isDark
  return {
    common: {
      fontFamily: '"Noto Sans SC", sans-serif',
      primaryColor: isDark ? '#d97a4a' : '#c2704a',
      primaryColorHover: isDark ? '#e08f63' : '#a85d38',
      primaryColorPressed: isDark ? '#c96a3e' : '#924e2d',
      primaryColorSuppl: isDark ? '#d97a4a' : '#c2704a',
      successColor: isDark ? '#34d399' : '#059669',
      warningColor: isDark ? '#fbbf24' : '#d97706',
      errorColor: isDark ? '#f87171' : '#dc2626',
      infoColor: isDark ? '#22d3ee' : '#0891b2',
      borderRadius: '8px',
      borderRadiusSmall: '6px',
      ...(isDark ? {
        borderColor: '#3d3631',
        dividerColor: '#332d28',
        inputColor: '#201c19',
        inputBorderColor: '#3d3631',
        inputColorFocus: '#26211e',
        inputBorderColorFocus: '#d97a4a',
        cardColor: '#26211e',
        tableHeaderColor: '#2d2824',
        tableColorHover: '#2d2824',
        popoverColor: '#26211e',
        modalColor: '#26211e',
        tagColor: '#2d2824',
        actionColor: '#2d2824',
        tabColor: '#1a1614',
        inputColorDisabled: '#2d2824',
        progressRailColor: '#3d3631',
        scrollbarColor: '#3d3631',
        scrollbarColorHover: '#78716c',
        skeletonColor: '#2d2824',
        skeletonColorEnd: '#332d28',
        clearColor: '#3d3631',
        clearColorHover: '#78716c',
        iconColor: '#a8a29e',
        iconColorHover: '#e7e5e4',
        iconColorPressed: '#e7e5e4',
      } : {
        borderColor: '#e5ddd5',
        dividerColor: '#efe9e2',
        inputColor: '#fafaf9',
        inputBorderColor: '#e5ddd5',
        inputColorFocus: '#ffffff',
        inputBorderColorFocus: '#c2704a',
        cardColor: '#ffffff',
        tableHeaderColor: '#faf8f5',
        popoverColor: '#ffffff',
        modalColor: '#ffffff',
        tagColor: '#f7f4f0',
        actionColor: '#faf8f5',
        tabColor: '#f7f4f0',
        inputColorDisabled: '#f7f4f0',
        progressRailColor: '#efe9e2',
        scrollbarColor: '#e5ddd5',
        scrollbarColorHover: '#a8a29e',
        skeletonColor: '#f7f4f0',
        skeletonColorEnd: '#efe9e2',
        clearColor: '#e5ddd5',
        clearColorHover: '#a8a29e',
        iconColor: '#a8a29e',
        iconColorHover: '#78716c',
        iconColorPressed: '#1c1917',
      })
    }
  }
})
</script>

<style scoped>
.theme-toggle {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
  width: 44px;
  height: 44px;
  border-radius: 12px;
  border: none;
  background: var(--color-surface, #ffffff);
  box-shadow: var(--shadow-md, 0 4px 12px rgba(28, 25, 23, 0.06));
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary, #78716c);
  transition: all 0.25s ease;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.theme-toggle:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg, 0 8px 30px rgba(28, 25, 23, 0.08));
  color: var(--color-primary, #c2704a);
}

.theme-toggle--dark {
  background: #2d2824;
  color: #d97a4a;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.35);
}

.theme-toggle--dark:hover {
  color: #e08f63;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.45);
}

.theme-toggle-icon {
  width: 22px;
  height: 22px;
  background-color: currentColor;
  -webkit-mask: var(--icon) center/contain no-repeat;
  mask: var(--icon) center/contain no-repeat;
}

.theme-toggle:not(.theme-toggle--dark) .theme-toggle-icon {
  --icon: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z' fill='currentColor'/%3E%3C/svg%3E");
}

.theme-toggle--dark .theme-toggle-icon {
  --icon: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M12 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10zm0-12V3m0 18v-2M4.22 4.22l1.42 1.42m12.72 12.72 1.42 1.42M3 12h2m14 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E");
}

@media (max-width: 767px) {
  .theme-toggle {
    bottom: 16px;
    right: 16px;
    width: 40px;
    height: 40px;
    border-radius: 10px;
  }
}
</style>
