package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ConversationRoutes(protected *gin.RouterGroup, conversationHandler *handlers.ConversationHandler) {
	conversationsGroup := protected.Group("/conversations")
	{
		conversationsGroup.GET("/between/:memberOneId/:memberTwoId", conversationHandler.FindConversation)

		conversationsGroup.POST("/ensure", conversationHandler.GetOrCreateConversation)
	}
}
