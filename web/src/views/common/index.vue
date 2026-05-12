<template>
  <div class="p-16">
    <n-h2>{{ route.meta.title }}</n-h2>
    <n-p>功能开发中...</n-p>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { NH2, NP } from 'naive-ui'

const route = useRoute()
</script>
