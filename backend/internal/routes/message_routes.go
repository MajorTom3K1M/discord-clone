package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func MessageRoutes(protected *gin.RouterGroup, messageHandler *handlers.MessageHandler) {
	messageGroup := protected.Group("/messages")
	{
		messageGroup.GET("", messageHandler.GetMessages)
	}
}
