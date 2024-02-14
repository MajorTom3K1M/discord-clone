package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	router.POST("/signup", authHandler.SignUp)
	router.POST("/signin", authHandler.SignIn)
	router.POST("/signout", authHandler.SignOut)
	router.GET("/refresh", authHandler.Refresh)
}
