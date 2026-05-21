import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../pages/index/index.vue') },
  { path: '/product/:id', component: () => import('../pages/product/detail.vue') },
  { path: '/forum', component: () => import('../pages/forum/list.vue') },
  { path: '/forum/:id', component: () => import('../pages/forum/detail.vue') },
  { path: '/barter', component: () => import('../pages/barter/list.vue') },
  { path: '/groupbuy', component: () => import('../pages/groupbuy/list.vue') },
  { path: '/rental', component: () => import('../pages/rental/list.vue') },
  { path: '/orders', component: () => import('../pages/user/orders.vue') },
  { path: '/user', component: () => import('../pages/user/profile.vue') },
  { path: '/login', component: () => import('../pages/user/login.vue') },
  { path: '/published', component: () => import('../pages/user/published.vue') },
  { path: '/publish', component: () => import('../pages/product/publish.vue') },
  { path: '/post/new', component: () => import('../pages/forum/publish.vue') },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
