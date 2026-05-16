<template>
  <div class="page-padding">
    <div class="toolbar-row mb-3">
      <n-h2 class="!my-0">字典管理</n-h2>
      <n-space wrap class="w-full md:w-auto">
        <n-input v-model:value="nameFilter" placeholder="搜索名称" clearable class="search-input" @keyup.enter="onSearch" @clear="onSearch" />
        <n-input v-model:value="typeFilter" placeholder="搜索类型标识" clearable class="search-input" @keyup.enter="onSearch" @clear="onSearch" />
        <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="状态" clearable class="search-select" @update:value="onSearch" />
        <n-button @click="onSearch">搜索</n-button>
        <n-button type="primary" @click="openDictModal()">新增字典</n-button>
      </n-space>
    </div>

    <n-card>
      <n-data-table
        :columns="dictColumns"
        :data="dictList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row: any) => row.id"
        :scroll-x="800"
        bordered
        single-column
        remote
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        @update:checked-row-keys="handleDictSelect"
      />
    </n-card>

    <!-- 字典编辑弹窗 -->
    <n-modal v-model:show="showDictModal" preset="card" :title="editingDictId ? '编辑字典' : '新增字典'" :style="{ width: '90vw', maxWidth: '500px' }">
      <n-form ref="dictFormRef" :model="dictForm" :rules="dictRules" label-placement="left" label-width="80">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="dictForm.name" placeholder="如：用户状态" />
        </n-form-item>
        <n-form-item label="类型标识" path="type">
          <n-input v-model:value="dictForm.type" placeholder="如：sys_user_status" :disabled="!!editingDictId" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch :value="dictForm.status === 1" @update:value="v => dictForm.status = v ? 1 : 2">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="dictForm.desc" type="textarea" placeholder="可选备注" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showDictModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSaveDict">保存</n-button>
      </template>
    </n-modal>

    <!-- 字典详情抽屉 -->
    <n-drawer v-model:show="showDetailDrawer" :width="drawerWidth" placement="right">
      <n-drawer-content :title="'字典值 - ' + (selectedDict?.name || '')" closable>
        <template #header>
          <div class="flex items-center justify-between w-full pr-8">
            <span>字典值 - {{ selectedDict?.name || '' }}</span>
            <n-button size="small" type="primary" @click="openDetailModal()">新增值</n-button>
          </div>
        </template>

        <n-data-table
          :columns="detailColumns"
          :data="detailList"
          :loading="detailLoading"
          :row-key="(row: any) => row.id"
          :scroll-x="600"
          bordered
          single-column
        />

        <!-- 详情编辑弹窗 -->
        <n-modal v-model:show="showDetailModal" preset="card" :title="editingDetailId ? '编辑字典值' : '新增字典值'" :style="{ width: '90vw', maxWidth: '500px' }">
          <n-form ref="detailFormRef" :model="detailForm" :rules="detailRules" label-placement="left" label-width="80">
            <n-form-item label="标签" path="label">
              <n-input v-model:value="detailForm.label" placeholder="如：正常" />
            </n-form-item>
            <n-form-item label="值" path="value">
              <n-input v-model:value="detailForm.value" placeholder="如：1" />
            </n-form-item>
            <n-form-item label="排序">
              <n-input-number v-model:value="detailForm.sort" :min="0" placeholder="排序号" />
            </n-form-item>
            <n-form-item label="状态">
              <n-switch :value="detailForm.status === 1" @update:value="v => detailForm.status = v ? 1 : 2">
                <template #checked>启用</template>
                <template #unchecked>禁用</template>
              </n-switch>
            </n-form-item>
            <n-form-item label="备注">
              <n-input v-model:value="detailForm.remark" type="textarea" placeholder="可选备注" />
            </n-form-item>
          </n-form>
          <template #footer>
            <n-button @click="showDetailModal = false">取消</n-button>
            <n-button type="primary" :loading="detailSaving" class="ml-8" @click="handleSaveDetail">保存</n-button>
          </template>
        </n-modal>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NH2, NCard, NSwitch, NTag, NSelect, NSpace, NDrawer, NDrawerContent, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns, SelectOption } from 'naive-ui'
