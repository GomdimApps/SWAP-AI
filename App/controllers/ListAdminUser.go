package controllers

import (
	"net/http"
	"time"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListAdminUser(c *gin.Context) {
	// Busca todos os usuários de administração
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}

	// Remover o campo "Clients" do retorno
	type adminUserResponse struct {
		ID           uuid.UUID `json:"ID"`
		Email        string    `json:"Email"`
		PasswordHash string    `json:"PasswordHash"`
		CreatedAt    time.Time `json:"CreatedAt"`
		UpdatedAt    time.Time `json:"UpdatedAt"`
	}
	var response []adminUserResponse
	for _, user := range users {
		response = append(response, adminUserResponse{
			ID:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"users": response})
}
