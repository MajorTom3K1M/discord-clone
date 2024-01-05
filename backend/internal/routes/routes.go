package routes

import (
	"discord-backend/internal/app/factory"
	"discord-backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, f *factory.Factory) {
	profileHandler := f.NewProfileHandler()
	authHandler := f.NewAuthHandler()

	AuthRoutes(router, authHandler)

	// AuthMiddleware
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware)

	ProfileRoutes(protected, profileHandler)
}