import { getDicts, getDict, getDictDetails, createDict, updateDict, deleteDict, createDictDetail, updateDictDetail, deleteDictDetail } from '@/api'
import type { DictResponse, DictDetailResponse, CreateDictRequest, UpdateDictRequest, CreateDictDetailRequest, UpdateDictDetailRequest } from '@/types'
import { useLayoutStore } from '@/stores/layout'

const message = useMessage()
const dialog = useDialog()
const layout = useLayoutStore()

const drawerWidth = computed<number | string>(() => layout.isMobile ? '100%' : 640)

// Dictionary list
const loading = ref(false)
const dictList = ref<DictResponse[]>([])
const selectedDict = ref<DictResponse | null>(null)
const nameFilter = ref('')
const typeFilter = ref('')
const statusFilter = ref<number | null>(null)
const pagination = reactive({ page: 1, pageSize: 10, pageCount: 1, itemCount: 0, pageSizes: [10, 20, 50, 100], showSizePicker: true })
const statusOptions: SelectOption[] = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 }
]

const dictColumns: DataTableColumns<DictResponse> = [
  { title: '序号', key: 'index', width: 70, render: (_row: DictResponse, index: number) => (pagination.page - 1) * pagination.pageSize + index + 1 },
  { title: '名称', key: 'name', width: 150 },
  { title: '类型标识', key: 'type', width: 180 },
  { title: '状态', key: 'status', width: 80, render(row) { return h(NTag, { type: row.status === 1 ? 'success' : 'error' }, { default: () => row.status === 1 ? '启用' : '禁用' }) } },
  { title: '描述', key: 'desc', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 200, fixed: 'right',
    render(row) {
      return [
        h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => openDetailDrawer(row) }, { default: () => '字典值' }),
        h(NButton, { size: 'tiny', quaternary: true, style: 'margin-left: 8px', onClick: () => openDictModal(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'tiny', quaternary: true, type: 'error', style: 'margin-left: 8px', onClick: () => handleDeleteDict(row) }, { default: () => '删除' })
      ]
    }
  }
]

const loadDicts = async () => {
  loading.value = true
  try {
    const res = await getDicts({
      page: pagination.page,
      page_size: pagination.pageSize,
      name: nameFilter.value || undefined,
      type: typeFilter.value || undefined,
      status: statusFilter.value || undefined
    })
    const p = res.data.data
    dictList.value = p.data || []
    pagination.page = p.page
    pagination.itemCount = p.total
    pagination.pageCount = p.total_page
  } catch (err: any) {
    message.error(err?.response?.data?.message || '获取字典列表失败')
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  pagination.page = 1
  loadDicts()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadDicts()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.page = 1
  pagination.pageSize = pageSize
  loadDicts()
}

const handleDictSelect = (keys: Array<string | number>) => {
  if (keys.length > 0) {
    selectedDict.value = dictList.value.find(d => d.id === keys[0]) || null
  } else {
    selectedDict.value = null
  }
}

// Dict modal
const showDictModal = ref(false)
const editingDictId = ref<string | null>(null)
const saving = ref(false)
const dictFormRef = ref<any>(null)
const dictForm = reactive<CreateDictRequest>({ name: '', type: '', status: 1, desc: '' })
const dictRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请输入类型标识', trigger: 'blur' }]
}

const openDictModal = (row?: DictResponse) => {
  if (row) {
    editingDictId.value = row.id
    dictForm.name = row.name
    dictForm.type = row.type
    dictForm.status = row.status
    dictForm.desc = row.desc || ''
  } else {
    editingDictId.value = null
    dictForm.name = ''
    dictForm.type = ''
    dictForm.status = 1
    dictForm.desc = ''
  }
  showDictModal.value = true
}

