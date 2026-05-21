import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res.data,
  err => Promise.reject(err)
)

export default api

export const getProducts = (params: any) => api.get('/products', { params })
export const getProduct = (id: number) => api.get(`/products/${id}`)
export const getPosts = (params: any) => api.get('/posts', { params })
export const getPost = (id: number) => api.get(`/posts/${id}`)
export const getGroupBuys = (params: any) => api.get('/groupbuys', { params })
export const joinGroupBuy = (id: number, data: any) => api.post(`/groupbuys/${id}/join`, data)
export const getRentalItems = (params: any) => api.get('/rentals', { params })
export const getBarterItems = (params: any) => api.get('/barters', { params })
export const getOrders = (params: any) => api.get('/orders', { params })
