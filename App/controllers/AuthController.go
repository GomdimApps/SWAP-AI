package controllers

import (
	"net/http"
	"time"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/NightGhost404/SWAP_AI/App/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	ExpiresAt int64     `json:"expiresAt"`
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token using the util function
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "user",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	tokenString, err := utils.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:     tokenString,
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// LoginUser handles user login with JWT username.
func LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token using the util function
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "user",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	tokenString, err := utils.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:     tokenString,
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}
