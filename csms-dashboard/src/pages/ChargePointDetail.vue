<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <!-- Header -->
        <div class="d-flex align-center justify-space-between mb-6">
          <div>
            <h1 class="text-h4 font-weight-bold mb-2">
              Charge Point Details
            </h1>
            <p class="text-subtitle-1 text-medium-emphasis">
              {{ chargePoint?.chargePointCode || 'Loading...' }}
            </p>
          </div>
          <v-btn
            color="primary"
            prepend-icon="mdi-arrow-left"
            @click="$router.push('/charge-points')"
          >
            Back to Charge Points
          </v-btn>
        </div>

        <!-- Loading State -->
        <v-skeleton-loader
          v-if="loading"
          type="card"
          class="mb-6"
        />

        <!-- Error State -->
        <v-alert
          v-else-if="error"
          type="error"
          class="mb-6"
        >
          {{ error }}
        </v-alert>

        <!-- Charge Point Details -->
        <div v-else-if="chargePoint" class="mb-6">
          <v-card>
            <v-card-title class="text-h5">
              <v-icon class="mr-2">mdi-ev-station</v-icon>
              Charge Point Information
            </v-card-title>
            <v-card-text>
              <v-row>
                <v-col cols="12" md="6">
                  <v-list>
                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-tag</v-icon>
                      </template>
                      <v-list-item-title>Charge Point Code</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.chargePointCode }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-barcode</v-icon>
                      </template>
                      <v-list-item-title>Serial Number</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.chargeBoxSerialNumber || 'N/A' }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-factory</v-icon>
                      </template>
                      <v-list-item-title>Vendor</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.chargePointVendor }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-cube</v-icon>
                      </template>
                      <v-list-item-title>Model</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.chargePointModel }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-chip</v-icon>
                      </template>
                      <v-list-item-title>Firmware Version</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.firmwareVersion }}</v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-col>

                <v-col cols="12" md="6">
                  <v-list>
                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-circle</v-icon>
                      </template>
                      <v-list-item-title>Status</v-list-item-title>
                      <v-list-item-subtitle>
                        <v-chip
                          :color="getStatusColor(chargePoint.status)"
                          size="small"
                        >
                          {{ chargePoint.status }}
                        </v-chip>
                      </v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-heart-pulse</v-icon>
                      </template>
                      <v-list-item-title>Last Heartbeat</v-list-item-title>
                      <v-list-item-subtitle>{{ formatDate(chargePoint.lastHeartbeat) }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-boot</v-icon>
                      </template>
                      <v-list-item-title>Last Boot</v-list-item-title>
                      <v-list-item-subtitle>{{ formatDate(chargePoint.lastBootNotification) }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-calendar</v-icon>
                      </template>
                      <v-list-item-title>Created</v-list-item-title>
                      <v-list-item-subtitle>{{ formatDate(chargePoint.createdAt) }}</v-list-item-subtitle>
                    </v-list-item>

                    <v-list-item>
                      <template #prepend>
                        <v-icon>mdi-update</v-icon>
                      </template>
                      <v-list-item-title>Last Updated</v-list-item-title>
                      <v-list-item-subtitle>{{ formatDate(chargePoint.updatedAt) }}</v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-col>
              </v-row>

              <!-- Additional Information -->
              <v-divider class="my-4"></v-divider>
              
              <v-row>
                <v-col cols="12" md="6">
                  <h3 class="text-h6 mb-3">Technical Details</h3>
                  <v-list density="compact">
                    <v-list-item>
                      <v-list-item-title>ICCID</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.iccid || 'N/A' }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>IMSI</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.imsi || 'N/A' }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Meter Type</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.meterType || 'N/A' }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Meter Serial Number</v-list-item-title>
                      <v-list-item-subtitle>{{ chargePoint.meterSerialNumber || 'N/A' }}</v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </div>

        <!-- Connectors Section -->
        <div v-if="chargePoint?.connectors?.length" class="mb-6">
          <v-card>
            <v-card-title class="text-h5">
              <v-icon class="mr-2">mdi-power-plug</v-icon>
              Connectors ({{ chargePoint.connectors.length }})
            </v-card-title>
            <v-card-text>
              <v-table>
                <thead>
                  <tr>
                    <th>Connector ID</th>
                    <th>Status</th>
                    <th>Error Code</th>
                    <th>Info</th>
                    <th>Vendor ID</th>
                    <th>Last Updated</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="connector in chargePoint.connectors" :key="connector.id">
                    <td>
                      <v-chip size="small" color="primary">
                        {{ connector.connectorId }}
                      </v-chip>
                    </td>
                    <td>
                      <v-chip
                        :color="getStatusColor(connector.status)"
                        size="small"
                      >
                        {{ connector.status }}
                      </v-chip>
                    </td>
                    <td>
                      <v-chip
                        :color="getErrorColor(connector.errorCode)"
                        size="small"
                      >
                        {{ connector.errorCode }}
                      </v-chip>
                    </td>
                    <td>{{ connector.info || 'N/A' }}</td>
                    <td>{{ connector.vendorId || 'N/A' }}</td>
                    <td>{{ formatDate(connector.updatedAt) }}</td>
                  </tr>
                </tbody>
              </v-table>
            </v-card-text>
          </v-card>
        </div>

        <!-- No Connectors Message -->
        <v-card v-else-if="chargePoint && (!chargePoint.connectors || chargePoint.connectors.length === 0)">
          <v-card-text class="text-center py-8">
            <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-power-plug-off</v-icon>
            <h3 class="text-h6 text-grey">No Connectors Found</h3>
            <p class="text-body-2 text-grey">This charge point doesn't have any connectors configured.</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useDashboardStore } from '../stores/dashboard'
import api from '../services/api'

const route = useRoute()
const dashboard = useDashboardStore()

const chargePoint = ref(null)
const loading = ref(true)
const error = ref(null)

// Get charge point ID from route params
const chargePointId = route.params.id

onMounted(async () => {
  try {
    loading.value = true
    error.value = null
    
    // Fetch charge point details
    const response = await api.get(`/api/v1/charge-points/${chargePointId}`)
    chargePoint.value = response.data
  } catch (err) {
    console.error('Error fetching charge point details:', err)
    error.value = 'Failed to load charge point details. Please try again.'
  } finally {
    loading.value = false
  }
})

// Helper functions
function formatDate(dateString) {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function getStatusColor(status) {
  const colors = {
    'Available': 'success',
    'Preparing': 'warning',
    'Charging': 'primary',
    'SuspendedEVSE': 'warning',
    'SuspendedEV': 'warning',
    'Finishing': 'info',
    'Reserved': 'info',
    'Unavailable': 'error'
  }
  return colors[status] || 'grey'
}

function getErrorColor(errorCode) {
  if (errorCode === 'NoError') return 'success'
  if (errorCode.includes('Error')) return 'error'
  return 'warning'
}
</script>

<style scoped>
.v-list-item {
  padding: 8px 0;
}
</style> 