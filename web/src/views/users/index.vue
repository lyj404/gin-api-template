<template>
  <div class="page-padding">
    <div class="toolbar-row mb-3">
      <n-h2 class="!my-0">用户管理</n-h2>
      <n-space wrap class="w-full md:w-auto">
        <n-input v-model:value="keyword" placeholder="搜索用户名/邮箱" clearable class="search-input" @keyup.enter="onSearch" @clear="onSearch" />
        <n-button @click="onSearch">搜索</n-button>
        <n-button type="primary" @click="openModal()">新增用户</n-button>
      </n-space>
    </div>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="data"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row: any) => row.id"
        :scroll-x="900"
        bordered
        single-column
        remote
        @update:page="handlePageChange"
      />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑用户' : '新增用户'" :style="{ width: '90vw', maxWidth: '520px' }">
      <n-form :model="form" :rules="rules" label-placement="left" label-width="90" ref="formRef">
        <n-form-item label="姓名" path="name">
          <n-input v-model:value="form.name" placeholder="请输入姓名" maxlength="50" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" maxlength="100" />
        </n-form-item>
        <n-form-item label="密码" path="password">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            :placeholder="editingId ? '留空则不修改' : '请输入密码'"
            maxlength="32"
          />
        </n-form-item>
        <n-form-item label="角色" path="role_ids">
          <n-select
            v-model:value="form.role_ids"
            multiple
            :options="roleOptions"
            placeholder="选择角色"
            clearable
          />
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
import {
  NButton, NDataTable, NModal, NForm, NFormItem, NInput, NH2, NCard,
  NTag, NSpace, NSelect, useMessage, useDialog
} from 'naive-ui'
import type { DataTableColumns, FormRules } from 'naive-ui'
import { getUsers, createUser, updateUser, deleteUser, getRoles } from '@/api'
import type { UserResponse, RoleResponse } from '@/types'

const message = useMessage()
const dialog = useDialog()
const formRef = ref()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const keyword = ref('')

const form = reactive<{
  name: string
  email: string
  password: string
  role_ids: number[]
}>({ name: '', email: '', password: '', role_ids: [] })

const data = ref<UserResponse[]>([])
const roleOptions = ref<{ label: string; value: number }[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 1,
  itemCount: 0,
  pageSlots: 5
})

const rules: FormRules = {
  name: { required: true, message: '请输入姓名', trigger: 'blur' },
  email: { required: true, type: 'email', message: '请输入有效邮箱', trigger: 'blur' },
  password: {
    required: false,
    trigger: 'blur',
    validator: (_rule, value) => {
      if (!editingId.value && !value) return new Error('请输入密码')
      if (value && value.length < 6) return new Error('密码至少 6 位')
      return true
    }
  }
}

const columns: DataTableColumns<UserResponse> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '姓名', key: 'name' },
  { title: '邮箱', key: 'email' },
  {
    title: '角色',
    key: 'roles',
    render: (row) => h(NSpace, { size: 4 }, {
      default: () => (row.roles || []).map(name =>
        h(NTag, { size: 'small', type: 'info' }, { default: () => name })
      )
    })
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render: (row) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => openModal(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getUsers({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: keyword.value || undefined
    })
    const payload = res.data.data as any
    data.value = payload?.data || []
    pagination.itemCount = payload?.total || 0
    pagination.pageCount = payload?.total_page || 1
  } catch (e: any) {
    message.error(e?.response?.data?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await getRoles({ page: 1, page_size: 100 })
    const payload = res.data.data as any
    const list: RoleResponse[] = payload?.data || []
    roleOptions.value = list.map(r => ({ label: r.name, value: r.id }))
  } catch (e) {
    // ignore
  }
}

const handlePageChange = (page: number) => {
  pagination.page = page
  fetchData()
}

const onSearch = () => {
  pagination.page = 1
  fetchData()
}

const openModal = (row?: UserResponse) => {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.email = row.email
    form.password = ''
    form.role_ids = [...(row.role_ids || [])]
  } else {
    editingId.value = null
    form.name = ''
    form.email = ''
    form.password = ''
    form.role_ids = []
  }
  showModal.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      const payload: any = {
        name: form.name,
        email: form.email,
        role_ids: form.role_ids
      }
      if (form.password) payload.password = form.password
      await updateUser(editingId.value, payload)
      message.success('更新成功')
    } else {
      await createUser({
        name: form.name,
        email: form.email,
        password: form.password,
        role_ids: form.role_ids
      })
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

const handleDelete = (row: UserResponse) => {
  dialog.warning({
    title: '确认删除',
    content: `确定删除用户 "${row.email}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteUser(row.id)
        message.success('删除成功')
        fetchData()
      } catch (e: any) {
        message.error(e?.response?.data?.message || '删除失败')
      }
    }
  })
}

onMounted(() => {
  fetchRoles()
  fetchData()
})
</script>

<style scoped>
.search-input {
  width: 100%;
}
@media (min-width: 768px) {
  .search-input {
    width: 220px;
  }
}
</style>
