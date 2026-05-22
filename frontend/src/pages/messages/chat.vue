<template>
  <div class="chat-page">
    <!-- 顶部导航 -->
    <div class="top-bar">
      <span class="back" @click="$router.back()">‹</span>
      <div class="peer-info">
        <div class="peer-avatar">{{ peerNick?.charAt(0) || '?' }}</div>
        <span class="peer-name">{{ peerNick || '聊天' }}</span>
      </div>
      <span class="more">···</span>
    </div>

    <!-- 消息列表 -->
    <div class="msg-area" ref="msgArea">
      <div class="msg-row" v-for="m in messages" :key="m.id" :class="m.from_user_id === myId ? 'me' : 'peer'">
        <div class="bubble" :class="m.from_user_id === myId ? 'b-me' : 'b-peer'">{{ m.content }}</div>
      </div>
      <div class="tip" v-if="messages.length === 0">打个招呼吧~</div>
    </div>

    <!-- 底部输入栏 -->
    <div class="bottom-bar">
      <div class="input-wrap">
        <input v-model="text" placeholder="输入消息..." @keyup.enter="send" />
      </div>
      <span class="send-btn" @click="send">发送</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../../store'
import api from '../../api'

const route = useRoute()
const store = useUserStore()
const peerId = Number(route.params.id)
const myId = store.userInfo?.id || 0
const peerNick = ref('')
const messages = ref<any[]>([])
const text = ref('')
const msgArea = ref<HTMLElement | null>(null)

const scroll = () => {
  setTimeout(() => {
    if (msgArea.value) msgArea.value.scrollTop = msgArea.value.scrollHeight
  }, 100)
}

const send = async () => {
  const txt = text.value.trim()
  if (!txt) return
  messages.value.push({ id: Date.now(), from_user_id: myId, content: txt })
  text.value = ''
  scroll()
  try {
    await api.post('/messages', { to_user_id: peerId, content: txt })
  } catch {}
}

onMounted(async () => {
  try {
    api.post('/messages/read/' + peerId)  // mark as read
    const [msgRes, userRes] = await Promise.all([
      api.get('/messages/' + peerId),
      api.get('/users/' + peerId)
    ])
    if (msgRes.code === 0) messages.value = msgRes.data || []
    if (userRes.code === 0) peerNick.value = userRes.data.nickname
    scroll()
  } catch {}
})
</script>

<style scoped>
.chat-page { position: fixed; inset: 0; z-index: 1001; background: #EDF2FF; display: flex; flex-direction: column; max-width: 450px; margin: 0 auto; }

.top-bar { display: flex; align-items: center; gap: 10px; padding: 10px 12px; background: #fff; box-shadow: 0 1px 0 #F1F5F9; }
.back { font-size: 26px; color: #1D4ED8; cursor: pointer; line-height: 1; }
.peer-info { display: flex; align-items: center; gap: 8px; flex: 1; }
.peer-avatar { width: 36px; height: 36px; border-radius: 50%; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 15px; font-weight: 600; }
.peer-name { font-size: 15px; font-weight: 600; }
.more { font-size: 18px; color: #94a3b8; cursor: pointer; letter-spacing: 1px; }

.msg-area { flex: 1; overflow-y: auto; padding: 12px 14px; display: flex; flex-direction: column; gap: 10px; }
.msg-row { display: flex; }
.msg-row.me { justify-content: flex-end; }
.msg-row.peer { justify-content: flex-start; }

.bubble { max-width: 76%; padding: 10px 14px; font-size: 14px; line-height: 1.5; position: relative; }
.b-me { background: #1D4ED8; color: #fff; border-radius: 16px 4px 16px 16px; margin-left: 40px; }
.b-peer { background: #fff; color: #1E293B; border-radius: 4px 16px 16px 16px; margin-right: 40px; box-shadow: 0 1px 2px rgba(0,0,0,0.04); }

.tip { text-align: center; color: #94a3b8; font-size: 13px; margin-top: 40px; }

.bottom-bar { display: flex; align-items: center; gap: 8px; padding: 8px 12px 12px; background: #fff; box-shadow: 0 -1px 0 #F1F5F9; }
.icon-btn { width: 36px; height: 36px; border-radius: 50%; background: #F1F5F9; display: flex; align-items: center; justify-content: center; font-size: 18px; cursor: pointer; color: #64748b; flex-shrink: 0; }
.input-wrap { flex: 1; background: #F1F5F9; border-radius: 20px; padding: 0 14px; }
.input-wrap input { width: 100%; border: none; background: transparent; padding: 10px 0; font-size: 14px; outline: none; }
.send-btn { padding: 8px 16px; background: #1D4ED8; color: #fff; border-radius: 18px; font-size: 13px; font-weight: 600; cursor: pointer; flex-shrink: 0; }
</style>
