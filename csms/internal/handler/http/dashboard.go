package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/domain"
)

type DashboardHandler struct {
	chargePointService domain.ChargePointService
	transactionService domain.TransactionService
	userService        domain.UserService
}

func NewDashboardHandler(
	chargePointService domain.ChargePointService,
	transactionService domain.TransactionService,
	userService domain.UserService,
) *DashboardHandler {
	return &DashboardHandler{
		chargePointService: chargePointService,
		transactionService: transactionService,
		userService:        userService,
	}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	chargePoints, err := h.chargePointService.ListChargePoints(ctx, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get charge points"})
		return
	}

	onlineCount := 0
	offlineCount := 0
	for _, cp := range chargePoints {
		if cp.Status == "Available" || cp.Status == "Charging" {
			onlineCount++
		} else {
			offlineCount++
		}
	}

	today := time.Now().Truncate(24 * time.Hour)
	transactions, err := h.transactionService.ListTransactions(ctx, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	todayTransactions := 0
	for _, tx := range transactions {
		if tx.StartTime.After(today) {
			todayTransactions++
		}
	}

	stats := gin.H{
		"totalChargePoints":   len(chargePoints),
		"onlineChargePoints":  onlineCount,
		"offlineChargePoints": offlineCount,
		"todayTransactions":   todayTransactions,
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DashboardHandler) GetWeeklyChart(c *gin.Context) {
	ctx := c.Request.Context()

	transactions, err := h.transactionService.ListTransactions(ctx, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	days := []string{"Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"}
	sessions := make([]int, 7)
	energy := make([]float64, 7)

	now := time.Now()
	for i := 0; i < 7; i++ {
		dayStart := now.AddDate(0, 0, -6+i).Truncate(24 * time.Hour)
		dayEnd := dayStart.Add(24 * time.Hour)

		for _, tx := range transactions {
			if tx.StartTime.After(dayStart) && tx.StartTime.Before(dayEnd) {
				sessions[i]++
				energy[i] += tx.EnergyConsumed
			}
		}
	}

	chartData := gin.H{
		"labels": days,
		"datasets": []gin.H{
			{
				"label":       "Sesi",
				"data":        sessions,
				"borderColor": "#1976d2",
				"fill":        false,
			},
			{
				"label":       "Energi (kWh)",
				"data":        energy,
				"borderColor": "#43a047",
				"fill":        false,
			},
		},
	}

	c.JSON(http.StatusOK, chartData)
}
