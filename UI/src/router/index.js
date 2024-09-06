import { createRouter, createWebHistory } from 'vue-router'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Request Monitor',
      component: () => import('../views/RequestsView.vue')
    },
    {
      path: '/mapping',
      name: 'Mapping',
      component: () => import('../views/MappingView.vue')
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
