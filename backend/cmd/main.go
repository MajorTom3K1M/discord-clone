package main

import (
	"discord-backend/internal/app/factory"
	"discord-backend/internal/db"
	"discord-backend/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configure CORS to allow your frontend domain, e.g., http://localhost:3000
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000"}
	// config.AllowCredentials = true // Important for cookies
	// r.Use(cors.New(config))

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
