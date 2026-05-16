<template>
  <div class="logo" :class="{ 'logo--collapsed': collapsed }">
    <div class="logo-icon">
      <span class="i-material-symbols:shield-person-outline" />
    </div>
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
  height: var(--header-height, 64px);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid var(--color-border-light, #efe9e2);
  overflow: hidden;
  white-space: nowrap;
  flex-shrink: 0;
  transition: border-color 0.3s ease;
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #c2704a, #d97706);
  border-radius: 8px;
  color: #fff;
  font-size: 18px;
  flex-shrink: 0;
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #c2704a 0%, #d97706 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.logo--collapsed .logo-icon {
  width: 36px;
  height: 36px;
  font-size: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
