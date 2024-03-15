package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ChannelRoutes(protected *gin.RouterGroup, channelHandler *handlers.ChannelHandler) {
	channelGroup := protected.Group("/channels")
	{
		channelGroup.POST("/servers/:serverId", channelHandler.CreateChannel)
	}
}
