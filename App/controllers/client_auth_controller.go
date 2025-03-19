package controllers

import (
	"net/http"
	"time"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ClientLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ClientAuthResponse struct {
	Token     string    `json:"token"`
	ClientID  uuid.UUID `json:"clientId"`
	Email     string    `json:"email"`
	ExpiresAt int64     `json:"expiresAt"`
}

func LoginClient(c *gin.Context) {
	var req ClientLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Client
	if err := database.DB.Where("email = ?", req.Email).First(&client).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Validar a senha
	if err := bcrypt.CompareHashAndPassword([]byte(client.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Gerar token (aqui, usamos o UUID do cliente como token)
	token := client.ID.String()

	c.JSON(http.StatusOK, ClientAuthResponse{
		Token:     token,
		ClientID:  client.ID,
		Email:     client.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}