const handleSaveDict = async () => {
  try {
    await dictFormRef.value?.validate()
  } catch { return }
  saving.value = true
  try {
    if (editingDictId.value) {
      await updateDict(editingDictId.value, dictForm as UpdateDictRequest)
      message.success('更新成功')
    } else {
      await createDict(dictForm)
      message.success('创建成功')
    }
    showDictModal.value = false
    await loadDicts()
  } catch (err: any) {
    message.error(err?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const handleDeleteDict = (row: DictResponse) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除字典「${row.name}」吗？关联的字典值将一同被删除。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteDict(row.id)
        message.success('删除成功')
        await loadDicts()
      } catch (err: any) {
        message.error(err?.response?.data?.message || '删除失败')
      }
    }
  })
}

// Detail drawer
const showDetailDrawer = ref(false)
const detailLoading = ref(false)
const detailList = ref<DictDetailResponse[]>([])

const detailColumns: DataTableColumns<DictDetailResponse> = [
  { title: '标签', key: 'label', width: 120 },
  { title: '值', key: 'value', width: 120 },
  { title: '排序', key: 'sort', width: 60 },
  { title: '状态', key: 'status', width: 80, render(row) { return h(NTag, { type: row.status === 1 ? 'success' : 'error' }, { default: () => row.status === 1 ? '启用' : '禁用' }) } },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 160, fixed: 'right',
    render(row) {
      return [
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openDetailModal(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'tiny', quaternary: true, type: 'error', style: 'margin-left: 8px', onClick: () => handleDeleteDetail(row) }, { default: () => '删除' })
      ]
    }
  }
]

const openDetailDrawer = async (row: DictResponse) => {
  selectedDict.value = row
  showDetailDrawer.value = true
  detailLoading.value = true
  try {
    const res = await getDictDetails(row.id)
    detailList.value = res.data.data || []
  } catch {
    detailList.value = []
  } finally {
    detailLoading.value = false
  }
}

// Detail modal
const showDetailModal = ref(false)
const editingDetailId = ref<string | null>(null)
const detailSaving = ref(false)
const detailFormRef = ref<any>(null)
const detailForm = reactive<CreateDictDetailRequest>({ dict_id: '', label: '', value: '', sort: 0, status: 1, remark: '' })
const detailRules = {
  label: [{ required: true, message: '请输入标签', trigger: 'blur' }],
  value: [{ required: true, message: '请输入值', trigger: 'blur' }]
}

const openDetailModal = (row?: DictDetailResponse) => {
  if (row) {
    editingDetailId.value = row.id
    detailForm.dict_id = selectedDict.value!.id
    detailForm.label = row.label
    detailForm.value = row.value
    detailForm.sort = row.sort
    detailForm.status = row.status
    detailForm.remark = row.remark || ''
  } else {
    editingDetailId.value = null
    detailForm.dict_id = selectedDict.value!.id
    detailForm.label = ''
    detailForm.value = ''
    detailForm.sort = 0
    detailForm.status = 1
    detailForm.remark = ''
  }
  showDetailModal.value = true
}

const handleSaveDetail = async () => {
  try {
    await detailFormRef.value?.validate()
  } catch { return }
  detailSaving.value = true
  try {
    if (editingDetailId.value) {
      await updateDictDetail(selectedDict.value!.id, editingDetailId.value, detailForm as UpdateDictDetailRequest)
      message.success('更新成功')
    } else {
      await createDictDetail(selectedDict.value!.id, detailForm)
      message.success('创建成功')
    }
    showDetailModal.value = false
    const res = await getDictDetails(selectedDict.value!.id)
    detailList.value = res.data.data || []
  } catch (err: any) {
    message.error(err?.response?.data?.message || '操作失败')
  } finally {
    detailSaving.value = false
  }
}

const handleDeleteDetail = (row: DictDetailResponse) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除字典值「${row.label}」吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteDictDetail(selectedDict.value!.id, row.id)
        message.success('删除成功')
        const res = await getDictDetails(selectedDict.value!.id)
        detailList.value = res.data.data || []
      } catch (err: any) {
        message.error(err?.response?.data?.message || '删除失败')
      }
    }
  })
}

onMounted(loadDicts)
</script>

<style scoped>
.search-input { width: 100%; }
.search-select { width: 100%; }
@media (min-width: 768px) {
  .search-input { width: 180px; }
  .search-select { width: 130px; }
}
</style>
