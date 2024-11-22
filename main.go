package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yantology/go-csv-to-postgres-coroutine/config"
	"github.com/yantology/go-csv-to-postgres-coroutine/handlers"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize router
	router := gin.Default()

	// Initialize app config
	config.InitAppConfig()

	// Serve static frontend files
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Initialize services
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Setup routes
	setupRoutes(router, db)

	// Start server
	router.Run(config.PORT)
}

func setupRoutes(router *gin.Engine, db *config.Database) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/upload", handlers.HandleFileUpload(db))
		v1.GET("/health", handlers.HealthCheck)
	}
}
