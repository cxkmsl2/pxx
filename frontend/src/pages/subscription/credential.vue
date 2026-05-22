<template>
  <div class="cred" v-if="data">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="nav-title">凭证提取</span>
      <span class="sos" @click="showSOS = true">🆘</span>
    </div>

    <!-- 倒计时 -->
    <div class="timer" :class="{ expired: remaining <= 0 }">
      {{ remaining > 0 ? '剩余使用时间 ' + fmtTime(remaining) : '已过期' }}
    </div>

    <!-- 凭证区 -->
    <div class="card">
      <div class="field">
        <div class="f-label">账号</div>
        <div class="f-value" v-if="revealed" @click="copy(data.account)">{{ data.account }} <span class="copy">复制</span></div>
        <div class="f-value masked" v-else @touchstart.prevent="startReveal" @touchend.prevent="stopReveal" @mousedown="startReveal" @mouseup="stopReveal" @mouseleave="stopReveal">•••••••• 长按显示</div>
      </div>
      <div class="field">
        <div class="f-label">密码</div>
        <div class="f-value" v-if="revealed" @click="copy(data.password)">{{ data.password }} <span class="copy">复制</span></div>
        <div class="f-value masked" v-else @touchstart.prevent="startReveal" @touchend.prevent="stopReveal" @mousedown="startReveal" @mouseup="stopReveal" @mouseleave="stopReveal">•••••••• 长按显示</div>
      </div>
      <div class="hint" v-if="!revealed">🔐 长按任意密码区域显示凭证，松手即隐藏</div>
    </div>

    <!-- SOS Bottom Sheet -->
    <div class="overlay" v-if="showSOS" @click="showSOS = false"></div>
    <div class="bottomsheet" v-if="showSOS">
      <div class="bs-title">🚨 维权申诉</div>
      <div class="bs-opts">
        <div class="bs-opt" @click="doDispute('密码错误')">🔑 密码错误</div>
        <div class="bs-opt" @click="doDispute('被人顶号')">👥 被人顶号</div>
        <div class="bs-opt" @click="doDispute('账号无法登录')">❌ 账号无法登录</div>
      </div>
      <div class="bs-cancel" @click="showSOS = false">取消</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../../api'

const route = useRoute()
const data = ref<any>(null)
const revealed = ref(false)
const showSOS = ref(false)
const remaining = ref(0)
let revealTimer = 0
let countdownTimer = 0

const orderId = Number(route.params.id)

const fmtTime = (s: number) => {
  const h = Math.floor(s / 3600), m = Math.floor(s / 3600 * 60) % 60, sec = s % 60
  return `${String(h).padStart(2,'0')}:${String(m).padStart(2,'0')}:${String(sec).padStart(2,'0')}`
}

const startReveal = () => {
  revealed.value = true
  if (navigator.vibrate) navigator.vibrate(50)
  clearTimeout(revealTimer)
}
const stopReveal = () => {
  revealTimer = window.setTimeout(() => { revealed.value = false }, 8000)
}

const copy = (text: string) => {
  navigator.clipboard.writeText(text).then(() => window.$toast('已复制', 'success'))
}

const doDispute = async (reason: string) => {
  try {
    await api.post('/subscriptions/orders/' + orderId + '/dispute', { reason })
    window.$toast('申诉已提交', 'info')
  } catch { window.$toast('申诉失败', 'error') }
  showSOS.value = false
}

onMounted(async () => {
  try {
    const res = await api.get('/subscriptions/orders/' + orderId + '/credential')
    if (res.code === 0) {
      data.value = res.data
    } else {
      window.$toast(res.msg || '获取凭证失败', 'error')
    }
  } catch {
    window.$toast('请先登录或订单不存在', 'error')
  }
  // Get real expire time
  if (data.value?.rent_end) {
    const expireAt = new Date(data.value.rent_end).getTime()
    remaining.value = Math.max(0, Math.floor((expireAt - Date.now()) / 1000))
    countdownTimer = window.setInterval(() => {
      remaining.value = Math.max(0, Math.floor((expireAt - Date.now()) / 1000))
    }, 1000)
  }
})

onUnmounted(() => {
  clearInterval(countdownTimer)
  clearTimeout(revealTimer)
})
</script>

<style scoped>
.cred { background: #0F172A; min-height: 100vh; color: #fff; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #1E293B; }
.back { color: #3B82F6; font-size: 14px; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }
.sos { font-size: 22px; cursor: pointer; }

.timer { text-align: center; padding: 16px; font-size: 14px; font-weight: 600; background: rgba(239,68,68,0.15); color: #EF4444; margin: 12px 14px; border-radius: 12px; }
.timer.expired { background: rgba(100,116,139,0.2); color: #64748b; }

.card { margin: 12px 14px; padding: 20px; background: #1E293B; border-radius: 16px; }
.field { margin-bottom: 16px; }
.f-label { font-size: 12px; color: #94a3b8; margin-bottom: 6px; }
.f-value { font-size: 15px; font-family: monospace; background: #0F172A; padding: 12px; border-radius: 10px; color: #34D399; position: relative; }
.masked { color: #64748b; cursor: pointer; user-select: none; }
.copy { position: absolute; right: 10px; top: 50%; transform: translateY(-50%); font-size: 11px; color: #3B82F6; cursor: pointer; font-family: sans-serif; }
.hint { text-align: center; font-size: 12px; color: #64748b; margin-top: 8px; }

.overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 200; }
.bottomsheet { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 450px; background: #1E293B; border-radius: 20px 20px 0 0; padding: 20px; z-index: 201; }
.bs-title { font-size: 16px; font-weight: 700; margin-bottom: 16px; }
.bs-opts { display: flex; flex-direction: column; gap: 8px; }
.bs-opt { padding: 14px; background: #0F172A; border-radius: 12px; font-size: 14px; cursor: pointer; }
.bs-cancel { text-align: center; padding: 14px; color: #94a3b8; margin-top: 8px; cursor: pointer; }
</style>