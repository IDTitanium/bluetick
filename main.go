package main

import (
	"bluetick/database"
	"bluetick/middlewares"
	"bluetick/models"
	"bluetick/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MySQL
	database.ConnectDB()

	// Auto-migrate tables
	database.DB.AutoMigrate(&models.User{}, &models.Message{})

	// Create Gin router
	r := gin.Default()

	// Enable CORS
	middlewares.HandleCors(r)

	// Register routes
	routes.RegisterAuthRoutes(r)
	routes.RegisterMessageRoutes(r)
	routes.RegisterWebSocketRoutes(r)
	routes.RegisterUserStatusRoutes(r)

	// Seed data
	if os.Getenv("SEED_DATA") == "true" {
		database.Seed()
	}

	fmt.Println("Server running on port 8050")
	r.Run(":8050")
}
