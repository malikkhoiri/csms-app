<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <TableTransactions />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import TableTransactions from '../components/tables/TableTransactions.vue'

const dashboard = useDashboardStore()

onMounted(async () => {
  // Fetch data on component mount
  await dashboard.fetchTransactions()
})
</script> 