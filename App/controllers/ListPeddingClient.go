package controllers

import (
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListPendingClientFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado ou acesso negado"})
		return
	}

	var pendingFiles []models.FileDocument
	if err := database.DB.Where("client_id = ? AND is_processed = false", clientID).Find(&pendingFiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar arquivos pendentes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pending_files": pendingFiles})
}

func ListProcessedClientFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado ou acesso negado"})
		return
	}

	var processedFiles []models.FileDocument
	if err := database.DB.Where("client_id = ? AND is_processed = true", clientID).Find(&processedFiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar arquivos processados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"processed_files": processedFiles})
}
