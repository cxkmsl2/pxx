<template>
  <div class="publish">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="title">{{ editId ? '编辑商品' : '发布商品' }}</span>
      <span class="submit" @click="doSubmit" :class="{ dim: submitting }">
        {{ submitting ? '发布中' : (editId ? '保存' : '发布') }}
      </span>
    </div>
    <div class="form">
      <input v-model="form.title" placeholder="商品标题" class="input" />
      <textarea v-model="form.desc" placeholder="商品描述" class="textarea" rows="3"></textarea>
      <div class="row">
        <select v-model="form.category" class="select">
          <option value="">选择分类</option>
          <option>考研资料</option><option>数码外设</option><option>日用百货</option>
          <option>美妆护肤</option><option>服饰鞋包</option><option>其他</option>
        </select>
        <input v-model="form.tag" placeholder="标签" class="input half" />
      </div>
      <div class="label-text">💰 售价（元）</div>
      <input v-model="priceYuan" type="number" placeholder="例如：45.00" class="input" />
      <div class="label-text">🏷 原价（元）<span class="hint">可不填</span></div>
      <input v-model="originalYuan" type="number" placeholder="例如：99.00" class="input" />
      <input v-model="form.campus" placeholder="所在校区" class="input" />
      <div class="btn" @click="doSubmit" :class="{ dim: submitting }">
        {{ submitting ? '提交中...' : (editId ? '保存修改' : '确认发布') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../../store'
import api, { getProduct } from '../../api'

const router = useRouter()
const route = useRoute()
const store = useUserStore()
const editId = Number(route.query.edit) || 0
const submitting = ref(false)
const priceYuan = ref('')
const originalYuan = ref('')

const form = ref({
  title: '', desc: '', category: '', tag: '',
  price: 0, originalPrice: 0, campus: ''
})

onMounted(async () => {
  if (editId) {
    const res: any = await getProduct(editId)
    if (res.code === 0) {
      const p = res.data
      form.value = {
        title: p.title, desc: p.desc || '', category: p.category, tag: p.tag || '',
        price: p.price, originalPrice: p.original_price || 0, campus: p.campus || ''
      }
      priceYuan.value = String((p.price || 0) / 100)
      originalYuan.value = String((p.original_price || 0) / 100)
    }
  }
})

const doSubmit = async () => {
  if (!store.token) { window.$toast('请先登录'); router.push('/login'); return }
  if (!form.value.title) { window.$toast('请输入标题'); return }
  submitting.value = true
  try {
    const p = Number(priceYuan.value) || 0
    const op = Number(originalYuan.value) || 0
    const data = {
      ...form.value,
      price: Math.round(p * 100),
      original_price: Math.round(op * 100),
    }
    const res = editId
      ? await api.put('/products/' + editId, data)
      : await api.post('/products', data)
    if (res.code === 0) {
      window.$toast(editId ? '已更新' : '发布成功', 'success')
      router.push('/')
    } else window.$toast(res.msg, 'error')
  } catch { window.$toast('操作失败', 'error') }
  submitting.value = false
}
</script>

<style scoped>
.publish { background: #f5f5f5; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { color: #333; font-size: 14px; cursor: pointer; }
.title { font-size: 15px; font-weight: 600; }
.submit { color: #ff4d4f; font-size: 14px; font-weight: 600; cursor: pointer; }
.submit.dim { opacity: 0.5; }
.form { padding: 16px; display: flex; flex-direction: column; gap: 8px; }
.label-text { font-size: 13px; color: #666; font-weight: 500; margin: 4px 0 0; }
.hint { font-size: 11px; color: #bbb; font-weight: 400; margin-left: 6px; }
.input, .textarea, .select {
  width: 100%; padding: 12px; border: 1px solid #e0e0e0; border-radius: 10px;
  font-size: 14px; outline: none; box-sizing: border-box; background: #fff;
}
.textarea { resize: vertical; min-height: 80px; }
.row { display: flex; gap: 10px; }
.half { flex: 1; }
.select { color: #333; }
.btn {
  background: linear-gradient(135deg, #ff6b6b, #ee5a24); color: #fff;
  text-align: center; padding: 14px; border-radius: 10px; font-size: 16px; font-weight: 600;
  cursor: pointer; margin-top: 12px;
}
.btn.dim { opacity: 0.6; }
</style>
