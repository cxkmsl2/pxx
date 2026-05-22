<template>
  <div v-if="!isLoggedIn" class="login-gate">
      <div class="gate-text">请先登录</div>
      <div class="gate-btn" @click="$router.push('/login')">立即登录</div>
    </div>
    <div class="vote" v-else-if="c">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">案件审判</span></div>
    <div class="case-header"><div class="case-title">{{ c.title }}</div></div>
    <div class="split">
      <div class="side red"><div class="side-label">原告 · {{ c.plaintiff?.nickname }}</div><div class="side-desc">{{ c.plaintiff_desc }}</div></div>
      <div class="side blue"><div class="side-label">被告 · {{ c.defendant?.nickname }}</div><div class="side-desc">{{ c.defendant_desc }}</div></div>
    </div>
    <div class="vote-bar" v-if="c.status === 1">
    <div class="vote-bar" v-if="c.status === 0"><div class="vote-hint">⏳ 等待评审员受理此案件</div></div>
      <div class="vote-btns">
        <div class="v-btn v-red" @click="doVote('plaintiff')">支持原告<br/><small>{{ c.plaintiff?.nickname }}</small></div>
        <div class="v-btn v-blue" @click="doVote('defendant')">支持被告<br/><small>{{ c.defendant?.nickname }}</small></div>
      </div>
      <div class="vote-hint" v-if="!voted">请选择你认为有理的一方投票</div>
      <div class="vote-hint" v-else>你已投票</div>
    </div>
    <div class="result" v-else>
      <div class="result-text">{{ c.winner === 'plaintiff' ? '原告胜诉' : '被告胜诉' }}</div>
      <div class="score">{{ c.plaintiff_votes }} : {{ c.defendant_votes }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../api'

const router = useRouter()
const isLoggedIn = !!localStorage.getItem('token')
const route = useRoute(); const c = ref<any>(null); onMounted(async () => { if (!localStorage.getItem('token')) { window.$toast('请先登录'); setTimeout(() => router.push('/login'), 1000); return }
  const res = await api.get('/disputes?page_size=50')
  if (res.code === 0) c.value = res.data.items.find((x: any) => x.id === Number(route.params.id))
})
const voted = ref(false)
const doVote = async (v: string) => {
  if (voted.value) { window.$toast('你已经投过票了'); return }
  const res = await api.post('/disputes/' + c.value.id + '/vote', { vote_for: v })
  if (res.code === 0) {
    voted.value = true
    window.$toast('投票成功！', 'success')
    c.value = { ...c.value, plaintiff_votes: c.value.plaintiff_votes + (v === 'plaintiff' ? 1 : 0), defendant_votes: c.value.defendant_votes + (v === 'defendant' ? 1 : 0), status: res.data.result?.includes('winner') ? 2 : 1, winner: res.data.result?.includes('plaintiff') ? 'plaintiff' : res.data.result?.includes('defendant') ? 'defendant' : '' }
  } else window.$toast(res.msg, 'error')
}
</script>

<style scoped>
.vote { background: #0F172A; min-height: 100vh; color: #E2E8F0; }
.nav { display: flex; gap: 10px; padding: 12px; }
.back { color: #3B82F6; font-size: 18px; cursor: pointer; }
.nav-title { font-weight: 600; }
.case-header { padding: 16px; } .case-title { font-size: 16px; font-weight: 700; }
.split { display: flex; margin: 0 14px; gap: 2px; }
.side { flex: 1; padding: 16px; min-height: 200px; }
.side.red { background: rgba(239,68,68,0.1); border-radius: 12px 0 0 12px; }
.side.blue { background: rgba(59,130,246,0.1); border-radius: 0 12px 12px 0; }
.side-label { font-size: 13px; font-weight: 600; margin-bottom: 8px; }
.side.red .side-label { color: #EF4444; }
.side.blue .side-label { color: #3B82F6; }
.side-desc { font-size: 12px; color: #94a3b8; line-height: 1.6; }
.vote-bar { padding: 20px 14px; }
.vote-btns { display: flex; gap: 16px; }
.v-btn { flex: 1; text-align: center; padding: 16px; border-radius: 14px; font-size: 15px; font-weight: 700; cursor: pointer; color: #fff; }
.v-btn small { font-size: 11px; font-weight: 400; opacity: 0.8; display: block; margin-top: 4px; }
.v-red { background: #EF4444; }
.v-blue { background: #3B82F6; }
.vote-hint { text-align: center; margin-top: 12px; font-size: 12px; color: #64748b; }
.result { text-align: center; padding: 30px; }
.result-text { font-size: 20px; font-weight: 700; color: #FBBF24; }
.score { font-size: 36px; font-weight: 800; margin-top: 8px; }
.login-gate { display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh; }
.gate-text { font-size: 16px; color: #94a3b8; margin-bottom: 16px; }
.gate-btn { padding: 10px 40px; background: #3B82F6; color: #fff; border-radius: 20px; font-size: 14px; cursor: pointer; }
</style>