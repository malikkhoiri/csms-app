package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/domain"
)

type ChargePointHandler struct {
	chargePointService domain.ChargePointService
}

func NewChargePointHandler(chargePointService domain.ChargePointService) *ChargePointHandler {
	return &ChargePointHandler{
		chargePointService: chargePointService,
	}
}

func (h *ChargePointHandler) GetChargePoints(c *gin.Context) {
	ctx := c.Request.Context()

	limit := 100
	offset := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	chargePoints, err := h.chargePointService.ListChargePoints(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get charge points"})
		return
	}

	c.JSON(http.StatusOK, chargePoints)
}

func (h *ChargePointHandler) GetChargePoint(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid charge point ID"})
		return
	}

	chargePoint, err := h.chargePointService.GetChargePoint(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Charge point not found"})
		return
	}

	c.JSON(http.StatusOK, chargePoint)
}

func (h *ChargePointHandler) UpdateChargePointStatus(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid charge point ID"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.chargePointService.UpdateChargePointStatus(ctx, uint(id), request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *ChargePointHandler) SendRemoteCommand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid charge point ID"})
		return
	}

	var request struct {
		Command string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// TODO: Implement remote command logic
	c.JSON(http.StatusOK, gin.H{
		"message":       "Command sent successfully",
		"command":       request.Command,
		"chargePointId": id,
	})
}
