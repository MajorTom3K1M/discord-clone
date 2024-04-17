package handlers

import (
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectMessageHandler struct {
	DirectMessageService *services.DirectMessageService
}

func NewDirectMessageHandler(directMessageService *services.DirectMessageService) *DirectMessageHandler {
	return &DirectMessageHandler{DirectMessageService: directMessageService}
}

func (h *DirectMessageHandler) GetDirectMessages(c *gin.Context) {
	conversationIDStr := c.Query("conversationId")
	cursor := c.Query("cursor")
	if conversationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing conversationId"})
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
		return
	}

	directMessages, nextCursor, err := h.DirectMessageService.GetDirectMessages(conversationID, cursor)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Get direct messages successfully",
		"items":      directMessages,
		"nextCursor": nextCursor,
	})
}
