<template>
  <div class="app">
    <div class="page-content"><router-view v-slot="{ Component }"><transition name="fade" mode="out-in"><component :is="Component" /></transition></router-view></div>
    <div class="tabbar">
      <div v-for="tab in tabs" :key="tab.path" class="tab-item" :class="{ active: $route.path === tab.path || $route.path.startsWith(tab.match) }" @click="$router.push(tab.path)">
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-text">{{ tab.text }}<i v-if="tab.badge && tab.badge > 0" class="tab-badge">{{ tab.badge }}</i></span>
      </div>
    </div>
    <div class="toast" v-if="toast.visible" :class="toast.type">{{ toast.msg }}</div>
    <div class="sheet-overlay" v-if="sheet.visible" @click="closeSheet"></div>
    <div class="sheet" v-if="sheet.visible">
      <div class="sheet-title" v-if="sheet.title">{{ sheet.title }}</div>
      <div class="sheet-opts">
        <div class="sheet-opt" v-for="(opt,i) in sheet.options" :key="i" @click="selectSheet(i)">{{ opt }}</div>
      </div>
      <div class="sheet-cancel" @click="closeSheet">取消</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, provide, ref, onMounted, computed } from 'vue'

const tabs = computed(() => [
  { path: '/', match: '/', icon: '🏠', text: '首页' },
  { path: '/tasks', match: '/tasks', icon: '🏃', text: '跑腿' },
  { path: '/messages', match: '/messages', icon: '💬', text: '消息', badge: unread.value },
  { path: '/user', match: '/user', text: '我的', icon: '👤' },
])

const unread = ref(0)
const toast = reactive({ visible: false, msg: '', type: 'info' })
const sheet = reactive({ visible: false, title: '', options: [] as string[], cb: null as any })
let timer = 0

const showToast = (msg: string, type = 'info') => {
  toast.visible = true; toast.msg = msg; toast.type = type
  clearTimeout(timer); timer = window.setTimeout(() => { toast.visible = false }, 2000)
}
const showSheet = (title: string, options: string[], cb: any) => {
  sheet.visible = true; sheet.title = title; sheet.options = options; sheet.cb = cb
}
const closeSheet = () => { sheet.visible = false }
const selectSheet = (i: number) => { closeSheet(); sheet.cb && sheet.cb(i) }

provide('toast', showToast); provide('sheet', showSheet)
window.$toast = showToast; window.$sheet = showSheet

const fetchUnread = async () => {
  if (!localStorage.getItem('token')) return
  try {
    const res = await (await fetch('/api/v1/messages/conversations', { headers: { Authorization: 'Bearer ' + localStorage.getItem('token') } })).json()
    if (res.code === 0) unread.value = (res.data || []).reduce((s: number, c: any) => s + (c.unread || 0), 0)
  } catch {}
}
onMounted(() => { fetchUnread(); setInterval(fetchUnread, 15000) })
</script>

<style>
:root {
  --primary: #1D4ED8; --primary-light: #3B82F6; --success: #10B981; --danger: #EF4444;
  --bg: #F2F2F6; --card: #FFFFFF; --text: #1C1C1E; --text-secondary: #8E8E93;
  --radius: 16px; --shadow: 0 8px 24px rgba(149,157,165,0.1);
}
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body { font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", "PingFang SC", sans-serif; background: var(--bg); -webkit-font-smoothing: antialiased; }
#app { max-width: 450px; margin: 0 auto; background: var(--bg); min-height: 100vh; position: relative; }
@media (min-width: 451px) { #app { box-shadow: 0 0 60px rgba(0,0,0,0.1); } }
.app { display: flex; flex-direction: column; min-height: 100vh; }
.page-content { flex: 1; padding-bottom: 64px; }

.tabbar { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 450px; height: 64px; display: flex; background: rgba(255,255,255,0.85); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px); z-index: 999; box-shadow: 0 -2px 12px rgba(0,0,0,0.04); }
.tab-item { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #8E8E93; font-size: 11px; cursor: pointer; user-select: none; transition: 0.15s; position: relative; }
.tab-item.active { color: var(--primary); }
.tab-icon { font-size: 22px; line-height: 1; margin-bottom: 2px; }
.tab-text { font-size: 10px; font-weight: 500; }
.tab-badge { position: absolute; top: 2px; right: calc(50% - 28px); background: var(--danger); color: #fff; font-size: 10px; font-style: normal; padding: 1px 5px; border-radius: 8px; min-width: 16px; text-align: center; }

.toast { position: fixed; top: 60px; left: 50%; transform: translateX(-50%); padding: 10px 24px; border-radius: 8px; font-size: 14px; color: #fff; z-index: 9999; white-space: nowrap; pointer-events: none; animation: toastIn 0.3s ease; backdrop-filter: blur(10px); }
.toast.info { background: rgba(30,41,59,0.88); }
.toast.success { background: rgba(16,185,129,0.88); }
.toast.error { background: rgba(239,68,68,0.88); }
@keyframes toastIn { from { opacity: 0; transform: translateX(-50%) translateY(-12px); } to { opacity: 1; transform: translateX(-50%) translateY(0); } }

.sheet-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); z-index: 3000; }
.sheet { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 450px; background: #fff; border-radius: 16px 16px 0 0; z-index: 3001; animation: slideUp 0.3s ease; padding-bottom: env(safe-area-inset-bottom); box-shadow: 0 -8px 32px rgba(0,0,0,0.1); }
@keyframes slideUp { from { transform: translateX(-50%) translateY(100%); } to { transform: translateX(-50%) translateY(0); } }
.sheet-title { text-align: center; padding: 16px; font-size: 13px; color: var(--text-secondary); }
.sheet-opt { text-align: center; padding: 16px; font-size: 16px; color: var(--primary); cursor: pointer; border-bottom: 0.5px solid #E5E5EA; }
.sheet-opt:active { background: #F8FAFC; }
.sheet-cancel { text-align: center; padding: 16px; font-size: 15px; color: var(--danger); cursor: pointer; margin-top: 6px; border-top: 6px solid var(--bg); }

.fade-enter-active, .fade-leave-active { transition: opacity 0.15s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
