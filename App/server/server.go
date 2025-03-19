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

func SetupRouter() *gin.Engine {
	files.StartCleanupRoutine()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

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

	uploadConfig := files.GetDefaultConfig()
	r.MaxMultipartMemory = int64(uploadConfig.MaxSizeMB) << 20

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.RegisterUser)
			auth.POST("/login", controllers.LoginUser)
			auth.POST("/client/login", controllers.LoginClient)
		}

		clients := api.Group("/clients")
		clients.Use(authMiddleware())
		{
			clients.GET("", controllers.GetClients)
			clients.POST("", controllers.CreateClient)
			clients.GET("/:id", controllers.GetClientByID)
			clients.PUT("/:id", controllers.UpdateClient)
			clients.DELETE("/:id", controllers.DeleteClient)
			clients.POST("/:id/upload", controllers.UploadClientFile)
			clients.GET("/:id/files", controllers.GetClientFiles)
		}

		admin := api.Group("/admin")
		admin.Use(authMiddleware())
		{
			admin.POST("/:id/upload", controllers.AdminUploadClientFile)
			admin.GET("/:id/files", controllers.GetAdminProcessedFiles)
			admin.GET("/:id/unprocessed-files", controllers.GetUnprocessedFiles)
			admin.DELETE("/:id/delete-unprocessed", controllers.DeleteUnprocessedFiles)
			admin.DELETE("/:id/delete-unprocessed/:fileID", controllers.DeleteUnprocessedFiles)
		}
	}

	uploadDirs := []string{
		uploadConfig.BaseDir,
		files.GetAdminUploadConfig().BaseDir,
	}

	for _, dir := range uploadDirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			log.Printf("Warning: Failed to create upload directory %s: %v", dir, err)
		} else {
			log.Printf("Upload directory configured at: %s", dir)
		}
	}

	log.Printf("Maximum upload size set to: %dMB", uploadConfig.MaxSizeMB)

	return r
}

func StartServer() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := SetupRouter()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port 8080 with extended timeouts for large file uploads")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		userID, err := uuid.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

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

		c.Set("userID", userID)
		c.Next()
	}
}
