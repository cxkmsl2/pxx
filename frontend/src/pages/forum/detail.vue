<template>
  <div class="detail" v-if="post">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="nav-title">帖子详情</span>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '../../api'

const route = useRoute()
const post = ref<any>(null)

onMounted(async () => {
  const id = Number(route.params.id)
  const res: any = await getPost(id)
  if (res.code === 0) post.value = res.data
})
</script>

<style scoped>
.detail { background: #f5f5f5; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }

.card { margin: 12px; padding: 20px; background: #fff; border-radius: 12px; }
.header { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
.badge { padding: 3px 12px; border-radius: 12px; font-size: 12px; color: #fff; font-weight: 600; }
.wanted { background: #ff6b6b; }
.share { background: #22a6b3; }
.tags { font-size: 12px; color: #999; }
.title { font-size: 20px; font-weight: 700; color: #222; margin-bottom: 12px; line-height: 1.4; }
.content { font-size: 15px; color: #555; line-height: 1.8; white-space: pre-wrap; }
.meta { margin-top: 20px; padding-top: 14px; border-top: 1px solid #f0f0f0; display: flex; gap: 16px; font-size: 12px; color: #aaa; }
</style>
