<template>
  <div class="order-detail">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">订单详情</span></div>
    <div class="card">
      <div class="ord-no">{{ order.order_no }}</div>
      <div class="ord-amount">¥{{ (order.amount / 100).toFixed(2) }}</div>
      <div class="ord-status" :class="'s'+order.status">{{ statusText }}</div>
    </div>
    <div class="timeline">
      <div class="tl-item" v-for="(t,i) in steps" :key="i" :class="{ done: t.done, current: t.current }">
        <div class="tl-dot"></div>
        <div class="tl-content">
          <div class="tl-title">{{ t.label }}</div>
          <div class="tl-time" v-if="t.time">{{ t.time }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../../api'
const route = useRoute()
const order = ref<any>({})
const statusText = computed(() => ['','待付款','已付款','已发货','已完成','已取消'][order.value.status] || '未知')
const steps = computed(() => {
  const s = order.value.status || 0
  return [
    { label: '下单', done: s >= 1, current: s === 1, time: order.value.created_at },
    { label: '付款', done: s >= 2, current: s === 2, time: order.value.paid_at },
    { label: '发货', done: s >= 3, current: s === 3 },
    { label: '收货', done: s >= 4, current: s === 4, time: order.value.completed_at },
    { label: '评价', done: false, current: false },
  ]
})
onMounted(async () => {
  // For now, use orders list to find one
  const res = await api.get('/orders?page_size=50')
  if (res.code === 0) order.value = (res.data.items || []).find((o: any) => o.id === Number(route.params.id)) || {}
})
</script>

<style scoped>
.order-detail { background: var(--bg, #F2F2F6); min-height: 100vh; }
.nav { display: flex; gap: 10px; padding: 12px; background: #fff; }
.back { color: var(--primary, #1D4ED8); font-size: 16px; cursor: pointer; }
.nav-title { font-weight: 600; }
.card { margin: 12px 14px; padding: 20px; background: #fff; border-radius: 16px; text-align: center; box-shadow: var(--shadow, 0 8px 24px rgba(0,0,0,0.05)); }
.ord-no { font-size: 12px; color: #8E8E93; }
.ord-amount { font-size: 32px; font-weight: 800; color: var(--primary, #1D4ED8); margin: 8px 0; }
.ord-status { font-size: 13px; padding: 4px 16px; border-radius: 12px; display: inline-block; }
.s1, .s5 { background: #F1F5F9; color: #8E8E93; }
.s2, .s3 { background: #EFF6FF; color: #1D4ED8; }
.s4 { background: #D1FAE5; color: #059669; }
.timeline { margin: 12px 14px; padding: 20px; background: #fff; border-radius: 16px; }
.tl-item { display: flex; gap: 12px; padding: 12px 0; position: relative; }
.tl-item::before { content: ''; position: absolute; left: 9px; top: 30px; bottom: 0; width: 2px; background: #E5E5EA; }
.tl-item:last-child::before { display: none; }
.tl-item.done::before { background: var(--primary, #1D4ED8); }
.tl-dot { width: 20px; height: 20px; border-radius: 50%; background: #E5E5EA; flex-shrink: 0; position: relative; z-index: 1; }
.tl-item.done .tl-dot { background: var(--primary, #1D4ED8); }
.tl-item.current .tl-dot { background: #fff; border: 3px solid var(--primary, #1D4ED8); animation: pulse 1.5s infinite; }
@keyframes pulse { 0%,100% { box-shadow: 0 0 0 0 rgba(29,78,216,0.4); } 50% { box-shadow: 0 0 0 8px rgba(29,78,216,0); } }
.tl-title { font-size: 14px; font-weight: 600; color: #8E8E93; }
.tl-item.done .tl-title, .tl-item.current .tl-title { color: #1C1C1E; }
.tl-time { font-size: 11px; color: #8E8E93; margin-top: 2px; }
</style>