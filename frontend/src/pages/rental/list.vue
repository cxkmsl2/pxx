<template>
  <div class="page">
    <div class="nav">⏰ 物品租赁</div>
    <div v-if="loading">加载中...</div>
    <div v-else>
      <div class="card" v-for="item in items" :key="item.id" @click="router.push('/rental/' + item.id)">
        <div class="card-thumb" :class="'bg' + (item.id % 4)"></div>
        <div class="card-title">{{ item.title }}</div>
        <div class="card-desc">{{ item.desc }}</div>
        <div class="card-prices">
          <span class="price">¥{{ (item.daily_price / 100).toFixed(2) }}<em>/天</em></span>
          <span class="price2" v-if="item.weekly_price">¥{{ (item.weekly_price / 100).toFixed(2) }}<em>/周</em></span>
          <span class="deposit">押金 ¥{{ (item.deposit / 100).toFixed(2) }}</span>
        </div>
        <div class="card-meta">{{ item.owner?.nickname }} · {{ item.campus }}</div>
      </div>
      <div class="empty" v-if="items.length === 0">暂无租赁物品</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getRentalItems } from '../../api'

const router = useRouter()
const items = ref<any[]>([])
const loading = ref(true)

const goDetail = (id: number) => router.push('/product/' + id)

onMounted(async () => {
  const res: any = await getRentalItems({ page: 1 })
  if (res.code === 0) items.value = res.data.items
  loading.value = false
})
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; }
.nav { padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #94a3b8; }
.card { margin: 10px 14px; padding: 14px; background: #fff; border-radius: 14px; cursor: pointer; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.card-thumb { width: 100%; height: 120px; border-radius: 10px; margin-bottom: 12px; }
.bg0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.bg1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.bg2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.bg3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.card-title { font-size: 15px; font-weight: 600; color: #1E293B; }
.card-desc { font-size: 12px; color: #94a3b8; margin-top: 4px; }
.card-prices { display: flex; gap: 12px; margin: 10px 0; align-items: baseline; }
.price { font-size: 20px; font-weight: 700; color: #1D4ED8; }
.price em { font-size: 11px; font-style: normal; color: #94a3b8; margin-left: 2px; }
.price2 { font-size: 14px; color: #1D4ED8; font-weight: 600; }
.price2 em { font-size: 10px; font-style: normal; color: #94a3b8; }
.deposit { font-size: 12px; color: #64748b; }
.card-meta { font-size: 11px; color: #94a3b8; margin-top: 8px; }
</style>