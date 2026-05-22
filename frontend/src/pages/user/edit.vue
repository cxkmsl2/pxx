<template>
  <div class="page">
    <div class="nav"><span class="back" @click="$router.back()">← 返回</span><span class="nav-title">编辑资料</span>
      <span class="save" @click="save">保存</span></div>
    <div class="form">
      <div class="avatar-row" @click="toast('头像功能开发中')"><span class="lbl">头像</span><span class="av">😊</span></div>
      <div class="row"><span class="lbl">昵称</span><input v-model="nick" placeholder="你的昵称" class="inp" /></div>
      <div class="row"><span class="lbl">校区</span><input v-model="campus" placeholder="所在校区" class="inp" /></div>
      <div class="row"><span class="lbl">手机号</span><input v-model="phone" placeholder="绑定手机号" class="inp" /></div>
      <div class="row"><span class="lbl">默认宿舍</span><input v-model="dorm" placeholder="如：12栋301" class="inp" /></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../store'
import api from '../../api'

const store = useUserStore()
const nick = ref(''), campus = ref(''), phone = ref(''), dorm = ref('')

onMounted(() => {
  nick.value = store.userInfo?.nickname || ''
  campus.value = store.userInfo?.campus || ''
  phone.value = store.userInfo?.phone || ''
  dorm.value = store.userInfo?.dormitory || ''
})

const save = async () => {
  try {
    const res = await api.put('/user/profile', {
      nickname: nick.value, campus: campus.value, phone: phone.value, dormitory: dorm.value
    })
    if (res.code === 0) {
      store.setUserInfo({ ...store.userInfo, nickname: nick.value, campus: campus.value, dormitory: dorm.value })
      window.$toast('已保存', 'success')
    }
  } catch { window.$toast('保存失败', 'error') }
}
const toast = (m: string) => window.$toast(m)
</script>

<style scoped>
.page { background: #F2F2F6; min-height: 100vh; }
.nav { display: flex; align-items: center; gap: 10px; padding: 12px; background: #fff; }
.back { font-size: 14px; color: #1D4ED8; cursor: pointer; }
.nav-title { flex: 1; font-size: 15px; font-weight: 600; }
.save { color: #1D4ED8; font-size: 14px; font-weight: 600; cursor: pointer; }
.form { background: #fff; margin: 16px 14px; border-radius: 12px; overflow: hidden; }
.avatar-row { display: flex; justify-content: space-between; align-items: center; padding: 14px; border-bottom: 0.5px solid #E5E5EA; cursor: pointer; }
.av { font-size: 40px; }
.row { display: flex; align-items: center; padding: 0 14px; border-bottom: 0.5px solid #E5E5EA; min-height: 48px; }
.lbl { width: 80px; font-size: 14px; color: #1C1C1E; }
.inp { flex: 1; border: none; font-size: 14px; outline: none; text-align: right; padding: 14px 0; color: #1C1C1E; }
</style>