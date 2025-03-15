package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/NightGhost404/SWAP_AI/App/controllers"
	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

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
			auth.POST("/register", controllers.RegisterUser)
			auth.POST("/login", controllers.LoginUser)
		}

		// Client routes (protected)
		clients := api.Group("/clients")
		clients.Use(authMiddleware())
		{
			clients.GET("", controllers.GetClients)
			clients.POST("", controllers.CreateClient)
			clients.GET("/:id", controllers.GetClientByID)
			clients.PUT("/:id", controllers.UpdateClient)
			clients.DELETE("/:id", controllers.DeleteClient)
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

	// Setup router
	router := SetupRouter()

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Middleware for authentication
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// For testing purposes, we'll treat the token as the user ID directly
		// In a real app, this would involve proper JWT validation
		userID, err := uuid.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
