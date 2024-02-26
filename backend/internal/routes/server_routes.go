package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ServerRoutes(protected *gin.RouterGroup, serverHandler *handlers.ServerHandler) {
	serverGroup := protected.Group("/server")
	{
		serverGroup.POST("", serverHandler.CreateServer)
	}
}
