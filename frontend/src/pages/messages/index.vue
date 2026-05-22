<template>
  <div class="msgs">
    <div class="nav">💬 消息</div>
    <div class="empty" v-if="!store.token">请先登录</div>
    <div class="sk-list" v-if="loading">
      <div class="sk-conv" v-for="i in 3" :key="i"><div class="sk-avatar"></div><div class="sk-body"><div class="sk-line w80"></div><div class="sk-line w60"></div></div></div>
    </div>
    <div class="empty" v-else-if="convs.length === 0 && !loading">暂无消息</div>
    <div v-else>
      <div class="conv-item" v-for="c in convs" :key="c.user_id" @click="$router.push('/chat/' + c.user_id)">
        <div class="c-avatar">{{ c.nickname?.charAt(0) || '?' }}</div>
        <div class="c-body">
          <div class="c-top"><span class="c-name">{{ c.nickname }}</span><span class="c-time">{{ fmt(c.last_time) }}</span></div>
          <div class="c-bottom"><span class="c-msg">{{ c.last_msg || '开始聊天' }}</span><span class="c-badge" v-if="c.unread">{{ c.unread }}</span></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../store'
import api from '../../api'

const store = useUserStore()
const convs = ref<any[]>([])
const loading = ref(true)

const fmt = (t: string) => t ? new Date(t).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) : ''

onMounted(async () => {
  if (!store.token) { loading.value = false; return }
  try {
    const res: any = await api.get('/messages/conversations')
    if (res.code === 0) convs.value = res.data
  } catch {}
  loading.value = false
})
</script>

<style scoped>
.msgs { background: #F8FAFC; min-height: 100vh; }
.nav { padding: 14px; background: #fff; font-size: 16px; font-weight: 700; color: #1E293B; }
.empty { text-align: center; padding: 60px; color: #94a3b8; }

.conv-item { display: flex; align-items: center; gap: 12px; padding: 14px; background: #fff; cursor: pointer; border-bottom: 0.5px solid #F1F5F9; }
.c-avatar { width: 48px; height: 48px; border-radius: 50%; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: 600; flex-shrink: 0; }
.c-body { flex: 1; min-width: 0; }
.c-top { display: flex; justify-content: space-between; }
.c-name { font-size: 15px; font-weight: 600; color: #1E293B; }
.c-time { font-size: 11px; color: #94a3b8; }
.c-bottom { display: flex; justify-content: space-between; align-items: center; margin-top: 4px; }
.c-msg { font-size: 13px; color: #94a3b8; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 80%; }
.c-badge { background: #EF4444; color: #fff; font-size: 10px; padding: 2px 6px; border-radius: 10px; min-width: 18px; text-align: center; }

.sk-list { }
.sk-conv { display: flex; gap: 12px; padding: 14px; background: #fff; align-items: center; }
.sk-avatar { width: 48px; height: 48px; border-radius: 50%; background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; flex-shrink: 0; }
.sk-body { flex: 1; }
.sk-line { height: 14px; border-radius: 7px; background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; margin-bottom: 6px; }
.w80 { width: 80%; } .w60 { width: 60%; }
@keyframes shimmer { 0% { background-position: -200% 0; } 100% { background-position: 200% 0; } }
</style>