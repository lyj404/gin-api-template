import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserMenus, getUserPermissions } from '@/api'
import type { MenuTreeNode, UserPermissions } from '@/types'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref<MenuTreeNode[]>([])
  const permissions = ref<UserPermissions | null>(null)

  const fetchMenus = async () => {
    const res = await getUserMenus()
    menus.value = res.data.data?.menus || []
  }

  const fetchPermissions = async () => {
    const res = await getUserPermissions()
    permissions.value = res.data.data || null
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
  }

  return { menus, permissions, fetchMenus, fetchPermissions, hasPermission, clear }
})
