<template>
  <n-config-provider :locale="zhCN" :theme-overrides="themeOverrides">
    <router-view />
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, type GlobalThemeOverrides } from 'naive-ui'
import zhCN from 'naive-ui/es/locales/common/zhCN'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    fontFamily: '"Noto Sans SC", sans-serif'
  }
}
</script>
