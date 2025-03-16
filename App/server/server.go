package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NightGhost404/SWAP_AI/App/controllers"
	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/files"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	// Aumentar tamanho máximo do corpo da requisição
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() // Use New() em vez de Default() para configuração manual

	// Configurar middleware personalizado para lidar com timeouts longos para uploads grandes
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

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

	// Set maximum upload size
	uploadConfig := files.GetDefaultConfig()
	r.MaxMultipartMemory = int64(uploadConfig.MaxSizeMB) << 20 // MB to bytes

	// Note: A configuração do timeout é feita na função StartServer
	// Os timeouts aqui são comentados para evitar a variável não utilizada
	// Os valores serão: ReadTimeout: 30min, WriteTimeout: 30min, IdleTimeout: 120s

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

			// Client file management routes
			clients.POST("/:id/upload", controllers.UploadClientFile)
			clients.GET("/:id/files", controllers.GetClientFiles)
		}
	}

	// Create upload directory if it doesn't exist
	uploadDir := uploadConfig.BaseDir
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		log.Printf("Warning: Failed to create upload directory %s: %v", uploadDir, err)
	} else {
		log.Printf("Upload directory configured at: %s", uploadDir)
		log.Printf("Maximum upload size set to: %dMB", uploadConfig.MaxSizeMB)
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

	// Configure custom HTTP server with longer timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	// Start server with the custom configuration
	log.Printf("Server starting on port 8080 with extended timeouts for large file uploads")
	if err := server.ListenAndServe(); err != nil {
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

		// Parse the UUID from the token
		userID, err := uuid.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Verify user exists in database
		var count int64
		if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}

		if count == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user token"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
