<template>
  <div class="forum">
    <div class="nav"><span>💬 校园论坛</span><span class="pub-link" @click="$router.push('/post/new')">＋ 发帖</span></div>
    <div class="tabs">
      <span class="tab" :class="{ active: tab === 0 }" @click="switchTab(0)">求购悬赏</span>
      <span class="tab" :class="{ active: tab === 1 }" @click="switchTab(1)">好物种草</span>
    </div>
    <div class="skeleton" v-if="loading">
      <div class="sk-item" v-for="i in 3" :key="i"><div class="sk-line w80"></div><div class="sk-line w60"></div></div>
    </div>
    <div v-else>
      <div class="post-card" v-for="p in posts" :key="p.id" @click="$router.push('/forum/'+p.id)">
        <div class="post-badge" :class="p.type === 1 ? 'wanted' : 'share'">
          {{ p.type === 1 ? '求' : '荐' }}
        </div>
        <div class="post-body">
          <div class="post-title">{{ p.title }}</div>
          <div class="post-meta">{{ p.user?.nickname }} · {{ p.created_at?.slice(0,10) }}</div>
        </div>
      </div>
      <div class="empty" v-if="posts.length === 0">暂无帖子</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPosts } from '../../api'

const tab = ref(0)
const posts = ref<any[]>([])
const loading = ref(true)

const fetch = async () => {
  loading.value = true; posts.value = []
  try {
    const res: any = await getPosts({ page: 1, type: tab.value + 1 })
    if (res.code === 0) posts.value = res.data.items
  } catch (e) { console.error(e) }
  loading.value = false
}

const switchTab = (t: number) => { tab.value = t; fetch() }

onMounted(() => fetch())
</script>

<style scoped>
.forum { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; justify-content: space-between; align-items: center; padding: 12px; background: #fff; font-size: 15px; font-weight: 600; }
.pub-link { color: #1D4ED8; font-size: 13px; cursor: pointer; }
.tabs { display: flex; background: #fff; padding: 0 12px 10px; gap: 10px; }
.tab { padding: 6px 16px; border-radius: 16px; font-size: 13px; background: #F8FAFC; color: #666; cursor: pointer; }
.tab.active { background: #1D4ED8; color: #fff; }
.loading { text-align: center; padding: 40px; color: #999; }
.post-card { display: flex; gap: 12px; margin: 8px 12px; padding: 14px; background: #fff; border-radius: 10px; cursor: pointer; align-items: flex-start; }
.post-badge { width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 12px; font-weight: 700; flex-shrink: 0; }
.wanted { background: #3B82F6; }
.share { background: #22a6b3; }
.post-body { flex: 1; min-width: 0; }
.post-title { font-size: 14px; font-weight: 600; color: #333; }
.post-meta { font-size: 11px; color: #999; margin-top: 6px; }
.empty { text-align: center; padding: 40px; color: #999; }

.skeleton { padding: 14px; }
.sk-item { background: #fff; border-radius: 14px; padding: 14px; margin-bottom: 10px; }
.sk-line { height: 14px; border-radius: 7px; background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; margin-bottom: 8px; }
.w80 { width: 80%; } .w60 { width: 60%; }
@keyframes shimmer { 0% { background-position: -200% 0; } 100% { background-position: 200% 0; } }
</style>
