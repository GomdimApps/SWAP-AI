package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/files"
	"github.com/NightGhost404/SWAP_AI/App/models"
	"github.com/gin-gonic/gin"
)

// FileInfo represents information regarding the file.
type FileInfo struct {
	FilePath string
}

// SendFile envia o arquivo com cabeçalhos para download.
func SendFile(c *gin.Context, file *os.File) {
	fileStat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do arquivo"})
		return
	}
	fileSize := fileStat.Size()
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(file.Name())))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileSize))

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar o arquivo"})
		return
	}
}

// GetUnprocessedFileInfo fetches the unprocessed file record from the DB
// and computes the absolute path by joining with the UPLOAD_DIR.
func GetUnprocessedFileInfo(clientID, projectID, fileID string) (FileInfo, error) {
	var file models.FileDocument
	if err := database.DB.
		Where("id = ? AND client_id = ? AND project_id = ? AND is_processed = false", fileID, clientID, projectID).
		First(&file).Error; err != nil {
		return FileInfo{}, errors.New("arquivo não encontrado")
	}
	finalPath := path.Join(files.GetDefaultConfig().BaseDir, file.FilePath)
	return FileInfo{FilePath: finalPath}, nil
}

// GetProcessedFileInfo fetches the processed file record from the DB
// and computes the absolute path using ADMIN_UPLOAD_DIR, folderID, projectID and fileID with extension.
func GetProcessedFileInfo(clientID, projectID, fileID string) (FileInfo, error) {
	var file models.FileDocument
	if err := database.DB.
		Where("id = ? AND client_id = ? AND project_id = ? AND is_processed = true", fileID, clientID, projectID).
		First(&file).Error; err != nil {
		return FileInfo{}, errors.New("arquivo não encontrado")
	}
	// Form the filename using file.ID and the extension from OriginalName.
	finalPath := path.Join(files.GetAdminUploadConfig().BaseDir,
		file.FolderID.String(),
		file.ProjectID.String(),
		file.ID.String()+path.Ext(file.OriginalName))
	return FileInfo{FilePath: finalPath}, nil
}

// DownloadUnprocessedFile abre o arquivo e envia usando SendFile.
func DownloadUnprocessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileID := c.Param("fileid")
	fileInfo, err := GetUnprocessedFileInfo(clientID, projectID, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	if _, err := os.Stat(fileInfo.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo inexistente no servidor"})
		return
	}
	f, err := os.Open(fileInfo.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo"})
		return
	}
	defer f.Close()
	SendFile(c, f)
}

// DownloadProcessedFile abre o arquivo e envia usando SendFile.
func DownloadProcessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileID := c.Param("fileid")
	fileInfo, err := GetProcessedFileInfo(clientID, projectID, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	if _, err := os.Stat(fileInfo.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo inexistente no servidor"})
		return
	}
	f, err := os.Open(fileInfo.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo"})
		return
	}
	defer f.Close()
	SendFile(c, f)
}
