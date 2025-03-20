package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	// Corrigindo o caminho de importação para ferramentas locais

	"github.com/NightGhost404/SWAP_AI/App/utils"
	"github.com/google/uuid"
)

const (
	DefaultBaseDir         = "/var/upload/clients/pedding"
	AdminProcessedBaseDir  = "/var/upload/clients/processed"
	DefaultMaxUploadSizeMB = 110 // 110MB
)

// FileUploadConfig holds configuration for file uploads
type FileUploadConfig struct {
	BaseDir     string
	MaxSizeMB   int
	AllowedExts []string
}

func getDefaultAllowedExtensions() []string {
	return []string{".mov", ".mkv", ".mp4", ".jpg", ".jpeg", ".png"}
}

func createUploadConfig(dirEnvVar, defaultDir string) *FileUploadConfig {
	return &FileUploadConfig{
		BaseDir:     getEnv(dirEnvVar, defaultDir),
		MaxSizeMB:   utils.GetEnvAsInt("MAX_UPLOAD_SIZE_MB", DefaultMaxUploadSizeMB),
		AllowedExts: getDefaultAllowedExtensions(),
	}
}

func GetDefaultConfig() *FileUploadConfig {
	return createUploadConfig("UPLOAD_DIR", DefaultBaseDir)
}

func GetAdminUploadConfig() *FileUploadConfig {
	return createUploadConfig("ADMIN_UPLOAD_DIR", AdminProcessedBaseDir)
}

// SaveUploadedFile saves an uploaded file to disk with a unique name
func SaveUploadedFile(file *multipart.FileHeader, clientID, folderID, projectID uuid.UUID, config *FileUploadConfig) (string, string, error) {
	// Cria a pasta base e a subpasta do projeto
	if err := os.MkdirAll(config.BaseDir, 0750); err != nil {
		return "", "", fmt.Errorf("failed to create base upload directory: %w", err)
	}
	// Novo caminho: folderID/projectID/files
	projectPath := filepath.Join(config.BaseDir, folderID.String(), projectID.String(), "files")
	if err := os.MkdirAll(projectPath, 0750); err != nil {
		return "", "", fmt.Errorf("failed to create project directory: %w", err)
	}

	originalName := sanitizeFileName(file.Filename)
	fileExt := filepath.Ext(originalName)

	uniqueID := uuid.New().String()
	timestamp := time.Now().UnixNano()
	uniqueFileName := fmt.Sprintf("%d-%s%s", timestamp, uniqueID, fileExt)

	filePath := filepath.Join(projectPath, uniqueFileName)

	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	buffer := make([]byte, 8*1024*1024)

	written, err := io.CopyBuffer(dst, src, buffer)
	if err != nil {
		os.Remove(filePath)
		return "", "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	if written != file.Size {
		os.Remove(filePath)
		return "", "", fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", file.Size, written)
	}

	if err = dst.Sync(); err != nil {
		return "", "", fmt.Errorf("failed to sync file: %w", err)
	}

	storagePath := filepath.Join(folderID.String(), projectID.String(), "files", uniqueFileName)
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

func sanitizeFileName(name string) string {
	return filepath.Base(name)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
