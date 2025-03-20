package controllers

import (
	"fmt"
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/files"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadClientFile handles file uploads for a specific client
func UploadClientFile(c *gin.Context) {
	clientIP := c.ClientIP()
	fmt.Printf("Início de upload - Cliente IP: %s\n", clientIP)

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found or access denied"})
		return
	}

	config := files.GetDefaultConfig()
	// Verifica se os dois arquivos foram enviados
	file1, err1 := c.FormFile("file1")
	file2, err2 := c.FormFile("file2")
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both file1 and file2 must be uploaded"})
		return
	}

	// Gera um novo projectID para o grupo de arquivos
	projectID := uuid.New()

	// Definir tamanho máximo permitido para cada arquivo
	maxAllowed := int64(config.MaxSizeMB) << 20
	if file1.Size > maxAllowed || file2.Size > maxAllowed {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":     "One or both files exceed the allowed size",
			"maxSizeMB": config.MaxSizeMB,
			"file1Size": float64(file1.Size) / (1 << 20),
			"file2Size": float64(file2.Size) / (1 << 20),
		})
		return
	}

	if !files.ValidateFileExtension(file1.Filename, config.AllowedExts) || !files.ValidateFileExtension(file2.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "File type not allowed",
			"allowedTypes": config.AllowedExts,
		})
		return
	}

	// Salva os dois arquivos usando a nova estrutura
	storagePath1, originalName1, err := files.SaveUploadedFile(file1, clientID, client.FolderID, projectID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file1: %v", err)})
		return
	}
	storagePath2, originalName2, err := files.SaveUploadedFile(file2, clientID, client.FolderID, projectID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file2: %v", err)})
		return
	}

	// Cria registros dos arquivos com o novo campo project_id
	var fileDoc1, fileDoc2 models.FileDocument
	fileDoc1 = models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName1,
		FilePath:     storagePath1,
		FileSize:     file1.Size,
		FileType:     http.DetectContentType([]byte(file1.Filename)),
		ProjectID:    projectID,
	}
	fileDoc2 = models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName2,
		FilePath:     storagePath2,
		FileSize:     file2.Size,
		FileType:     http.DetectContentType([]byte(file2.Filename)),
		ProjectID:    projectID,
	}

	if err := database.DB.Create(&fileDoc1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file1 information"})
		return
	}
	if err := database.DB.Create(&fileDoc2).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file2 information"})
		return
	}

	// Cria registro na nova tabela de projeto usando apenas o projectID e o clientID
	projectRecord := models.Project{
		ID:       projectID, // Usa o projectID gerado anteriormente
		ClientID: clientID,
	}
	if err := database.DB.Create(&projectRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project record"})
		return
	}

	fmt.Printf("Upload concluído - File1: %s, File2: %s, ProjectID: %s\n", originalName1, originalName2, projectID)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Files uploaded successfully",
		"file1_id":   fileDoc1.ID,
		"file2_id":   fileDoc2.ID,
		"project_id": projectID,
		"clientId":   clientID,
	})
}

// GetClientFiles returns all files for a specific client
func GetClientFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found or access denied"})
		return
	}

	var files []models.FileDocument
	if err := database.DB.Where("client_id = ?", clientID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
