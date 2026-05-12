<template>
  <n-layout has-sider class="h-full">
    <n-layout-sider
      bordered
      :width="220"
      :collapsed-width="64"
      :collapsed="collapsed"
      show-trigger
      @update:collapsed="collapsed = $event"
    >
      <div class="h-full flex flex-col">
        <div class="h-60px flex-center font-bold text-lg px16">
          <span v-if="!collapsed">Admin</span>
          <span v-else>A</span>
        </div>
        <n-menu
          v-model:value="activeKey"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
          :render-label="renderLabel"
          :default-expand-all="true"
        />
      </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="h-60px flex items-center px16">
        <div class="flex-1 font-medium">{{ route.meta.title || '' }}</div>
        <n-dropdown :options="userOptions" @select="handleUserCommand">
          <n-button text class="flex items-center">
            <span class="i-material-symbols:person-outline mr-8" />
            {{ userName }}
          </n-button>
        </n-dropdown>
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
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NButton } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { getUserInfo } from '@/utils/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const permission = usePermissionStore()

const collapsed = ref(false)
const activeKey = ref<string>()
const userName = ref(getUserInfo()?.name || 'Admin')

const menuOptions = computed<MenuOption[]>(() => {
  return permission.menus.map(menu => ({
    key: String(menu.id),
    label: menu.name,
    icon: () => h('span', { class: menu.icon || 'i-material-symbols:circle' }),
    children: menu.children?.map(child => ({
      key: String(child.id),
      label: child.name,
      path: child.path,
      icon: () => h('span', { class: child.icon || 'i-material-symbols:circle' })
    }))
  }))
})

const renderLabel = (option: MenuOption) => {
  if ((option as any).path) {
    return h('a', {
      onClick: () => router.push((option as any).path)
    }, option.label as string)
  }
  return option.label as string
}

const userOptions = [
  { label: '退出登录', key: 'logout' }
]

const handleUserCommand = (key: string) => {
  if (key === 'logout') {
    auth.logout()
  }
}
</script>
