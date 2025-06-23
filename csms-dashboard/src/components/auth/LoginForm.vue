<template>
  <v-card>
    <v-card-title>Login</v-card-title>
    <v-card-text>
      <v-form @submit.prevent="onLogin">
        <v-text-field v-model="email" label="Email" type="email" required />
        <v-text-field v-model="password" label="Password" type="password" required />
        <v-btn :loading="auth.loading" type="submit" color="primary" block>Login</v-btn>
        <v-alert v-if="auth.error" type="error" class="mt-2">{{ auth.error }}</v-alert>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const auth = useAuthStore()
const router = useRouter()

const onLogin = async () => {
  await auth.login(email.value, password.value)
  if (auth.token) {
    router.push('/')
  }
}
</script> 