package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ChannelRoutes(protected *gin.RouterGroup, channelHandler *handlers.ChannelHandler) {
	channelGroup := protected.Group("/channels")
	{
		channelGroup.GET("/:channelId", channelHandler.GetChannel)

		channelGroup.POST("/servers/:serverId", channelHandler.CreateChannel)
		channelGroup.POST("/:channelId/servers/:serverId", channelHandler.UpdateChannel)

		channelGroup.DELETE("/:channelId/servers/:serverId", channelHandler.DeleteChannel)
	}
}
