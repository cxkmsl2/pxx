<template>
  <div class="publish">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="title">发布拼卡</span>
      <span class="submit" @click="doSubmit">{{ submitting ? '发布中' : '发布' }}</span>
    </div>
    <div class="form">
      <select v-model="form.brand" class="input">
        <option value="">选择品牌</option>
        <option>Netflix</option><option>B站</option><option>腾讯视频</option>
        <option>百度网盘</option><option>ChatGPT</option><option>Spotify</option>
        <option>饿了么</option><option>网易云音乐</option>
      </select>
      <input v-model="form.title" placeholder="标题（如：B站大会员月卡）" class="input" />
      <textarea v-model="form.desc" placeholder="描述（如：独享账号，稳定不掉）" class="input textarea" rows="2"></textarea>
      <div class="warn">⚠ 请确保出租结束后修改密码，保护账号安全</div>
      <input v-model="form.account" placeholder="账号（加密存储）" class="input" />
      <input v-model="form.password" placeholder="密码（加密存储）" class="input" type="password" />
      <div class="label">💰 价格</div>
      <input v-model="priceDay" placeholder="按天价格（元）" type="number" class="input" />
      <input v-model="priceWeek" placeholder="按周价格（元，可选）" type="number" class="input" />
      <div class="btn" @click="doSubmit">{{ submitting ? '发布中...' : '确认发布' }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../api'

const router = useRouter()
const submitting = ref(false)
const priceDay = ref('')
const priceWeek = ref('')
const form = ref({ brand: '', title: '', desc: '', account: '', password: '' })

const doSubmit = async () => {
  if (!form.value.title || !form.value.account || !form.value.password) { window.$toast('请完整填写'); return }
  submitting.value = true
  try {
    const res = await api.post('/subscriptions', {
      ...form.value,
      price_per_day: Math.round((Number(priceDay.value) || 0) * 100),
      price_per_week: Math.round((Number(priceWeek.value) || 0) * 100),
    })
    if (res.code === 0) { window.$toast('发布成功', 'success'); router.push('/subscription') }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('发布失败', 'error') }
  submitting.value = false
}
</script>

<style scoped>
.publish { background: #0F172A; min-height: 100vh; color: #fff; }
.nav { display: flex; justify-content: space-between; align-items: center; padding: 12px; background: #1E293B; }
.back, .submit { color: #3B82F6; font-size: 14px; cursor: pointer; }
.title { font-size: 15px; font-weight: 600; }
.form { padding: 16px; display: flex; flex-direction: column; gap: 10px; }
.input { padding: 12px; border: 1px solid #334155; border-radius: 10px; background: #1E293B; color: #F1F5F9; font-size: 14px; outline: none; width: 100%; box-sizing: border-box; }
.textarea { resize: vertical; min-height: 60px; }
.warn { font-size: 12px; color: #FBBF24; background: rgba(251,191,36,0.1); padding: 10px; border-radius: 8px; }
.label { font-size: 13px; color: #94a3b8; }
.btn { background: linear-gradient(135deg, #3B82F6, #8B5CF6); color: #fff; text-align: center; padding: 14px; border-radius: 10px; font-size: 16px; font-weight: 600; cursor: pointer; margin-top: 10px; }
</style>