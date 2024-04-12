package routes

import (
	"discord-backend/internal/app/handlers"
	"discord-backend/internal/app/websocket"

	"github.com/gin-gonic/gin"
)

func SocketRoutes(router *gin.RouterGroup) {
	wsHub := websocket.NewHub()
	go wsHub.Run()

	router.GET("/ws", handlers.WebSocketHandler(wsHub))
}
