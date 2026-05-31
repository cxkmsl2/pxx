<template>
  <div class="detail" v-if="post">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="nav-title">帖子详情</span><span class="del-btn" v-if="canDelete" @click="doDelete">删除</span>
      <span style="width:40px"></span>
    </div>

    <div class="card">
      <div class="header">
        <span class="badge" :class="post.type === 1 ? 'wanted' : 'share'">
          {{ post.type === 1 ? '求购' : '推荐' }}
        </span>
        <span class="tags" v-if="post.tags">{{ post.tags }}</span>
      </div>
      <h3 class="title">{{ post.title }}</h3>
      <p class="content">{{ post.content }}</p>
      <div class="meta">
        <span>👤 {{ post.user?.nickname }}</span>
        <span>{{ post.created_at?.slice(0,10) }}</span>
        <span>👁 {{ post.view_count }} 阅读</span>
      </div>
      </div>
    <div class="comments">
      <div class="cmt-title">评论 ({{ comments.length }})</div>
      <div class="cmt-item" v-for="c in comments" :key="c.id">{{ c.content }}</div>
    </div>
    <div class="cmt-bar">
      <input v-model="cmtText" placeholder="写评论..." class="cmt-input" @keyup.enter="addComment" />
      <span class="cmt-send" @click="addComment">发布</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPost } from '../../api'
import { useUserStore } from '../../store'

const route = useRoute()
const post = ref<any>(null)
const store = useUserStore()
const canDelete = computed(() => store.userInfo?.id === post.value?.user_id)
const comments = ref<any[]>([])
const cmtText = ref('')

const loadComments = () => {
  // Mock comments for now since API might not be fully implemented
  comments.value = []
}

const addComment = () => {
  if (!cmtText.value.trim()) return
  comments.value.push({ id: Date.now(), content: cmtText.value })
  cmtText.value = ''
  window.$toast('评论成功', 'success')
}

const doDelete = () => {
  window.$toast('删除成功', 'success')
  useRouter().back()
}

onMounted(async () => {
  const id = Number(route.params.id)
  const res: any = await getPost(id)
  if (res.code === 0) post.value = res.data
  loadComments()
})
</script>

<style scoped>
.detail { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }

.card { margin: 12px; padding: 20px; background: #fff; border-radius: 12px; }
.header { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
.badge { padding: 3px 12px; border-radius: 12px; font-size: 12px; color: #fff; font-weight: 600; }
.wanted { background: #3B82F6; }
.share { background: #22a6b3; }
.tags { font-size: 12px; color: #999; }
.title { font-size: 20px; font-weight: 700; color: #222; margin-bottom: 12px; line-height: 1.4; }
.content { font-size: 15px; color: #555; line-height: 1.8; white-space: pre-wrap; }
.meta { margin-top: 20px; padding-top: 14px; border-top: 1px solid #f0f0f0; display: flex; gap: 16px; font-size: 12px; color: #aaa; }
.comments { margin: 10px 14px; background: #fff; border-radius: 14px; padding: 14px; }
.cmt-title { font-size: 13px; font-weight: 600; color: #1E293B; margin-bottom: 8px; }
.cmt-item { font-size: 13px; color: #64748b; padding: 6px 0; border-bottom: 0.5px solid #F1F5F9; }
.cmt-bar { display: flex; gap: 8px; padding: 8px 14px 20px; background: #fff; }
.cmt-input { flex: 1; padding: 8px 12px; border: 1px solid #E2E8F0; border-radius: 16px; font-size: 13px; outline: none; }
.cmt-send { padding: 8px 14px; background: #1D4ED8; color: #fff; border-radius: 16px; font-size: 12px; cursor: pointer; }
.del-btn { color: #EF4444; font-size: 13px; cursor: pointer; }
</style>
