<template>
  <div class="notif">
    <div class="nav">🔔 通知中心</div>
    <div class="loading" v-if="loading">加载中...</div>
    <div class="empty" v-if="!loading && items.length === 0">暂无通知</div>
    <div class="item" v-for="n in items" :key="n.id">
      <div class="n-icon">{{ n.from_user_id === 0 ? '📢' : '💬' }}</div>
      <div class="n-body">
        <div class="n-text">{{ n.content }}</div>
        <div class="n-time">{{ fmt(n.created_at) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../api'
const items = ref<any[]>([]); const loading = ref(true)
const fmt = (t: string) => t ? new Date(t).toLocaleString('zh-CN') : ''
onMounted(async () => {
  // Get system messages (from_user_id = 0)
  const res = await api.get('/messages/sys')
  if (res.code === 0) items.value = (res.data || []).reverse()
  loading.value = false
})
</script>

<style scoped>
.notif { background: var(--bg, #F2F2F6); min-height: 100vh; }
.nav { padding: 14px; background: #fff; font-size: 16px; font-weight: 700; }
.loading, .empty { text-align: center; padding: 40px; color: #8E8E93; }
.item { display: flex; gap: 10px; padding: 14px; background: #fff; margin: 2px 0; }
.n-icon { font-size: 20px; width: 32px; text-align: center; }
.n-body { flex: 1; }
.n-text { font-size: 14px; color: #1C1C1E; line-height: 1.4; }
.n-time { font-size: 11px; color: #8E8E93; margin-top: 4px; }
</style>