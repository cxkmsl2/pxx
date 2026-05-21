<template>
  <div class="home">
    <div class="search-bar">
      <div class="publish-btn" @click="$router.push('/publish')">＋ 卖闲置</div>
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input v-model="keyword" class="search-input" placeholder="搜搜闲置好物" @keyup.enter="onSearch" />
      </div>
    </div>

    <div class="swiper" @touchstart="onTouchStart" @touchend="onTouchEnd">
      <div class="swiper-track" :style="{ transform: 'translateX(-' + currentBanner * 100 + '%)' }">
        <div class="banner" v-for="(b, i) in banners" :key="i" :class="'banner-' + i" @click="goBanner(i)">{{ b }}</div>
      </div>
      <div class="swiper-dots">
        <span v-for="(b, i) in banners" :key="i" :class="{ active: i === currentBanner }"></span>
      </div>
    </div>

    <div class="cats">
      <div class="cat-item" v-for="c in cats" :key="c.text" @click="$router.push(c.to)">
        <span class="cat-icon">{{ c.icon }}</span>
        <span class="cat-text">{{ c.text }}</span>
      </div>
    </div>

    <div class="list-header" @click="doRefresh">热门推荐<span class="refresh-icon"> ↻</span></div>
    <div class="loading" v-if="loading">加载中...</div>
    <div class="products" v-else>
      <div class="product-card" v-for="item in products" :key="item.id" @click="$router.push('/product/'+item.id)">
        <div class="card-pic" :class="'pic-' + (item.id % 4)"></div>
        <div class="card-right">
          <div class="card-title">{{ item.title }}</div>
          <div class="card-tags">
            <span class="tag tag-cat">{{ item.category }}</span>
            <span class="tag tag-sub">{{ item.tag }}</span>
          </div>
          <div class="card-footer">
            <div class="card-prices">
              <span class="price-now">¥{{ (item.price / 100).toFixed(2) }}</span>
              <span class="price-old" v-if="item.original_price > 0">¥{{ (item.original_price / 100).toFixed(2) }}</span>
            </div>
            <span class="card-campus">{{ item.campus }}</span>
          </div>
        </div>
      </div>
    </div>
    <div class="load-more" v-if="!finished && !loading" @click="loadMore">点击加载更多</div>
    <div class="end" v-if="finished">— 已经到底了 —</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProducts } from '../../api'

const router = useRouter()
const keyword = ref('')
const products = ref<any[]>([])
const loading = ref(true)
const finished = ref(false)
let page = 1

const banners = ['🎓 毕业季大促 · 全场低至 1 折', '👥 宿舍拼团 · 满 3 人成团', '🔄 以物换物 · 闲置焕新生']
const currentBanner = ref(0)
let bannerTimer = 0
let touchStartX = 0

