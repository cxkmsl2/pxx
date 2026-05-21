<template>
  <div class="app">
    <div class="page-content"><router-view /></div>
    <div class="tabbar">
      <div v-for="tab in tabs" :key="tab.path" class="tab-item" :class="{ active: $route.path === tab.path }" @click="$router.push(tab.path)">
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-text">{{ tab.text }}</span>
      </div>
    </div>
    <div class="toast" v-if="toast.visible" :class="toast.type">{{ toast.msg }}</div>
  </div>
</template>

<script setup lang="ts">
import { reactive, provide } from 'vue'

const tabs = [
  { path: '/', icon: '🏠', text: '首页' },
  { path: '/forum', icon: '💬', text: '论坛' },
  { path: '/orders', icon: '📦', text: '订单' },
  { path: '/user', icon: '👤', text: '我的' },
]

const toast = reactive({ visible: false, msg: '', type: 'info' })
let timer = 0

const showToast = (msg: string, type = 'info') => {
  toast.visible = true; toast.msg = msg; toast.type = type
  clearTimeout(timer)
  timer = window.setTimeout(() => { toast.visible = false }, 2000)
}

provide('toast', showToast)
window.$toast = showToast
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body { font-family: -apple-system, BlinkMacSystemFont, "PingFang SC", "Helvetica Neue", sans-serif; background: #f5f5f5; -webkit-font-smoothing: antialiased; }
#app { max-width: 450px; margin: 0 auto; background: #f5f5f5; min-height: 100vh; position: relative; }
@media (min-width: 451px) { #app { box-shadow: 0 0 30px rgba(0,0,0,0.12); } }

.app { display: flex; flex-direction: column; min-height: 100vh; }
.page-content { flex: 1; padding-bottom: 56px; }

.tabbar {
  position: fixed; bottom: 0; left: 50%; transform: translateX(-50%);
  width: 100%; max-width: 450px; height: 56px;
  display: flex; background: #fff; border-top: 1px solid #eee;
  z-index: 999; padding-bottom: env(safe-area-inset-bottom);
}
.tab-item {
  flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center;
  color: #999; font-size: 11px; cursor: pointer; user-select: none;
}
.tab-item.active { color: #ff4d4f; }
.tab-icon { font-size: 22px; line-height: 1; margin-bottom: 2px; }
.tab-text { font-size: 10px; }

.toast {
  position: fixed; top: 60px; left: 50%; transform: translateX(-50%);
  padding: 10px 24px; border-radius: 20px; font-size: 14px; color: #fff;
  z-index: 9999; white-space: nowrap; pointer-events: none;
  animation: fadeIn 0.2s ease;
}
.toast.info { background: rgba(0,0,0,0.75); }
.toast.success { background: rgba(82, 196, 26, 0.9); }
.toast.error { background: rgba(255, 77, 79, 0.9); }

@keyframes fadeIn { from { opacity: 0; transform: translateX(-50%) translateY(-10px); } to { opacity: 1; transform: translateX(-50%) translateY(0); } }
</style>
