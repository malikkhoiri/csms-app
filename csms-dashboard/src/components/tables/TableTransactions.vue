<template>
  <v-card>
    <v-card-title class="d-flex align-center justify-space-between">
      <span>Transactions Management</span>
      <v-btn color="primary" @click="refreshData" prepend-icon="mdi-refresh">
        Refresh
      </v-btn>
    </v-card-title>
    
    <v-card-text>
      <v-data-table
        :headers="headers"
        :items="transactions"
        :loading="loading"
        :search="search"
        class="elevation-1"
      >
        <template v-slot:top>
          <v-text-field
            v-model="search"
            label="Search Transactions"
            prepend-inner-icon="mdi-magnify"
            single-line
            hide-details
            class="mx-4"
          />
        </template>
        
        <template v-slot:item.actions="{ item }">
          <v-icon size="small" @click="viewDetails(item)">
            mdi-eye
          </v-icon>
        </template>
        
        <template v-slot:item.id="{ item }">
          <span class="font-weight-medium">#{{ item?.id }}</span>
        </template>
        
        <template v-slot:item.chargePointId="{ item }">
          <span class="font-weight-medium">{{ item?.chargePoint?.chargePointCode || `CP-${item?.chargePointId}` }}</span>
        </template>
        
        <template v-slot:item.connectorId="{ item }">
          <v-chip size="small" color="info" variant="outlined">
            Connector {{ item?.connectorId }}
          </v-chip>
        </template>
        
        <template v-slot:item.idTag="{ item }">
          <div class="d-flex flex-column">
            <span class="font-weight-medium">{{ item?.idTag?.tag || 'N/A' }}</span>
          </div>
        </template>
        
        <template v-slot:item.status="{ item }">
          <v-chip
            :color="getStatusColor(item?.status)"
            :text="item?.status"
            size="small"
          />
        </template>
        
        <template v-slot:item.startTime="{ item }">
          {{ formatDateTime(item?.startTime) }}
        </template>
        
        <template v-slot:item.stopTime="{ item }">
          {{ formatDateTime(item?.stopTime) }}
        </template>
        
        <template v-slot:item.energyConsumed="{ item }">
          <span class="font-weight-medium">{{ item?.energyConsumed?.toFixed(2) || '0.00' }} kWh</span>
        </template>
        
        <template v-slot:item.currentMeterValue="{ item }">
          <span class="font-weight-medium">{{ item?.currentMeterValue?.toFixed(2) || '0.00' }}</span>
        </template>
        
        <template v-slot:item.totalCost="{ item }">
          <span class="font-weight-medium text-success">{{ formatCurrency(item?.totalCost || 0) }}</span>
        </template>
        
        <template v-slot:item.createdAt="{ item }">
          {{ formatDate(item?.createdAt) }}
        </template>
      </v-data-table>
    </v-card-text>

    <!-- Transaction Details Dialog -->
    <v-dialog v-model="dialog" max-width="700px">
      <v-card>
        <v-card-title>
          <span class="text-h5">Transaction Details</span>
        </v-card-title>

        <v-card-text>
          <v-container v-if="selectedTransaction">
            <v-row>
              <v-col cols="6">
                <v-list>
                  <v-list-item>
                    <v-list-item-title>Transaction ID</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.transactionId }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Charge Point</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.chargePoint.chargePointCode }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Connector</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.connectorId }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>ID Tag</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.idTag.tag }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-col>
              <v-col cols="6">
                <v-list>
                  <v-list-item>
                    <v-list-item-title>Status</v-list-item-title>
                    <v-list-item-subtitle>
                      <v-chip
                        :color="getStatusColor(selectedTransaction.status)"
                        :text="selectedTransaction.status"
                        size="small"
                      />
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Start Time</v-list-item-title>
                    <v-list-item-subtitle>{{ formatDateTime(selectedTransaction.startTime) }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Stop Time</v-list-item-title>
                    <v-list-item-subtitle>{{ formatDateTime(selectedTransaction.stopTime) }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Duration</v-list-item-title>
                    <v-list-item-subtitle>{{ calculateDuration(selectedTransaction.startTime, selectedTransaction.stopTime) }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-col>
            </v-row>
            
            <v-divider class="my-4" />
            
            <v-row>
              <v-col cols="12">
                <h3>Energy & Cost Details</h3>
                <v-list>
                  <v-list-item>
                    <v-list-item-title>Total Energy</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.totalEnergy }} kWh</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Total Cost</v-list-item-title>
                    <v-list-item-subtitle>{{ formatCurrency(selectedTransaction.totalCost) }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Meter Start</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.meterStart }}</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Meter Stop</v-list-item-title>
                    <v-list-item-subtitle>{{ selectedTransaction.meterStop }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer />
          <v-btn color="blue-darken-1" variant="text" @click="dialog = false">
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useDashboardStore } from '../../stores/dashboard'

const dashboard = useDashboardStore()
const search = ref('')
const dialog = ref(false)
const selectedTransaction = ref(null)

const headers = [
  { title: 'ID', key: 'id', sortable: true },
  { title: 'Charge Point', key: 'chargePointId', sortable: true },
  { title: 'Connector', key: 'connectorId', sortable: true },
  { title: 'ID Tag', key: 'idTag', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Start Time', key: 'startTime', sortable: true },
  { title: 'Stop Time', key: 'stopTime', sortable: true },
  { title: 'Current Meter', key: 'currentMeterValue', sortable: true },
  { title: 'Energy (kWh)', key: 'energyConsumed', sortable: true },
  { title: 'Cost (Rp)', key: 'totalCost', sortable: true },
  { title: 'Created At', key: 'createdAt', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

// Computed properties to get data from store
const transactions = computed(() => dashboard.transactions)
const loading = computed(() => dashboard.transactionsLoading)

onMounted(() => {
  fetchTransactions()
})

async function fetchTransactions() {
  await dashboard.fetchTransactions()
}

function refreshData() {
  fetchTransactions()
}

function viewDetails(item) {
  selectedTransaction.value = item
  dialog.value = true
}

function getStatusColor(status) {
  switch (status) {
    case 'Active': return 'success'
    case 'Completed': return 'primary'
    case 'Cancelled': return 'warning'
    case 'Failed': return 'error'
    case 'Pending': return 'info'
    default: return 'grey'
  }
}

function formatDateTime(dateString) {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return '-'
    return date.toLocaleString()
  } catch (error) {
    console.warn('Error formatting date time:', error)
    return '-'
  }
}

function formatDate(dateString) {
  if (!dateString) return ''
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return ''
    return date.toLocaleDateString()
  } catch (error) {
    console.warn('Error formatting date:', error)
    return ''
  }
}

function formatCurrency(amount) {
  if (typeof amount !== 'number') {
    amount = parseFloat(amount) || 0
  }
  
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount)
}

function calculateDuration(startTime, stopTime) {
  if (!startTime || !stopTime) return '-'
  
  try {
    const start = new Date(startTime)
    const stop = new Date(stopTime)
    
    if (isNaN(start.getTime()) || isNaN(stop.getTime())) return '-'
    
    const diffMs = stop - start
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMins / 60)
    const remainingMins = diffMins % 60
    
    if (diffHours > 0) {
      return `${diffHours}h ${remainingMins}m`
    } else {
      return `${remainingMins}m`
    }
  } catch (error) {
    console.warn('Error calculating duration:', error)
    return '-'
  }
}
</script> 