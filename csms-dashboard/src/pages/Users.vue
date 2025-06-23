<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="6">
        <TableUsers />
      </v-col>
      <v-col cols="12" md="6">
        <TableIdTags />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import TableUsers from '../components/tables/TableUsers.vue'
import TableIdTags from '../components/tables/TableIdTags.vue'

const dashboard = useDashboardStore()

onMounted(async () => {
  // Fetch data on component mount
  await Promise.all([
    dashboard.fetchUsers(),
    dashboard.fetchIDTags()
  ])
})
</script> 