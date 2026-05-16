<template>
  <div class="page-padding">
    <div class="toolbar-row mb-3">
      <n-h2 class="!my-0">菜单管理</n-h2>
      <n-button type="primary" @click="openModal()">新增菜单</n-button>
    </div>

    <n-card>
      <n-tree
        :data="treeData"
        :loading="loading"
        key-field="id"
        label-field="label"
        children-field="children"
        default-expand-all
        selectable
        :render-label="renderNodeLabel"
      />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑菜单' : '新增菜单'" :style="{ width: '90vw', maxWidth: '600px' }">
      <n-form :model="form" label-placement="left" label-width="90">
        <n-form-item label="菜单名称">
          <n-input v-model:value="form.name" placeholder="请输入菜单名称" />
        </n-form-item>
        <n-form-item label="上级菜单">
          <n-select v-model:value="form.parent_id" :options="menuOptions" placeholder="请选择上级菜单" clearable />
        </n-form-item>
        <n-form-item label="路由路径">
          <n-input v-model:value="form.path" placeholder="如: /users" />
        </n-form-item>
        <n-form-item label="图标">
          <div class="flex gap-8 items-center w-full">
            <Icon v-if="form.icon" :icon="toIconifyName(form.icon)" style="font-size: 1.25rem" />
            <n-button @click="openIconPicker">选择</n-button>
          </div>
        </n-form-item>
        <n-form-item label="排序">
          <n-input-number v-model:value="form.order_num" :min="0" />
        </n-form-item>
        <n-form-item label="关联资源">
          <n-select v-model:value="form.resource_ids" multiple :options="allResourceOptions" placeholder="选择资源" :loading="resourcesLoading" />
        </n-form-item>
        <n-form-item label="是否显示">
          <n-switch v-model:value="form.is_visible" />
        </n-form-item>
        <n-form-item label="状态">
          <n-select v-model:value="form.status" :options="statusOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSave">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showIconPicker" preset="card" title="选择图标" :style="{ width: '95vw', maxWidth: '720px' }">
      <n-input v-model:value="iconSearch" placeholder="搜索图标名称" clearable style="margin-bottom: 12px" />
      <div style="font-size: 12px; color: #999; margin-bottom: 8px">共 {{ filteredIcons.length }} 个，显示前 {{ Math.min(filteredIcons.length, iconDisplayLimit) }} 个</div>
      <div @scroll="onIconGridScroll" class="icon-grid">
        <div v-for="name in displayedIcons" :key="name"
          :title="name"
          class="icon-cell"
          :style="form.icon === `i-material-symbols:${name}` ? { borderColor: '#18a058', color: '#18a058' } : {}"
          @click="selectIcon(`i-material-symbols:${name}`)">
          <Icon :icon="`material-symbols:${name}`" style="font-size: 1.5rem" />
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed } from 'vue'
import { NButton, NTree, NCard, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NSwitch, NTag, NH2, NSpace, useMessage } from 'naive-ui'
import type { TreeOption, SelectOption } from 'naive-ui'
import { Icon } from '@iconify/vue'
import { getMenuTree, getMenu, createMenu, updateMenu, deleteMenu, getResources, bindMenuResource, unbindMenuResource } from '@/api'
import { usePermissionStore } from '@/stores/permission'
import type { MenuTreeNode, MenuResponse, ResourceResponse } from '@/types'
import { useDict } from '@/composables/useDict'

const message = useMessage()
const permission = usePermissionStore()
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const showIconPicker = ref(false)
const editingId = ref<string | null>(null)

const toIconifyName = (icon?: string) => (icon || '').replace(/^i-/, '')

const allIconNames = ref<string[]>([])
const iconSearch = ref('')
const iconDisplayLimit = ref(240)

const filteredIcons = computed(() => {
  const q = iconSearch.value.trim().toLowerCase()
  if (!q) return allIconNames.value
  return allIconNames.value.filter(n => n.includes(q))
})
const displayedIcons = computed(() => filteredIcons.value.slice(0, iconDisplayLimit.value))

const openIconPicker = () => {
  iconSearch.value = ''
  iconDisplayLimit.value = 240
  showIconPicker.value = true
}

const onIconGridScroll = (e: Event) => {
  const el = e.target as HTMLElement
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 40 && iconDisplayLimit.value < filteredIcons.value.length) {
    iconDisplayLimit.value = Math.min(iconDisplayLimit.value + 240, filteredIcons.value.length)
  }
}

const form = reactive({ name: '', parent_id: null as string | null, path: '', icon: '', order_num: 0, resource_ids: [] as string[], is_visible: true, status: 'enabled' })
const data = ref<MenuTreeNode[]>([])

const { options: statusOptions, lookup: statusLabel, load: loadStatusOptions } = useDict('menu_status')

const buildTreeOptions = (nodes: MenuTreeNode[]): TreeOption[] => {
  return (nodes || []).map(n => ({
    id: n.id,
    label: n.name,
    key: n.id,
    icon: n.icon,
    path: n.path,
    order_num: n.order_num,
    is_visible: n.is_visible,
    status: n.status,
    children: n.children ? buildTreeOptions(n.children) : undefined
  }))
}

