<template>
  <div class="login">
    <div class="nav">🔐 登录 PXX</div>
    <div class="form">
      <div class="label">学号 / 校园卡号</div>
      <input v-model="openId" class="input" placeholder="输入你的学号或校园卡号" />
      <div class="hint">测试用：输入任意内容即可登录</div>
      <div class="btn" @click="doLogin" :class="{ disabled: loading }">
        {{ loading ? '登录中...' : '立即登录' }}
      </div>
      <div class="skip" @click="$router.back()">暂不登录，先逛逛 →</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store'
import api from '../../api'

const router = useRouter()
const store = useUserStore()
const openId = ref('')
const loading = ref(false)

const doLogin = async () => {
  if (!openId.value.trim()) { window.$toast('请输入账号'); return }
  loading.value = true
  try {
    const res = await api.post('/login', { open_id: openId.value.trim() })
    if (res.code === 0) {
      store.setToken(res.data.token)
      store.setUserInfo(res.data.user)
      localStorage.setItem('token', res.data.token)
      window.$toast('登录成功！')
      router.back()
    } else {
      window.$toast(res.msg)
    }
  } catch { window.$toast('网络错误') }
  loading.value = false
}
</script>

<style scoped>
.login { min-height: 100vh; background: #F8FAFC; }
.nav { padding: 14px; background: #fff; font-size: 16px; font-weight: 600; }
.form { padding: 30px 20px; }
.label { font-size: 14px; color: #333; margin-bottom: 8px; font-weight: 500; }
.input {
  width: 100%; padding: 12px 14px; border: 1px solid #e0e0e0; border-radius: 10px;
  font-size: 15px; outline: none; box-sizing: border-box;
}
.input:focus { border-color: #1D4ED8; }
.hint { font-size: 11px; color: #aaa; margin: 8px 0 24px; }
.btn {
  background: linear-gradient(135deg, #3B82F6, #2563EB); color: #fff;
  text-align: center; padding: 12px; border-radius: 10px; font-size: 15px; font-weight: 600; cursor: pointer;
}
.btn.disabled { opacity: 0.6; }
.skip { text-align: center; margin-top: 20px; font-size: 13px; color: #1D4ED8; cursor: pointer; }
</style>
