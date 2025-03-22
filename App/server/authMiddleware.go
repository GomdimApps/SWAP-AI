package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

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

		// Debug: log dos claims para identificar campos presentes
		fmt.Printf("DEBUG: JWT claims: %+v\n", claims)

		tokenType, ok := claims["type"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}
		// Armazena token_type no contexto para uso posterior
		c.Set("token_type", tokenType)

		// Se o endpoint for de clientes, permite token do tipo "client" ou "user"
		if strings.HasPrefix(c.FullPath(), "/api/v1/clients") {
			if tokenType == "client" {
				clientIDStr, ok := claims["client_id"].(string)
				if !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client token payload"})
					c.Abort()
					return
				}
				clientID, err := uuid.Parse(clientIDStr)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client ID in token"})
					c.Abort()
					return
				}
				c.Set("userID", clientID)
			} else if tokenType == "user" {
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
				var user models.User
				if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
					c.Abort()
					return
				}
				if time.Now().After(user.ExpiresAt) || user.Status != "ativo" {
					c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive or expired"})
					c.Abort()
					return
				}
				c.Set("userID", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized token type for clients"})
				c.Abort()
				return
			}
		} else {
			// Para endpoints não pertencentes ao grupo clients, somente tokens do tipo "user" são aceitos
			if tokenType != "user" {
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
			var user models.User
			if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
			if time.Now().After(user.ExpiresAt) || user.Status != "ativo" {
				c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive or expired"})
				c.Abort()
				return
			}
			c.Set("userID", userID)
		}

		c.Next()
	}
}
