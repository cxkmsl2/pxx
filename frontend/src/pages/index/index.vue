<template>
  <div class="home">
    <div class="top-bar">
      <div class="search-box" @click="$router.push('/search')">
        <span>🔍</span>
        <span class="placeholder-text">搜搜闲置好物</span>
      </div>
      <div class="publish-btn" @click="$router.push('/publish')">＋ 卖闲置</div>
    </div>

    <div class="auction-banner" v-if="auctionItem" @click="$router.push('/product/' + auctionItem.id)">
      <div class="auction-label">⚡ 荷兰拍 · 极速清仓</div>
      <div class="auction-timer">每 10 秒自动降价 5%</div>
      <div class="auction-progress"><div class="auction-bar" :style="{ width: pct + '%' }"></div></div>
      <div class="auction-price">¥{{ (auctionItem.price / 100).toFixed(2) }}<em v-if="auctionItem.original_price > 0">原价 ¥{{ (auctionItem.original_price / 100).toFixed(2) }}</em></div>
    </div>

    <div class="cats">
      <div class="cat-item" v-for="c in cats" :key="c.text" @click="c.to ? $router.push(c.to) : $router.push('/search?q=' + encodeURIComponent(c.kw || c.text))">
        <span class="cat-icon">{{ c.icon }}</span><span class="cat-text">{{ c.text }}</span>
      </div>
    </div>

    <div class="section-head" @click="doRefresh">热门推荐 <span class="rf">↻</span></div>
    <div class="waterfall" v-if="!loading">
      <div class="wf-card" v-for="item in products" :key="item.id" @click="$router.push('/product/'+item.id)">
        <div class="wf-pic" :class="(item.images && item.images !== '[]' && item.images !== '') ? '' : 'bg' + (item.id % 4)" :style="{ height: (120 + (item.id % 3) * 40) + 'px', backgroundImage: getImg(item) }"></div>
        <div class="wf-info"><div class="wf-title">{{ item.title }}</div><div class="wf-price">¥{{ (item.price / 100).toFixed(2) }}</div><div class="wf-meta">{{ item.campus }}</div></div>
      </div>
    </div>
    <div class="loading" v-else>加载中...</div>
    <div class="empty-state" v-if="!loading && products.length === 0">📭 暂无商品<br/><span>成为第一个发布的人吧</span></div>
    <div class="load-more" v-if="!finished && !loading" @click="loadMore">加载更多</div>
    <div class="end" v-if="finished">— 已经到底了 —</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProducts } from '../../api'

const router = useRouter()
const products = ref<any[]>([])
const auctionItem = ref<any>(null)
const loading = ref(true)
const finished = ref(false)
let page = 1

const cats = [
  { icon: '🔄', text: '置换', to: '/barter' },
  { icon: '🏃', text: '跑腿', to: '/tasks' },
  { icon: '📦', text: '拼团', to: '/groupbuy' },
  { icon: '⏰', text: '租赁', to: '/rental' },
  { icon: '🎫', text: '拼卡', to: '/subscription' },
  { icon: '💬', text: '论坛', to: '/forum' },
  { icon: '⚖️', text: '法庭', to: '/tribunal' },
]

const getImg = (item: any) => { 
  if (item.images && typeof item.images === 'string' && item.images !== '') { 
    try { 
      const imgs = JSON.parse(item.images); 
      if (Array.isArray(imgs) && imgs.length > 0) {
        return 'url(' + imgs[0] + ')';
      }
    } catch (e) {
      console.error("Image parse error:", e);
    }
  } 
  return ''; 
}

const pct = ref(50)

