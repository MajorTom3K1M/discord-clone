package routes

import (
	"discord-backend/internal/app/factory"
	"discord-backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, f *factory.Factory) {
	profileHandler := f.NewProfileHandler()
	authHandler := f.NewAuthHandler()
	serverHandler := f.NewServerHandler()
	memberHandler := f.NewMemberHandler()
	channelHandler := f.NewChannelHandler()
	converstaionHandler := f.NewConversationHandler()
	websocketHandler := f.NewWebsocketHandler()
	messageHandler := f.NewMessageHandler()
	directMessageHandler := f.NewDirectMessageHandler()

	AuthRoutes(router, authHandler)

	// AuthMiddleware
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware)

	SocketRoutes(protected, websocketHandler)
	ProfileRoutes(protected, profileHandler)
	ServerRoutes(protected, serverHandler)
	MemberRoutes(protected, memberHandler)
	ChannelRoutes(protected, channelHandler)
	ConversationRoutes(protected, converstaionHandler)
	MessageRoutes(protected, messageHandler)
	DirectMessageRoutes(protected, directMessageHandler)
}
