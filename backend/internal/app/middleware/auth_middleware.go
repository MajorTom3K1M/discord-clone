package middleware

import (
	"discord-backend/internal/app/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	splitToken := strings.Split(tokenString, "Bearer ")
	if len(splitToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		c.Abort()
		return
	}

	tokenString = splitToken[1]

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
