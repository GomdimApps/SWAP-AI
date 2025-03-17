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

	if c.Request.ContentLength > 0 && c.Request.ContentLength > int64(config.MaxSizeMB)<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":      "File size exceeds limit",
			"maxSizeMB":  config.MaxSizeMB,
			"receivedMB": float64(c.Request.ContentLength) / (1 << 20),
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file upload: %v", err)})
		return
	}

	fmt.Printf("Upload em progresso - Arquivo: %s, Tamanho: %.2f MB\n",
		file.Filename, float64(file.Size)/(1<<20))

	if !files.ValidateFileExtension(file.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "File type not allowed",
			"allowedTypes": config.AllowedExts,
		})
		return
	}

	storagePath, originalName, err := files.SaveUploadedFile(file, clientID, client.FolderID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}

	fileDoc := models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName,
		FilePath:     storagePath,
		FileSize:     file.Size,
		FileType:     http.DetectContentType([]byte(file.Filename)),
	}

	if err := database.DB.Create(&fileDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file information"})
		return
	}

	fmt.Printf("Upload concluído - Arquivo: %s, ID: %s\n", originalName, fileDoc.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":      "File uploaded successfully",
		"id":           fileDoc.ID,
		"originalName": fileDoc.OriginalName,
		"size":         fileDoc.FileSize,
		"sizeInMB":     float64(fileDoc.FileSize) / (1 << 20),
		"clientId":     clientID,
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
