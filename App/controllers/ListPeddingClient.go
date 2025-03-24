package controllers

import (
	"fmt"
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileResponse struct {
	models.FileDocument
	DownloadURL string `json:"download_url"`
}

type ProjectFilesGroup struct {
	ProjectID uuid.UUID      `json:"project_id"`
	Files     []FileResponse `json:"files"`
}

func ListPendingClientFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	// Recupera informações do token
	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return
	}
	var client models.Client
	if tokenType == "client" {
		// Usuário cliente: busca somente pelo id usando tokenID
		if err := database.DB.Where("id = ?", tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	} else {
		// Token de admin: utiliza o parâmetro de rota e verifica se o client foi criado pelo admin
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
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

		var fileResponses []FileResponse
		for _, file := range files {
			// Alterado: utilizando file.ID como 'fileid' para download
			downloadLink := fmt.Sprintf("http://192.168.10.3:8080/download/unprocessed/%s/%s/%s", file.ClientID, file.ProjectID, file.ID)
			fileResponses = append(fileResponses, FileResponse{
				FileDocument: file,
				DownloadURL:  downloadLink,
			})
		}

		group := ProjectFilesGroup{
			ProjectID: pid,
			Files:     fileResponses,
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

	// Recupera informações do token
	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return
	}
	var client models.Client
	if tokenType == "client" {
		// Usuário cliente: busca somente pelo id usando tokenID
		if err := database.DB.Where("id = ?", tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	} else {
		// Token de admin: utiliza o parâmetro de rota e verifica se o client foi criado pelo admin
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
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

		var fileResponses []FileResponse
		for _, file := range files {
			// Alterado: utilizando file.ID como 'fileid' para download
			downloadLink := fmt.Sprintf("http://192.168.10.3:8080/download/processed/%s/%s/%s", file.ClientID, file.ProjectID, file.ID)
			fileResponses = append(fileResponses, FileResponse{
				FileDocument: file,
				DownloadURL:  downloadLink,
			})
		}

		group := ProjectFilesGroup{
			ProjectID: pid,
			Files:     fileResponses,
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, gin.H{"processed_projects": groups})
}
