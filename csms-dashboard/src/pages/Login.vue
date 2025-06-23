<template>
  <div class="login-page">
    <div class="login-container">
      <v-card class="login-card elevation-12">
        <v-toolbar color="primary" dark flat>
          <v-toolbar-title class="text-center w-100">
            <v-icon size="32" class="me-2">mdi-ev-station</v-icon>
            CSMS
          </v-toolbar-title>
        </v-toolbar>
        
        <v-card-text class="pa-6">
          <v-form @submit.prevent="handleLogin" ref="form">
            <v-text-field
              v-model="credentials.email"
              label="Email"
              name="email"
              prepend-inner-icon="mdi-email"
              type="email"
              :rules="emailRules"
              required
              variant="outlined"
              class="mb-4"
            />
            
            <v-text-field
              v-model="credentials.password"
              label="Password"
              name="password"
              prepend-inner-icon="mdi-lock"
              type="password"
              :rules="passwordRules"
              required
              variant="outlined"
              class="mb-6"
            />
          </v-form>
        </v-card-text>
        
        <v-card-actions class="pa-6 pt-0">
          <v-btn
            color="primary"
            :loading="auth.loading"
            @click="handleLogin"
            block
            size="large"
            class="text-body-1"
          >
            <v-icon start>mdi-login</v-icon>
            Login
          </v-btn>
        </v-card-actions>
        
        <v-alert
          v-if="error"
          type="error"
          variant="tonal"
          class="ma-6 mt-0"
          closable
          @click:close="error = ''"
        >
          {{ error }}
        </v-alert>
      </v-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const form = ref(null)
const error = ref('')

const credentials = reactive({
  email: '',
  password: ''
})

const emailRules = [
  v => !!v || 'Email is required',
  v => /.+@.+\..+/.test(v) || 'Email must be valid'
]

const passwordRules = [
  v => !!v || 'Password is required',
  v => v.length >= 6 || 'Password must be at least 6 characters'
]

async function handleLogin() {
  error.value = ''
  
  const { valid } = await form.value.validate()
  if (!valid) return
  
  const result = await auth.login(credentials)
  
  if (result.success) {
    router.push('/')
  } else {
    error.value = result.error
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  border-radius: 12px;
  overflow: hidden;
}

.v-toolbar-title {
  font-size: 1.5rem;
  font-weight: 600;
}
</style> 