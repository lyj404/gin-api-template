<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>角色管理</n-h2>
      <n-button type="primary" @click="openModal()">新增角色</n-button>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" />
    </n-card>

    <!-- 新建/编辑模态框 -->
    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑角色' : '新增角色'" style="width: 500px">
      <n-form :model="form" label-placement="left" label-width="80">
        <n-form-item label="角色名称">
          <n-input v-model:value="form.name" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.description" type="textarea" placeholder="请输入描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSave">保存</n-button>
      </template>
    </n-modal>

    <!-- 权限配置抽屉 -->
    <n-drawer v-model:show="showDrawer" :width="640" placement="right">
      <n-drawer-content :title="'权限配置 - ' + (configRole?.name || '')" closable>
        <n-tabs type="line" animated>
          <n-tab-pane name="menus" tab="菜单权限">
            <template v-if="menuLoading">
              <n-spin size="large" />
            </template>
            <template v-else>
              <n-space class="mb-12">
                <n-button size="tiny" @click="selectAllMenus(true)">全选</n-button>
                <n-button size="tiny" @click="selectAllMenus(false)">取消全选</n-button>
                <span class="text-sm text-gray-400 ml-8">已选 {{ checkedMenuIds.length }} / {{ totalMenuCount }}</span>
              </n-space>
              <n-tree
                :data="menuTree"
                key-field="id"
                label-field="name"
                children-field="children"
                selectable
                checkable
                cascade
                :checked-keys="checkedMenuIds"
                @update:checked-keys="checkedMenuIds = $event"
              />
              <n-space class="mt-16">
                <n-button type="primary" :loading="menuSaving" @click="saveMenus">保存菜单权限</n-button>
                <n-button @click="loadRoleConfig(configRole!.id)">重置</n-button>
              </n-space>
            </template>
          </n-tab-pane>
          <n-tab-pane name="resources" tab="资源权限">
            <template v-if="resLoading">
              <n-spin size="large" />
            </template>
            <template v-else>
              <n-space class="mb-12">
                <n-button size="tiny" @click="selectAllResources(true)">全选</n-button>
                <n-button size="tiny" @click="selectAllResources(false)">取消全选</n-button>
                <span class="text-sm text-gray-400 ml-8">已选 {{ selectedResCount }} / {{ allResources.length }}</span>
              </n-space>
              <n-data-table :columns="resColumns" :data="allResources" :loading="resLoading" />
              <n-space class="mt-16">
                <n-button type="primary" :loading="resSaving" @click="saveResources">保存资源权限</n-button>
                <n-button @click="loadRoleConfig(configRole!.id)">重置</n-button>
              </n-space>
            </template>
          </n-tab-pane>
        </n-tabs>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, computed, onMounted } from 'vue'
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NH2, NCard, NTag, NSpace, NDrawer, NDrawerContent, NTabs, NTabPane, NTree, NSpin, NSwitch, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns, TreeOption } from 'naive-ui'
import { getRoles, getRoleDetail, createRole, updateRole, deleteRole, getMenuTree, getResources, bindRoleMenu, unbindRoleMenu, bindRoleResource, unbindRoleResource } from '@/api'
import type { RoleResponse, ResourceResponse, ResourceBrief, MenuTreeNode } from '@/types'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({ name: '', description: '' })
const data = ref<RoleResponse[]>([])
const pagination = reactive({ page: 1, pageSize: 10, pageCount: 1, itemCount: 0 })

// 权限配置
const showDrawer = ref(false)
const configRole = ref<RoleResponse | null>(null)
const menuTree = ref<TreeOption[]>([])
const checkedMenuIds = ref<number[]>([])
const menuLoading = ref(false)
const menuSaving = ref(false)
const allResources = ref<ResourceResponse[]>([])
const resLoading = ref(false)
const resSaving = ref(false)
// 资源勾选状态: key=resource_id, value={checked, is_write}
const resourceChecked = ref<Record<number, { checked: boolean; is_write: boolean }>>({})

const columns: DataTableColumns<RoleResponse> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '角色名称', key: 'name' },
  { title: '描述', key: 'description' },
  { title: '系统角色', key: 'is_system', render: (row: RoleResponse) => h(NTag, { type: row.is_system ? 'success' : 'default', size: 'small' }, { default: () => row.is_system ? '是' : '否' }) },
  { title: '操作', key: 'actions', width: 200, render: (row: RoleResponse) => h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'warning', onClick: () => openConfig(row) }, { default: () => '权限配置' }),
      h(NButton, { size: 'small', type: 'error', disabled: row.is_system, onClick: () => handleDelete(row) }, { default: () => '删除' })
    ]
  }) }
]

