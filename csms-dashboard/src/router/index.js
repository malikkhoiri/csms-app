/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('../pages/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/charge-points',
    name: 'ChargePoints',
    component: () => import('../pages/ChargePoints.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/charge-points/:id',
    name: 'ChargePointDetail',
    component: () => import('../pages/ChargePointDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/transactions',
    name: 'Transactions',
    component: () => import('../pages/Transactions.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: () => import('../pages/Users.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()
  
  // Initialize auth on first load
  if (!auth.user && auth.token) {
    auth.initAuth()
  }
  
  // Close any open dialogs by dispatching a custom event
  window.dispatchEvent(new CustomEvent('close-dialogs'))
  
  if (to.meta.requiresAuth) {
    if (!auth.isAuthenticated) {
      next('/login')
    } else {
      // Check if user data is valid
      const isValid = await auth.checkAuth()
      if (!isValid) {
        next('/login')
      } else {
        next()
      }
    }
  } else {
    // For login page, redirect to dashboard if already authenticated
    if (to.path === '/login' && auth.isAuthenticated) {
      next('/')
    } else {
      next()
    }
  }
})

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (localStorage.getItem('vuetify:dynamic-reload')) {
      console.error('Dynamic import error, reloading page did not fix it', err)
    } else {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

export default router
