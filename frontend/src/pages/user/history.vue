<template>
  <div class="page">
    <div class="nav"><span class="back" @click="$router.back()">← 返回</span><span class="nav-title">浏览历史</span></div>
    <div class="empty" v-if="items.length === 0">暂无浏览记录</div>
    <div class="card" v-for="item in items" :key="item.id" @click="$router.push('/product/'+item.id)">
      <div class="c-title">{{ item.title }}</div>
      <div class="c-meta">{{ item.time }} · {{ item.campus }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../api'

const items = ref<any[]>([])

onMounted(async () => {
  const ids = JSON.parse(localStorage.getItem('pxx_history') || '[]')
  if (ids.length === 0) return
  for (const id of ids.slice(-20)) {
    try {
      const res: any = await api.get('/products/' + id)
      if (res.code === 0 && res.data) items.value.unshift({ ...res.data, time: new Date().toLocaleDateString() })
    } catch {}
  }
})
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; align-items: center; gap: 10px; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }
.empty { text-align: center; padding: 40px; color: #94a3b8; }
.card { margin: 8px 14px; padding: 14px; background: #fff; border-radius: 12px; cursor: pointer; }
.c-title { font-size: 14px; font-weight: 600; color: #1E293B; }
.c-meta { font-size: 12px; color: #94a3b8; margin-top: 4px; }
</style>