const getResourcePerms = (res: ResourceResponse): { read: boolean; write: boolean } => {
  if (res.type === 'api') {
    if (res.method === 'GET' || res.method === 'HEAD') return { read: true, write: false }
    if (res.method === 'POST' || res.method === 'PUT' || res.method === 'DELETE' || res.method === 'PATCH') return { read: false, write: true }
    return { read: true, write: true } // * 或空
  }
  if (res.type === 'entity') {
    if (res.action === 'read') return { read: true, write: false }
    if (res.action === 'write' || res.action === 'delete') return { read: false, write: true }
    return { read: true, write: true } // * 或空
  }
  return { read: true, write: true }
}

const toggleResourcePerm = (row: ResourceResponse) => {
  const perms = getResourcePerms(row)
  if (!perms.read || !perms.write) return // 只有单一模式，不可切换
  if (!resourceChecked.value[row.id]) resourceChecked.value[row.id] = { checked: true, is_write: true }
  resourceChecked.value[row.id].is_write = !resourceChecked.value[row.id].is_write
  resourceChecked.value = { ...resourceChecked.value }
}

const resColumns: DataTableColumns<ResourceResponse> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '类型', key: 'type', width: 80, render: (row: ResourceResponse) => h(NTag, { type: row.type === 'api' ? 'info' : 'warning', size: 'small' }, { default: () => row.type }) },
  { title: '模式', key: 'pattern' },
  { title: '方法/操作', key: 'method', width: 90, render: (row: ResourceResponse) => row.method || row.action || '-' },
  { title: '选中', key: 'checked', width: 70, render: (row: ResourceResponse) => h('div', { style: 'display:flex;justify-content:center' }, [
    h('input', { type: 'checkbox', checked: resourceChecked.value[row.id]?.checked || false, onChange: (e: any) => {
      const perms = getResourcePerms(row)
      if (!resourceChecked.value[row.id]) resourceChecked.value[row.id] = { checked: false, is_write: perms.write }
      resourceChecked.value[row.id].checked = e.target.checked
      resourceChecked.value[row.id].is_write = perms.write
      resourceChecked.value = { ...resourceChecked.value }
    } })
  ]) },
  { title: '权限', key: 'perm', width: 120, render: (row: ResourceResponse) => {
    const checked = resourceChecked.value[row.id]?.checked
    if (!checked) return null
    const perms = getResourcePerms(row)
    const isWrite = resourceChecked.value[row.id]?.is_write ?? perms.write
    const canToggle = perms.read && perms.write
    if (perms.read && !perms.write) return h(NTag, { type: 'info', size: 'small' }, { default: () => '只读' })
    if (!perms.read && perms.write) return h(NTag, { type: 'warning', size: 'small' }, { default: () => '只写' })
    // 读写都支持 → 可点击切换
    return h(NButton, { size: 'tiny', tertiary: true, type: isWrite ? 'warning' : 'info', onClick: () => toggleResourcePerm(row) }, { default: () => isWrite ? '读写' : '只读' })
  } }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getRoles({ page: pagination.page, page_size: pagination.pageSize })
    const p = res.data.data
    data.value = p.data || []
    pagination.page = p.page
    pagination.itemCount = p.total
    pagination.pageCount = p.total_page
  } finally {
    loading.value = false
  }
}

const openModal = (row?: RoleResponse) => {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.description = row.description
  } else {
    editingId.value = null
    form.name = ''
    form.description = ''
  }
  showModal.value = true
}

