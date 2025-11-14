import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Implants from '@/views/Implants.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/implants',
      name: 'implants',
      component: Implants,
    },
  ],
})

export default router
