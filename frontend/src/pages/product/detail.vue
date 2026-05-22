<template>
  <div class="detail" v-if="product">
    <div class="nav"><span class="back" @click="$router.back()">←</span><span class="nav-title">商品详情</span></div>
    <div class="pic" :style="{ backgroundImage: product.images ? 'url(' + (JSON.parse(product.images)[0] || '') + ')' : '' }" :class="product.images ? '' : 'pic-' + (product.id % 4)"></div>
    <div class="price-row">
      <span class="price">¥{{ (product.price / 100).toFixed(2) }}</span>
      <span class="original" v-if="product.original_price > 0">¥{{ (product.original_price / 100).toFixed(2) }}</span>
      <span class="discount" v-if="product.original_price > 0">{{ ((1 - product.price / product.original_price) * 100).toFixed(0) }}% OFF</span>
    </div>
    <div class="section"><div class="title">{{ product.title }}</div><p class="desc">{{ product.desc }}</p></div>
    <div class="section"><div class="info-row"><span>👤 {{ product.seller?.nickname }} <b class="credit-badge">{{ creditText }}</b></span><span>📍 {{ product.campus }}</span><span>👁 {{ product.view_count }} 浏览</span></div></div>
    <div class="footer">
      <span class="btn btn-plain" @click="toggleFav">{{ faved ? '❤️ 已收藏' : '🤍 收藏' }}</span>
      <span class="btn btn-plain" @click="goChat">💬 私信</span>
      <span class="btn btn-plain" @click="doBargain">🔪 砍价</span>
      <span class="btn btn-primary" @click="goPay">{{ buying ? '支付中' : '立即购买' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { getProduct } from '../../api'
import { useUserStore } from '../../store'

const route = useRoute(); const router = useRouter(); const store = useUserStore()
const product = ref<any>(null); const buying = ref(false); const faved = ref(false)

const sellerScore = computed(() => product.value?.seller?.credit || 60)
const creditText = computed(() => {
  const s = sellerScore.value
  if (s >= 90) return '信用极好'
  if (s >= 80) return '信用优秀'
  if (s >= 70) return '信用良好'
  if (s >= 60) return '信用一般'
  return '信用差'
})

onMounted(async () => {
  const id = Number(route.params.id)
  const res = await getProduct(id)
  if (res.code === 0) {
    product.value = res.data
    const favs = JSON.parse(localStorage.getItem('pxx_favs') || '[]')
    faved.value = favs.includes(id)
    const hist = JSON.parse(localStorage.getItem('pxx_history') || '[]')
    if (!hist.includes(id)) { hist.push(id); localStorage.setItem('pxx_history', JSON.stringify(hist.slice(-20))) }
  }
})

const toggleFav = () => {
  const favs = JSON.parse(localStorage.getItem('pxx_favs') || '[]')
  const pid = product.value.id
  if (favs.includes(pid)) { favs.splice(favs.indexOf(pid), 1); window.$toast('已取消') }
  else { favs.push(pid); window.$toast('已收藏', 'success') }
  localStorage.setItem('pxx_favs', JSON.stringify(favs))
  faved.value = !faved.value
}

const goChat = () => {
  if (product.value?.seller?.id) router.push('/chat/' + product.value.seller.id)
  else window.$toast('无法获取卖家信息')
}

const doBargain = () => {
  if (!store.token) { window.$toast('请先登录'); router.push('/login'); return }
  window.$sheet('输入报价', ['¥10', '¥20', '¥30', '¥50', '自定义'], async (i: number) => {
    const prices = ['10','20','30','50','custom']
    let price = prices[i]
    if (price === 'custom') price = prompt('输入报价（元）') || '0'
    if (!price || Number(price) <= 0) return
    await api.post('/messages', { to_user_id: product.value.seller.id, content: '🔪 砍价：' + store.userInfo.nickname + ' 出价 ¥' + price + '元' })
    window.$toast('砍价已发送', 'success')
  })
}

const goPay = () => {
  if (!store.token) { router.push('/login'); return }
  window.$sheet('支付方式', ['💚 微信支付', '💙 支付宝', '💰 余额'], async (i: number) => {
    buying.value = true
    try {
      const res = await api.post('/orders', { product_id: product.value.id })
      if (res.code === 0) { window.$toast('支付成功', 'success'); router.push('/orders') }
      else window.$toast(res.msg, 'error')
    } catch { window.$toast('支付失败', 'error') }
    buying.value = false
  })
}
</script>

<style scoped>
.detail { background: var(--bg, #F2F2F6); min-height: 100vh; padding-bottom: 130px; }
.nav { display: flex; gap: 10px; padding: 12px; background: #fff; }
.back { font-size: 16px; color: #1D4ED8; cursor: pointer; }
.nav-title { font-weight: 600; }
.pic { height: 260px; background-size: cover; background-position: center; }
.pic-0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.pic-1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.pic-2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.pic-3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.price-row { display: flex; align-items: center; gap: 10px; padding: 14px; background: #fff; }
.price { font-size: 24px; font-weight: 800; color: #1D4ED8; }
.original { font-size: 13px; color: #94a3b8; text-decoration: line-through; }
.discount { font-size: 12px; background: #D1FAE5; color: #059669; padding: 2px 8px; border-radius: 4px; margin-left: auto; }
.section { background: #fff; margin: 8px 14px; padding: 14px; border-radius: 14px; }
.title { font-size: 16px; font-weight: 700; margin-bottom: 6px; }
.desc { font-size: 13px; color: #64748b; line-height: 1.6; }
.info-row { display: flex; gap: 12px; font-size: 13px; color: #64748b; align-items: center; }
.credit-badge { font-size: 10px; background: #D1FAE5; color: #059669; padding: 1px 6px; border-radius: 8px; margin-left: 4px; }
.footer { position: fixed; bottom: 56px; left: 50%; transform: translateX(-50%); width: 100%; max-width: 450px; display: flex; gap: 6px; padding: 8px 12px; background: #fff; z-index: 1000; box-shadow: 0 -1px 4px rgba(0,0,0,0.06); }
.btn { flex: 1; text-align: center; padding: 10px 6px; border-radius: 20px; font-size: 12px; cursor: pointer; }
.btn-plain { background: #F1F5F9; color: #64748b; }
.btn-primary { flex: 1.5; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; font-weight: 600; }
</style>
