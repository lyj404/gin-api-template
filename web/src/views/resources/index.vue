<template>
  <div class="p-16">
    <div class="flex justify-between items-center mb-16">
      <n-h2>资源管理</n-h2>
    </div>

    <n-card>
      <n-data-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" :row-key="(row: any) => row.id" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, h, reactive, onMounted } from 'vue'
import { NDataTable, NCard, NH2, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getResources } from '@/api'
import type { ResourceResponse } from '@/types'

const loading = ref(false)
const data = ref<ResourceResponse[]>([])
const pagination = reactive({ page: 1, pageSize: 20 })

const columns: DataTableColumns<ResourceResponse> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', key: 'name' },
  { title: '类型', key: 'type', render: (row: ResourceResponse) => h(NTag, { type: row.type === 'api' ? 'info' : 'warning', size: 'small' }, { default: () => row.type }) },
  { title: '模式', key: 'pattern' },
  { title: '方法', key: 'method', render: (row: ResourceResponse) => row.method ? h(NTag, { size: 'small', bordered: false }, { default: () => row.method }) : null },
  { title: '实体', key: 'entity' },
  { title: '操作', key: 'action' },
  { title: '描述', key: 'description' }
]

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getResources({ page: pagination.page, page_size: pagination.pageSize })
    data.value = res.data.data || []
    pagination.page = res.data.page
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
