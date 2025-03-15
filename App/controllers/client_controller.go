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

// GetClients returns all clients for a user
func GetClients(c *gin.Context) {
	// In a real application, get user ID from JWT token
	// For testing, we'll use a parameter
	userID, err := uuid.Parse(c.Query("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find user's folder
	var folder models.Folder
	if err := database.DB.Where("user_id = ?", userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User folder not found"})
		return
	}

	// Get all clients for this folder
	var clients []models.Client
	if err := database.DB.Where("folder_id = ?", folder.ID).Find(&clients).Error; err != nil {
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

	// In a real application, get user ID from JWT token
	// For testing, we'll use a parameter
	userID, err := uuid.Parse(c.Query("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find user's folder
	var folder models.Folder
	if err := database.DB.Where("user_id = ?", userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User folder not found"})
		return
	}

	// Create client
	client := models.Client{
		FolderID: folder.ID,
		Name:     req.Name,
		Email:    req.Email,
		WhatsApp: req.WhatsApp,
		Status:   req.Status,
	}

	if err := database.DB.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// GetClientByID returns a client by ID
func GetClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
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

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
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

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	if err := database.DB.Delete(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
