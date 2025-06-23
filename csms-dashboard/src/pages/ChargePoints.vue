<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex justify-space-between align-center">
            <span>Charge Points</span>
            <v-btn 
              color="primary" 
              @click="refreshData"
              :loading="dashboard.chargePointsLoading"
            >
              <v-icon>mdi-refresh</v-icon>
              Refresh
            </v-btn>
          </v-card-title>
          <v-card-text>
            <div v-if="dashboard.chargePointsError" class="mb-4">
              <v-alert type="error">{{ dashboard.chargePointsError }}</v-alert>
            </div>
            <TableChargePoints 
              :items="dashboard.chargePoints"
              :loading="dashboard.chargePointsLoading"
              @status-update="handleStatusUpdate"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <RemoteCommandPanel @send-command="handleRemoteCommand" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import TableChargePoints from '../components/tables/TableChargePoints.vue'
import RemoteCommandPanel from '../components/common/RemoteCommandPanel.vue'

const dashboard = useDashboardStore()

onMounted(async () => {
  await dashboard.fetchChargePoints()
})

const refreshData = async () => {
  await dashboard.fetchChargePoints()
}

const handleStatusUpdate = async (chargePointId, status) => {
  try {
    await dashboard.updateChargePointStatus(chargePointId, status)
  } catch (error) {
    console.error('Error updating status:', error)
  }
}

const handleRemoteCommand = async (chargePointId, command) => {
  try {
    await dashboard.sendRemoteCommand(chargePointId, command)
  } catch (error) {
    console.error('Error sending remote command:', error)
  }
}
</script> 