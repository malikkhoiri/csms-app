package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/domain"
)

type TransactionHandler struct {
	transactionService domain.TransactionService
}

func NewTransactionHandler(transactionService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
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

	transactions, err := h.transactionService.ListTransactions(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionService.GetTransaction(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
