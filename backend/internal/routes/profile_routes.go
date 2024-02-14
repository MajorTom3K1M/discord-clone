package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(protected *gin.RouterGroup, profileHandler *handlers.ProfileHandler) {
	profileGroup := protected.Group("/profile")
	{
		profileGroup.GET("/auth/me", profileHandler.GetMyProfile)
		profileGroup.GET("/:id", profileHandler.GetProfile)
		profileGroup.POST("/:id", profileHandler.UpdateProfile)
	}
}
