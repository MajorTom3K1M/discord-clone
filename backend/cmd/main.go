package main

import (
	"discord-backend/internal/app/factory"
	"discord-backend/internal/db"
	"discord-backend/internal/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // Fallback to default if not set
	}

	// Split the origins into a slice
	origins := strings.Split(allowedOrigins, ",")

	// Configure CORS to allow your frontend domain, e.g., http://localhost:3000
	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowCredentials = true // Important for cookies
	r.Use(cors.New(config))

	database, err := db.ConnectToDB()

	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("AutoMigrate failed: ", err)
	}

	appFactory := factory.NewFactory(database)

	routes.SetupRoutes(r, appFactory)

	r.Run()
}
