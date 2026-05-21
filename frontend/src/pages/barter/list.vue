<template>
  <div class="page">
    <div class="nav">🔄 以物换物</div>
    <div v-if="loading">加载中...</div>
    <div v-else>
      <div class="card" v-for="item in items" :key="item.id">
        <div class="card-thumb" :class="'bg' + (item.id % 4)"></div>
        <div class="card-title">📤 {{ item.title }}</div>
        <div class="card-want">🔁 想换：{{ item.want_item }}</div>
        <div class="card-footer">
          <span class="card-meta">{{ item.user?.nickname }} · {{ item.created_at?.slice(0,10) }}</span>
          <span class="agree-btn" @click="doAgree(item.id)" v-if="item.user_id !== myId">同意置换</span>
        </div>
      </div>
      <div class="empty" v-if="items.length === 0">暂无置换物品</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api, { getBarterItems } from '../../api'

const router = useRouter()
const items = ref<any[]>([])
const loading = ref(true)
const myId = ref(0)

const loadData = async () => {
  const res: any = await getBarterItems({ page: 1 })
  if (res.code === 0) items.value = res.data.items
  loading.value = false
}

const doAgree = async (id: number) => {
  if (!localStorage.getItem('token')) { window.$toast('请先登录'); router.push('/login'); return }
  const place = prompt('约定见面地点（如：图书馆一楼）', '图书馆一楼')
  if (!place) return
  try {
    const res = await api.post('/barters/' + id + '/agree', {
      user_id: 0, meet_place: place, meet_time: new Date().toISOString().slice(0,19).replace("T"," ")
    })
    if (res.code === 0) window.$toast('置换订单已生成！', 'success')
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('操作失败', 'error') }
}

onMounted(() => loadData())
</script>

<style scoped>
.page { background: #f5f5f5; min-height: 100vh; }
.nav { padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #999; font-size: 13px; }

.card { margin: 10px 12px; padding: 14px; background: #fff; border-radius: 12px; }
.card-thumb { width: 100%; height: 100px; border-radius: 10px; margin-bottom: 12px; }
.card-title { font-size: 15px; font-weight: 600; color: #222; }
.card-want { font-size: 13px; color: #ff4d4f; margin-top: 6px; font-weight: 500; }
.card-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 12px; }
.card-meta { font-size: 11px; color: #aaa; }
.agree-btn {
  padding: 6px 14px; background: linear-gradient(135deg, #22a6b3, #7ed6df);
  color: #fff; border-radius: 14px; font-size: 12px; font-weight: 600; cursor: pointer;
}

.bg0 { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.bg1 { background: linear-gradient(135deg, #4834d4, #686de0); }
.bg2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.bg3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }
</style>
