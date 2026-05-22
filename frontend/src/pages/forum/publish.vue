<template>
  <div class="publish">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="title">发布帖子</span>
      <span class="submit" @click="doSubmit" :class="{ dim: submitting }">
        {{ submitting ? '发布中' : '发布' }}
      </span>
    </div>
    <div class="form">
      <div class="tabs">
        <span class="tab" :class="{ active: form.type === 1 }" @click="form.type = 1">求购悬赏</span>
        <span class="tab" :class="{ active: form.type === 2 }" @click="form.type = 2">好物种草</span>
      </div>
      <input v-model="form.title" placeholder="标题" class="input" />
      <textarea v-model="form.content" placeholder="详细描述..." class="textarea" rows="4"></textarea>
      <input v-model="form.tags" placeholder="标签（逗号分隔）" class="input" />
      <div class="btn" @click="doSubmit" :class="{ dim: submitting }">
        {{ submitting ? '发布中...' : '确认发布' }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store'
import api from '../../api'

const router = useRouter()
const store = useUserStore()
const submitting = ref(false)
const form = ref({ type: 1, title: '', content: '', tags: '' })

const doSubmit = async () => {
  if (!store.token) { window.$toast('请先登录'); router.push('/login'); return }
  if (!form.value.title) { window.$toast('请输入标题'); return }
  submitting.value = true
  try {
    const res = await api.post('/posts', form.value)
    if (res.code === 0) {
      window.$toast('发布成功！')
      router.push('/forum')
    } else { window.$toast(res.msg) }
  } catch { window.$toast('发布失败') }
  submitting.value = false
}
</script>

<style scoped>
.publish { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { color: #333; font-size: 14px; cursor: pointer; }
.title { font-size: 15px; font-weight: 600; }
.submit { color: #1D4ED8; font-size: 14px; font-weight: 600; cursor: pointer; }
.submit.dim { opacity: 0.5; }

.form { padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.tabs { display: flex; gap: 10px; margin-bottom: 4px; }
.tab { padding: 8px 20px; border-radius: 16px; font-size: 13px; background: #f0f0f0; color: #666; cursor: pointer; }
.tab.active { background: #1D4ED8; color: #fff; }

.input, .textarea {
  width: 100%; padding: 12px; border: 1px solid #e0e0e0; border-radius: 10px;
  font-size: 14px; outline: none; box-sizing: border-box; background: #fff;
}
.textarea { resize: vertical; min-height: 100px; }

.btn {
  background: linear-gradient(135deg, #3B82F6, #2563EB); color: #fff;
  text-align: center; padding: 12px; border-radius: 10px; font-size: 15px; font-weight: 600; cursor: pointer; margin-top: 10px;
}
.btn.dim { opacity: 0.6; }
</style>