const handleSave = async () => {
  saving.value = true
  try {
    if (editingId.value) {
      await updateRole(editingId.value, form)
      message.success('更新成功')
    } else {
      await createRole(form)
      message.success('创建成功')
    }
    showModal.value = false
    fetchData()
  } catch (e: any) {
    message.error(e?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = (row: RoleResponse) => {
  dialog.warning({
    title: '确认删除',
    content: `确定删除角色 "${row.name}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteRole(row.id)
        message.success('删除成功')
        fetchData()
      } catch (e: any) {
        message.error(e?.response?.data?.message || '删除失败')
      }
    }
  })
}

// 权限配置
const openConfig = async (row: RoleResponse) => {
  configRole.value = row
  showDrawer.value = true
  await loadRoleConfig(row.id)
}

const loadRoleConfig = async (roleId: number) => {
  checkedMenuIds.value = []
  resourceChecked.value = {}

  // 加载菜单树
  menuLoading.value = true
  try {
    const menuRes = await getMenuTree()
    menuTree.value = buildTreeOptions(menuRes.data.data || [])
  } finally {
    menuLoading.value = false
  }

  // 加载资源列表
  resLoading.value = true
  try {
    const resRes = await getResources({ page: 1, page_size: 200 })
    allResources.value = resRes.data.data?.data || []
  } finally {
    resLoading.value = false
  }

  // 加载角色当前绑定
  try {
    const detail = await getRoleDetail(roleId)
    const d = detail.data.data
    if (d.menus) {
      checkedMenuIds.value = d.menus.map(m => m.menu_id)
    }
    if (d.resources) {
      const rc: Record<number, { checked: boolean; is_write: boolean }> = {}
      for (const r of d.resources) {
        const res = allResources.value.find(res => res.id === r.resource_id)
        const perms = res ? getResourcePerms(res) : { read: true, write: false }
        rc[r.resource_id] = { checked: true, is_write: r.is_write && perms.write }
      }
      resourceChecked.value = rc
    }
  } catch (e: any) {
    message.error('加载角色权限失败')
  }
}

const selectedResCount = computed(() => Object.values(resourceChecked.value).filter(v => v.checked).length)

const selectAllResources = (select: boolean) => {
  const rc: Record<number, { checked: boolean; is_write: boolean }> = {}
  for (const res of allResources.value) {
    const perms = getResourcePerms(res)
    rc[res.id] = { checked: select, is_write: perms.write }
  }
  resourceChecked.value = rc
}

const collectMenuIds = (nodes: TreeOption[]): number[] => {
  const ids: number[] = []
  const walk = (list: TreeOption[]) => {
    for (const n of list) {
      ids.push(n.id as number)
      if (n.children) walk(n.children)
    }
  }
  walk(nodes)
  return ids
}

const totalMenuCount = computed(() => collectMenuIds(menuTree.value).length)

const selectAllMenus = (select: boolean) => {
  checkedMenuIds.value = select ? collectMenuIds(menuTree.value) : []
}

const buildTreeOptions = (nodes: MenuTreeNode[]): TreeOption[] => {
  return (nodes || []).map(n => ({
    id: n.id,
    name: n.name,
    label: n.name,
    key: n.id,
    children: n.children ? buildTreeOptions(n.children) : undefined
  }))
}

const saveMenus = async () => {
  if (!configRole.value) return
  menuSaving.value = true
  try {
    const detail = await getRoleDetail(configRole.value.id)
    const currentIds: number[] = detail.data.data.menus?.map(m => m.menu_id) || []
    const newIds = checkedMenuIds.value

    const toAdd = newIds.filter(id => !currentIds.includes(id))
    const toRemove = currentIds.filter(id => !newIds.includes(id))

    for (const id of toRemove) {
      await unbindRoleMenu(configRole.value.id, id)
    }
    for (const id of toAdd) {
      await bindRoleMenu(configRole.value.id, { menu_id: id })
    }
    message.success('菜单权限保存成功')
  } catch (e: any) {
    message.error(e?.response?.data?.message || '保存失败')
  } finally {
    menuSaving.value = false
  }
}

const saveResources = async () => {
  if (!configRole.value) return
  resSaving.value = true
  try {
    const detail = await getRoleDetail(configRole.value.id)
    const currentMap: Record<number, { is_write: boolean }> = {}
    for (const r of detail.data.data.resources || []) {
      currentMap[r.resource_id] = { is_write: r.is_write }
    }

    const entryArr = Object.entries(resourceChecked.value) as [string, { checked: boolean; is_write: boolean }][]
    for (const [resIdStr, val] of entryArr) {
      const resId = parseInt(resIdStr)
      const res = allResources.value.find(r => r.id === resId)
      const perms = res ? getResourcePerms(res) : { read: true, write: false }
      const isWrite = val.is_write && perms.write // 只有资源支持写才有效

      if (val.checked && !currentMap[resId]) {
        await bindRoleResource(configRole.value.id, resId, { is_write: isWrite })
      } else if (!val.checked && currentMap[resId]) {
        await unbindRoleResource(configRole.value.id, resId)
      } else if (val.checked && currentMap[resId] && currentMap[resId].is_write !== isWrite) {
        await unbindRoleResource(configRole.value.id, resId)
        await bindRoleResource(configRole.value.id, resId, { is_write: isWrite })
      }
    }
    message.success('资源权限保存成功')
  } catch (e: any) {
    message.error(e?.response?.data?.message || '保存失败')
  } finally {
    resSaving.value = false
  }
}

onMounted(fetchData)
</script>
