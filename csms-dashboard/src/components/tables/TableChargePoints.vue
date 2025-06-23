<template>
  <div>
    <v-data-table 
      :headers="headers" 
      :items="items" 
      :loading="loading"
      class="elevation-1"
    >
      <template #item.status="{ item }">
        <v-chip 
          :color="getStatusColor(item.status)" 
          :text="item.status"
          @click="updateStatus(item.id, item.status)"
        />
      </template>
      <template #item.actions="{ item }">
        <v-btn 
          size="small" 
          color="primary" 
          @click="viewDetails(item)"
        >
          Details
        </v-btn>
      </template>
    </v-data-table>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['status-update'])

const headers = [
  { title: 'ID', value: 'id' },
  { title: 'Code', value: 'chargePointCode' },
  { title: 'Model', value: 'chargePointModel' },
  { title: 'Vendor', value: 'chargePointVendor' },
  { title: 'Status', value: 'status' },
  { title: 'Last Heartbeat', value: 'lastHeartbeat' },
  { title: 'Actions', value: 'actions', sortable: false }
]

const getStatusColor = (status) => {
  switch (status?.toLowerCase()) {
    case 'available':
    case 'online':
      return 'success'
    case 'charging':
      return 'warning'
    case 'faulted':
    case 'offline':
      return 'error'
    default:
      return 'grey'
  }
}

const updateStatus = (id, currentStatus) => {
  // Toggle status or show status update dialog
  const newStatus = currentStatus === 'Available' ? 'Unavailable' : 'Available'
  emit('status-update', id, newStatus)
}

const viewDetails = (item) => {
  router.push(`/charge-points/${item.id}`)
}
</script> 