<template>
  <div class="page">
    <div class="nav"><span class="back" @click="$router.back()">← 返回</span><span class="nav-title">我的收藏</span></div>
    <div class="loading" v-if="loading">加载中...</div>
    <div class="empty" v-if="!loading && items.length === 0">还没有收藏商品</div>
    <div v-else>
      <div class="card" v-for="item in items" :key="item.id" @click="$router.push('/product/'+item.id)">
        <div class="card-pic" :class="'bg' + (item.id % 4)"></div>
        <div class="card-right">
          <div class="card-title">{{ item.title }}</div>
          <div class="card-price">¥{{ (item.price / 100).toFixed(2) }}</div>
          <div class="card-meta">{{ item.campus }} · {{ item.seller?.nickname }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../api'

const items = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  const ids: number[] = JSON.parse(localStorage.getItem('pxx_favs') || '[]')
  if (ids.length === 0) { loading.value = false; return }
  for (const id of ids.slice(-30)) {
    try {
      const res: any = await api.get('/products/' + id)
      if (res.code === 0 && res.data) items.value.push(res.data)
    } catch {}
  }
  loading.value = false
})
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; align-items: center; gap: 10px; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #94a3b8; }

.card { display: flex; gap: 12px; margin: 10px 14px; padding: 12px; background: #fff; border-radius: 14px; cursor: pointer; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.card-pic { width: 80px; height: 80px; border-radius: 10px; flex-shrink: 0; }
.bg0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.bg1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.bg2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.bg3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.card-right { flex: 1; min-width: 0; }
.card-title { font-size: 14px; font-weight: 600; color: #1E293B; }
.card-price { font-size: 17px; font-weight: 700; color: #10B981; margin-top: 4px; }
.card-meta { font-size: 11px; color: #94a3b8; margin-top: 2px; }
</style>