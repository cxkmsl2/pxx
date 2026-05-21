<template>
  <div class="page">
    <div class="nav">⏰ 物品租赁</div>
    <div v-if="loading">加载中...</div>
    <div v-else>
      <div class="card" v-for="item in items" :key="item.id">
        <div class="card-thumb" :class="'bg' + (item.id % 4)"></div>
        <div class="card-title">{{ item.title }}</div>
        <div class="card-desc">{{ item.desc }}</div>
        <div class="card-prices">
          <span class="price">¥{{ (item.daily_price / 100).toFixed(2) }}<em>/天</em></span>
          <span class="price2">¥{{ (item.weekly_price / 100).toFixed(2) }}<em>/周</em></span>
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
import { getRentalItems } from '../../api'
const items = ref<any[]>([])
const loading = ref(true)
onMounted(async () => {
  const res: any = await getRentalItems({ page: 1 })
  if (res.code === 0) items.value = res.data.items
  loading.value = false
})
</script>

<style scoped>
.page { background: #f5f5f5; min-height: 100vh; }
.nav { padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #999; font-size: 13px; }

.card { margin: 10px 12px; padding: 14px; background: #fff; border-radius: 12px; }
.card-thumb { width: 100%; height: 120px; border-radius: 10px; margin-bottom: 12px; }
.card-title { font-size: 15px; font-weight: 600; color: #222; }
.card-desc { font-size: 12px; color: #999; margin-top: 4px; line-height: 1.5; }
.card-prices { display: flex; gap: 12px; margin: 10px 0; align-items: baseline; }
.price { font-size: 20px; font-weight: 700; color: #ff4d4f; }
.price em { font-size: 11px; font-style: normal; color: #999; font-weight: 400; margin-left: 2px; }
.price2 { font-size: 14px; color: #ff4d4f; font-weight: 600; }
.price2 em { font-size: 10px; font-style: normal; color: #999; font-weight: 400; margin-left: 2px; }
.deposit { font-size: 12px; color: #666; }
.card-meta { font-size: 11px; color: #aaa; margin-top: 8px; }

.bg0 { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.bg1 { background: linear-gradient(135deg, #4834d4, #686de0); }
.bg2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.bg3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }
</style>
