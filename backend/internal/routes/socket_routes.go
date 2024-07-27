package routes

import (
	"discord-backend/internal/app/handlers"
	"discord-backend/internal/app/websocket"

	"github.com/gin-gonic/gin"
)

func SocketRoutes(router *gin.RouterGroup, socketHandler *handlers.WebsocketHandler) {
	wsHub := websocket.NewHub()
	go wsHub.Run()

	router.GET("/ws", socketHandler.WebSocketHandler(wsHub))
	router.POST("/ws/messages", socketHandler.WebSocketMessageHandler(wsHub))
	router.PATCH("/ws/messages/:messageId", socketHandler.WebScoketEditMessageHandler(wsHub))
	router.DELETE("/ws/messages/:messageId", socketHandler.WebScoketDeleteMessageHandler(wsHub))

	router.GET("/ws/servers/:serverId/participants", socketHandler.WebSocketGetParticipants(wsHub))

	router.POST("/ws/direct-messages", socketHandler.WebSocketDirectMessageHandler(wsHub))
	router.PATCH("/ws/direct-messages/:directMessageId", socketHandler.WebSocketEditDirectMessageHandler(wsHub))
	router.DELETE("/ws/direct-messages/:directMessageId", socketHandler.WebSocketDeleteDirectMessageHandler(wsHub))
}
