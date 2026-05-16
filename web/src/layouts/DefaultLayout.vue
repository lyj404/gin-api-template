<template>
  <n-layout has-sider class="h-full">
    <n-layout-sider
      v-if="!layout.isMobile"
      bordered
      :width="240"
      :collapsed-width="68"
      :collapsed="layout.sidebarCollapsed"
      collapse-mode="width"
      :native-scrollbar="false"
      :show-trigger="false"
      class="sider"
      @update:collapsed="layout.sidebarCollapsed = $event"
    >
      <Logo :collapsed="layout.sidebarCollapsed" />
      <SideMenu :collapsed="layout.sidebarCollapsed" />
    </n-layout-sider>

    <n-drawer
      v-else
      v-model:show="layout.mobileDrawerOpen"
      :width="280"
      placement="left"
      :block-scroll="true"
      :auto-focus="false"
    >
      <n-drawer-content :native-scrollbar="false" body-content-style="padding:0">
        <Logo :collapsed="false" />
        <SideMenu :collapsed="false" @select="layout.closeMobileDrawer()" />
      </n-drawer-content>
    </n-drawer>

    <n-layout>
      <n-layout-header bordered class="layout-header">
        <button class="header-btn" @click="layout.toggleSidebar()">
          <span :class="[hamburgerIcon, 'text-xl']" />
        </button>

        <n-breadcrumb class="header-breadcrumb">
          <n-breadcrumb-item v-for="crumb in breadcrumbs" :key="crumb.path">
            <span class="breadcrumb-dot" />
            {{ crumb.title }}
          </n-breadcrumb-item>
        </n-breadcrumb>

        <div class="header-actions">
          <n-tooltip trigger="hover">
            <template #trigger>
              <button class="header-btn hide-on-mobile" @click="reload">
                <span class="i-material-symbols:refresh text-xl" />
              </button>
            </template>
            刷新
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <button class="header-btn hide-on-mobile" @click="toggleFullscreen">
                <span :class="[isFullscreen ? 'i-material-symbols:fullscreen-exit' : 'i-material-symbols:fullscreen', 'text-xl']" />
              </button>
            </template>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </n-tooltip>

          <n-dropdown
            :options="userOptions"
            placement="bottom-end"
            trigger="click"
            @select="handleUserCommand"
          >
            <div class="user-box">
              <n-avatar round :size="30" class="user-avatar">
                {{ userInitial }}
              </n-avatar>
              <span class="user-name hidden md:inline">{{ userName }}</span>
              <span class="i-material-symbols:keyboard-arrow-down text-sm text-gray-400 hidden md:inline" />
            </div>
          </n-dropdown>
        </div>
      </n-layout-header>

      <n-layout-content class="app-content">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NDrawer, NDrawerContent,
  NDropdown, NAvatar, NBreadcrumb, NBreadcrumbItem, NTooltip
} from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useLayoutStore } from '@/stores/layout'
import { getUserInfo } from '@/utils/auth'
import Logo from './components/Logo.vue'
import SideMenu from './components/SideMenu.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const layout = useLayoutStore()

const userName = ref(getUserInfo()?.name || 'Admin')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())
const isFullscreen = ref(false)

const hamburgerIcon = computed(() => {
  if (layout.isMobile) {
    return layout.mobileDrawerOpen ? 'i-material-symbols:close' : 'i-material-symbols:menu'
  }
  return layout.sidebarCollapsed ? 'i-material-symbols:menu' : 'i-material-symbols:menu-open'
})

const breadcrumbs = computed(() => {
  return route.matched
    .filter(r => r.meta?.title)
    .map(r => ({ path: r.path, title: r.meta.title as string }))
})

const userOptions: DropdownOption[] = [
  {
    key: 'profile',
    label: '个人信息',
    icon: () => h('span', { class: 'i-material-symbols:person-outline' })
  },
  { type: 'divider', key: 'd1' },
  {
    key: 'logout',
    label: '退出登录',
    icon: () => h('span', { class: 'i-material-symbols:logout' })
  }
]

const handleUserCommand = (key: string) => {
  if (key === 'profile') {
    router.push('/profile')
  } else if (key === 'logout') {
    auth.logout()
  }
}

const reload = () => {
  window.location.reload()
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}
</script>

<style scoped>
.layout-header {
  height: var(--header-height, 64px);
  display: flex;
  align-items: center;
  padding: 0 16px 0 10px;
  background: var(--color-header-bg, rgba(255,255,255,0.78));
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-header-border, rgba(229,221,213,0.6));
  gap: 6px;
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

@media (min-width: 768px) {
  .layout-header {
    padding: 0 24px 0 14px;
    gap: 10px;
  }
}

.header-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  color: var(--color-text-secondary, #78716c);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.header-btn:hover {
  background: var(--color-primary-soft, rgba(194, 112, 74, 0.08));
  color: var(--color-primary, #c2704a);
}

.header-breadcrumb {
  flex: 1;
  min-width: 0;
  margin-left: 4px;
  overflow: hidden;
  white-space: nowrap;
}

.breadcrumb-dot {
  display: inline-block;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--color-primary, #c2704a);
  margin-right: 6px;
  vertical-align: middle;
  opacity: 0.6;
}

.header-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

@media (min-width: 768px) {
  .header-actions {
    gap: 8px;
  }
}

.user-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-left: 4px;
}

.user-box:hover {
  background: var(--color-primary-soft, rgba(194, 112, 74, 0.08));
}

.user-avatar {
  background: linear-gradient(135deg, #c2704a, #d97706) !important;
  color: #fff !important;
  font-weight: 700 !important;
  font-size: 13px !important;
  flex-shrink: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text, #1c1917);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
