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
            <span class="edit-btn" @click.stop="$router.push('/publish?edit='+item.id)">编辑</span>
            <span class="card-status" :class="item.status === 1 ? 'on' : 'off'">{{ item.status === 1 ? '在售' : '已下架' }}</span>
          </div>
        </div>
      </div>
      <div class="empty" v-if="products.length === 0">还没发布过商品</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../store'
import { getProducts } from '../../api'

const store = useUserStore()
const products = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  if (!store.token || !store.userInfo) { loading.value = false; return }
  const res: any = await getProducts({ seller_id: store.userInfo.id, page_size: 50 })
  if (res.code === 0) products.value = res.data.items
  loading.value = false
})
</script>

<style scoped>
.page { background: #f5f5f5; min-height: 100vh; }
.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }
.loading, .empty { text-align: center; padding: 40px; color: #999; font-size: 13px; }

.product-card { display: flex; gap: 12px; background: #fff; border-radius: 12px; padding: 12px; margin: 10px 12px; cursor: pointer; }
.card-pic { width: 80px; height: 80px; border-radius: 10px; flex-shrink: 0; }
.pic-0 { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.pic-1 { background: linear-gradient(135deg, #4834d4, #686de0); }
.pic-2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.pic-3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }
.card-right { flex: 1; display: flex; flex-direction: column; justify-content: space-between; min-width: 0; }
.card-title { font-size: 14px; font-weight: 600; color: #222; }
.card-tags { display: flex; gap: 6px; margin: 4px 0; }
.tag { font-size: 10px; padding: 2px 8px; border-radius: 3px; }
.tag-cat { background: #fff0f0; color: #ff4d4f; }
.tag-sub { background: #f0f0f8; color: #666; }
.card-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 4px; }
.price-now { font-size: 17px; font-weight: 700; color: #ff4d4f; }
.card-status { font-size: 11px; padding: 2px 8px; border-radius: 4px; }
.on { background: #e8f8e8; color: #52c41a; }
.off { background: #f0f0f0; color: #999; }
.edit-btn { font-size: 11px; padding: 3px 10px; background: #f0f0f0; color: #666; border-radius: 10px; cursor: pointer; }
</style>
