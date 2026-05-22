<template>
  <div class="page" v-if="item">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">置换详情</span></div>
    <div class="pic" :class="'bg' + (item.id % 4)"></div>
    <div class="section">
      <div class="title">📤 {{ item.title }}</div>
      <p class="desc">{{ item.desc }}</p>
    </div>
    <div class="section want">
      <div class="want-label">🔁 他想换</div>
      <div class="want-text">{{ item.want_item }}</div>
    </div>
    <div class="section">
      <div class="info"><span>👤 {{ item.user?.nickname }}</span><span class="credit-badge">信用一般</span></div>
    </div>
    <div class="btn" @click="agree">同意置换</div>
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
  const res = await api.get('/barters/' + id)
  if (res.code === 0) item.value = res.data
})
const agree = async () => {
  const place = prompt('约定见面地点', '图书馆一楼')
  if (!place) return
  try {
    const r = await api.post('/barters/' + item.value.id + '/agree', { meet_place: place, meet_time: new Date().toISOString().slice(0,19).replace('T',' ') })
    if (r.code === 0) window.$toast('置换订单已生成', 'success')
    else window.$toast(r.msg, 'error')
  } catch { window.$toast('操作失败', 'error') }
}
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; padding-bottom: 80px; }
.nav { display: flex; gap: 10px; padding: 12px; background: #fff; }
.back { color: #1D4ED8; font-size: 16px; cursor: pointer; }
.nav-title { font-weight: 600; }
.pic { height: 200px; }
.bg0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.bg1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.bg2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.bg3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.section { background: #fff; margin: 10px 14px; padding: 14px; border-radius: 14px; }
.title { font-size: 16px; font-weight: 700; color: #1E293B; }
.desc { font-size: 13px; color: #64748b; margin-top: 4px; line-height: 1.6; }
.want { background: #FEF2F2; }
.want-label { font-size: 12px; color: #EF4444; }
.want-text { font-size: 16px; font-weight: 700; color: #EF4444; margin-top: 4px; }
.info { display: flex; gap: 16px; font-size: 13px; color: #64748b; }
.credit-badge { font-size: 10px; background: #F1F5F9; color: #64748b; padding: 1px 6px; border-radius: 8px; }
.btn { margin: 20px 14px; background: #1D4ED8; color: #fff; text-align: center; padding: 14px; border-radius: 12px; font-size: 16px; font-weight: 600; cursor: pointer; }
</style>