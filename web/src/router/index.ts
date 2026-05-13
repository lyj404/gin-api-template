import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
  }
}

export const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'i-material-symbols:dashboard-outline' }
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('@/views/roles/index.vue'),
        meta: { title: '角色管理', icon: 'i-material-symbols:manage-accounts-outline' }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('@/views/menus/index.vue'),
        meta: { title: '菜单管理', icon: 'i-material-symbols:menu-outline' }
      },
      {
        path: 'orgs',
        name: 'Orgs',
        component: () => import('@/views/orgs/index.vue'),
        meta: { title: '组织管理', icon: 'i-material-symbols:corporate-fare-outline' }
      },
      {
        path: 'resources',
        name: 'Resources',
        component: () => import('@/views/resources/index.vue'),
        meta: { title: '资源管理', icon: 'i-material-symbols:security-outline' }
      },
      {
        path: 'audit-logs',
        name: 'AuditLogs',
        component: () => import('@/views/audit-logs/index.vue'),
        meta: { title: '审计日志', icon: 'i-material-symbols:history-outline' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const whiteList = ['/login']

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()
  const token = authStore.token

  if (token) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      if (permissionStore.menus.length === 0) {
        try {
          await Promise.all([
            permissionStore.fetchMenus(),
            permissionStore.fetchPermissions()
          ])
          next({ ...to, replace: true })
        } catch (error) {
          authStore.clearToken()
          next(`/login?redirect=${to.path}`)
        }
      } else {
        next()
      }
    }
  } else {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router
