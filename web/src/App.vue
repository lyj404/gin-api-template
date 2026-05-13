<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <router-view />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, NNotificationProvider, type GlobalThemeOverrides } from 'naive-ui'
import zhCN from 'naive-ui/es/locales/common/zhCN'
import dateZhCN from 'naive-ui/es/locales/date/zhCN'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    fontFamily: '"Noto Sans SC", sans-serif'
  }
}
</script>
