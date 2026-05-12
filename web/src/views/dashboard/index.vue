<template>
  <div class="p-16">
    <n-grid :cols="4" :x-gap="16" :y-gap="16">
      <n-gi>
        <n-card>
          <div class="text-gray-500 text-sm">总用户数</div>
          <div class="text-2xl font-bold mt-8">1,234</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <div class="text-gray-500 text-sm">活跃用户</div>
          <div class="text-2xl font-bold mt-8">892</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <div class="text-gray-500 text-sm">角色数</div>
          <div class="text-2xl font-bold mt-8">12</div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card>
          <div class="text-gray-500 text-sm">操作日志</div>
          <div class="text-2xl font-bold mt-8">5,678</div>
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
import { NGrid, NGi, NCard, NH3, NP } from 'naive-ui'
</script>
