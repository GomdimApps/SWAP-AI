package utils

import (
	"net/http"
	"strings"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTokenInfo extrai informações do token do header "Authorization".
// Espera o header no formato "Bearer <token>", onde <token> é "tokenID:tokenType".
func GetTokenInfo(c *gin.Context) (uuid.UUID, string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return uuid.Nil, "", false
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return uuid.Nil, "", false
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	parts := strings.Split(tokenStr, ":")
	if len(parts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return uuid.Nil, "", false
	}
	tokenID, err := uuid.Parse(parts[0])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token ID"})
		return uuid.Nil, "", false
	}
	tokenType := parts[1]
	return tokenID, tokenType, true
}

// GetClientFromToken recupera o client a partir do token presente no contexto.
// Para token do tipo "client", verifica se o parâmetro "id" é igual ao tokenID.
func GetClientFromToken(c *gin.Context) (models.Client, bool) {
	// Recupera informações do token
	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return models.Client{}, false
	}

	var client models.Client
	if tokenType == "client" {
		// Verifica se o parâmetro de rota corresponde ao tokenID
		paramID, err := uuid.Parse(c.Param("id"))
		if err != nil || paramID != tokenID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return models.Client{}, false
		}
		if err := database.DB.Where("id = ?", tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return models.Client{}, false
		}
	} else {
		// Token de admin: utiliza o parâmetro de rota e verifica se o client foi criado pelo admin
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return models.Client{}, false
		}
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return models.Client{}, false
		}
	}

	return client, true
}
