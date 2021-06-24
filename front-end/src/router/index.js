import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [ 
  {
    path: '/', 
    name: 'AddInfo', 
    component: () => import('../views/AddInfo.vue')
  }, 
  {
    path: '/buyers', 
    name: 'Buyers', 
    component: () => import('../views/Buyers.vue')
  }, 
  {
    path: '/buyer/:buyerId', 
    name: 'SingleBuyer', 
    component: () => import('../views/BuyerInfo.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
