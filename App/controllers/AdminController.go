package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/files"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminUploadClientFile handles file uploads for a specific client by an admin
func AdminUploadClientFile(c *gin.Context) {
	clientIP := c.ClientIP()
	fmt.Printf("Início de upload administrativo - Admin IP: %s\n", clientIP)

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
	if err := database.DB.Where("id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Use a configuração de upload admin
	config := files.GetAdminUploadConfig()
	if c.Request.ContentLength > 0 && c.Request.ContentLength > int64(config.MaxSizeMB)<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":      "File size exceeds limit",
			"maxSizeMB":  config.MaxSizeMB,
			"receivedMB": float64(c.Request.ContentLength) / (1 << 20),
		})
		return
	}

	// Obtém o project_id do formulário, se fornecido; senão, gera um novo
	projectIDStr := c.PostForm("project_id")
	var projectID uuid.UUID
	if projectIDStr != "" {
		pid, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
			return
		}
		projectID = pid
	} else {
		projectID = uuid.New()
	}

	// Validação do tamanho e extensão do arquivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file upload: %v", err)})
		return
	}

	fmt.Printf("Upload admin em progresso - Arquivo: %s, Tamanho: %.2f MB\n",
		file.Filename, float64(file.Size)/(1<<20))

	if !files.ValidateFileExtension(file.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "File type not allowed",
			"allowedTypes": config.AllowedExts,
		})
		return
	}

	// Salva o arquivo admin usando a configuração ADMIN_UPLOAD_DIR
	storagePath, originalName, err := files.SaveUploadedFile(file, clientID, client.FolderID, projectID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}

	// Cria registro do arquivo processado
	fileDoc := models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName,
		FilePath:     storagePath,
		FileSize:     file.Size,
		FileType:     http.DetectContentType([]byte(file.Filename)),
		ProcessedBy:  userID,
		IsProcessed:  true,
		ProjectID:    projectID,
	}
	if err := database.DB.Create(&fileDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file information"})
		return
	}

	// Obtém configuração do diretório padrão (do cliente)
	defaultConfig := files.GetDefaultConfig()

	// Apaga os registros e arquivos não processados pertencentes a esse project_id
	var unprocessedFiles []models.FileDocument
	if err := database.DB.Where("client_id = ? AND project_id = ? AND is_processed = false", clientID, projectID).Find(&unprocessedFiles).Error; err == nil {
		for _, uf := range unprocessedFiles {
			// caminho completo usando a base do cliente
			path := filepath.Join(defaultConfig.BaseDir, uf.FilePath)
			if err := os.Remove(path); err != nil {
				fmt.Printf("Warning: Could not delete file %s: %v\n", path, err)
			}
			// Remove o registro do banco
			database.DB.Delete(&uf)
		}
	}

	// Apaga a pasta do projeto do cliente
	clientProjectFolder := filepath.Join(defaultConfig.BaseDir, client.FolderID.String(), projectID.String())
	if err := os.RemoveAll(clientProjectFolder); err != nil {
		fmt.Printf("Warning: Failed to delete client's project folder %s: %v\n", clientProjectFolder, err)
	}

	fmt.Printf("Upload admin concluído - Arquivo: %s, ID: %s, ProjectID: %s\n", originalName, fileDoc.ID, projectID)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Admin file uploaded successfully",
		"id":           fileDoc.ID,
		"originalName": fileDoc.OriginalName,
		"size":         fileDoc.FileSize,
		"sizeInMB":     float64(fileDoc.FileSize) / (1 << 20),
		"clientId":     clientID,
		"processedBy":  userID,
		"project_id":   projectID,
	})
}

// GetAdminProcessedFiles returns all processed files for a specific client
func GetAdminProcessedFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	var files []models.FileDocument
	if err := database.DB.Where("client_id = ? AND is_processed = true", clientID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client processed files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// GetUnprocessedFiles returns all unprocessed files for a specific client
func GetUnprocessedFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if !verifyClientAccess(c, clientID) {
		return
	}

	var files []models.FileDocument
	if err := database.DB.Where("client_id = ? AND is_processed = false", clientID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client unprocessed files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// DeleteUnprocessedFiles deletes unprocessed files for a specific client
func DeleteUnprocessedFiles(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var targetFileID *uuid.UUID
	if fileIDParam := c.Param("fileID"); fileIDParam != "" {
		parsedID, err := uuid.Parse(fileIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}
		targetFileID = &parsedID
	}

	if !verifyClientAccess(c, clientID) {
		return
	}

	query := database.DB.Where("client_id = ? AND is_processed = false", clientID)
	if targetFileID != nil {
		query = query.Where("id = ?", targetFileID)
	}

	var filesToDelete []models.FileDocument
	if err := query.Find(&filesToDelete).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files for deletion"})
		return
	}

	if len(filesToDelete) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":      "No matching files found to delete",
			"deletedCount": 0,
		})
		return
	}

	config := files.GetDefaultConfig()
	deletedFiles := deleteFiles(filesToDelete, config.BaseDir)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Files deleted successfully",
		"deletedCount": len(deletedFiles),
		"deletedFiles": deletedFiles,
	})
}

func verifyClientAccess(c *gin.Context, clientID uuid.UUID) bool {
	_, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return false
	}

	var client models.Client
	if err := database.DB.Where("id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return false
	}

	return true
}

func deleteFiles(files []models.FileDocument, basePath string) []uuid.UUID {
	var deletedFiles []uuid.UUID

	for _, file := range files {
		fullPath := filepath.Join(basePath, file.FilePath)
		if err := os.Remove(fullPath); err != nil {
			fmt.Printf("Warning: Could not delete file %s: %v\n", fullPath, err)
		}

		if err := database.DB.Delete(&file).Error; err != nil {
			fmt.Printf("Warning: Could not delete file record %s: %v\n", file.ID, err)
		} else {
			deletedFiles = append(deletedFiles, file.ID)
		}
	}

	return deletedFiles
}
