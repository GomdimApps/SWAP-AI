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
	// Get client ID from URL parameter
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	// Get current user ID
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Verify client belongs to current user
	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found or access denied"})
		return
	}

	// Get upload configuration
	config := files.GetDefaultConfig()

	// Check file size limit
	if c.Request.ContentLength > int64(config.MaxSizeMB)<<20 { // Convert MB to bytes
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	// Validate file extension
	if !files.ValidateFileExtension(file.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "File type not allowed",
			"allowedTypes": config.AllowedExts,
		})
		return
	}

	// Save file to disk
	storagePath, originalName, err := files.SaveUploadedFile(file, clientID, client.FolderID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}

	// Create file document record in database
	fileDoc := models.FileDocument{
		ClientID:     clientID,
		FolderID:     client.FolderID,
		OriginalName: originalName,
		FilePath:     storagePath,
		FileSize:     file.Size,
		FileType:     http.DetectContentType([]byte(file.Filename)), // Not accurate, but a placeholder
	}

	if err := database.DB.Create(&fileDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file information"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":      "File uploaded successfully",
		"id":           fileDoc.ID,
		"originalName": fileDoc.OriginalName,
		"size":         fileDoc.FileSize,
		"clientId":     clientID,
	})
}

// GetClientFiles returns all files for a specific client
func GetClientFiles(c *gin.Context) {
	// Get client ID from URL parameter
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	// Get current user ID
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Verify client belongs to current user
	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found or access denied"})
		return
	}

	// Get all files for this client
	var files []models.FileDocument
	if err := database.DB.Where("client_id = ?", clientID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
