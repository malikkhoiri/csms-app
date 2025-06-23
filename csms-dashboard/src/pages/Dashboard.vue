<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="3">
        <StatCard 
          title="Charge Points" 
          :value="dashboard.stats.totalChargePoints" 
          icon="mdi-ev-station" 
        />
      </v-col>
      <v-col cols="12" md="3">
        <StatCard 
          title="Online" 
          :value="dashboard.stats.onlineChargePoints" 
          icon="mdi-check-circle" 
          color="success" 
        />
      </v-col>
      <v-col cols="12" md="3">
        <StatCard 
          title="Offline" 
          :value="dashboard.stats.offlineChargePoints" 
          icon="mdi-close-circle" 
          color="error" 
        />
      </v-col>
      <v-col cols="12" md="3">
        <StatCard 
          title="Transaksi Hari Ini" 
          :value="dashboard.stats.todayTransactions" 
          icon="mdi-cash" 
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" md="8">
        <ChartLine :chart-data="dashboard.weeklyChart" />
      </v-col>
      <v-col cols="12" md="4">
        <StatusIndicator status="online" label="Sistem Normal" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import StatCard from '../components/common/StatCard.vue'
import ChartLine from '../components/charts/ChartLine.vue'
import StatusIndicator from '../components/common/StatusIndicator.vue'

const dashboard = useDashboardStore()

onMounted(async () => {
  // Fetch dashboard data on component mount
  await Promise.all([
    dashboard.fetchDashboardStats(),
    dashboard.fetchWeeklyChart()
  ])
})
</script> 