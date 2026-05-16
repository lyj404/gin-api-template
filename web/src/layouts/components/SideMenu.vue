<template>
  <n-menu
    v-model:value="activeKey"
    :collapsed="collapsed"
    :collapsed-width="64"
    :collapsed-icon-size="22"
    :icon-size="20"
    :indent="18"
    :options="menuOptions"
    :default-expand-all="!collapsed"
    class="side-menu"
    @update:value="onSelect"
  />
</template>

<script setup lang="ts">
import { computed, h, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, type MenuOption } from 'naive-ui'
import { usePermissionStore } from '@/stores/permission'
import type { MenuTreeNode } from '@/types'

defineProps<{
  collapsed?: boolean
}>()

const emit = defineEmits<{
  (e: 'select', key: string): void
}>()

const router = useRouter()
const route = useRoute()
const permission = usePermissionStore()

const activeKey = ref<string>(route.path)

watch(() => route.path, (newPath) => {
  activeKey.value = newPath
})

const safelistedFallback = (hasChildren?: boolean) =>
  hasChildren ? 'i-material-symbols:folder-outline' : 'i-material-symbols:circle-outline'

const resolveIconClass = (icon: string | undefined, hasChildren: boolean): string => {
  const raw = (icon || '').trim()
  if (!raw) return safelistedFallback(hasChildren)
  return raw.startsWith('i-') ? raw : `i-${raw}`
}

const renderIcon = (iconClass: string) => () => h('span', {
  class: iconClass,
  style: 'display: inline-block; width: 20px; height: 20px; flex-shrink: 0; background-color: currentColor;',
})

const toMenuOption = (menu: MenuTreeNode): MenuOption => {
  const hasChildren = !!(menu.children && menu.children.length)
  const iconClass = resolveIconClass(menu.icon, hasChildren)
  return {
    key: menu.path || String(menu.id),
    label: menu.name,
    icon: renderIcon(iconClass),
    children: hasChildren ? menu.children!.map(toMenuOption) : undefined,
  }
}

const menuOptions = computed<MenuOption[]>(() => permission.menus.map(toMenuOption))

const onSelect = (key: string) => {
  if (key && key.startsWith('/')) {
    router.push(key)
  }
  emit('select', key)
}
</script>
