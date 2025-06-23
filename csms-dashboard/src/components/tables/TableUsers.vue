<template>
  <v-card>
    <v-card-title class="d-flex align-center justify-space-between">
      <span>Users Management</span>
      <v-btn color="primary" @click="dialog = true" prepend-icon="mdi-plus">
        Add User
      </v-btn>
    </v-card-title>
    
    <v-card-text>
      <v-data-table
        :headers="headers"
        :items="users"
        :loading="loading"
        :search="search"
        class="elevation-1"
      >
        <template v-slot:top>
          <v-text-field
            v-model="search"
            label="Search Users"
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
                  v-model="editedItem.name"
                  label="Name"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="editedItem.email"
                  label="Email"
                  type="email"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="editedItem.phone"
                  label="Phone"
                  required
                />
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="editedItem.role"
                  :items="roles"
                  label="Role"
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
        <v-card-title class="text-h5">Are you sure you want to delete this user?</v-card-title>
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
  { title: 'Name', key: 'name', sortable: true },
  { title: 'Email', key: 'email', sortable: true },
  { title: 'Phone', key: 'phone', sortable: true },
  { title: 'Role', key: 'role', sortable: true },
  { title: 'Created At', key: 'createdAt', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false }
]

const roles = ['admin', 'customer', 'operator']

const editedIndex = ref(-1)
const editedItem = ref({
  id: null,
  name: '',
  email: '',
  phone: '',
  role: 'customer'
})

const defaultItem = ref({
  id: null,
  name: '',
  email: '',
  phone: '',
  role: 'customer'
})

// Computed properties to get data from store
const users = computed(() => dashboard.users)
const loading = computed(() => dashboard.usersLoading)

const formTitle = computed(() => {
  return editedIndex.value === -1 ? 'New User' : 'Edit User'
})

onMounted(() => {
  fetchUsers()
  window.addEventListener('close-dialogs', cleanupDialogs)
})

async function fetchUsers() {
  await dashboard.fetchUsers()
}

function editItem(item) {
  editedIndex.value = users.value.indexOf(item)
  editedItem.value = Object.assign({}, item)
  dialog.value = true
}

function deleteItem(item) {
  editedIndex.value = users.value.indexOf(item)
  editedItem.value = Object.assign({}, item)
  dialogDelete.value = true
}

async function deleteItemConfirm() {
  try {
    await dashboard.deleteUser(editedItem.value.id)
    closeDelete()
  } catch (error) {
    console.error('Error deleting user:', error)
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
    if (editedIndex.value > -1) {
      // Update existing user
      await dashboard.updateUser(editedItem.value.id, editedItem.value)
    } else {
      // Create new user
      await dashboard.createUser(editedItem.value)
    }
    close()
  } catch (error) {
    console.error('Error saving user:', error)
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