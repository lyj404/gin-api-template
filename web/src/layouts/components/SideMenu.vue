<template>
  <n-menu
    v-model:value="activeKey"
    :collapsed="collapsed"
    :collapsed-width="64"
    :collapsed-icon-size="20"
    :indent="20"
    :options="menuOptions"
    :default-expand-all="!collapsed"
    @update:value="onSelect"
  />
</template>

<script setup lang="ts">
import { computed, h, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, type MenuOption } from 'naive-ui'
import { Icon } from '@iconify/vue'
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

const fallbackIcon = (icon?: string, hasChildren?: boolean) => {
  const raw = (icon || '').replace(/^i-/, '').trim()
  if (raw) return raw
  return hasChildren ? 'material-symbols:folder-outline' : 'material-symbols:circle-outline'
}

const toMenuOption = (menu: MenuTreeNode): MenuOption => {
  const hasChildren = !!(menu.children && menu.children.length)
  const iconName = fallbackIcon(menu.icon, hasChildren)
  return {
    key: menu.path || String(menu.id),
    label: menu.name,
    icon: () => h(Icon, { icon: iconName, class: 'text-lg' }),
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
