import api from './api'

export const dashboardService = {
  // Dashboard Overview
  async getDashboardStats() {
    const response = await api.get('/api/v1/dashboard/stats')
    return response.data
  },

  async getWeeklyChart() {
    const response = await api.get('/api/v1/dashboard/weekly-chart')
    return response.data
  },

  // Charge Points
  async getChargePoints(params = {}) {
    const response = await api.get('/api/v1/charge-points', { params })
    return response.data
  },

  async getChargePoint(id) {
    const response = await api.get(`/api/v1/charge-points/${id}`)
    return response.data
  },

  async updateChargePointStatus(id, status) {
    const response = await api.patch(`/api/v1/charge-points/${id}/status`, { status })
    return response.data
  },

  // Transactions
  async getTransactions(params = {}) {
    const response = await api.get('/api/v1/transactions', { params })
    return response.data
  },

  async getTransaction(id) {
    const response = await api.get(`/api/v1/transactions/${id}`)
    return response.data
  },

  // Users
  async getUsers(params = {}) {
    const response = await api.get('/api/v1/users', { params })
    return response.data
  },

  async getUser(id) {
    const response = await api.get(`/api/v1/users/${id}`)
    return response.data
  },

  async createUser(userData) {
    const response = await api.post('/api/v1/users', userData)
    return response.data
  },

  async updateUser(id, userData) {
    const response = await api.put(`/api/v1/users/${id}`, userData)
    return response.data
  },

  async deleteUser(id) {
    const response = await api.delete(`/api/v1/users/${id}`)
    return response.data
  },

  // ID Tags
  async getIDTags(params = {}) {
    const response = await api.get('/api/v1/id-tags', { params })
    return response.data
  },

  async getIDTag(id) {
    const response = await api.get(`/api/v1/id-tags/${id}`)
    return response.data
  },

  async createIDTag(idTagData) {
    const response = await api.post('/api/v1/id-tags', idTagData)
    return response.data
  },

  async updateIDTag(id, idTagData) {
    const response = await api.put(`/api/v1/id-tags/${id}`, idTagData)
    return response.data
  },

  async deleteIDTag(id) {
    const response = await api.delete(`/api/v1/id-tags/${id}`)
    return response.data
  },

  // Remote Commands
  async sendRemoteCommand(chargePointId, command) {
    const response = await api.post(`/api/v1/charge-points/${chargePointId}/commands`, command)
    return response.data
  }
} 