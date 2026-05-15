<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>审计日志</n-h2>
    </div>

    <n-card>
      <div class="mb-16">
        <div class="flex gap-8 items-center mb-12">
          <span class="text-sm text-gray-500 whitespace-nowrap">查询方式：</span>
          <n-select v-model:value="searchType" :options="searchTypeOptions" style="width: 140px" />
          <n-button @click="searchAll">刷新</n-button>
        </div>
        <div class="flex gap-8 items-center">
          <template v-if="searchType === 'operator'">
            <n-input v-model:value="operatorId" placeholder="请输入操作者ID" clearable style="width: 220px" />
            <n-button type="primary" @click="searchByOperator">查询</n-button>
          </template>
          <template v-else-if="searchType === 'target'">
            <n-select v-model:value="targetType" :options="targetTypeOptions" placeholder="目标类型" clearable style="width: 130px" />
            <n-input v-model:value="targetId" placeholder="目标ID" clearable style="width: 180px" />
            <n-button type="primary" @click="searchByTarget">查询</n-button>
          </template>
          <template v-else-if="searchType === 'time'">
            <n-date-picker v-model:value="timeRange" type="daterange" clearable style="width: 260px" />
            <n-button type="primary" @click="searchByTime">查询</n-button>
          </template>
          <span v-else class="text-sm text-gray-400">点击刷新按钮获取最新日志</span>
        </div>
      </div>

      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" bordered single-column />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { NButton, NDataTable, NCard, NInput, NSelect, NDatePicker, NH2, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getAuditLogs, getAuditLogsByTarget, getAuditLogsByTime } from '@/api'

const message = useMessage()
const searchType = ref('all')
const loading = ref(false)
const data = ref<any[]>([])
const operatorId = ref('')
const targetType = ref('')
const targetId = ref('')
const timeRange = ref<[number, number] | null>(null)
const pagination = reactive({ page: 1, pageSize: 10, itemCount: 0, pageCount: 1 })

const searchTypeOptions = [
  { label: '全部', value: 'all' },
  { label: '按操作者', value: 'operator' },
  { label: '按目标', value: 'target' },
  { label: '按时间', value: 'time' }
]

const targetTypeOptions = [
  { label: '角色', value: 'role' },
  { label: '用户', value: 'user' },
  { label: '菜单', value: 'menu' },
  { label: '组织', value: 'org_unit' },
  { label: '资源', value: 'resource' }
]

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '操作者', key: 'operator_name' },
  { title: '操作', key: 'action' },
  { title: '目标类型', key: 'target_type' },
  { title: '目标ID', key: 'target_id', width: 80 },
  { title: '描述', key: 'description' },
  { title: '时间', key: 'created_at', render: (row: any) => new Date(row.created_at).toLocaleString() }
]

const applyPage = (p: any) => {
  data.value = p.data || []
  pagination.page = p.page
  pagination.pageCount = p.total_page
  pagination.itemCount = p.total
}

const searchAll = async () => {
  loading.value = true
  try {
    const res = await getAuditLogs({ page: pagination.page, page_size: pagination.pageSize })
    applyPage(res.data.data)
  } catch { message.error('查询失败') } finally { loading.value = false }
}

const searchByOperator = async () => {
  if (!operatorId.value) { message.warning('请输入操作者ID'); return }
  loading.value = true
  try {
    const res = await getAuditLogs({ page: pagination.page, page_size: pagination.pageSize, operator_id: operatorId.value })
    applyPage(res.data.data)
  } catch { message.error('查询失败') } finally { loading.value = false }
}

const searchByTarget = async () => {
  if (!targetType.value || !targetId.value) { message.warning('请输入目标类型和ID'); return }
  loading.value = true
  try {
    const res = await getAuditLogsByTarget({ target_type: targetType.value, target_id: targetId.value })
    data.value = res.data.data?.data || []
  } catch { message.error('查询失败') } finally { loading.value = false }
}

const searchByTime = async () => {
  if (!timeRange.value) { message.warning('请选择时间范围'); return }
  loading.value = true
  try {
    const [start, end] = timeRange.value
    const fmt = (d: number) => new Date(d).toISOString().split('T')[0]
    const res = await getAuditLogsByTime({ start_time: fmt(start), end_time: fmt(end) })
    data.value = res.data.data?.data || []
  } catch { message.error('查询失败') } finally { loading.value = false }
}

onMounted(searchAll)
</script>
