<template>
  <div class="settings">
    <div class="nav"><span class="back" @click="$router.back()">← 返回</span><span class="nav-title">设置</span></div>

    <!-- 账号与安全 -->
    <div class="group">
      <div class="group-title">账号与安全</div>
      <div class="card">
        <div class="item" @click="toast('请联系学校信息中心认证')"><span>学生身份认证</span><span class="right">未认证 ›</span></div>
        <div class="item" @click="bindPhone"><span>手机号绑定</span><span class="right">未绑定 ›</span></div>
        <div class="item" @click="changePass"><span>修改密码</span><span class="right">›</span></div>
        <div class="item last" @click="toast('即将上线')"><span class="danger">注销账号</span><span class="right">›</span></div>
      </div>
    </div>

    <!-- 消息与通知 -->
    <div class="group">
      <div class="group-title">消息与通知</div>
      <div class="card">
        <div class="item"><span>交易提醒</span><span class="switch" :class="{ on: notifyTrade }" @click="notifyTrade=!notifyTrade"><i></i></span></div>
        <div class="item"><span>降价订阅通知</span><span class="switch" :class="{ on: notifyPrice }" @click="notifyPrice=!notifyPrice"><i></i></span></div>
        <div class="item last"><span>夜间勿扰 (23:00-07:00)</span><span class="switch" :class="{ on: dnd }" @click="dnd=!dnd"><i></i></span></div>
      </div>
    </div>

    <!-- 隐私 -->
    <div class="group">
      <div class="group-title">隐私</div>
      <div class="card">
        <div class="item"><span>对同城隐藏位置</span><span class="switch" :class="{ on: hideLoc }" @click="hideLoc=!hideLoc"><i></i></span></div>
        <div class="item" @click="$router.push('/tribunal')"><span>纠纷记录</span><span class="right">›</span></div>
        <div class="item last"><span>允许陌生人私信</span><span class="switch" :class="{ on: allowMsg }" @click="allowMsg=!allowMsg"><i></i></span></div>
      </div>
    </div>

    <!-- 通用 -->
    <div class="group">
      <div class="group-title">通用</div>
      <div class="card">
        <div class="item" @click="clearCache">
          <span>清除缓存</span><span class="right">{{ cacheSize }}</span>
        </div>
        <div class="item"><span>当前版本</span><span class="right">V1.0.0</span></div>
        <div class="item last" @click="$router.push('/')"><span>关于 PXX</span><span class="right">›</span></div>
      </div>
    </div>

    <!-- 退出登录 -->
    <div class="logout-btn" @click="doLogout">退出登录</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store'
import api from '../../api'

const router = useRouter()
const store = useUserStore()
const notifyTrade = ref(true)
const notifyPrice = ref(true)
const dnd = ref(false)
const hideLoc = ref(false)
const allowMsg = ref(true)
const cacheSize = ref('24.5 MB')

const toast = (m: string) => window.$toast(m)

const clearCache = () => {
  localStorage.clear()
  cacheSize.value = '0 KB'
  window.$toast('缓存已清理', 'success')
}

const bindPhone = () => {
  const phone = prompt('输入手机号')
  if (!phone) return
  api.put('/user/profile', { phone }).then(() => window.$toast('绑定成功', 'success'))
}
const changePass = () => {
  window.$sheet('修改密码', ['更改登录密码', '设置支付密码'], async (i: number) => {
    if (i === 0) {
      const pass = prompt('输入新密码')
      if (!pass) return
      await api.put('/user/profile', { password: pass })
      window.$toast('修改成功', 'success')
    }
  })
}

const doLogout = () => {
  if (confirm('确定退出登录？')) {
    store.logout()
    localStorage.removeItem('token')
    router.push('/')
    window.$toast('已退出')
  }
}
</script>

<style scoped>
.settings { background: #F2F2F6; min-height: 100vh; padding-bottom: 40px; }
.nav { display: flex; align-items: center; gap: 10px; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #1D4ED8; cursor: pointer; }
.nav-title { font-size: 15px; font-weight: 600; }

.group { margin-top: 16px; }
.group-title { font-size: 12px; color: #8E8E93; padding: 0 20px 6px; text-transform: uppercase; letter-spacing: 0.5px; }

.card { background: #fff; margin: 0 14px; border-radius: 12px; overflow: hidden; }

.item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 16px; font-size: 14px; color: #1C1C1E; cursor: pointer;
  border-bottom: 0.5px solid #E5E5EA; transition: background 0.1s;
  min-height: 48px;
}
.item:active { background: #F9F9F9; }
.item.last { border-bottom: none; }
.right { color: #8E8E93; font-size: 13px; }
.danger { color: #FF3B30; }

/* Switch Toggle */
.switch { width: 50px; height: 30px; border-radius: 15px; background: #E5E5EA; position: relative; transition: 0.2s; cursor: pointer; }
.switch i { position: absolute; top: 2px; left: 2px; width: 26px; height: 26px; border-radius: 50%; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,0.2); transition: 0.2s; }
.switch.on { background: #34C759; }
.switch.on i { left: 22px; }

.logout-btn {
  margin: 24px 14px; background: #fff; text-align: center; padding: 14px;
  border-radius: 12px; color: #FF3B30; font-size: 15px; font-weight: 500; cursor: pointer;
}
.logout-btn:active { background: #F9F9F9; }
</style>
