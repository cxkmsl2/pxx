<template>
  <div class="orders">
    <div class="nav">📦 我的订单</div>
    <div class="tabs">
      <span class="tab" :class="{ active: tab === 0 }" @click="tab = 0">全部</span>
      <span class="tab" :class="{ active: tab === 1 }" @click="tab = 1">待付款</span>
      <span class="tab" :class="{ active: tab === 2 }" @click="tab = 2">已完成</span>
    </div>
    <div class="empty" v-if="filtered.length === 0">暂无订单</div>
    <div class="order-card" v-for="o in filtered" :key="o.id" @click="$router.push('/order/'+o.id)">
      <div class="order-head">
        <span class="order-no">{{ o.order_no }}</span>
        <span class="order-status" :class="'s' + o.status">{{ st(o.status) }}</span>
      </div>
      <div class="order-amount">¥{{ (o.amount / 100).toFixed(2) }}</div>
      <div class="order-time">{{ o.created_at?.slice(0,10) }} {{ o.created_at?.slice(11,19) }}</div>

      <div class="order-detail" v-if="expanded === o.id">
        <div class="dl"><span>买家ID</span><span>{{ o.buyer_id }}</span></div>
        <div class="dl"><span>卖家ID</span><span>{{ o.seller_id }}</span></div>
        <div class="dl"><span>商品ID</span><span>{{ o.product_id }}</span></div>
        <div class="dl"><span>金额</span><span>¥{{ (o.amount / 100).toFixed(2) }}</span></div>
        <div class="dl"><span>支付时间</span><span>{{ o.paid_at || '未支付' }}</span></div>
        <div class="dl"><span>完成时间</span><span>{{ o.completed_at || '未完成' }}</span></div>
        <div class="actions" v-if="o.status >= 2"><span class="act tribunal-act" @click.stop="toTribunal(o.id)">申请小法庭</span></div>
        <div class="actions" v-if="o.status === 1">
          <span class="act pay" @click.stop="doPay(o.id)">去支付</span>
          <span class="act cancel" @click.stop="doCancel(o.id)">取消订单</span>
        </div>
        <div class="actions" v-if="o.status === 2">
          <span class="act pay" @click.stop="doShip(o.id)">标记发货</span>
        </div>
        <div class="actions" v-if="o.status === 4"><span class="act rate-btn" @click.stop="doRate(o.id)">评价</span></div>
        <div class="actions" v-if="o.status === 3">
          <span class="act pay" @click.stop="doConfirm(o.id)">确认收货</span>
        </div>
        <div class="actions-old" v-if="o.status === 1" style="display:none">
          <span class="act pay" @click.stop="doPay(o.id)">去支付</span>
          <span class="act cancel" @click.stop="doCancel(o.id)">取消订单</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api, { getOrders } from '../../api'

const tab = ref(0)
const orders = ref<any[]>([])
const expanded = ref(0)

const st = (s: number) => ['', '待付款', '已付款', '已发货', '已完成', '已取消'][s] || '未知'

const filtered = computed(() => {
  if (tab.value === 0) return orders.value
  if (tab.value === 1) return orders.value.filter(o => o.status === 1)
  return orders.value.filter(o => o.status >= 4)
})

const toggleExpand = (id: number) => {
  expanded.value = expanded.value === id ? 0 : id
}

const doPay = async (id: number) => {
  try {
    const res = await api.put('/orders/' + id + '/status', { status: 2 },
      )
    if (res.code === 0) { window.$toast('支付成功', 'success'); loadOrders() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('支付失败', 'error') }
}

const doCancel = async (id: number) => {
  try {
    const res = await api.put('/orders/' + id + '/status', { status: 5 },
      )
    if (res.code === 0) { window.$toast('已取消', 'success'); loadOrders() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('取消失败', 'error') }
}

const loadOrders = async () => {
  const res: any = await getOrders({ page: 1 })
  if (res.code === 0) orders.value = res.data.items
}

const doShip = async (id: number) => {
  try {
    const res = await api.put('/orders/' + id + '/status', { status: 3 })
    if (res.code === 0) { window.$toast('已发货', 'success'); loadOrders() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('操作失败', 'error') }
}

const doConfirm = async (id: number) => {
  try {
    const res = await api.put('/orders/' + id + '/status', { status: 4 })
    if (res.code === 0) { window.$toast('收货成功', 'success'); loadOrders() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('操作失败', 'error') }
}

const showMsg = (m: string) => window.$toast(m)

onMounted(() => loadOrders())
</script>

<style scoped>
.orders { background: #F8FAFC; min-height: 100vh; }
.nav { padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.tabs { display: flex; background: #fff; padding: 0 12px 10px; gap: 10px; }
.tab { padding: 6px 16px; border-radius: 16px; font-size: 13px; background: #F8FAFC; color: #666; cursor: pointer; }
.tab.active { background: #1D4ED8; color: #fff; }
.empty { text-align: center; padding: 40px; color: #999; }

.order-card { margin: 8px 12px; padding: 14px; background: #fff; border-radius: 10px; cursor: pointer; }
.order-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.order-no { font-size: 12px; color: #999; }
.order-status { font-size: 11px; padding: 2px 8px; border-radius: 4px; }
.s1, .s5 { background: #F8FAFC; color: #999; }
.s2, .s3 { background: #EFF6FF; color: #1D4ED8; }
.s4 { background: #e8f8e8; color: #52c41a; }
.order-amount { font-size: 18px; font-weight: 700; color: #1D4ED8; }
.order-time { font-size: 11px; color: #bbb; margin-top: 4px; }

.order-detail { margin-top: 10px; padding-top: 10px; border-top: 1px solid #f0f0f0; }
.dl { display: flex; justify-content: space-between; padding: 4px 0; font-size: 12px; color: #666; }
.dl span:first-child { color: #999; }
.actions { display: flex; gap: 10px; margin-top: 10px; }
.act { padding: 5px 16px; border-radius: 14px; font-size: 12px; cursor: pointer; }
.pay { background: #1D4ED8; color: #fff; }
.cancel { background: #f0f0f0; color: #999; }

.tribunal-act { background: #FBBF24; color: #1E293B; }

.rate-btn { background: #FBBF24; color: #1E293B; }
</style>
