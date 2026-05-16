<template>
  <div class="not-found">
    <div class="not-found-bg" />
    <div class="not-found-content">
      <div class="not-found-illustration">
        <div class="compass">
          <div class="compass-ring">
            <div v-for="i in 8" :key="i" class="compass-tick" :style="{ transform: `rotate(${i * 45}deg)` }" />
          </div>
          <div class="compass-needle needle-n" />
          <div class="compass-needle needle-s" />
          <div class="compass-center" />
          <div class="compass-text">404</div>
        </div>
        <div class="compass-shadow" />
      </div>
      <h1 class="not-found-title">页面迷失了方向</h1>
      <p class="not-found-desc">您访问的页面不存在，可能已被移除或地址有误</p>
      <div class="not-found-actions">
        <button class="btn btn-primary" @click="router.push('/')">
          <span class="btn-icon i-material-symbols:home-outline" />
          返回首页
        </button>
        <button class="btn btn-ghost" @click="router.push('/dashboard')">
          <span class="btn-icon i-material-symbols:dashboard-outline" />
          去仪表盘
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--color-bg, #f7f4f0);
}

.not-found-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 60% 50% at 50% 40%, rgba(194, 112, 74, 0.06) 0%, transparent 60%);
  pointer-events: none;
}

.not-found-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 20px;
  animation: fadeInUp 0.6s ease both;
}

.not-found-illustration {
  margin-bottom: 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.compass {
  width: 140px;
  height: 140px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.compass-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 3px solid var(--color-border-light, #efe9e2);
  animation: spin 20s linear infinite;
}

.dark .compass-ring {
  border-color: var(--color-border, #3d3631);
}

.compass-tick {
  position: absolute;
  top: -3px;
  left: 50%;
  width: 2px;
  height: 10px;
  margin-left: -1px;
  background: var(--color-text-muted, #a8a29e);
  transform-origin: center 73px;
  border-radius: 1px;
}

.compass-tick:nth-child(4n+1) {
  height: 14px;
  background: var(--color-primary, #c2704a);
}

.compass-needle {
  position: absolute;
  width: 4px;
  height: 48px;
  border-radius: 2px;
  transform-origin: center bottom;
}

.needle-n {
  top: 22px;
  background: linear-gradient(to top, var(--color-primary, #c2704a), #d97706);
  transform: rotate(-15deg);
  animation: wobble 3s ease-in-out infinite;
  z-index: 2;
}

.needle-s {
  top: 22px;
  background: var(--color-text-muted, #a8a29e);
  transform: rotate(165deg);
  animation: wobble 3s ease-in-out 0.1s infinite;
  z-index: 1;
  opacity: 0.5;
}

.compass-center {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--color-surface, #ffffff);
  border: 2.5px solid var(--color-primary, #c2704a);
  z-index: 3;
  position: absolute;
}

.compass-text {
  position: absolute;
  font-size: 28px;
  font-weight: 800;
  color: var(--color-primary, #c2704a);
  letter-spacing: -1px;
  opacity: 0.15;
  z-index: 0;
}

.compass-shadow {
  width: 100px;
  height: 12px;
  border-radius: 50%;
  background: rgba(194, 112, 74, 0.1);
  margin-top: 8px;
  filter: blur(6px);
  animation: pulse 2s ease-in-out infinite;
}

.dark .compass-shadow {
  background: rgba(217, 122, 74, 0.15);
}

.not-found-title {
  margin: 0 0 10px;
  font-size: 26px;
  font-weight: 700;
  color: var(--color-text, #1c1917);
}

.not-found-desc {
  margin: 0 0 32px;
  font-size: 15px;
  color: var(--color-text-secondary, #78716c);
  line-height: 1.6;
}

.not-found-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  text-decoration: none;
}

.btn-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.btn-primary {
  background: var(--color-primary, #c2704a);
  color: #fff;
  box-shadow: 0 2px 8px rgba(194, 112, 74, 0.25);
}

.btn-primary:hover {
  background: var(--color-primary-hover, #a85d38);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(194, 112, 74, 0.35);
}

.btn-primary:active {
  transform: translateY(0);
}

.btn-ghost {
  background: var(--color-surface-2, #faf8f5);
  color: var(--color-text-secondary, #78716c);
  border: 1.5px solid var(--color-border-light, #efe9e2);
}

.btn-ghost:hover {
  border-color: var(--color-primary-soft, rgba(194, 112, 74, 0.3));
  color: var(--color-primary, #c2704a);
  background: var(--color-primary-soft, rgba(194, 112, 74, 0.05));
  transform: translateY(-2px);
}

.dark .btn-ghost {
  background: var(--color-surface, #26211e);
  border-color: var(--color-border, #3d3631);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes wobble {
  0%, 100% { transform: rotate(-15deg); }
  50% { transform: rotate(5deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.1); }
}
</style>
