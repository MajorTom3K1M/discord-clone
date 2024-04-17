package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func DirectMessageRoutes(protected *gin.RouterGroup, directMessageHandler *handlers.DirectMessageHandler) {
	messageGroup := protected.Group("/direct-messages")
	{
		messageGroup.GET("", directMessageHandler.GetDirectMessages)
	}
}
