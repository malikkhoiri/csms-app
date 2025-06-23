package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malikkhoiri/csms/internal/domain"
)

type IDTagHandler struct {
	idTagService domain.IDTagService
}

func NewIDTagHandler(idTagService domain.IDTagService) *IDTagHandler {
	return &IDTagHandler{
		idTagService: idTagService,
	}
}

func (h *IDTagHandler) GetIDTags(c *gin.Context) {
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

	idTags, err := h.idTagService.ListIDTags(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID tags"})
		return
	}

	c.JSON(http.StatusOK, idTags)
}

func (h *IDTagHandler) GetIDTagsByUser(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	idTags, err := h.idTagService.ListByUser(ctx, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID tags for user"})
		return
	}

	c.JSON(http.StatusOK, idTags)
}

func (h *IDTagHandler) GetIDTag(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID tag ID"})
		return
	}

	idTag, err := h.idTagService.GetIDTag(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID tag not found"})
		return
	}

	c.JSON(http.StatusOK, idTag)
}

func (h *IDTagHandler) CreateIDTag(c *gin.Context) {
	ctx := c.Request.Context()

	var idTag domain.IDTag
	if err := c.ShouldBindJSON(&idTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "msg": err.Error()})
		return
	}

	err := h.idTagService.CreateIDTag(ctx, &idTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ID tag"})
		return
	}

	c.JSON(http.StatusCreated, idTag)
}

func (h *IDTagHandler) UpdateIDTag(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID tag ID"})
		return
	}

	var idTag domain.IDTag
	if err := c.ShouldBindJSON(&idTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	idTag.ID = uint(id)
	err = h.idTagService.UpdateIDTag(ctx, &idTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ID tag"})
		return
	}

	c.JSON(http.StatusOK, idTag)
}

func (h *IDTagHandler) DeleteIDTag(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID tag ID"})
		return
	}

	err = h.idTagService.DeleteIDTag(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ID tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ID tag deleted successfully"})
}
