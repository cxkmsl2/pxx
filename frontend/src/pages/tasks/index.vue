<template>
  <div class="tasks">
    <div class="nav">🏃 微跑腿大厅</div>
    <div class="hero">
      <div class="hero-text">校园最后 500 米 · 顺手赚零花</div>
      <div class="hero-sub">帮同学带饭、取快递、代打印</div>
    </div>
    <div class="quick-actions">
      <span class="action" :class="{ on: filter === f.val }" v-for="f in filters" :key="f.val" @click="filter = f.val; loadData()">{{ f.icon }} {{ f.label }}</span>
      <span class="action" :class="{ on: filter === 'my' && myTab === 0 }" @click="filter='my'; myTab=0; loadMyTasks()">📤 我发的</span>
      <span class="action" :class="{ on: filter === 'my' && myTab === 1 }" @click="filter='my'; myTab=1; loadMyTasks()">📥 我接的</span>
    </div>
    <div class="task-list">
      <div class="task-card" v-for="t in items" :key="t.id">
        <div class="task-reward" :class="t.urgent ? 'hot' : ''">¥{{ t.reward }}</div>
        <div class="task-body">
          <div class="task-title">{{ t.title }}</div>
          <div class="task-detail">{{ t.from }} → {{ t.to }}</div>
        </div>
        <div class="task-grab" v-if="filter!=='my'" @click="takeTask(t.id)">接单</div>
          <div class="task-grab done" v-else @click="completeTask(t.id)">完成</div>
      </div>
    </div>
    <div class="fab" @click="showMsg('发单功能开发中')">＋</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../api'
import { useUserStore } from '../../store'

const router = useRouter()
const store = useUserStore()
const filter = ref('all')
const items = ref<any[]>([])
const loading = ref(true)

const filters = [
  { val: 'all', icon: '🔥', label: '全部' },
  { val: 'food', icon: '🍱', label: '带饭' },
  { val: 'express', icon: '📦', label: '快递' },
  { val: 'print', icon: '🖨', label: '打印' },
  { val: 'book', icon: '📚', label: '还书' },
]

const taskTypes: Record<string, string> = { food: 'food', express: 'express', print: 'print', book: 'book' }

const loadData = async () => {
  loading.value = true
  let url = '/tasks?page_size=30'
  if (filter.value !== 'all') url += '&type=' + (taskTypes[filter.value] || filter.value)
  const res: any = await api.get(url)
  if (res.code === 0) items.value = res.data.items
  loading.value = false
}

const takeTask = async (id: number) => {
  if (!localStorage.getItem('token')) { window.$toast('请先登录'); router.push('/login'); return }
  try {
    const res = await api.post('/tasks/' + id + '/take')
    if (res.code === 0) { window.$toast('接单成功！', 'success'); loadData() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('接单失败', 'error') }
}

const showMsg = (m: string) => window.$toast(m)
onMounted(() => loadData())
</script>

<style scoped>
.tasks { background: #F8FAFC; min-height: 100vh; }
.nav { padding: 14px; background: rgba(255,255,255,0.9); backdrop-filter: blur(20px); font-size: 16px; font-weight: 700; color: #1D4ED8; position: sticky; top: 0; z-index: 10; }
.hero { padding: 20px 16px; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; }
.hero-text { font-size: 18px; font-weight: 700; }
.hero-sub { font-size: 13px; opacity: 0.8; margin-top: 4px; }
.quick-actions { display: flex; gap: 8px; padding: 12px 16px; overflow-x: auto; }
.action { padding: 6px 14px; border-radius: 16px; font-size: 12px; background: #fff; color: #64748b; white-space: nowrap; cursor: pointer; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
.action.on { background: #1D4ED8; color: #fff; }
.task-list { padding: 0 16px; }
.task-card { display: flex; align-items: center; gap: 12px; background: #fff; border-radius: 14px; padding: 14px; margin-bottom: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.task-reward { width: 44px; height: 44px; border-radius: 12px; background: #EFF6FF; color: #1D4ED8; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; flex-shrink: 0; }
.task-reward.hot { background: #FEF2F2; color: #EF4444; }
.task-body { flex: 1; min-width: 0; }
.task-title { font-size: 14px; font-weight: 600; color: #1E293B; }
.task-detail { font-size: 12px; color: #94a3b8; margin-top: 4px; }
.task-grab { padding: 6px 16px; background: linear-gradient(135deg, #10B981, #059669); color: #fff; border-radius: 16px; font-size: 12px; font-weight: 600; cursor: pointer; }
.fab { position: fixed; bottom: 100px; right: 20px; width: 56px; height: 56px; border-radius: 50%; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; font-size: 28px; display: flex; align-items: center; justify-content: center; box-shadow: 0 8px 24px rgba(29,78,216,0.3); cursor: pointer; z-index: 100; }

.done { background: #059669; }
</style>