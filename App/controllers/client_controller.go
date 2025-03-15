package controllers

import (
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"`
	WhatsApp string `json:"whatsapp"`
	Status   string `json:"status"`
}

// GetUserID extracts user ID from context (set by auth middleware)
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	// In a real app, this would be extracted from the JWT token
	// For testing, we simulate this by using a context value set by authMiddleware
	userIDValue, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := userIDValue.(uuid.UUID)
	return userID, ok
}

// GetClients returns all clients for the authenticated user
func GetClients(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get all clients for this user
	var clients []models.Client
	if err := database.DB.Where("user_id = ?", userID).Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

// CreateClient creates a new client
func CreateClient(c *gin.Context) {
	var req ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	tx := database.DB.Begin()

	// Create a new folder for the client
	folder := models.Folder{}
	if err := tx.Create(&folder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	// Create client
	client := models.Client{
		UserID:   userID,
		FolderID: folder.ID,
		Name:     req.Name,
		Email:    req.Email,
		WhatsApp: req.WhatsApp,
		Status:   req.Status,
	}

	if err := tx.Create(&client).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, client)
}

// GetClientByID returns a client by ID
func GetClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
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
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// UpdateClient updates a client
func UpdateClient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var req ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Update fields
	client.Name = req.Name
	client.Email = req.Email
	client.WhatsApp = req.WhatsApp
	if req.Status != "" {
		client.Status = req.Status
	}

	if err := database.DB.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// DeleteClient deletes a client
func DeleteClient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
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
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Get folder ID before deleting the client
	folderID := client.FolderID

	tx := database.DB.Begin()

	// Delete client
	if err := tx.Delete(&client).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	// Delete associated folder
	if err := tx.Delete(&models.Folder{}, folderID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
