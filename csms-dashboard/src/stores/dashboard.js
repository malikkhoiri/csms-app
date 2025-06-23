import { defineStore } from 'pinia'
import { dashboardService } from '../services/dashboard'

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    // Dashboard Overview
    stats: {
      totalChargePoints: 0,
      onlineChargePoints: 0,
      offlineChargePoints: 0,
      todayTransactions: 0
    },
    weeklyChart: {
      labels: [],
      datasets: []
    },

    // Charge Points
    chargePoints: [],
    chargePointsLoading: false,
    chargePointsError: null,

    // Transactions
    transactions: [],
    transactionsLoading: false,
    transactionsError: null,

    // Users
    users: [],
    usersLoading: false,
    usersError: null,

    // ID Tags
    idTags: [],
    idTagsLoading: false,
    idTagsError: null
  }),

  actions: {
    // Dashboard Overview
    async fetchDashboardStats() {
      try {
        const stats = await dashboardService.getDashboardStats()
        this.stats = stats
      } catch (error) {
        console.error('Error fetching dashboard stats:', error)
      }
    },

    async fetchWeeklyChart() {
      try {
        const chartData = await dashboardService.getWeeklyChart()
        this.weeklyChart = chartData
      } catch (error) {
        console.error('Error fetching weekly chart:', error)
      }
    },

    // Charge Points
    async fetchChargePoints(params = {}) {
      this.chargePointsLoading = true
      this.chargePointsError = null
      try {
        const chargePoints = await dashboardService.getChargePoints(params)
        this.chargePoints = chargePoints
      } catch (error) {
        this.chargePointsError = error.message
        console.error('Error fetching charge points:', error)
      } finally {
        this.chargePointsLoading = false
      }
    },

    async updateChargePointStatus(id, status) {
      try {
        await dashboardService.updateChargePointStatus(id, status)
        // Refresh charge points data
        await this.fetchChargePoints()
      } catch (error) {
        console.error('Error updating charge point status:', error)
        throw error
      }
    },

    // Transactions
    async fetchTransactions(params = {}) {
      this.transactionsLoading = true
      this.transactionsError = null
      try {
        const transactions = await dashboardService.getTransactions(params)
        this.transactions = transactions
      } catch (error) {
        this.transactionsError = error.message
        console.error('Error fetching transactions:', error)
      } finally {
        this.transactionsLoading = false
      }
    },

    // Users
    async fetchUsers(params = {}) {
      this.usersLoading = true
      this.usersError = null
      try {
        const users = await dashboardService.getUsers(params)
        this.users = users
      } catch (error) {
        this.usersError = error.message
        console.error('Error fetching users:', error)
      } finally {
        this.usersLoading = false
      }
    },

    async createUser(userData) {
      try {
        await dashboardService.createUser(userData)
        await this.fetchUsers()
      } catch (error) {
        console.error('Error creating user:', error)
        throw error
      }
    },

    async updateUser(id, userData) {
      try {
        await dashboardService.updateUser(id, userData)
        await this.fetchUsers()
      } catch (error) {
        console.error('Error updating user:', error)
        throw error
      }
    },

    async deleteUser(id) {
      try {
        await dashboardService.deleteUser(id)
        await this.fetchUsers()
      } catch (error) {
        console.error('Error deleting user:', error)
        throw error
      }
    },

    // ID Tags
    async fetchIDTags(params = {}) {
      this.idTagsLoading = true
      this.idTagsError = null
      try {
        const idTags = await dashboardService.getIDTags(params)
        this.idTags = idTags
      } catch (error) {
        this.idTagsError = error.message
        console.error('Error fetching ID tags:', error)
      } finally {
        this.idTagsLoading = false
      }
    },

    async createIDTag(idTagData) {
      try {
        await dashboardService.createIDTag(idTagData)
        await this.fetchIDTags()
      } catch (error) {
        console.error('Error creating ID tag:', error)
        throw error
      }
    },

    async updateIDTag(id, idTagData) {
      try {
        await dashboardService.updateIDTag(id, idTagData)
        await this.fetchIDTags()
      } catch (error) {
        console.error('Error updating ID tag:', error)
        throw error
      }
    },

    async deleteIDTag(id) {
      try {
        await dashboardService.deleteIDTag(id)
        await this.fetchIDTags()
      } catch (error) {
        console.error('Error deleting ID tag:', error)
        throw error
      }
    },

    // Remote Commands
    async sendRemoteCommand(chargePointId, command) {
      try {
        await dashboardService.sendRemoteCommand(chargePointId, command)
        // Refresh charge points data
        await this.fetchChargePoints()
      } catch (error) {
        console.error('Error sending remote command:', error)
        throw error
      }
    }
  }
}) 