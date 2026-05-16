<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="themeStore.isDark ? darkTheme : undefined" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <router-view />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
    <button class="theme-toggle" @click="themeStore.toggle">
      <span :class="themeStore.isDark ? 'i-material-symbols:light-mode' : 'i-material-symbols:dark-mode'" />
    </button>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, NNotificationProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import zhCN from 'naive-ui/es/locales/common/zhCN'
import dateZhCN from 'naive-ui/es/locales/date/zhCN'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const themeOverrides: GlobalThemeOverrides = {
  common: {
    fontFamily: '"Noto Sans SC", sans-serif'
  }
}
</script>

<style scoped>
.theme-toggle {
  position: fixed;
  bottom: 16px;
  right: 16px;
  z-index: 9999;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid #e2e8f0;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #334155;
  transition: all 0.2s;
}
@media (min-width: 768px) {
  .theme-toggle {
    bottom: 24px;
    right: 24px;
    width: 44px;
    height: 44px;
    font-size: 22px;
  }
}
.theme-toggle:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}
.dark .theme-toggle {
  background: #1e293b;
  border-color: #334155;
  color: #e2e8f0;
}
</style>
