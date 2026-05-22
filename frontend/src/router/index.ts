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
  { path: '/tasks', component: () => import('../pages/tasks/index.vue') },
  { path: '/search', component: () => import('../pages/index/search.vue') },
  { path: '/subscription', component: () => import('../pages/subscription/list.vue') },
  { path: '/subscription/publish', component: () => import('../pages/subscription/publish.vue') },
  { path: '/subscription/order/:id', component: () => import('../pages/subscription/credential.vue') },
  { path: '/notifications', component: () => import('../pages/messages/notifications.vue') },
  { path: '/messages', component: () => import('../pages/messages/index.vue') },
  { path: '/chat/:id', component: () => import('../pages/messages/chat.vue') },
  { path: '/rental/:id', component: () => import('../pages/rental/detail.vue') },
  { path: '/barter/:id', component: () => import('../pages/barter/detail.vue') },
  { path: '/tribunal', component: () => import('../pages/tribunal/list.vue') },
  { path: '/tribunal/:id', component: () => import('../pages/tribunal/detail.vue') },
  { path: '/login', component: () => import('../pages/user/login.vue') },
  { path: '/published', component: () => import('../pages/user/published.vue') },
  { path: '/history', component: () => import('../pages/user/history.vue') },
  { path: '/order/:id', component: () => import('../pages/user/order-detail.vue') },
  { path: '/qrcode/:id', component: () => import('../pages/user/qrcode.vue') },
  { path: '/settings', component: () => import('../pages/user/settings.vue') },
  { path: '/favorites', component: () => import('../pages/user/favorites.vue') },
  { path: '/edit-profile', component: () => import('../pages/user/edit.vue') },
  { path: '/publish', component: () => import('../pages/product/publish.vue') },
  { path: '/post/new', component: () => import('../pages/forum/publish.vue') },
]

const router = createRouter({
  scrollBehavior: () => ({ top: 0 }),
  history: createWebHashHistory(),
  routes,
})

export default router
