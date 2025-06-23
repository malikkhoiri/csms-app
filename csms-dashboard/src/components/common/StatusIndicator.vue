<template>
  <v-card class="d-flex align-center pa-4" outlined>
    <v-icon :color="statusColor" size="32">mdi-circle</v-icon>
    <span class="ml-4 text-h6">{{ label }}</span>
  </v-card>
</template>

<script setup>
const props = defineProps({
  status: String,
  label: String,
})
const statusColor = {
  online: 'success',
  offline: 'error',
  faulted: 'warning',
}[props.status] || 'grey'
</script> 