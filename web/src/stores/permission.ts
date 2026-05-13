import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserMenus, getUserPermissions } from '@/api'
import type { MenuTreeNode, UserPermissions } from '@/types'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref<MenuTreeNode[]>([])
  const permissions = ref<UserPermissions | null>(null)
  const loaded = ref(false)

  const fetchMenus = async () => {
    const res = await getUserMenus()
    menus.value = res.data.data?.menus || []
  }

  const fetchPermissions = async () => {
    const res = await getUserPermissions()
    permissions.value = res.data.data || null
  }

  const loadAll = async () => {
    await Promise.all([fetchMenus(), fetchPermissions()])
    loaded.value = true
  }

  const hasPermission = (resource: string, method: string = 'GET') => {
    if (!permissions.value) return false
    return permissions.value.permissions.some(
      (p: { resource: string; method: string }) => p.resource === resource && (p.method === method || p.method === '*')
    )
  }

  const clear = () => {
    menus.value = []
    permissions.value = null
    loaded.value = false
  }

  return { menus, permissions, loaded, fetchMenus, fetchPermissions, loadAll, hasPermission, clear }
})
