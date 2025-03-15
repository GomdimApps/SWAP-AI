package server

import (
	"log"
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API Routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", registerUser)
			auth.POST("/login", loginUser)
		}

		// Client routes (protected)
		clients := api.Group("/clients")
		clients.Use(authMiddleware())
		{
			clients.GET("", getClients)
			clients.POST("", createClient)
			clients.GET("/:id", getClientByID)
			clients.PUT("/:id", updateClient)
			clients.DELETE("/:id", deleteClient)
		}
	}

	return r
}

// StartServer configures and starts the HTTP server
func StartServer() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate schemas
	if err := database.DB.AutoMigrate(&models.User{}, &models.Folder{}, &models.Client{}); err != nil {
		log.Fatalf("Failed to migrate database schemas: %v", err)
	}

	// Setup router
	router := SetupRouter()

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Placeholder route handlers
func registerUser(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func loginUser(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func getClients(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func createClient(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func getClientByID(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func updateClient(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func deleteClient(c *gin.Context) {
	// Implementation will be added later
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT authentication will be implemented later
		c.Next()
	}
}
