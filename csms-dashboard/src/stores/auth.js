import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import api from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('auth_token') || null)
  const user = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  // Watch for token changes and sync with localStorage
  watch(token, (newToken) => {
    if (newToken) {
      localStorage.setItem('auth_token', newToken)
    } else {
      localStorage.removeItem('auth_token')
    }
  })

  async function login(credentials) {
    loading.value = true
    try {
      const response = await api.post('/api/v1/auth/login', credentials)
      const { accessToken, user: userData } = response.data
      
      token.value = accessToken
      user.value = userData
      
      localStorage.setItem('user', JSON.stringify(userData))
      
      return { success: true }
    } catch (error) {
      console.error('Login error:', error)
      return { 
        success: false, 
        error: error.response?.data?.message || 'Login failed' 
      }
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await api.post('/api/v1/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuthState()
    }
  }

  function clearAuthState() {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user')
  }

  async function checkAuth() {
    if (!token.value) return false
    
    try {
      const response = await api.get('/api/v1/auth/me')
      user.value = response.data
      return true
    } catch (error) {
      console.error('Auth check error:', error)
      clearAuthState()
      return false
    }
  }

  function initAuth() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      user.value = JSON.parse(savedUser)
    }
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    login,
    logout,
    checkAuth,
    initAuth,
    clearAuthState
  }
}) 