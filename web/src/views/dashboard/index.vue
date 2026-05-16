<template>
  <div class="dashboard">
    <div class="stats-grid">
      <div
        v-for="card in statsCards"
        :key="card.label"
        class="stat-card"
      >
        <div class="stat-card-body">
            <div class="stat-icon" :class="`stat-icon--${card.color}`">
              <span :class="[card.icon, 'icon-inner']" />
            </div>
          <div class="stat-info">
            <div class="stat-value">
              <n-spin v-if="loading" size="small" />
              <span v-else>{{ card.value.toLocaleString() }}</span>
            </div>
            <div class="stat-label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card chart-card--wide">
        <div class="chart-card-header">
          <span class="chart-title">
            <span class="chart-dot chart-dot--primary" />
            审计日志趋势（近 7 天）
          </span>
        </div>
        <div class="chart-card-body">
          <v-chart class="chart" :option="trendOption" autoresize />
        </div>
      </div>
      <div class="chart-card">
        <div class="chart-card-header">
          <span class="chart-title">
            <span class="chart-dot chart-dot--amber" />
            资源类型分布
          </span>
        </div>
        <div class="chart-card-body">
          <v-chart class="chart" :option="pieOption" autoresize />
        </div>
      </div>
    </div>

    <div class="welcome-card">
      <div class="welcome-content">
        <div class="welcome-icon">
          <span class="i-material-symbols:rocket-launch-outline text-2xl" />
        </div>
        <div>
          <div class="welcome-title">欢迎使用后台管理系统</div>
          <div class="welcome-desc">基于 Vue 3 + Naive UI + Go Gin 构建</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NSpin } from 'naive-ui'
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
import { useThemeStore } from '@/stores/theme'

use([CanvasRenderer, LineChart, PieChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent])

const themeStore = useThemeStore()

const loading = ref(true)
const stats = ref({ user_count: 0, role_count: 0, menu_count: 0, resource_count: 0 })
const trend = ref<{ date: string; count: number }[]>([])

const statsCards = computed(() => [
  { label: '用户总数', value: stats.value.user_count, icon: 'i-material-symbols:group-outline', color: 'primary' },
  { label: '角色总数', value: stats.value.role_count, icon: 'i-material-symbols:manage-accounts-outline', color: 'success' },
  { label: '菜单总数', value: stats.value.menu_count, icon: 'i-material-symbols:menu', color: 'amber' },
  { label: '资源总数', value: stats.value.resource_count, icon: 'i-material-symbols:shield-outline', color: 'info' }
])

const isDark = computed(() => themeStore.isDark)

const trendOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: isDark.value ? '#26211e' : '#ffffff',
    borderColor: isDark.value ? '#3d3631' : '#efe9e2',
    textStyle: { color: isDark.value ? '#e7e5e4' : '#1c1917', fontSize: 12 }
  },
  grid: { left: 40, right: 16, top: 20, bottom: 28 },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    axisLine: { lineStyle: { color: isDark.value ? '#3d3631' : '#efe9e2' } },
    axisTick: { show: false },
    axisLabel: { color: isDark.value ? '#a8a29e' : '#78716c', fontSize: 11 },
    data: trend.value.map(t => t.date)
  },
  yAxis: {
    type: 'value',
    minInterval: 1,
    splitLine: { lineStyle: { color: isDark.value ? '#332d28' : '#f7f4f0', type: 'dashed' } },
    axisLabel: { color: isDark.value ? '#a8a29e' : '#78716c', fontSize: 11 }
  },
  series: [
    {
      type: 'line',
      smooth: true,
      lineStyle: { width: 2.5, color: '#c2704a' },
      itemStyle: { color: '#c2704a' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: isDark.value ? 'rgba(217, 122, 74, 0.25)' : 'rgba(194, 112, 74, 0.15)' },
            { offset: 1, color: isDark.value ? 'rgba(217, 122, 74, 0.02)' : 'rgba(194, 112, 74, 0.02)' }
          ]
        }
      },
      symbol: 'circle',
      symbolSize: 6,
      showSymbol: false,
      data: trend.value.map(t => t.count)
    }
  ]
}))

const pieOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    backgroundColor: isDark.value ? '#26211e' : '#ffffff',
    borderColor: isDark.value ? '#3d3631' : '#efe9e2',
    textStyle: { color: isDark.value ? '#e7e5e4' : '#1c1917', fontSize: 12 }
  },
  legend: {
    bottom: 0,
    textStyle: { color: isDark.value ? '#a8a29e' : '#78716c', fontSize: 11 },
    itemWidth: 10,
    itemHeight: 10
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '68%'],
      center: ['50%', '42%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: isDark.value ? '#26211e' : '#ffffff', borderWidth: 2 },
      label: { show: false },
      emphasis: {
        itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.15)' }
      },
      data: [
        { value: stats.value.user_count, name: '用户', itemStyle: { color: '#c2704a' } },
        { value: stats.value.role_count, name: '角色', itemStyle: { color: '#059669' } },
        { value: stats.value.menu_count, name: '菜单', itemStyle: { color: '#d97706' } }
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
.dashboard {
  width: 100%;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

@media (min-width: 640px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.stat-card {
  background: var(--color-surface, #ffffff);
  border-radius: 12px;
  border: 1px solid var(--color-border-light, #efe9e2);
  transition: all 0.25s ease;
}

.stat-card:hover {
  box-shadow: var(--shadow-md, 0 4px 12px rgba(28, 25, 23, 0.06));
  transform: translateY(-2px);
}

.stat-card-body {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-inner {
  width: 20px;
  height: 20px;
  display: inline-block;
  flex-shrink: 0;
}

.stat-icon--primary {
  background: rgba(194, 112, 74, 0.1);
  color: #c2704a;
}

.stat-icon--success {
  background: rgba(5, 150, 105, 0.1);
  color: #059669;
}

.stat-icon--amber {
  background: rgba(217, 119, 6, 0.1);
  color: #d97706;
}

.stat-icon--info {
  background: rgba(8, 145, 178, 0.1);
  color: #0891b2;
}

.dark .stat-icon--primary {
  background: rgba(217, 122, 74, 0.15);
  color: #d97a4a;
}

.dark .stat-icon--success {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.dark .stat-icon--amber {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.dark .stat-icon--info {
  background: rgba(34, 211, 238, 0.15);
  color: #22d3ee;
}

.stat-info {
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text, #1c1917);
  line-height: 1.2;
  margin-bottom: 2px;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-secondary, #78716c);
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  margin-bottom: 20px;
}

@media (min-width: 1024px) {
  .charts-grid {
    grid-template-columns: 2fr 1fr;
  }
}

.chart-card {
  background: var(--color-surface, #ffffff);
  border-radius: 12px;
  border: 1px solid var(--color-border-light, #efe9e2);
  overflow: hidden;
  transition: box-shadow 0.25s ease;
}

.chart-card:hover {
  box-shadow: var(--shadow-sm, 0 1px 2px rgba(28, 25, 23, 0.04));
}

.chart-card-header {
  padding: 16px 20px 0;
}

.chart-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text, #1c1917);
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.chart-dot--primary {
  background: #c2704a;
}

.chart-dot--amber {
  background: #d97706;
}

.chart-card-body {
  padding: 12px 8px 8px;
}

.chart {
  height: 300px;
  width: 100%;
}

@media (max-width: 767px) {
  .chart {
    height: 220px;
  }
}

.welcome-card {
  background: linear-gradient(135deg, rgba(194, 112, 74, 0.06), rgba(217, 119, 6, 0.04));
  border: 1px solid var(--color-border-light, #efe9e2);
  border-radius: 12px;
  padding: 20px 24px;
}

.dark .welcome-card {
  background: linear-gradient(135deg, rgba(217, 122, 74, 0.08), rgba(251, 191, 36, 0.04));
  border-color: var(--color-border, #3d3631);
}

.welcome-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.welcome-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--color-primary-soft, rgba(194, 112, 74, 0.08));
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-primary, #c2704a);
  flex-shrink: 0;
}

.welcome-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text, #1c1917);
  margin-bottom: 2px;
}

.welcome-desc {
  font-size: 13px;
  color: var(--color-text-secondary, #78716c);
}
</style>
