<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>角色管理</n-h2>
      <n-button type="primary" @click="openModal()">新增角色</n-button>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" />
    </n-card>

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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NH2, NCard, NTag, NSpace, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getRoles, createRole, updateRole, deleteRole } from '@/api'
import type { RoleResponse } from '@/types'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({ name: '', description: '' })
const data = ref<RoleResponse[]>([])
const pagination = reactive({ page: 1, pageSize: 10, pageCount: 1, itemCount: 0 })

const columns: DataTableColumns<RoleResponse> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '角色名称', key: 'name' },
  { title: '描述', key: 'description' },
  { title: '系统角色', key: 'is_system', render: (row: RoleResponse) => h(NTag, { type: row.is_system ? 'success' : 'default', size: 'small' }, { default: () => row.is_system ? '是' : '否' }) },
  { title: '操作', key: 'actions', width: 160, render: (row: RoleResponse) => h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', disabled: row.is_system, onClick: () => handleDelete(row) }, { default: () => '删除' })
    ]
  }) }
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

onMounted(fetchData)
</script>
