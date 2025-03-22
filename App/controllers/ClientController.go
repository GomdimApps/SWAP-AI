package controllers

import (
	"net/http"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ClientRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	WhatsApp string `json:"whatsapp" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Status   string `json:"status"`
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	userID, ok := userIDValue.(uuid.UUID)
	return userID, ok
}

// Nova função auxiliar para obter ID e token_type do contexto.
func GetTokenInfo(c *gin.Context) (uuid.UUID, string, bool) {
	idVal, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, "", false
	}
	tokenTypeVal, exists := c.Get("token_type")
	if !exists {
		return uuid.Nil, "", false
	}
	id, ok := idVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, "", false
	}
	tokenType, ok := tokenTypeVal.(string)
	if !ok {
		return uuid.Nil, "", false
	}
	return id, tokenType, true
}

func GetClients(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var clients []models.Client
	if err := database.DB.Where("user_id = ?", userID).Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

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

	// O campo UserID do client armazena o ID do admin que o criou.
	tx := database.DB.Begin()

	folder := models.Folder{}
	if err := tx.Create(&folder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	client := models.Client{
		UserID:       userID,
		FolderID:     folder.ID,
		Name:         req.Name,
		Email:        req.Email,
		WhatsApp:     req.WhatsApp,
		PasswordHash: string(hashedPassword),
		Status:       req.Status,
	}

	if err := tx.Create(&client).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, client)
}

func GetClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return
	}

	var client models.Client
	if tokenType == "client" {
		// O cliente só pode acessar a si mesmo.
		if id != tokenID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		if err := database.DB.Where("id = ?", id).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	} else {
		// Usuário admin: valida que o client foi criado por ele.
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	}

	c.JSON(http.StatusOK, client)
}

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

	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return
	}

	var client models.Client
	if tokenType == "client" {
		if id != tokenID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		if err := database.DB.Where("id = ?", id).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	} else {
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	}

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

func DeleteClient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	tokenID, tokenType, ok := GetTokenInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token info not found"})
		return
	}

	var client models.Client
	if tokenType == "client" {
		if id != tokenID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		if err := database.DB.Where("id = ?", id).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	} else {
		if err := database.DB.Where("id = ? AND user_id = ?", id, tokenID).First(&client).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
	}

	folderID := client.FolderID
	tx := database.DB.Begin()

	if err := tx.Delete(&client).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	if err := tx.Delete(&models.Folder{}, folderID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
