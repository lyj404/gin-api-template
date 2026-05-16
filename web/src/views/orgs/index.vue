<template>
  <div class="page-padding">
    <div class="toolbar-row mb-3">
      <n-h2 class="!my-0">组织管理</n-h2>
      <n-button type="primary" @click="openModal(null)">新增组织</n-button>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="treeData" :loading="loading" :pagination="false" :scroll-x="700" bordered single-column />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑组织' : '新增组织'" :style="{ width: '90vw', maxWidth: '500px' }">
      <n-form :model="form" label-placement="left" label-width="80">
        <n-form-item label="组织名称">
          <n-input v-model:value="form.name" placeholder="请输入组织名称" />
        </n-form-item>
        <n-form-item label="上级组织">
          <n-select v-model:value="form.parent_id" :options="orgOptions" placeholder="请选择上级组织" clearable />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSave">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { NButton, NCard, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NH2, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns, SelectOption } from 'naive-ui'
import { getOrgTree, createOrgUnit, updateOrgUnit, deleteOrgUnit } from '@/api'
import type { OrgUnitResponse } from '@/types'

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({ name: '', parent_id: null as number | null })
const allOrgs = ref<OrgUnitResponse[]>([])

interface FlatNode {
  id: number
  name: string
  parent_id: number | null
  level: number
}

const buildFlat = (units: OrgUnitResponse[]): FlatNode[] => {
  return units.map(u => ({ id: u.id, name: u.name, parent_id: u.parent_id, level: u.level }))
}

const treeData = computed(() => buildFlat(allOrgs.value))
const orgOptions = computed<SelectOption[]>(() => [{ label: '顶级组织', value: 0 }, ...allOrgs.value.map((o: OrgUnitResponse) => ({ label: o.name, value: o.id }))])

const columns = computed<DataTableColumns<FlatNode>>(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '名称',
    key: 'name',
    render: (row: FlatNode) => {
      return h('span', { style: { marginLeft: `${row.level * 24}px` } }, row.name)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 280,
    render: (row: FlatNode) => {
      return h(NSpace, { size: 4 }, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '新增子级' }),
          h(NButton, { size: 'small', onClick: () => openModal(row, null, row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
        ]
      })
    }
  }
])

import { h } from 'vue'

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOrgTree()
    allOrgs.value = res.data.data || []
  } catch {} finally {
    loading.value = false
  }
}

const openModal = (node: FlatNode | null, _parent?: any, editNode?: FlatNode) => {
  if (editNode) {
    editingId.value = editNode.id
    form.name = editNode.name
    form.parent_id = editNode.parent_id || null
  } else {
    editingId.value = null
    form.name = ''
    form.parent_id = node ? (node.parent_id || node.id) : null
  }
  showModal.value = true
}

const handleSave = async () => {
  saving.value = true
  try {
    const payload: any = { name: form.name }
    if (form.parent_id) payload.parent_id = form.parent_id === 0 ? null : form.parent_id
    if (editingId.value) {
      await updateOrgUnit(editingId.value, payload)
    } else {
      await createOrgUnit(payload)
    }
    message.success(editingId.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchData()
  } catch (e: any) {
    message.error(e?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = (node: FlatNode) => {
  deleteOrgUnit(node.id).then(() => { message.success('删除成功'); fetchData() }).catch((e: any) => message.error(e?.response?.data?.message || '删除失败'))
}

onMounted(fetchData)
</script>
