<template>
  <div class="search-page">
    <div class="search-top">
      <span class="back" @click="$router.back()">←</span>
      <div class="search-input-wrap">
        <input v-model="kw" ref="inputRef" placeholder="搜搜闲置好物" class="search-input" @input="onInput" @keyup.enter="doSearch" />
        <span class="clear" v-if="kw" @click="kw='';doSearch()">✕</span>
      </div>
      <span class="search-btn" @click="doSearch">搜索</span>
    </div>

    <!-- 热门搜索 -->
    <div class="hot-wrap" v-if="!searched">
      <div class="hot-title">🔥 热门搜索</div>
      <div class="hot-tags">
        <span class="hot-tag" v-for="h in hots" :key="h" @click="kw=h;doSearch()">{{ h }}</span>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div v-else>
      <div class="result-head" v-if="!loading">共找到 {{ total }} 件商品</div>
      <div class="loading" v-if="loading">搜索中...</div>
      <div class="empty" v-if="!loading && products.length === 0">暂无相关商品</div>
      <div class="products" v-if="products.length > 0">
        <div class="product-card" v-for="item in products" :key="item.id" @click="$router.push('/product/'+item.id)">
          <div class="card-pic" :class="'bg' + (item.id % 4)"></div>
          <div class="card-info">
            <div class="card-title">{{ item.title }}</div>
            <div class="card-price">¥{{ (item.price / 100).toFixed(2) }}</div>
            <div class="card-meta">{{ item.campus }} · {{ item.seller?.nickname }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import api from '../../api'

const kw = ref('')
let debounceTimer = 0
const onInput = () => { clearTimeout(debounceTimer); debounceTimer = window.setTimeout(() => doSearch(), 300) }
const searched = ref(false)
const products = ref<any[]>([])
const total = ref(0)
const loading = ref(false)
const inputRef = ref()

const hots = ['iPad', '考研', '鼠标', '台灯', '手环', '自行车', '显示器', '收纳']

const doSearch = async () => {
  if (!kw.value.trim()) { searched.value = false; return }
  searched.value = true
  loading.value = true
  try {
    const res: any = await api.get('/products?keyword=' + encodeURIComponent(kw.value.trim()) + '&page_size=30')
    if (res.code === 0) { products.value = res.data.items; total.value = res.data.total }
  } catch {}
  loading.value = false
}

onMounted(() => nextTick(() => inputRef.value?.focus()))
</script>

<style scoped>
.search-page { background: #F8FAFC; min-height: 100vh; }

.search-top { display: flex; align-items: center; gap: 10px; padding: 10px 12px; background: #fff; position: sticky; top: 0; z-index: 10; }
.back { font-size: 18px; color: #333; cursor: pointer; }
.search-input-wrap { flex: 1; position: relative; }
.search-input { width: 100%; padding: 9px 30px 9px 14px; border: none; background: #F1F5F9; border-radius: 20px; font-size: 14px; outline: none; }
.clear { position: absolute; right: 10px; top: 50%; transform: translateY(-50%); color: #94a3b8; cursor: pointer; font-size: 14px; }
.search-btn { color: #1D4ED8; font-size: 14px; font-weight: 600; cursor: pointer; }

.hot-wrap { padding: 20px 14px; }
.hot-title { font-size: 14px; font-weight: 600; color: #1E293B; margin-bottom: 12px; }
.hot-tags { display: flex; flex-wrap: wrap; gap: 10px; }
.hot-tag { padding: 8px 16px; background: #F1F5F9; border-radius: 20px; font-size: 13px; color: #64748b; cursor: pointer; }
.hot-tag:active { background: #DBEAFE; color: #1D4ED8; }

.result-head { padding: 12px 14px 0; font-size: 13px; color: #94a3b8; }
.loading, .empty { text-align: center; padding: 40px; color: #94a3b8; }

.products { padding: 12px 14px; }
.product-card { display: flex; gap: 12px; background: #fff; border-radius: 14px; padding: 12px; margin-bottom: 10px; cursor: pointer; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.card-pic { width: 100px; height: 100px; border-radius: 10px; flex-shrink: 0; }
.bg0 { background: linear-gradient(135deg, #DBEAFE, #BFDBFE); }
.bg1 { background: linear-gradient(135deg, #D1FAE5, #A7F3D0); }
.bg2 { background: linear-gradient(135deg, #FEF3C7, #FDE68A); }
.bg3 { background: linear-gradient(135deg, #EDE9FE, #DDD6FE); }
.card-info { flex: 1; display: flex; flex-direction: column; justify-content: center; min-width: 0; }
.card-title { font-size: 14px; font-weight: 600; color: #1E293B; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-price { font-size: 18px; font-weight: 700; color: #10B981; margin-top: 4px; }
.card-meta { font-size: 11px; color: #94a3b8; margin-top: 4px; }
</style>