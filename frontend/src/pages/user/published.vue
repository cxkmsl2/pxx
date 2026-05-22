<template>
  <div class="page">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="nav-title">我发布的</span>
      <span style="width:40px"></span>
    </div>
    <div class="empty" v-if="!store.token">请先登录</div>
    <div class="loading" v-else-if="loading">加载中...</div>
    <div v-else>
      <div class="product-card" v-for="item in products" :key="item.id" @click="$router.push('/product/'+item.id)">
        <div class="card-pic" :class="'pic-' + (item.id % 4)"></div>
        <div class="card-right">
          <div class="card-title">{{ item.title }}</div>
          <div class="card-tags">
            <span class="tag tag-cat">{{ item.category }}</span>
            <span class="tag tag-sub">{{ item.tag }}</span>
          </div>
          <div class="card-footer">
            <span class="price-now">¥{{ (item.price / 100).toFixed(2) }}</span>
            <span class="card-status" :class="item.status === 1 ? 'on' : 'off'">{{ item.status === 1 ? '在售' : '已下架' }}</span>
            <span class="edit-btn" @click.stop="$router.push('/publish?edit='+item.id)">编辑</span>
            <span class="del-btn" @click.stop="doDelete(item.id)">删除</span>
          </div>
        </div>
      </div>
      <div class="empty" v-if="products.length === 0 && !loading">还没发布过商品</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../store'
import api, { getProducts } from '../../api'

const store = useUserStore()
const products = ref<any[]>([])
const loading = ref(true)

const loadData = async () => {
  loading.value = true
  if (!store.token || !store.userInfo) { loading.value = false; return }
  const res: any = await getProducts({ seller_id: store.userInfo.id, page_size: 50 })
  if (res.code === 0) products.value = res.data.items
  loading.value = false
}

const doDelete = async (id: number) => {
  if (!confirm('确定删除？')) return
  try {
    const res = await api.delete('/products/' + id)
    if (res.code === 0) { window.$toast('已删除', 'success'); loadData() }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('删除失败', 'error') }
}

onMounted(() => loadData())
</script>

<style scoped>
.page { background: #F8FAFC; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #94a3b8; font-size: 13px; }

.product-card { display: flex; gap: 12px; background: #fff; border-radius: 14px; padding: 12px; margin: 10px 14px; cursor: pointer; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.card-pic { width: 80px; height: 80px; border-radius: 10px; flex-shrink: 0; }
.pic-0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.pic-1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.pic-2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.pic-3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.card-right { flex: 1; display: flex; flex-direction: column; justify-content: space-between; min-width: 0; }
.card-title { font-size: 14px; font-weight: 600; color: #1E293B; }
.card-tags { display: flex; gap: 6px; margin: 4px 0; }
.tag { font-size: 10px; padding: 2px 8px; border-radius: 3px; }
.tag-cat { background: #EFF6FF; color: #1D4ED8; }
.tag-sub { background: #F1F5F9; color: #64748b; }
.card-footer { display: flex; align-items: center; gap: 6px; margin-top: 4px; }
.price-now { font-size: 17px; font-weight: 700; color: #10B981; }
.card-status { font-size: 11px; padding: 2px 8px; border-radius: 4px; }
.on { background: #D1FAE5; color: #059669; }
.off { background: #F1F5F9; color: #94a3b8; }
.edit-btn { font-size: 11px; padding: 3px 10px; background: #F1F5F9; color: #64748b; border-radius: 10px; cursor: pointer; margin-left: auto; }
.del-btn { font-size: 11px; padding: 3px 10px; background: #FEF2F2; color: #EF4444; border-radius: 10px; cursor: pointer; }
</style>