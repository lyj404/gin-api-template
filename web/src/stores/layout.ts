import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import { useBreakpoints, breakpointsTailwind } from '@vueuse/core'

const STORAGE_KEY = 'sidebar-collapsed'

export const useLayoutStore = defineStore('layout', () => {
  const breakpoints = useBreakpoints(breakpointsTailwind)
  const isMobile = breakpoints.smaller('md')          // < 768px
  const isTablet = breakpoints.between('md', 'lg')    // 768-1023px
  const isDesktop = breakpoints.greaterOrEqual('lg')  // >= 1024px

  const sidebarCollapsed = ref<boolean>(localStorage.getItem(STORAGE_KEY) === 'true')
  const mobileDrawerOpen = ref(false)

  watch(sidebarCollapsed, (val) => {
    localStorage.setItem(STORAGE_KEY, String(val))
  })

  watch(isTablet, (val) => {
    if (val) sidebarCollapsed.value = true
  }, { immediate: true })

  watch(isMobile, (val) => {
    if (!val) mobileDrawerOpen.value = false
  })

  const device = computed<'mobile' | 'tablet' | 'desktop'>(() => {
    if (isMobile.value) return 'mobile'
    if (isTablet.value) return 'tablet'
    return 'desktop'
  })

  const toggleSidebar = () => {
    if (isMobile.value) {
      mobileDrawerOpen.value = !mobileDrawerOpen.value
    } else {
      sidebarCollapsed.value = !sidebarCollapsed.value
    }
  }

  const closeMobileDrawer = () => {
    mobileDrawerOpen.value = false
  }

  return {
    isMobile,
    isTablet,
    isDesktop,
    device,
    sidebarCollapsed,
    mobileDrawerOpen,
    toggleSidebar,
    closeMobileDrawer,
  }
})
