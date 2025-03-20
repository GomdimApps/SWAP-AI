package controllers

import (
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
)

func ListAdminUser(c *gin.Context) {
	// Busca todos os usuários de administração
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