const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await getProducts({ page, page_size: 10 })
    if (res.code === 0) {
      products.value.push(...res.data.items)
      if (res.data.items.length < 10) finished.value = true
      page++
      if (!auctionItem.value) {
        const d = res.data.items.filter((p: any) => p.original_price > p.price)
        auctionItem.value = d[0] || res.data.items[0]
        if (auctionItem.value && auctionItem.value.original_price > 0)
          pct.value = Math.round((1 - auctionItem.value.price / auctionItem.value.original_price) * 100)
      }
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

const loadMore = () => fetchData()
const doRefresh = () => { page = 1; products.value = []; finished.value = false; if (!auctionItem.value) auctionItem.value = null; fetchData() }
onMounted(() => fetchData())
</script>

<style scoped>
.home { background: #F8FAFC; }
.top-bar { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: rgba(255,255,255,0.9); backdrop-filter: blur(20px); position: sticky; top: 0; z-index: 10; }
.search-box { display: flex; align-items: center; gap: 6px; flex: 1; background: #F1F5F9; border-radius: 20px; padding: 9px 14px; cursor: pointer; }
.placeholder-text { color: #94a3b8; font-size: 13px; }
.publish-btn { padding: 9px 14px; background: linear-gradient(135deg, #1D4ED8, #3B82F6); color: #fff; border-radius: 20px; font-size: 13px; font-weight: 600; white-space: nowrap; cursor: pointer; }

.auction-banner { margin: 10px 12px; padding: 16px; background: linear-gradient(135deg, #1E293B, #334155); border-radius: 16px; color: #fff; cursor: pointer; }
.auction-label { font-size: 13px; color: #FBBF24; font-weight: 600; }
.auction-timer { font-size: 11px; color: #94a3b8; margin-top: 4px; }
.auction-progress { height: 4px; background: rgba(255,255,255,0.15); border-radius: 2px; margin: 10px 0; }
.auction-bar { height: 100%; background: linear-gradient(90deg, #10B981, #34D399); border-radius: 2px; }
.auction-price { font-size: 28px; font-weight: 800; }
.auction-price em { font-size: 12px; color: #94a3b8; font-style: normal; text-decoration: line-through; margin-left: 8px; font-weight: 400; }

.cats { display: grid; grid-template-columns: repeat(3, 1fr); background: #fff; padding: 10px 12px; }
.cat-item { display: flex; flex-direction: column; align-items: center; padding: 8px 0; gap: 4px; cursor: pointer; }
.cat-icon { font-size: 24px; } .cat-text { font-size: 11px; color: #64748b; }

.section-head { padding: 16px 14px 8px; font-size: 16px; font-weight: 700; color: #1E293B; cursor: pointer; }
.rf { color: #1D4ED8; margin-left: 4px; }
.loading, .end { text-align: center; padding: 24px; color: #94a3b8; font-size: 12px; }

.waterfall { column-count: 2; column-gap: 10px; padding: 0 12px 80px; }
.wf-card { break-inside: avoid; display: inline-block; width: 100%; background: #fff; border-radius: 14px; overflow: hidden; margin-bottom: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); cursor: pointer; }
.wf-pic { width: 100%; border-radius: 14px 14px 0 0; }
.bg0 { background-size: cover; background-position: center; background-repeat: no-repeat; background-color: #DBEAFE; }
.bg1 { background-color: #D1FAE5; }
.bg2 { background-color: #FEF3C7; }
.bg3 { background-color: #EDE9FE; }
.wf-info { padding: 10px; }
.wf-title { font-size: 13px; font-weight: 600; color: #1E293B; line-height: 1.3; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.wf-price { font-size: 17px; font-weight: 700; color: #10B981; margin-top: 4px; }
.wf-meta { font-size: 10px; color: #94a3b8; margin-top: 2px; }
.load-more { text-align: center; padding: 14px; color: #1D4ED8; font-size: 13px; cursor: pointer; }
.skeleton { padding: 0 12px; }
.sk-card { display: flex; gap: 12px; background: #fff; border-radius: 16px; padding: 12px; margin-bottom: 10px; }
.sk-pic { width: 110px; height: 110px; border-radius: 10px; background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
.sk-body { flex: 1; display: flex; flex-direction: column; gap: 8px; justify-content: center; }
.sk-line { height: 14px; border-radius: 7px; background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
.w80 { width: 80%; } .w60 { width: 60%; } .w40 { width: 40%; }
@keyframes shimmer { 0% { background-position: -200% 0; } 100% { background-position: 200% 0; } }

.empty-state { text-align: center; padding: 40px 20px; color: #8E8E93; font-size: 15px; }
.empty-state span { font-size: 12px; color: #C7C7CC; }</style>