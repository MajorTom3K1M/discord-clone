package middleware

import (
	"discord-backend/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		c.Abort()
		return
	}

	claims, err := services.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	c.Set("profile_id", claims["profile_id"])
	c.Set("name", claims["name"])

	c.Next()
}
