<template>
  <div class="page" v-if="item">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">租赁详情</span></div>
    <div class="pic" :class="'bg' + (item.id % 4)"></div>
    <div class="section">
      <div class="title">{{ item.title }}</div>
      <p class="desc">{{ item.desc }}</p>
    </div>
    <div class="section">
      <div class="prices">
        <div class="price-item"><span class="val">¥{{ (item.daily_price / 100).toFixed(2) }}</span><span class="unit">/天</span></div>
        <div class="price-item" v-if="item.weekly_price"><span class="val">¥{{ (item.weekly_price / 100).toFixed(2) }}</span><span class="unit">/周</span></div>
        <div class="price-item"><span class="val dep">¥{{ (item.deposit / 100).toFixed(2) }}</span><span class="unit">押金</span></div>
      </div>
    </div>
    <div class="section">
      <div class="info"><span>👤 {{ item.owner?.nickname }}</span><span>📍 {{ item.campus }}</span></div>
    </div>
    <div class="btn" @click="rent">立即租用</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../../api'

const route = useRoute()
const item = ref<any>(null)

onMounted(async () => {
  const id = Number(route.params.id)
  const res = await api.get('/rentals/' + id)
  if (res.code === 0) item.value = res.data
})

const rent = () => {
  if (!localStorage.getItem('token')) { window.$toast('请先登录'); return }
  window.$sheet('选择租用天数', ['1天', '3天', '7天', '14天', '30天'], async (i: number) => {
    const days = [1,3,7,14,30][i]
    try {
      const res = await api.post('/rentals/' + item.value.id + '/order', { days })
      if (res.code === 0) { window.$toast('租用成功！', 'success') }
      else window.$toast(res.msg, 'error')
    } catch { window.$toast('租用失败', 'error') }
  })
}
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; padding-bottom: 80px; }
.nav { display: flex; gap: 10px; padding: 12px; background: #fff; }
.back { color: #1D4ED8; font-size: 16px; cursor: pointer; }
.nav-title { font-weight: 600; }
.pic { height: 240px; }
.bg0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.bg1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.bg2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.bg3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.section { background: #fff; margin: 10px 14px; padding: 14px; border-radius: 14px; }
.title { font-size: 16px; font-weight: 700; color: #1E293B; }
.desc { font-size: 13px; color: #64748b; margin-top: 4px; line-height: 1.6; }
.prices { display: flex; gap: 20px; }
.price-item { text-align: center; }
.val { font-size: 22px; font-weight: 700; color: #1D4ED8; }
.dep { color: #64748b; font-size: 16px; }
.unit { font-size: 11px; color: #94a3b8; display: block; }
.info { display: flex; gap: 16px; font-size: 13px; color: #64748b; }
.btn { margin: 20px 14px; background: #1D4ED8; color: #fff; text-align: center; padding: 14px; border-radius: 12px; font-size: 16px; font-weight: 600; cursor: pointer; }
</style>