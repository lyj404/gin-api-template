<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>资源管理</n-h2>
      <n-button type="primary" @click="openModal()">新增资源</n-button>
    </div>

    <n-card>
      <n-alert type="info" closable class="mb-16">
        <template #header>
          资源 = API 权限 或 实体权限。角色绑定资源后，该角色的用户即可访问。
        </template>
        API 类型控制接口访问，实体类型控制业务数据操作。
      </n-alert>

      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑资源' : '新增资源'" style="width: 580px">
      <n-alert type="info" closable class="mb-12">
        <template #header>
          {{ isApiType ? 'API 类型 — 控制接口访问权限' : '实体类型 — 控制业务数据操作权限' }}
        </template>
        <div v-if="isApiType">
          <div>名称建议: <n-tag size="small" bordered>模块:操作</n-tag> 如 <n-tag size="small" bordered>user:create</n-tag></div>
          <div class="mt-4">模式: API 路径，用 <n-tag size="small" bordered>/*</n-tag> 匹配子路径，如 <n-tag size="small" bordered>/users/*</n-tag></div>
        </div>
        <div v-else>
          <div>名称建议: <n-tag size="small" bordered>entity:{实体}:{操作}</n-tag> 如 <n-tag size="small" bordered>entity:user:read</n-tag></div>
          <div class="mt-4">由 <n-tag size="small" bordered>CheckEntityPermission</n-tag> 在代码中调用检查</div>
        </div>
      </n-alert>

      <n-form :model="form" label-placement="left" label-width="70">
        <n-form-item label="名称">
          <n-input v-model:value="form.name"
            :placeholder="isApiType ? '如: user:create' : '如: entity:user:read'"
            @update:value="autoFillExample" />
        </n-form-item>
        <n-form-item label="类型">
          <n-select v-model:value="form.type" :options="typeOptions" @update:value="onTypeChange" />
        </n-form-item>
        <n-form-item v-if="isApiType" label="方法">
          <n-select v-model:value="form.method" :options="methodOptions" clearable placeholder="GET / POST / PUT / DELETE / *" />
        </n-form-item>
        <n-form-item v-if="isApiType" label="模式">
          <n-input v-model:value="form.pattern" placeholder="如: /users、/users/*、/roles/:id" />
        </n-form-item>
        <n-form-item v-if="!isApiType" label="实体">
          <n-input v-model:value="form.entity" placeholder="如: user、role、order" />
        </n-form-item>
        <n-form-item v-if="!isApiType" label="操作">
          <n-input v-model:value="form.action" placeholder="如: read、write、delete、*" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.description" type="textarea" rows="3"
            :placeholder="isApiType ? '如: 创建用户、查看用户列表' : '如: 读取用户数据权限'" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space vertical class="w-full">
          <n-space justify="center" v-if="!editingId">
            <n-button size="tiny" quaternary @click="fillExample('api')">API 示例: 用户管理</n-button>
            <n-button size="tiny" quaternary @click="fillExample('api-role')">API 示例: 角色管理</n-button>
            <n-button size="tiny" quaternary @click="fillExample('entity')">实体示例: 用户</n-button>
          </n-space>
          <n-space justify="end">
            <n-button @click="showModal = false">取消</n-button>
            <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
          </n-space>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, computed, onMounted } from 'vue'
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NH2, NCard, NTag, NSpace, NAlert, useMessage, useDialog } from 'naive-ui'
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

const isApiType = computed(() => form.type === 'api')

const typeOptions: SelectOption[] = [
  { label: 'API — 接口权限', value: 'api' },
  { label: '实体 — 数据权限', value: 'entity' }
]

const methodOptions: SelectOption[] = [
  { label: 'GET (读取)', value: 'GET' },
  { label: 'POST (创建)', value: 'POST' },
  { label: 'PUT (更新)', value: 'PUT' },
  { label: 'DELETE (删除)', value: 'DELETE' },
  { label: '* (任意)', value: '*' }
]

const columns: DataTableColumns<ResourceResponse> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name', width: 150 },
  { title: '类型', key: 'type', width: 80, render: (row: ResourceResponse) => h(NTag, { type: row.type === 'api' ? 'info' : 'warning', size: 'small' }, { default: () => row.type === 'api' ? 'API' : '实体' }) },
  { title: '模式', key: 'pattern' },
  { title: '方法', key: 'method', width: 70, render: (row: ResourceResponse) => row.method ? h(NTag, { size: 'small', bordered: false }, { default: () => row.method }) : null },
  { title: '实体', key: 'entity', width: 80 },
  { title: '操作', key: 'action', width: 70 },
  { title: '描述', key: 'description', ellipsis: true },
  { title: '操作', key: 'actions', width: 120, render: (row: ResourceResponse) => h(NSpace, null, {
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
    resetForm()
  }
  showModal.value = true
}

const resetForm = () => {
  form.name = ''
  form.type = 'api'
  form.pattern = ''
  form.method = ''
  form.entity = ''
  form.action = ''
  form.description = ''
}

const onTypeChange = () => {
  if (form.type === 'api') {
    form.entity = ''
    form.action = ''
  } else {
    form.method = ''
    form.pattern = '*'
  }
}

const autoFillExample = (val: string) => {
  // 用户输入时清除旧提示
}

const fillExample = (type: string) => {
  if (type === 'api') {
    form.type = 'api'
    form.name = 'user:create'
    form.pattern = '/users'
    form.method = 'POST'
    form.entity = ''
    form.action = ''
    form.description = '创建用户'
  } else if (type === 'api-role') {
    form.type = 'api'
    form.name = 'role:manage'
    form.pattern = '/roles/*'
    form.method = '*'
    form.entity = ''
    form.action = ''
    form.description = '角色管理（所有操作）'
  } else if (type === 'entity') {
    form.type = 'entity'
    form.name = 'entity:user:read'
    form.pattern = '*'
    form.method = ''
    form.entity = 'user'
    form.action = 'read'
    form.description = '读取用户数据'
  }
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
