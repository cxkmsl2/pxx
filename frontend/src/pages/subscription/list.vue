<template>
  <div class="subs">
    <div class="nav">🎫 拼卡局</div>
    <div class="brands">
      <span class="brand" :class="{ on: brand === b.val }" v-for="b in brands" :key="b.val" @click="onBrand(b.val)">{{ b.label }}</span>
    </div>
    <div class="loading" v-if="loading">加载中...</div>
    <div v-else>
      <div class="sub-card" v-for="item in items" :key="item.id" :class="'brand-' + (item.brand || '').replace(/[^a-zA-Z]/g,'').toLowerCase()">
        <div class="card-glow"></div>
        <div class="card-content">
          <div class="card-tag" :class="item.status === 1 ? 'free' : 'used'">
            <span class="dot"></span>{{ item.status === 1 ? '空闲中 · 秒发' : '使用中' }}
          </div>
          <div class="card-title">{{ item.title }}</div>
          <div class="card-desc">{{ item.desc }}</div>
          <div class="card-footer" @click="item.status === 1 && (expanded = expanded === item.id ? 0 : item.id)">
            <span class="price">¥{{ (item.price_per_day / 100).toFixed(2) }}<em>/天</em></span>
            <span class="price2" v-if="item.price_per_week">¥{{ (item.price_per_week / 100).toFixed(2) }}<em>/周</em></span>
            <span class="owner">{{ item.owner?.nickname }}</span>
          </div>
          <div class="expand-panel" v-if="expanded === item.id">
            <div class="exp-title">选择租用天数</div>
            <div class="exp-days">
              <span v-for="d in [1,3,7,14,30]" :key="d" class="exp-day" @click.stop="rent(item.id, d)">{{ d }}天</span>
            </div>
          </div>
        </div>
      </div>
      <div class="empty" v-if="items.length === 0">暂无拼卡</div>
    </div>
    <div class="fab" @click="$router.push('/subscription/publish')">＋</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../api'

const brand = ref(''); const items = ref<any[]>([]); const loading = ref(true); const expanded = ref(0)

const brands = [
  { val: '', label: '全部' }, { val: 'Netflix', label: 'Netflix' },
  { val: 'B站', label: 'B站' }, { val: '腾讯视频', label: '腾讯' },
  { val: '百度网盘', label: '网盘' }, { val: 'ChatGPT', label: 'ChatGPT' },
]

const onBrand = (b: string) => { brand.value = b; loadData() }

const rent = async (id: number, days: number) => {
  try {
    const res = await api.post('/subscriptions/' + id + '/rent', { days })
    if (res.code === 0) { window.$toast('租用成功', 'success'); expanded.value = 0; loadData() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('租用失败', 'error') }
}

const loadData = async () => {
  loading.value = true
  let url = '/subscriptions?page_size=30'
  if (brand.value) url += '&brand=' + encodeURIComponent(brand.value)
  const res = await api.get(url)
  if (res.code === 0) items.value = res.data.items
  loading.value = false
}

onMounted(() => loadData())
</script>

<style scoped>
.subs { background: #0F172A; min-height: 100vh; color: #fff; }
.nav { padding: 14px; font-size: 16px; font-weight: 700; }
.brands { display: flex; gap: 8px; padding: 0 14px 12px; overflow-x: auto; }
.brand { padding: 6px 14px; border-radius: 16px; font-size: 12px; background: #1E293B; color: #94a3b8; white-space: nowrap; cursor: pointer; }
.brand.on { background: #3B82F6; color: #fff; }
.loading, .empty { text-align: center; padding: 60px; color: #64748b; }
.sub-card { margin: 10px 14px; padding: 16px; border-radius: 16px; position: relative; overflow: hidden; background: rgba(30,41,59,0.8); backdrop-filter: blur(10px); border: 1px solid rgba(255,255,255,0.06); }
.card-glow { position: absolute; top: -20px; right: -20px; width: 100px; height: 100px; border-radius: 50%; filter: blur(40px); opacity: 0.3; }
.brand-netflix .card-glow { background: #E50914; }
.brand-b .card-glow { background: #FB7299; }
.card-content { position: relative; z-index: 1; }
.card-tag { font-size: 11px; margin-bottom: 8px; display: flex; align-items: center; gap: 6px; }
.card-tag .dot { width: 6px; height: 6px; border-radius: 50%; animation: breathe 1.5s infinite; }
.free { color: #34D399; } .free .dot { background: #34D399; }
.used { color: #F87171; } .used .dot { background: #F87171; }
@keyframes breathe { 0%,100% { opacity: 1; } 50% { opacity: 0.3; } }
.card-title { font-size: 16px; font-weight: 700; color: #F1F5F9; }
.card-desc { font-size: 12px; color: #94a3b8; margin-top: 4px; }
.card-footer { display: flex; align-items: center; gap: 10px; margin-top: 10px; cursor: pointer; }
.price { font-size: 20px; font-weight: 700; color: #34D399; }
.price em { font-size: 11px; font-style: normal; color: #64748b; margin-left: 2px; }
.price2 { font-size: 14px; color: #34D399; font-weight: 600; }
.price2 em { font-size: 10px; font-style: normal; color: #64748b; }
.owner { margin-left: auto; font-size: 11px; color: #64748b; }
.expand-panel { border-top: 1px solid rgba(255,255,255,0.1); padding-top: 12px; margin-top: 8px; }
.exp-title { font-size: 12px; color: #94a3b8; margin-bottom: 8px; }
.exp-days { display: flex; gap: 8px; }
.exp-day { padding: 8px 16px; background: rgba(59,130,246,0.2); color: #60A5FA; border-radius: 20px; font-size: 13px; font-weight: 600; cursor: pointer; }
.exp-day:active { background: #3B82F6; color: #fff; }
.fab { position: fixed; bottom: 100px; right: 20px; width: 56px; height: 56px; border-radius: 50%; background: linear-gradient(135deg, #3B82F6, #8B5CF6); color: #fff; font-size: 28px; display: flex; align-items: center; justify-content: center; box-shadow: 0 8px 24px rgba(59,130,246,0.4); cursor: pointer; z-index: 100; }
</style>
