import type { Directive, DirectiveBinding } from 'vue'
import { usePermissionStore } from '@/stores/permission'

const vPermission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const permission = usePermissionStore()
    const { value } = binding
    if (value && !permission.hasPermission(value)) {
      el.style.display = 'none'
    }
  }
}

export default vPermission
