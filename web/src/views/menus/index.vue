<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>菜单管理</n-h2>
      <n-button type="primary" @click="openModal()">新增菜单</n-button>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="false" :row-key="(row: any) => row.id" />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑菜单' : '新增菜单'" style="width: 600px">
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
            <span v-if="form.icon" :class="[toIconClass(form.icon), 'text-xl inline-block']" />
            <n-button @click="showIconPicker = true">选择</n-button>
          </div>
        </n-form-item>
        <n-form-item label="排序">
          <n-input-number v-model:value="form.order_num" :min="0" />
        </n-form-item>
        <n-form-item label="关联资源">
          <n-select v-model:value="form.resource_id" :options="resourceOptions" placeholder="请选择资源" :loading="resourcesLoading" />
        </n-form-item>
        <n-form-item label="是否显示">
          <n-switch v-model:value="form.is_visible" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSave">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showIconPicker" preset="card" title="选择图标" style="width: 650px">
      <div class="grid grid-cols-8 gap-8 max-h-400 overflow-y-auto">
        <div v-for="icon in commonIcons" :key="icon"
          class="flex-center p-8 border border-gray-200 rounded cursor-pointer hover:border-primary hover:text-primary"
          :class="{ 'border-primary text-primary': form.icon === icon }"
          @click="selectIcon(icon)">
          <Icon :icon="toIconifyName(icon)" class="text-xl" />
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed } from 'vue'
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NSwitch, NH2, NCard, NTag, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns, SelectOption } from 'naive-ui'
import { Icon } from '@iconify/vue'
import { getMenuTree, createMenu, updateMenu, deleteMenu, getResources } from '@/api'
import { usePermissionStore } from '@/stores/permission'
import type { MenuTreeNode } from '@/types'

const message = useMessage()
const permission = usePermissionStore()
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const showIconPicker = ref(false)
const editingId = ref<number | null>(null)

const toIconifyName = (icon?: string) => (icon || '').replace(/^i-/, '')
const toIconClass = (icon?: string) => icon ? (icon.startsWith('i-') ? icon : `i-${icon}`) : ''

const form = reactive({ name: '', parent_id: null as number | null, path: '', icon: '', order_num: 0, resource_id: null as number | null, is_visible: true })
const data = ref<MenuTreeNode[]>([])

const menuOptions = computed<SelectOption[]>(() => [{ label: '顶级菜单', value: 0 }, ...data.value.map((m: MenuTreeNode) => ({ label: m.name, value: m.id }))])
const resourceOptions = ref<SelectOption[]>([])
const resourcesLoading = ref(false)

const columns: DataTableColumns<MenuTreeNode> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '菜单名称', key: 'name' },
  { title: '路由路径', key: 'path', render: (row: MenuTreeNode) => row.path || '-' },
  { title: '图标', key: 'icon', render: (row: MenuTreeNode) => h(Icon, { icon: toIconifyName(row.icon) || 'material-symbols:circle-outline', class: 'text-lg' }) },
  { title: '排序', key: 'order_num', width: 80 },
  { title: '操作', key: 'actions', width: 160, render: (row: MenuTreeNode) => h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
    ]
  }) }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMenuTree()
    data.value = res.data.data || []
  } finally {
    loading.value = false
  }
}

const openModal = (row?: MenuTreeNode) => {
  if (row) {
    editingId.value = row.id
    Object.assign(form, { name: row.name, parent_id: null, path: row.path, icon: row.icon, order_num: row.order_num, resource_id: null, is_visible: row.is_visible })
  } else {
    editingId.value = null
    Object.assign(form, { name: '', parent_id: null, path: '', icon: '', order_num: 0, resource_id: null, is_visible: true })
  }
  showModal.value = true
  fetchResources()
}

const handleSave = async () => {
  saving.value = true
  try {
    const payload: any = { name: form.name, path: form.path, icon: form.icon, order_num: form.order_num, is_visible: form.is_visible }
    if (form.parent_id) (payload as any).parent_id = form.parent_id === 0 ? null : form.parent_id
    if (form.resource_id) (payload as any).resource_id = form.resource_id
    if (editingId.value) {
      await updateMenu(editingId.value, payload)
    } else {
      await createMenu(payload)
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

const handleDelete = (row: MenuTreeNode) => {
  deleteMenu(row.id).then(() => { message.success('删除成功'); fetchData(); permission.fetchMenus() }).catch((e: any) => message.error(e?.response?.data?.message || '删除失败'))
}

const commonIcons = [
  'i-material-symbols:dashboard-outline',
  'i-material-symbols:group-outline',
  'i-material-symbols:person-outline',
  'i-material-symbols:manage-accounts-outline',
  'i-material-symbols:menu-outline',
  'i-material-symbols:corporate-fare-outline',
  'i-material-symbols:security-outline',
  'i-material-symbols:history-outline',
  'i-material-symbols:settings-outline',
  'i-material-symbols:home-outline',
  'i-material-symbols:notifications-outline',
  'i-material-symbols:mail-outline',
  'i-material-symbols:search',
  'i-material-symbols:edit-outline',
  'i-material-symbols:delete-outline',
  'i-material-symbols:add-circle-outline',
  'i-material-symbols:refresh',
  'i-material-symbols:download',
  'i-material-symbols:upload',
  'i-material-symbols:print',
  'i-material-symbols:share',
  'i-material-symbols:star-outline',
  'i-material-symbols:favorite-outline',
  'i-material-symbols:lock-outline',
  'i-material-symbols:visibility-outline',
  'i-material-symbols:map',
  'i-material-symbols:bar-chart',
  'i-material-symbols:pie-chart',
  'i-material-symbols:table',
  'i-material-symbols:calendar-month',
  'i-material-symbols:description',
  'i-material-symbols:folder-outline',
  'i-material-symbols:file-upload',
  'i-material-symbols:article-outline',
  'i-material-symbols:build-outline',
  'i-material-symbols:help-outline',
  'i-material-symbols:info-outline',
  'i-material-symbols:warning-outline',
  'i-material-symbols:check-circle-outline',
  'i-material-symbols:cloud-outline',
  'i-material-symbols:link',
  'i-material-symbols:tag',
  'i-material-symbols:label-outline',
  'i-material-symbols:category-outline',
  'i-material-symbols:layers',
]

const fetchResources = async () => {
  resourcesLoading.value = true
  try {
    const res = await getResources({ page: 1, page_size: 100 })
    resourceOptions.value = (res.data.data?.data || []).map((r: any) => ({ label: r.description, value: r.id }))
  } catch { /* ignore */ }
  finally { resourcesLoading.value = false }
}

const selectIcon = (icon: string) => {
  form.icon = icon
  showIconPicker.value = false
}

onMounted(fetchData)
</script>
