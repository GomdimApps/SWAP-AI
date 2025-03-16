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

// AdminUploadClientFile handles file uploads for a specific client by an admin
func AdminUploadClientFile(c *gin.Context) {
	// Iniciar log para monitoramento de upload administrativo
	clientIP := c.ClientIP()
	fmt.Printf("Início de upload administrativo - Admin IP: %s\n", clientIP)

	// Get client ID from URL parameter
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	// Get current user ID (admin)
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Verify client exists (without user_id check since admin can access any client)
	var client models.Client
	if err := database.DB.Where("id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Get upload configuration with processed directory for admin uploads
	config := files.GetAdminUploadConfig()

	// Check file size limit
	if c.Request.ContentLength > 0 && c.Request.ContentLength > int64(config.MaxSizeMB)<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":      "File size exceeds limit",
			"maxSizeMB":  config.MaxSizeMB,
			"receivedMB": float64(c.Request.ContentLength) / (1 << 20),
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file upload: %v", err)})
		return
	}

	// Log do tamanho do arquivo
	fmt.Printf("Upload admin em progresso - Arquivo: %s, Tamanho: %.2f MB\n",
		file.Filename, float64(file.Size)/(1<<20))

	// Validate file extension
	if !files.ValidateFileExtension(file.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "File type not allowed",
			"allowedTypes": config.AllowedExts,
		})
		return
	}

	// Save file to processed directory
	storagePath, originalName, err := files.SaveUploadedFile(file, clientID, client.FolderID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}

	// Create admin file document record in database
	fileDoc := models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName,
		FilePath:     storagePath,
		FileSize:     file.Size,
		FileType:     http.DetectContentType([]byte(file.Filename)), // Not accurate, but a placeholder
		ProcessedBy:  userID,                                        // Track which admin processed this file
		IsProcessed:  true,                                          // Mark as processed since admin uploaded it
	}

	if err := database.DB.Create(&fileDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file information"})
		return
	}

	// Log de conclusão
	fmt.Printf("Upload admin concluído - Arquivo: %s, ID: %s\n", originalName, fileDoc.ID)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":      "Admin file uploaded successfully",
		"id":           fileDoc.ID,
		"originalName": fileDoc.OriginalName,
		"size":         fileDoc.FileSize,
		"sizeInMB":     float64(fileDoc.FileSize) / (1 << 20),
		"clientId":     clientID,
		"processedBy":  userID,
	})
}

// GetAdminProcessedFiles returns all processed files for a specific client
func GetAdminProcessedFiles(c *gin.Context) {
	// Get client ID from URL parameter
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	// Verify client exists
	var client models.Client
	if err := database.DB.Where("id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Get all processed files for this client
	var files []models.FileDocument
	if err := database.DB.Where("client_id = ? AND is_processed = true", clientID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client processed files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
