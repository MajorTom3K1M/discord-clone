package handlers

import (
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageHandler struct {
	MessageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{MessageService: messageService}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	channelIDStr := c.Query("channelId")
	cursor := c.Query("cursor")
	if channelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing channelId"})
		return
	}

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Server UUID format"})
		return
	}

	messages, nextCursor, err := h.MessageService.GetMessages(channelID, cursor)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Get messages successfully",
		"items":      messages,
		"nextCursor": nextCursor,
	})
}
