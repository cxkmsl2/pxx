<template>
  <div class="page">
    <div class="nav">👥 宿舍拼团</div>
    <div v-if="loading">加载中...</div>
    <div v-else>
      <div class="card" v-for="gb in items" :key="gb.id">
        <div class="card-thumb" :class="'bg' + (gb.id % 4)"></div>
        <div class="card-title">{{ gb.title }}</div>
        <div class="desc">{{ gb.desc }}</div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: Math.min(gb.current_people / gb.min_people * 100, 100) + '%' }"></div>
        </div>
        <div class="card-footer">
          <div>
            <span class="price">¥{{ (gb.price_per_unit / 100).toFixed(2) }}</span>
            <span class="unit">/件</span>
          </div>
          <div class="people">{{ gb.current_people >= gb.min_people ? '🎉 已满' : '已拼 '+gb.current_people+'/'+gb.min_people+' 人' }}</div>
          <div class="join-btn" @click="doJoin(gb.id)" v-if="gb.sold_count < gb.total_stock">参团</div>
          <div class="join-btn sold-out" v-else>售罄</div>
        </div>
      </div>
      <div class="empty" v-if="items.length === 0">暂无拼团</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api, { getGroupBuys } from '../../api'

const router = useRouter()
const items = ref<any[]>([])
const loading = ref(true)

const loadData = async () => {
  const res: any = await getGroupBuys({ page: 1 })
  if (res.code === 0) items.value = res.data.items
  loading.value = false
}

const doJoin = async (id: number) => {
  if (!localStorage.getItem('token')) {
    window.$toast('请先登录', 'info')
    router.push('/login')
    return
  }
  try {
    const res = await api.post('/groupbuys/' + id + '/join', { quantity: 1 })
    if (res.code === 0) {
      window.$toast('参团成功！', 'success')
      // 延迟 500ms 等缓存失效后重新加载
      setTimeout(() => loadData(), 500)
    } else {
      window.$toast(res.msg, 'error')
    }
  } catch { window.$toast('参团失败', 'error') }
}

onMounted(() => loadData())
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; }
.nav { padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #999; font-size: 13px; }

.card { margin: 10px 12px; padding: 14px; background: #fff; border-radius: 12px; }
.card-thumb { width: 100%; height: 100px; border-radius: 10px; margin-bottom: 12px; }
.card-title { font-size: 15px; font-weight: 600; color: #222; }
.desc { font-size: 12px; color: #999; margin-top: 4px; }

.progress-bar {
  height: 6px; background: #f0f0f0; border-radius: 3px; margin: 10px 0; overflow: hidden;
}
.progress-fill { height: 100%; background: #1D4ED8; border-radius: 3px; transition: width 0.3s; }

.card-footer { display: flex; align-items: center; gap: 10px; margin-top: 6px; }
.price { font-size: 20px; font-weight: 700; color: #1D4ED8; }
.unit { font-size: 12px; color: #999; margin-left: 2px; }
.people { font-size: 12px; color: #666; }
.join-btn {
  margin-left: auto; padding: 6px 16px; background: linear-gradient(135deg, #3B82F6, #2563EB);
  color: #fff; border-radius: 16px; font-size: 12px; font-weight: 600; cursor: pointer;
}
.sold-out { background: #ccc; cursor: default; }

.bg0 { background: linear-gradient(135deg, #3B82F6, #2563EB); }
.bg1 { background: linear-gradient(135deg, #1E40AF, #686de0); }
.bg2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.bg3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }
</style>
