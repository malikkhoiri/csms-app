<template>
  <v-app>
    <!-- Navigation drawer - only show when not on login page -->
    <v-navigation-drawer v-model="drawer" app theme="dark" v-if="!isLoginPage">
      <v-list>
        <v-list-item title="CSMS" />
        <v-divider />
        <v-list-item v-for="item in menu" :key="item.to" :to="item.to" link>
          <template v-slot:prepend>
            <v-icon>{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    
    <!-- App bar - only show when not on login page -->
    <v-app-bar app v-if="!isLoginPage">
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <!-- <v-toolbar-title>{{ currentPageTitle }}</v-toolbar-title> -->
      <v-spacer />
      <v-btn v-if="auth.token" icon @click="logout" title="Logout">
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-app-bar>
    
    <!-- Main content -->
    <v-main :class="{ 'pa-0': isLoginPage }">
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from './stores/auth'
import { useRouter, useRoute } from 'vue-router'

const drawer = ref(false)
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

// Check if current page is login page
const isLoginPage = computed(() => {
  return route.path === '/login'
})

const menu = [
  { to: '/', title: 'Dashboard', icon: 'mdi-view-dashboard' },
  { to: '/charge-points', title: 'Charge Points', icon: 'mdi-ev-station' },
  { to: '/transactions', title: 'Transaksi', icon: 'mdi-cash' },
  { to: '/users', title: 'User & ID Tag', icon: 'mdi-account-group' },
]

// Dynamic page title based on current route
// const currentPageTitle = computed(() => {
//   const currentMenuItem = menu.find(item => item.to === route.path)
//   return currentMenuItem ? currentMenuItem.title : 'Dashboard'
// })

async function logout() {
  auth.clearAuthState()
  await router.push('/login')
}
</script>
