<template>
  <div class="tribunal">
    <div v-if="!isLoggedIn" class="login-gate">
      <div class="gate-text">请先登录</div>
      <div class="gate-btn" @click="$router.push('/login')">立即登录</div>
    </div>
    <div v-else>
    <div class="nav">⚖️ 校园小法庭</div>
    <div class="tabs">
      <span v-if="isReviewer" :class="{ on: tab === 0 }" @click="tab=0">待审核</span><span v-if="!isReviewer" :class="{ on: tab === 0 }" @click="tab=0">我的案件</span>
      <span :class="{ on: tab === 1 }" @click="tab=1">投票中</span>
      <span :class="{ on: tab === 2 }" @click="tab=2">已结案</span>
    </div>
    <div v-if="loading">加载中...</div>
    <div v-else>
      <div class="case-card" v-for="c in filtered" :key="c.id" @click="goDetail(c.id)">
        <div class="case-title">{{ c.title }}</div>
        <div class="case-parties">
          <span class="party p-red">{{ c.plaintiff?.nickname }}</span>
          <span class="vs">VS</span>
          <span class="party p-blue">{{ c.defendant?.nickname }}</span>
        </div>
        <div class="case-bar" v-if="c.status === 1 || c.status === 2">
          <div class="bar-red" :style="{ width: pct(c.plaintiff_votes, c.defendant_votes) }"></div>
          <div class="bar-blue" :style="{ width: pct(c.defendant_votes, c.plaintiff_votes) }"></div>
        </div>
        <div class="case-status" v-if="c.status === 0">⏳ 待评审员受理</div>
        <div class="case-status" v-if="c.status === 1">🔵 投票进行中</div>
        <div class="case-status" v-if="c.status === 2">{{ c.winner === 'plaintiff' ? '原告胜诉' : '被告胜诉' }}</div>
        <div class="accept-btn" v-if="isReviewer && c.status === 0" @click.stop="doAccept(c.id)">受理案件</div>
      </div>
      <div class="empty" v-if="filtered.length===0">暂无案件</div>
    </div>
  </div>
</div>
  </template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store'
import api from '../../api'

const isLoggedIn = !!localStorage.getItem('token')
const router = useRouter()

const store = useUserStore()
const isReviewer = store.userInfo?.role === 'admin'
const tab = ref(0); const cases = ref<any[]>([]); const loading = ref(true)
const pct = (a: number, b: number) => { const t = a + b; return t > 0 ? (a/t*100) + '%' : '0%' }
const myId = store.userInfo?.id
const filtered = computed(() => cases.value.filter(c => {
  if (tab.value === 0) {
    if (isReviewer) return c.status === 0
    return c.plaintiff_id === myId || c.defendant_id === myId
  }
  if (tab.value === 1) return c.status === 1
  return c.status === 2
}))

const showCreate = ref(false)
const newCase = ref({ title: '', order_id: 0, defendant_id: 0, plaintiff_desc: '', defendant_desc: '' })
const goDetail = (id: number) => router.push('/tribunal/' + id)

const doAccept = async (id: number) => {
  const res = await api.post('/disputes/' + id + '/accept')
  if (res.code === 0) { window.$toast('已受理，案件进入投票阶段', 'success'); load() }
  else window.$toast(res.msg, 'error')
}

const doCreate = async () => {
  const res = await api.post('/disputes', { ...newCase.value })
  if (res.code === 0) { window.$toast('案件已创建', 'success'); showCreate.value = false; load() }
  else window.$toast(res.msg, 'error')
}
const load = async () => {
  const res = await api.get('/disputes?page_size=30')
  if (res.code === 0) cases.value = res.data.items
  loading.value = false
}
onMounted(() => { if (!localStorage.getItem('token')) { router.push('/login') } else { load() } })
</script>

<style scoped>
.tribunal { background: #0F172A; min-height: 100vh; color: #E2E8F0; }
.nav { padding: 14px; font-size: 16px; font-weight: 700; }
.tabs { display: flex; gap: 12px; padding: 0 14px 12px; }
.tabs span { padding: 6px 16px; border-radius: 16px; font-size: 13px; background: #1E293B; color: #94a3b8; cursor: pointer; }
.tabs span.on { background: #3B82F6; color: #fff; }
.loading, .empty { text-align: center; padding: 60px; color: #64748b; }

.case-card { margin: 10px 14px; padding: 16px; background: #1E293B; border-radius: 14px; cursor: pointer; }
.case-title { font-size: 14px; font-weight: 600; }
.case-parties { display: flex; align-items: center; gap: 10px; margin: 10px 0; }
.party { padding: 4px 12px; border-radius: 8px; font-size: 12px; font-weight: 600; }
.p-red { background: rgba(239,68,68,0.2); color: #EF4444; }
.p-blue { background: rgba(59,130,246,0.2); color: #3B82F6; }
.vs { color: #64748b; font-size: 11px; }
.case-bar { display: flex; height: 6px; border-radius: 3px; overflow: hidden; }
.bar-red { background: #EF4444; height: 100%; transition: 0.5s; }
.bar-blue { background: #3B82F6; height: 100%; transition: 0.5s; }
.case-status { margin-top: 8px; font-size: 12px; color: #FBBF24; }
.admin-btn { font-size: 12px; padding: 4px 12px; background: #3B82F6; border-radius: 12px; cursor: pointer; margin-left: auto; color: #fff; font-weight: 400; }
.overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); z-index: 200; }
.create-form { position: fixed; top: 50%; left: 50%; transform: translate(-50%,-50%); background: #1E293B; padding: 20px; border-radius: 16px; width: 85%; max-width: 400px; z-index: 201; }
.form-title { font-size: 16px; font-weight: 700; margin-bottom: 12px; }
.inp { width: 100%; padding: 10px; margin-bottom: 8px; border: 1px solid #334155; border-radius: 8px; background: #0F172A; color: #E2E8F0; font-size: 13px; box-sizing: border-box; outline: none; }
.ta { height: 80px; resize: vertical; }
.btn { background: #3B82F6; color: #fff; text-align: center; padding: 10px; border-radius: 8px; font-weight: 600; cursor: pointer; margin-top: 4px; }
.cancel { text-align: center; padding: 10px; color: #64748b; cursor: pointer; }

.accept-btn { margin-top: 8px; padding: 6px 14px; background: #3B82F6; color: #fff; border-radius: 8px; font-size: 12px; font-weight: 600; cursor: pointer; display: inline-block; }
.login-gate { display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh; }
.gate-text { font-size: 16px; color: #94a3b8; margin-bottom: 16px; }
.gate-btn { padding: 10px 40px; background: #3B82F6; color: #fff; border-radius: 20px; font-size: 14px; cursor: pointer; }
</style>