const cats = [
  { icon: '📚', text: '考研资料', to: '/forum' },
  { icon: '💻', text: '数码外设', to: '/' },
  { icon: '🧴', text: '日用百货', to: '/' },
  { icon: '🔄', text: '以物换物', to: '/barter' },
  { icon: '📦', text: '宿舍拼团', to: '/groupbuy' },
  { icon: '⏰', text: '物品租赁', to: '/rental' },
  { icon: '💬', text: '论坛', to: '/forum' },
  { icon: '👤', text: '我的', to: '/user' },
]

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page, page_size: 10 }
    if (keyword.value) params.keyword = keyword.value
    const res: any = await getProducts(params)
    if (res.code === 0) {
      products.value.push(...res.data.items)
      if (res.data.items.length < 10) finished.value = true
      page++
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

const loadMore = () => fetchData()
const onSearch = () => { page = 1; products.value = []; finished.value = false; fetchData() }
const doRefresh = () => { page = 1; products.value = []; finished.value = false; fetchData() }


const startBannerAuto = () => {
  clearInterval(bannerTimer)
  bannerTimer = window.setInterval(() => { currentBanner.value = (currentBanner.value + 1) % banners.length }, 3500)
}
const onTouchStart = (e: TouchEvent) => { touchStartX = e.touches[0].clientX; clearInterval(bannerTimer) }
const onTouchEnd = (e: TouchEvent) => {
  const dx = touchStartX - e.changedTouches[0].clientX
  if (dx > 40) currentBanner.value = (currentBanner.value + 1) % banners.length
  if (dx < -40) currentBanner.value = (currentBanner.value - 1 + banners.length) % banners.length
  startBannerAuto()
}

onMounted(() => { fetchData(); startBannerAuto() })
const bannerLinks = ['/', '/groupbuy', '/barter']
const goBanner = (i: number) => router.push(bannerLinks[i] || '/')

onUnmounted(() => clearInterval(bannerTimer))
</script>

<style scoped>
.home { background: #f5f5f5; }

.search-bar { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: #fff; }
.publish-btn { padding: 9px 14px; background: linear-gradient(135deg, #ff6b6b, #ee5a24); color: #fff; border-radius: 20px; font-size: 13px; font-weight: 600; white-space: nowrap; cursor: pointer; flex-shrink: 0; }
.search-box { display: flex; align-items: center; gap: 8px; background: #f5f5f5; border-radius: 20px; padding: 9px 16px; flex: 1; }
.search-icon { font-size: 14px; }
.search-input { flex: 1; border: none; background: transparent; font-size: 13px; outline: none; color: #333; }

.swiper { overflow: hidden; margin: 10px 12px; border-radius: 10px; position: relative; }
.swiper-track { display: flex; transition: transform 0.4s ease; }
.banner { min-width: 100%; height: 85px; padding: 14px 16px; color: #fff; font-size: 14px; font-weight: 600; display: flex; align-items: center; text-shadow: 0 1px 2px rgba(0,0,0,0.3); flex-shrink: 0; cursor: pointer; }
.banner:nth-child(1) { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.banner:nth-child(2) { background: linear-gradient(135deg, #4834d4, #686de0); }
.banner:nth-child(3) { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.swiper-dots { position: absolute; bottom: 8px; left: 50%; transform: translateX(-50%); display: flex; gap: 6px; }
.swiper-dots span { width: 6px; height: 6px; border-radius: 50%; background: rgba(255,255,255,0.5); transition: 0.2s; }
.swiper-dots span.active { background: #fff; width: 16px; border-radius: 3px; }

.cats { display: grid; grid-template-columns: repeat(4, 1fr); background: #fff; padding: 8px 12px 12px; }
.cat-item { display: flex; flex-direction: column; align-items: center; padding: 8px 0; gap: 5px; cursor: pointer; }
.cat-icon { font-size: 26px; }
.cat-text { font-size: 11px; color: #666; }

.list-header { padding: 16px 12px 8px; font-size: 16px; font-weight: 700; color: #222; cursor: pointer; user-select: none; }
.refresh-icon { color: #ff4d4f; font-size: 18px; }
.loading, .end { text-align: center; padding: 24px; color: #bbb; font-size: 12px; }

.products { padding: 0 12px; }
.product-card { display: flex; gap: 12px; background: #fff; border-radius: 12px; padding: 12px; margin-bottom: 10px; cursor: pointer; }
.product-card:active { background: #f9f9f9; }

.card-pic { width: 110px; min-height: 110px; border-radius: 10px; flex-shrink: 0; }
.pic-0 { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.pic-1 { background: linear-gradient(135deg, #4834d4, #686de0); }
.pic-2 { background: linear-gradient(135deg, #22a6b3, #7ed6df); }
.pic-3 { background: linear-gradient(135deg, #f9ca24, #f0932b); }

.card-right { flex: 1; display: flex; flex-direction: column; justify-content: space-between; min-width: 0; }
.card-title { font-size: 15px; font-weight: 600; color: #222; line-height: 1.4; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; word-break: break-all; }
.card-tags { display: flex; gap: 6px; margin: 8px 0; flex-wrap: wrap; }
.tag { font-size: 10px; padding: 2px 8px; border-radius: 3px; line-height: 1.4; white-space: nowrap; }
.tag-cat { background: #fff0f0; color: #ff4d4f; }
.tag-sub { background: #f0f0f8; color: #666; }

.card-footer { display: flex; align-items: flex-end; justify-content: space-between; margin-top: auto; }
.card-prices { display: flex; align-items: baseline; gap: 6px; }
.price-now { font-size: 17px; font-weight: 700; color: #ff4d4f; }
.price-old { font-size: 11px; color: #bbb; text-decoration: line-through; }
.card-campus { font-size: 11px; color: #aaa; flex-shrink: 0; }

.load-more { text-align: center; padding: 14px; color: #ff4d4f; font-size: 13px; cursor: pointer; }
</style>
