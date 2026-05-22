<template>
  <div class="qr-page">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">当面交易核销</span></div>
    <div class="qr-card">
      <div class="qr-title">请向买家出示此码</div>
      <div class="qr-img">[二维码]</div>
      <div class="qr-code">{{ orderNo }}</div>
      <div class="qr-hint">买家扫码后自动确认收货</div>
    </div>
    <div class="qr-actions">
      <span class="qr-btn" @click="copyCode">复制核销码</span>
      <span class="qr-btn primary" @click="simulateScan">模拟扫码</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const orderNo = ref(route.params.id || 'PX20240522001')

const copyCode = () => { navigator.clipboard.writeText(orderNo.value); window.$toast('已复制') }
const simulateScan = () => {
  window.$toast('核销成功！订单已收货', 'success')
  setTimeout(() => (window as any).location.href = '/#/orders', 1000)
}
</script>

<style scoped>
.qr-page { background: #1C1C1E; min-height: 100vh; color: #fff; }
.nav { display: flex; gap: 10px; padding: 12px; }
.back { color: #3B82F6; font-size: 18px; cursor: pointer; }
.nav-title { font-weight: 600; }
.qr-card { text-align: center; padding: 40px 20px; }
.qr-title { font-size: 14px; color: #8E8E93; margin-bottom: 20px; }
.qr-img { width: 200px; height: 200px; background: #fff; margin: 0 auto; border-radius: 16px; display: flex; align-items: center; justify-content: center; color: #1C1C1E; font-size: 48px; }
.qr-code { font-size: 20px; font-weight: 700; margin: 16px 0; letter-spacing: 4px; }
.qr-hint { font-size: 12px; color: #8E8E93; }
.qr-actions { display: flex; gap: 12px; padding: 20px 14px; }
.qr-btn { flex: 1; text-align: center; padding: 14px; border-radius: 12px; font-size: 14px; background: #2C2C2E; color: #fff; cursor: pointer; }
.qr-btn.primary { background: #3B82F6; }
</style>