const treeData = computed(() => buildTreeOptions(data.value))

const menuOptions = computed(() => [{ label: '顶级菜单', value: '0' }, ...data.value.map((m: MenuTreeNode) => ({ label: m.name, value: m.id }))])
const allResourceOptions = ref<SelectOption[]>([])
const resourcesLoading = ref(false)

const renderNodeLabel = (info: { option: TreeOption }) => {
  const opt = info.option
  return h('div', { style: 'display: flex; align-items: center; gap: 12px; width: 100%; padding: 4px 0' }, [
    h(Icon, { icon: toIconifyName(opt.icon as string) || 'material-symbols:circle-outline', class: 'text-lg', style: 'flex-shrink: 0' }),
    h('span', { style: 'flex-shrink: 0' }, opt.label),
    opt.path ? h('span', { style: 'color: #999; font-size: 12px; flex-shrink: 0' }, opt.path as string) : null,
    h('span', { style: 'flex: 1' }),
    h(NTag, { size: 'tiny', type: opt.status === 'enabled' ? 'success' : 'error', bordered: false }, { default: () => statusLabel(opt.status as string || '') }),
    h(NButton, { size: 'tiny', quaternary: true, onClick: (e: Event) => { e.stopPropagation(); openModal(opt) } }, { default: () => '编辑' }),
    h(NButton, { size: 'tiny', quaternary: true, type: 'error', onClick: (e: Event) => { e.stopPropagation(); handleDelete(opt) } }, { default: () => '删除' })
  ])
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMenuTree()
    data.value = res.data.data || []
  } finally {
    loading.value = false
  }
}

const fetchAllResources = async () => {
  resourcesLoading.value = true
  try {
    const res = await getResources({ page: 1, page_size: 200 })
    allResourceOptions.value = (res.data.data?.data || []).map((r: ResourceResponse) => ({ label: `${r.description} (${r.name})`, value: r.id }))
  } catch { /* ignore */ }
  finally { resourcesLoading.value = false }
}

const openModal = async (row?: TreeOption) => {
  loadStatusOptions()
  if (row) {
    editingId.value = row.id as string
    form.name = row.label as string
    form.path = row.path as string || ''
    form.icon = row.icon as string || ''
    form.order_num = row.order_num as number || 0
    form.is_visible = row.is_visible as boolean ?? true
    form.status = (row.status as string) || 'enabled'

    try {
      const detail = await getMenu(row.id as string)
      const resources = (detail.data.data as MenuResponse).resources || []
      form.resource_ids = resources.map((r: any) => r.id)
    } catch {
      form.resource_ids = []
    }
  } else {
    editingId.value = null
    Object.assign(form, { name: '', parent_id: null, path: '', icon: '', order_num: 0, resource_ids: [], is_visible: true, status: 'enabled' })
  }
  showModal.value = true
  fetchAllResources()
}

const handleSave = async () => {
  saving.value = true
  try {
    const payload: any = { name: form.name, path: form.path, icon: form.icon, order_num: form.order_num, is_visible: form.is_visible, status: form.status }
    if (form.parent_id !== null && form.parent_id !== undefined) {
      payload.parent_id = form.parent_id === '0' ? null : form.parent_id
    }

    let menuId: string
    if (editingId.value) {
      await updateMenu(editingId.value, payload)
      menuId = editingId.value
    } else {
      const res = await createMenu(payload)
      menuId = res.data.data?.id
    }

    let currentIds: string[] = []
    try {
      const detail = await getMenu(menuId)
      currentIds = ((detail.data.data as MenuResponse).resources || []).map((r: any) => r.id)
    } catch { /* ignore */ }

    const newIds = form.resource_ids || []
    const toAdd = newIds.filter((id: string) => !currentIds.includes(id))
    const toRemove = currentIds.filter((id: string) => !newIds.includes(id))

    for (const id of toRemove) {
      await unbindMenuResource(menuId, id)
    }
    for (const id of toAdd) {
      await bindMenuResource(menuId, { resource_id: id })
    }

    message.success(editingId.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchData()
    permission.fetchMenus()
  } catch (e: any) {
    message.error(e?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = (opt: TreeOption) => {
  deleteMenu(opt.id as string).then(() => { message.success('删除成功'); fetchData(); permission.fetchMenus() }).catch((e: any) => message.error(e?.response?.data?.message || '删除失败'))
}

const selectIcon = (icon: string) => {
  form.icon = icon
  showIconPicker.value = false
}

onMounted(async () => {
  fetchData()
  try {
    const { default: iconsData } = await import('@iconify-json/material-symbols/icons.json')
    const names = Object.keys((iconsData as any).icons || {})
    const aliases = Object.keys((iconsData as any).aliases || {})
    allIconNames.value = [...names, ...aliases].sort()
  } catch (e) {
    console.error('load icons failed', e)
  }
})
</script>

<style scoped>
.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(48px, 1fr));
  gap: 6px;
  max-height: 60vh;
  overflow-y: auto;
  padding: 4px;
}
.icon-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  cursor: pointer;
}
</style>