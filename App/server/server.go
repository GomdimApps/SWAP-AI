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
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
			clients.GET("/:id/pending-files", controllers.ListPendingClientFiles)
			clients.GET("/:id/processed-files", controllers.ListProcessedClientFiles)
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

		tokenString := tokenParts[1]
		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			secretKey = "mysecret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check that the token is for a user (not a client)
		if claims["type"] != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized token type"})
			c.Abort()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
