<template>
  <div class="logo" :class="{ 'logo--collapsed': collapsed }">
    <span class="i-material-symbols:shield-person-outline text-2xl text-primary" />
    <transition name="fade">
      <span v-if="!collapsed" class="logo-text">Admin 后台</span>
    </transition>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  collapsed?: boolean
}>()
</script>

<style scoped>
.logo {
  height: var(--header-height, 60px);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid #efeff5;
  font-weight: 700;
  font-size: 16px;
  letter-spacing: 0.5px;
  overflow: hidden;
  white-space: nowrap;
  flex-shrink: 0;
}

.logo-text {
  background: linear-gradient(90deg, #18a058 0%, #36ad6a 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.text-primary {
  color: #18a058;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:global(.dark) .logo {
  border-bottom-color: #2a2a4a;
}
</style>
