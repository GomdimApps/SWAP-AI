package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultBaseDir        = "/var/upload"
	AdminProcessedBaseDir = "/var/upload/processed"
	MaxUploadSizeMB       = 110 // 110MB
)

// FileUploadConfig holds configuration for file uploads
type FileUploadConfig struct {
	BaseDir     string
	MaxSizeMB   int
	AllowedExts []string
}

// GetDefaultConfig returns the default upload configuration for clients
func GetDefaultConfig() *FileUploadConfig {
	return &FileUploadConfig{
		BaseDir:     getEnv("UPLOAD_DIR", DefaultBaseDir),
		MaxSizeMB:   MaxUploadSizeMB,
		AllowedExts: []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"},
	}
}

// GetAdminUploadConfig returns the upload configuration for admin processed files
func GetAdminUploadConfig() *FileUploadConfig {
	return &FileUploadConfig{
		BaseDir:     getEnv("ADMIN_UPLOAD_DIR", AdminProcessedBaseDir),
		MaxSizeMB:   MaxUploadSizeMB,
		AllowedExts: []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"},
	}
}

// SaveUploadedFile saves an uploaded file to disk with a unique name
func SaveUploadedFile(file *multipart.FileHeader, clientID, folderID uuid.UUID, config *FileUploadConfig) (string, string, error) {
	// Create base upload directory if it doesn't exist
	if err := os.MkdirAll(config.BaseDir, 0750); err != nil {
		return "", "", fmt.Errorf("failed to create base upload directory: %w", err)
	}

	// Create folder-specific directory using folder UUID
	folderPath := filepath.Join(config.BaseDir, folderID.String())
	if err := os.MkdirAll(folderPath, 0750); err != nil {
		return "", "", fmt.Errorf("failed to create folder directory: %w", err)
	}

	// Sanitize and make unique filename
	originalName := sanitizeFileName(file.Filename)
	fileExt := filepath.Ext(originalName)

	// Generate unique filename
	uniqueID := uuid.New().String()
	timestamp := time.Now().UnixNano()
	uniqueFileName := fmt.Sprintf("%d-%s%s", timestamp, uniqueID, fileExt)

	// Full path to save the file
	filePath := filepath.Join(folderPath, uniqueFileName)

	// Open source file with buffer
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file with buffer
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Usar buffer maior para cópia mais eficiente (8MB buffer)
	buffer := make([]byte, 8*1024*1024)

	// Copy file contents with progress tracking
	written, err := io.CopyBuffer(dst, src, buffer)
	if err != nil {
		// Em caso de erro, tenta remover o arquivo parcial
		os.Remove(filePath)
		return "", "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Verificar tamanho do arquivo copiado
	if written != file.Size {
		os.Remove(filePath)
		return "", "", fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", file.Size, written)
	}

	// Sync to ensure write is complete
	if err = dst.Sync(); err != nil {
		return "", "", fmt.Errorf("failed to sync file: %w", err)
	}

	// Return the relative storage path and original filename
	storagePath := filepath.Join(folderID.String(), uniqueFileName)
	return storagePath, originalName, nil
}

// ValidateFileExtension checks if the file has an allowed extension
func ValidateFileExtension(filename string, allowedExts []string) bool {
	if len(allowedExts) == 0 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExts {
		if strings.ToLower(allowedExt) == ext {
			return true
		}
	}

	return false
}

// sanitizeFileName removes potentially unsafe characters from filenames
func sanitizeFileName(name string) string {
	return filepath.Base(name)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
