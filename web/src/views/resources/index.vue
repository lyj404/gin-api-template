<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>资源管理</n-h2>
      <n-button type="primary" @click="openModal()">新增资源</n-button>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑资源' : '新增资源'" style="width: 550px">
      <n-form :model="form" label-placement="left" label-width="80">
        <n-form-item label="名称">
          <n-input v-model:value="form.name" placeholder="资源名称/标识符" />
        </n-form-item>
        <n-form-item label="类型">
          <n-select v-model:value="form.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="模式">
          <n-input v-model:value="form.pattern" placeholder="API路径或实体模式" />
        </n-form-item>
        <n-form-item label="方法">
          <n-select v-model:value="form.method" :options="methodOptions" clearable placeholder="仅API类型" />
        </n-form-item>
        <n-form-item label="实体">
          <n-input v-model:value="form.entity" placeholder="仅实体类型" />
        </n-form-item>
        <n-form-item label="操作">
          <n-input v-model:value="form.action" placeholder="仅实体类型" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.description" type="textarea" placeholder="资源描述" />
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
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NH2, NCard, NTag, NSpace, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns, SelectOption } from 'naive-ui'
import { getResources, createResource, updateResource, deleteResource } from '@/api'
import type { ResourceResponse } from '@/types'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({ name: '', type: 'api', pattern: '', method: '', entity: '', action: '', description: '' })
const data = ref<ResourceResponse[]>([])
const pagination = reactive({ page: 1, pageSize: 20, pageCount: 1, itemCount: 0 })

const typeOptions: SelectOption[] = [
  { label: 'API', value: 'api' },
  { label: '实体', value: 'entity' }
]

const methodOptions: SelectOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' }
]

const columns: DataTableColumns<ResourceResponse> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', key: 'name' },
  { title: '类型', key: 'type', render: (row: ResourceResponse) => h(NTag, { type: row.type === 'api' ? 'info' : 'warning', size: 'small' }, { default: () => row.type }) },
  { title: '模式', key: 'pattern' },
  { title: '方法', key: 'method', render: (row: ResourceResponse) => row.method ? h(NTag, { size: 'small', bordered: false }, { default: () => row.method }) : null },
  { title: '实体', key: 'entity' },
  { title: '操作', key: 'action' },
  { title: '描述', key: 'description' },
  { title: '操作', key: 'actions', width: 160, render: (row: ResourceResponse) => h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
    ]
  }) }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getResources({ page: pagination.page, page_size: pagination.pageSize })
    const p = res.data.data
    data.value = p.data || []
    pagination.page = p.page
    pagination.itemCount = p.total
    pagination.pageCount = p.total_page
  } finally {
    loading.value = false
  }
}

const openModal = (row?: ResourceResponse) => {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.type = row.type
    form.pattern = row.pattern
    form.method = row.method || ''
    form.entity = row.entity || ''
    form.action = row.action || ''
    form.description = row.description || ''
  } else {
    editingId.value = null
    form.name = ''
    form.type = 'api'
    form.pattern = ''
    form.method = ''
    form.entity = ''
    form.action = ''
    form.description = ''
  }
  showModal.value = true
}

const handleSave = async () => {
  saving.value = true
  try {
    if (editingId.value) {
      await updateResource(editingId.value, form as any)
      message.success('更新成功')
    } else {
      await createResource(form as any)
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

const handleDelete = (row: ResourceResponse) => {
  dialog.warning({
    title: '确认删除',
    content: `确定删除资源 "${row.name}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteResource(row.id)
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
