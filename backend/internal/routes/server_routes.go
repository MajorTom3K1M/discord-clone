package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ServerRoutes(protected *gin.RouterGroup, serverHandler *handlers.ServerHandler) {
	serversGroup := protected.Group("/servers")
	{
		serversGroup.GET("", serverHandler.GetServers)
		serversGroup.GET("/by-profile", serverHandler.GetServerByProfileID)
		serversGroup.GET("/invite-code/:inviteCode", serverHandler.GetServerByInviteCode)
		serversGroup.GET("/:serverId", serverHandler.GetServer)
		serversGroup.GET("/:serverId/details", serverHandler.GetServerDetails)

		serversGroup.POST("", serverHandler.CreateServer)
		serversGroup.POST("/invite-code/:inviteCode/members", serverHandler.UpdateServerMember)

		serversGroup.PATCH("/:serverId", serverHandler.UpdateServer)
		serversGroup.PATCH("/:serverId/leave", serverHandler.LeaveServer)
		serversGroup.PATCH("/:serverId/invite-code", serverHandler.UpdateServerInviteCode)
	}
}
