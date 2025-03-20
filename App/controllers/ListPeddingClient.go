package controllers

import (
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectFilesGroup struct {
	ProjectID uuid.UUID             `json:"project_id"`
	Files     []models.FileDocument `json:"files"`
}

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

	var projectIDs []uuid.UUID
	if err := database.DB.Model(&models.FileDocument{}).
		Where("client_id = ? AND is_processed = false", clientID).
		Distinct("project_id").
		Pluck("project_id", &projectIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar project IDs de arquivos pendentes"})
		return
	}

	var groups []ProjectFilesGroup
	for _, pid := range projectIDs {
		var files []models.FileDocument
		if err := database.DB.
			Where("client_id = ? AND is_processed = false AND project_id = ?", clientID, pid).
			Find(&files).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar arquivos para o project_id"})
			return
		}
		group := ProjectFilesGroup{
			ProjectID: pid,
			Files:     files,
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, gin.H{"pending_projects": groups})
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

	var projectIDs []uuid.UUID
	if err := database.DB.Model(&models.FileDocument{}).
		Where("client_id = ? AND is_processed = true", clientID).
		Distinct("project_id").
		Pluck("project_id", &projectIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar project IDs de arquivos processados"})
		return
	}

	var groups []ProjectFilesGroup
	for _, pid := range projectIDs {
		var files []models.FileDocument
		if err := database.DB.
			Where("client_id = ? AND is_processed = true AND project_id = ?", clientID, pid).
			Find(&files).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar arquivos para o project_id"})
			return
		}
		group := ProjectFilesGroup{
			ProjectID: pid,
			Files:     files,
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, gin.H{"processed_projects": groups})
}
