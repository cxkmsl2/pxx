<template>
  <div class="detail" v-if="product">
    <div class="nav">
      <span class="back" @click="$router.back()">← 返回</span>
      <span class="nav-title">商品详情</span>
      <span style="width:40px"></span>
    </div>

    <div class="pic" :class="'pic-' + (product.id % 4)"></div>

    <div class="price-row">
      <span class="price">¥{{ (product.price / 100).toFixed(2) }}</span>
      <span class="original" v-if="product.original_price > 0">¥{{ (product.original_price / 100).toFixed(2) }}</span>
      <span class="discount">{{ product.original_price > 0 ? ((1 - product.price / product.original_price) * 100).toFixed(0) + '% OFF' : '优惠' }}</span>
    </div>

    <div class="section">
      <div class="title">{{ product.title }}</div>
      <p class="desc">{{ product.desc }}</p>
    </div>

    <div class="section">
      <div class="info-row">
        <span>👤 {{ product.seller?.nickname }}</span>
        <span>📍 {{ product.campus }}</span>
        <span>👁 {{ product.view_count }} 浏览</span>
      </div>
    </div>

    <div class="footer">
      <span class="btn btn-plain" @click="msg('已收藏')">🤍 收藏</span>
      <span class="btn btn-plain" @click="msg('私信功能')">💬 私信</span>
      <span class="btn btn-primary" @click="doBuy">{{ buying ? '提交中...' : '立即购买' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { getProduct } from '../../api'
import { useUserStore } from '../../store'

const route = useRoute()
const router = useRouter()
const store = useUserStore()
const product = ref<any>(null)
const buying = ref(false)

onMounted(async () => {
  const id = Number(route.params.id)
  const res: any = await getProduct(id)
  if (res.code === 0) product.value = res.data
})

const doBuy = async () => {
  if (!store.token) { router.push('/login'); return }
  buying.value = true
  try {
    const res = await api.post('/orders', { product_id: product.value.id })
    if (res.code === 0) { window.$toast('下单成功', 'success'); router.push('/orders') }
    else window.$toast(res.msg, 'error')
  } catch { window.$toast('下单失败', 'error') }
  buying.value = false
}

const msg = (m: string) => window.$toast(m)
</script>

<style scoped>
.detail { background: #f5f5f5; min-height: 100vh; padding-bottom: 120px; }

.nav { display: flex; align-items: center; justify-content: space-between; padding: 12px; background: #fff; position: sticky; top: 0; z-index: 10; }
.back { font-size: 14px; color: #333; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }

.pic { height: 260px; }
.pic-0 { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.pic-1 { background: linear-gradient(135deg, #4834d4, #686de0); }
.pic-2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.pic-3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }

.price-row { display: flex; align-items: center; gap: 10px; padding: 14px; background: #fff; }
.price { font-size: 24px; font-weight: 800; color: #ff4d4f; }
.original { font-size: 13px; color: #bbb; text-decoration: line-through; }
.discount { font-size: 12px; background: #fff0f0; color: #ff4d4f; padding: 2px 8px; border-radius: 4px; margin-left: auto; }

.section { background: #fff; margin: 10px 12px; padding: 14px; border-radius: 10px; }
.title { font-size: 16px; font-weight: 700; color: #222; margin-bottom: 6px; }
.desc { font-size: 13px; color: #666; line-height: 1.6; }
.info-row { display: flex; gap: 16px; font-size: 13px; color: #666; }

.footer { position: fixed; bottom: 56px; left: 50%; transform: translateX(-50%); width: 100%; max-width: 450px; display: flex; gap: 10px; padding: 8px 12px; background: #fff; z-index: 1000; box-shadow: 0 -1px 4px rgba(0,0,0,0.06); }
.btn { flex: 1; text-align: center; padding: 10px; border-radius: 20px; font-size: 13px; cursor: pointer; }
.btn-plain { background: #f5f5f5; color: #666; }
.btn-primary { flex: 2; background: linear-gradient(135deg, #ff6b6b, #ee5a24); color: #fff; font-weight: 600; }
</style>
