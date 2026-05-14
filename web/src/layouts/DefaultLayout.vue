<template>
  <n-layout has-sider class="h-full">
    <n-layout-sider
      bordered
      :width="220"
      :collapsed-width="64"
      :collapsed="collapsed"
      :native-scrollbar="false"
      class="sider"
      @update:collapsed="collapsed = $event"
    >
      <div class="logo">
        <span class="i-material-symbols:shield-person-outline text-2xl text-primary" />
        <transition name="fade">
          <span v-if="!collapsed" class="logo-text">Admin 后台</span>
        </transition>
      </div>

      <n-menu
        v-model:value="activeKey"
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="20"
        :indent="20"
        :options="menuOptions"
        :render-label="renderLabel"
        :default-expand-all="true"
      />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="header">
        <button class="action-btn collapse-btn" @click="collapsed = !collapsed">
          <span :class="[collapsed ? 'i-material-symbols:menu' : 'i-material-symbols:menu-open', 'text-xl']" />
        </button>

        <n-breadcrumb class="ml-4">
          <n-breadcrumb-item v-for="crumb in breadcrumbs" :key="crumb.path">
            {{ crumb.title }}
          </n-breadcrumb-item>
        </n-breadcrumb>

        <div class="actions">
          <n-tooltip trigger="hover">
            <template #trigger>
              <button class="action-btn" @click="reload">
                <span class="i-material-symbols:refresh text-xl" />
              </button>
            </template>
            刷新
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <button class="action-btn" @click="toggleFullscreen">
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
              <n-avatar round :size="32" class="user-avatar">
                {{ userInitial }}
              </n-avatar>
              <span class="user-name">{{ userName }}</span>
              <span class="i-material-symbols:keyboard-arrow-down text-base text-gray-400" />
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
import { ref, computed, h, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent,
  NMenu, NDropdown, NAvatar, NBreadcrumb, NBreadcrumbItem, NTooltip
} from 'naive-ui'
import type { MenuOption, DropdownOption } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { getUserInfo } from '@/utils/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const permission = usePermissionStore()

const collapsed = ref(false)
const activeKey = ref<string>(route.path)
const userName = ref(getUserInfo()?.name || 'Admin')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())
const isFullscreen = ref(false)

watch(() => route.path, (newPath) => {
  activeKey.value = newPath
})

const breadcrumbs = computed(() => {
  return route.matched
    .filter(r => r.meta?.title)
    .map(r => ({ path: r.path, title: r.meta.title as string }))
})

const menuOptions = computed<MenuOption[]>(() => {
  return permission.menus.map(menu => ({
    key: menu.path || String(menu.id),
    label: menu.name,
    path: menu.path,
    icon: () => h('span', { class: `${menu.icon || 'i-material-symbols:circle-outline'} text-lg` }),
    children: menu.children?.map(child => ({
      key: child.path,
      label: child.name,
      path: child.path,
      icon: () => h('span', { class: `${child.icon || 'i-material-symbols:circle-outline'} text-lg` })
    }))
  }))
})

const renderLabel = (option: MenuOption) => {
  const menuKey = option.key as string
  if (menuKey && menuKey.startsWith('/')) {
    return h('a', {
      onClick: (e: Event) => {
        e.preventDefault()
        router.push(menuKey)
      }
    }, option.label as string)
  }
  return option.label as string
}

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
.sider :deep(.n-layout-sider-scroll-container) {
  display: flex;
  flex-direction: column;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid #efeff5;
  font-weight: 700;
  font-size: 16px;
  letter-spacing: 0.5px;
  overflow: hidden;
  white-space: nowrap;
  flex-shrink: 0;
}

.logo-text {
  background: linear-gradient(90deg, #18a058 0%, #36ad6a 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.text-primary {
  color: #18a058;
}

.header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 24px 0 12px;
  background-color: #fff;
  gap: 8px;
}

.collapse-btn {
  color: #444;
}

.ml-4 {
  margin-left: 4px;
}

.actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  color: #555;
  transition: background-color 0.2s, color 0.2s;
}

.action-btn:hover {
  background-color: rgba(0, 0, 0, 0.04);
  color: #18a058;
}

.user-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-left: 4px;
}

.user-box:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.user-avatar {
  background-color: #18a058 !important;
  color: #fff !important;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
