<template>
  <div class="page-padding">
    <n-grid cols="1 s:2 m:4" responsive="screen" :x-gap="16" :y-gap="16">
      <n-gi v-for="card in statsCards" :key="card.label">
        <n-card>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-gray-500 text-sm">{{ card.label }}</div>
              <div class="text-2xl font-bold mt-8">
                <n-spin v-if="loading" size="small" />
                <span v-else>{{ card.value.toLocaleString() }}</span>
              </div>
            </div>
            <span :class="['text-3xl', card.icon, card.color]" />
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid cols="1 m:3" responsive="screen" :x-gap="16" :y-gap="16" class="mt-16">
      <n-gi span="1 m:2">
        <n-card title="审计日志趋势（近 7 天）">
          <v-chart class="chart" :option="trendOption" autoresize />
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="资源类型分布">
          <v-chart class="chart" :option="pieOption" autoresize />
        </n-card>
      </n-gi>
    </n-grid>

    <n-card class="mt-16">
      <n-h3>欢迎使用后台管理系统</n-h3>
      <n-p depth="3">基于 Vue 3 + Naive UI + Go Gin 构建</n-p>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NGrid, NGi, NCard, NH3, NP, NSpin } from 'naive-ui'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent
} from 'echarts/components'
import { getDashboardStats, getAuditTrend } from '@/api'

use([CanvasRenderer, LineChart, PieChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent])

const loading = ref(true)
const stats = ref({ user_count: 0, role_count: 0, menu_count: 0, audit_log_count: 0 })
const trend = ref<{ date: string; count: number }[]>([])

const statsCards = computed(() => [
  { label: '用户总数', value: stats.value.user_count, icon: 'i-material-symbols:group-outline', color: 'text-blue-500' },
  { label: '角色总数', value: stats.value.role_count, icon: 'i-material-symbols:manage-accounts-outline', color: 'text-emerald-500' },
  { label: '菜单总数', value: stats.value.menu_count, icon: 'i-material-symbols:menu-outline', color: 'text-amber-500' },
  { label: '审计日志', value: stats.value.audit_log_count, icon: 'i-material-symbols:history-outline', color: 'text-violet-500' }
])

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 40, right: 20, top: 30, bottom: 30 },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: trend.value.map(t => t.date)
  },
  yAxis: { type: 'value', minInterval: 1 },
  series: [
    {
      name: '日志数',
      type: 'line',
      smooth: true,
      areaStyle: { opacity: 0.15 },
      lineStyle: { width: 2 },
      itemStyle: { color: '#2563eb' },
      data: trend.value.map(t => t.count)
    }
  ]
}))

const pieOption = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { bottom: 0 },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      data: [
        { value: stats.value.user_count, name: '用户' },
        { value: stats.value.role_count, name: '角色' },
        { value: stats.value.menu_count, name: '菜单' }
      ]
    }
  ]
}))

const loadStats = async () => {
  try {
    const res = await getDashboardStats()
    stats.value = res.data.data || stats.value
  } finally {
    loading.value = false
  }
}

const loadTrend = async () => {
  try {
    const res = await getAuditTrend()
    trend.value = res.data.data || []
  } catch (e) {
    trend.value = []
  }
}

onMounted(() => {
  loadStats()
  loadTrend()
})
</script>

<style scoped>
.chart {
  height: 320px;
  width: 100%;
}
@media (max-width: 767px) {
  .chart {
    height: 240px;
  }
}
</style>
