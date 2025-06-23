<template>
  <v-card>
    <v-card-title class="d-flex align-center justify-space-between">
      <span>ID Tags Management</span>
      <v-btn color="primary" @click="dialog = true" prepend-icon="mdi-plus">
        Add ID Tag
      </v-btn>
    </v-card-title>
    
    <v-card-text>
      <v-data-table
        :headers="headers"
        :items="idTags"
        :loading="loading"
        :search="search"
        class="elevation-1"
      >
        <template v-slot:top>
          <v-text-field
            v-model="search"
            label="Search ID Tags"
            prepend-inner-icon="mdi-magnify"
            single-line
            hide-details
            class="mx-4"
          />
        </template>
        
        <template v-slot:item.actions="{ item }">
          <v-icon size="small" class="me-2" @click="editItem(item)">
            mdi-pencil
          </v-icon>
          <v-icon size="small" @click="deleteItem(item)">
            mdi-delete
          </v-icon>
        </template>
        
        <template v-slot:item.userName="{ item }">
          {{ item?.user?.name || `User ID: ${item?.userId}` || 'N/A' }}
        </template>
        
        <template v-slot:item.status="{ item }">
          <v-chip
            :color="getStatusColor(item?.status)"
            :text="item?.status"
            size="small"
          />
        </template>
        
        <template v-slot:item.expiryDate="{ item }">
          {{ formatDate(item?.expiryDate) }}
        </template>
        
        <template v-slot:item.createdAt="{ item }">
          {{ formatDate(item?.createdAt) }}
        </template>
      </v-data-table>
    </v-card-text>

    <!-- Add/Edit Dialog -->
    <v-dialog v-model="dialog" max-width="500px" @click:outside="close" persistent>
      <v-card>
        <v-card-title>
          <span class="text-h5">{{ formTitle }}</span>
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="editedItem.tag"
                  label="Tag"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="editedItem.userId"
                  :items="users"
                  item-title="name"
                  item-value="id"
                  label="User"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="editedItem.status"
                  :items="statuses"
                  label="Status"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="editedItem.expiryDate"
                  label="Expiry Date"
                  type="datetime-local"
                  required
                />
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer />
          <v-btn color="blue-darken-1" variant="text" @click="close">
            Cancel
          </v-btn>
          <v-btn color="blue-darken-1" variant="text" @click="save">
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="dialogDelete" max-width="500px" @click:outside="closeDelete" persistent>
      <v-card>
        <v-card-title class="text-h5">Are you sure you want to delete this ID tag?</v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="blue-darken-1" variant="text" @click="closeDelete">Cancel</v-btn>
          <v-btn color="blue-darken-1" variant="text" @click="deleteItemConfirm">OK</v-btn>
          <v-spacer />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { useDashboardStore } from '../../stores/dashboard'

const dashboard = useDashboardStore()
const search = ref('')
const dialog = ref(false)
const dialogDelete = ref(false)

const headers = [
  { title: 'ID', key: 'id', sortable: true },
  { title: 'Tag', key: 'tag', sortable: true },
  { title: 'User', key: 'userName', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Expiry Date', key: 'expiryDate', sortable: true },
  { title: 'Created At', key: 'createdAt', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

const statuses = ['Accepted', 'Blocked', 'Expired', 'Invalid', "ConcurrentTx"]

const editedIndex = ref(-1)
const editedItem = ref({
  id: null,
  tag: '',
  userId: null,
  status: 'Accepted',
  expiryDate: ''
})

const defaultItem = ref({
  id: null,
  tag: '',
  userId: null,
  status: 'Accepted',
  expiryDate: ''
})

// Computed properties to get data from store
const idTags = computed(() => {
  return dashboard.idTags
})
const users = computed(() => dashboard.users)
const loading = computed(() => dashboard.idTagsLoading)

const formTitle = computed(() => {
  return editedIndex.value === -1 ? 'New Tag' : 'Edit Tag'
})

onMounted(() => {
  fetchIdTags()
  fetchUsers()
  window.addEventListener('close-dialogs', cleanupDialogs)
})

async function fetchIdTags() {
  await dashboard.fetchIDTags()
}

async function fetchUsers() {
  await dashboard.fetchUsers()
}

function editItem(item) {
  editedIndex.value = idTags.value.indexOf(item)
  
  const itemCopy = {
    id: item.id,
    tag: item.tag,
    userId: item.userId,
    status: item.status,
    expiryDate: item.expiryDate
  }
  
  if (itemCopy.expiryDate) {
    itemCopy.expiryDate = formatDateTimeForInput(itemCopy.expiryDate)
  }
  
  editedItem.value = itemCopy
  dialog.value = true
}

function formatDateTimeForInput(dateTimeString) {
  if (!dateTimeString) return ''
  
  try {
    const date = new Date(dateTimeString)
    if (isNaN(date.getTime())) return ''
    
    // Format to datetime-local input format: "2025-06-30T17:16"
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    
    return `${year}-${month}-${day}T${hours}:${minutes}`
  } catch (error) {
    console.warn('Error formatting datetime for input:', error)
    return ''
  }
}

function deleteItem(item) {
  editedIndex.value = idTags.value.indexOf(item)
  editedItem.value = Object.assign({}, item)
  dialogDelete.value = true
}

async function deleteItemConfirm() {
  try {
    await dashboard.deleteIDTag(editedItem.value.id)
    closeDelete()
  } catch (error) {
    console.error('Error deleting ID tag:', error)
  }
}

function close() {
  dialog.value = false
  editedIndex.value = -1
  editedItem.value = Object.assign({}, defaultItem.value)
}

function closeDelete() {
  dialogDelete.value = false
  editedIndex.value = -1
  editedItem.value = Object.assign({}, defaultItem.value)
}

async function save() {
  try {
    const formattedData = {
      id: editedItem.value.id,
      tag: editedItem.value.tag,
      userId: editedItem.value.userId,
      status: editedItem.value.status,
      expiryDate: formatDateTimeForBackend(editedItem.value.expiryDate)
    }

    if (editedIndex.value > -1) {
      // Update existing ID tag
      await dashboard.updateIDTag(formattedData.id, formattedData)
    } else {
      // Create new ID tag
      await dashboard.createIDTag(formattedData)
    }
    close()
  } catch (error) {
    console.error('Error saving ID tag:', error)
  }
}

function formatDateTimeForBackend(dateTimeString) {
  if (!dateTimeString) return null
  
  try {
    // datetime-local input format: "2025-06-30T17:16"
    // Convert to RFC3339 format: "2025-06-30T17:16:00Z"
    const date = new Date(dateTimeString)
    if (isNaN(date.getTime())) return null
    
    // Format to RFC3339 with UTC timezone
    return date.toISOString()
  } catch (error) {
    console.warn('Error formatting datetime for backend:', error)
    return null
  }
}

function getStatusColor(status) {
  switch (status) {
    case 'Accepted': return 'success'
    case 'Blocked': return 'error'
    case 'Expired': return 'warning'
    case 'Invalid': return 'error'
    case 'ConcurrentTx': return 'warning'
    default: return 'primary'
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

// Cleanup function to close all dialogs
function cleanupDialogs() {
  dialog.value = false
  dialogDelete.value = false
  editedIndex.value = -1
  editedItem.value = Object.assign({}, defaultItem.value)
}

// Cleanup on component unmount
onBeforeUnmount(() => {
  cleanupDialogs()
  window.removeEventListener('close-dialogs', cleanupDialogs)
})
</script> 