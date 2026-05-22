<template>
  <div class="user">
    <div class="header">
      <div class="avatar" @click="$router.push('/edit-profile')">😊</div>
      <div class="user-info">
        <div class="nick" @click="$router.push('/edit-profile')">{{ store.userInfo?.nickname || '校园用户' }}</div>
        <div class="campus">{{ store.userInfo?.campus || '未认证' }}</div>
      </div>
      <div class="login-btn" v-if="!store.token" @click="$router.push('/login')">立即登录</div>
      <div class="login-btn" v-else @click="logout">退出</div>
    </div>

    
    <div class="stats">
      <div class="stat"><b>0</b><span>收藏</span></div>
      <div class="stat"><b>0</b><span>浏览</span></div>
      <div class="stat"><b>0</b><span>交易</span></div>
    </div>

    <div class="menu">
      <div class="menu-item" @click="$router.push('/orders')">
        <span>📦 我的订单</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/groupbuy')">
        <span>👥 我的拼团</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/rental')">
        <span>⏰ 我的租赁</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/barter')">
        <span>🔄 以物换物</span><span class="arrow">›</span>
      </div>
    </div>

    <div class="menu">
      <div class="menu-item" @click="$router.push('/login')" v-if="!store.token">
        <span>🔐 立即登录</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/published')">
        <span>🔔 通知中心</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/published')">
        <span>📝 我发布的</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/favorites')">
        <span>⭐ 我的收藏</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" @click="$router.push('/history')">
        <span>👁 浏览历史</span><span class="arrow">›</span>
      </div>
      <div class="menu-item" ><span>⚖️ 小法庭</span><span class="arrow">›</span></div>
      <div class="menu-item" @click="$router.push('/settings')">
        <span>⚙️ 设置</span><span class="arrow">›</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '../../store'

const store = useUserStore()

const logout = () => {
  store.logout()
  localStorage.removeItem('token')
  window.$toast('已退出')
}

const showMsg = (m: string) => window.$toast(m)
</script>

<style scoped>
.user { background: #F8FAFC; min-height: 100vh; }

.header {
  display: flex; align-items: center; gap: 12px;
  padding: 24px 16px; background: linear-gradient(135deg, #3B82F6, #2563EB);
  color: #fff;
}
.avatar {
  width: 56px; height: 56px; border-radius: 50%; background: rgba(255,255,255,0.3);
  display: flex; align-items: center; justify-content: center; font-size: 28px;
}
.user-info { flex: 1; }
.nick { font-size: 18px; font-weight: 600; }
.campus { font-size: 12px; opacity: 0.8; margin-top: 2px; }
.login-btn {
  padding: 6px 16px; border: 1px solid rgba(255,255,255,0.8);
  border-radius: 16px; font-size: 12px; cursor: pointer;
}

.stats { display: flex; background: #fff; padding: 16px 0; margin-bottom: 10px; }
.stat { flex: 1; text-align: center; display: flex; flex-direction: column; gap: 4px; }
.stat b { font-size: 20px; color: #333; }
.stat span { font-size: 12px; color: #999; }

.menu { background: #fff; margin-bottom: 10px; }
.menu-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 16px; font-size: 14px; color: #333; cursor: pointer;
  border-bottom: 1px solid #F8FAFC;
}
.menu-item:last-child { border-bottom: none; }
.arrow { color: #ccc; font-size: 18px; }










.cl-text { font-size: 15px; font-weight: 700; color: #059669; }
.cl-score { font-size: 12px; color: #94a3b8; margin-top: 2px; }
</style>
