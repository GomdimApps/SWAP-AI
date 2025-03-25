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

// GetUnprocessedFileInfo busca o arquivo não processado pelo ID e retorna o caminho completo
func GetUnprocessedFileInfo(clientID, projectID, fileID string) (FileInfo, error) {
	var file models.FileDocument
	if err := database.DB.
		Where("id = ?", fileID).
		First(&file).Error; err != nil {
		return FileInfo{}, errors.New("arquivo não encontrado")
	}

	// O caminho base dos arquivos não processados
	baseDir := "/var/upload/clients/pedding"

	// Monta o caminho completo usando o FilePath armazenado no banco
	finalPath := path.Join(baseDir, file.FilePath)
	return FileInfo{FilePath: finalPath}, nil
}

// GetProcessedFileInfo busca o arquivo processado pelo ID
func GetProcessedFileInfo(clientID, projectID, fileID string) (FileInfo, error) {
	var file models.FileDocument
	if err := database.DB.
		Where("id = ?", fileID).
		First(&file).Error; err != nil {
		return FileInfo{}, errors.New("arquivo não encontrado")
	}

	// O caminho base dos arquivos processados
	baseDir := "/var/upload/clients/processed"

	// Monta o caminho usando o ID do arquivo e a extensão do arquivo original
	finalPath := path.Join(baseDir,
		file.FolderID.String(),
		file.ProjectID.String(),
		file.ID.String()+path.Ext(file.OriginalName))

	return FileInfo{FilePath: finalPath}, nil
}

// DownloadUnprocessedFile serve o arquivo não processado
func DownloadUnprocessedFile(c *gin.Context) {
	// Extrai os parâmetros da URL
	fileID := c.Param("fileid")
	clientID := c.Param("clientid")   // mantido para compatibilidade com a rota
	projectID := c.Param("projectid") // mantido para compatibilidade com a rota

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

// DownloadProcessedFile serve o arquivo processado
func DownloadProcessedFile(c *gin.Context) {
	// Extrai os parâmetros da URL
	fileID := c.Param("fileid")
	clientID := c.Param("clientid")   // mantido para compatibilidade com a rota
	projectID := c.Param("projectid") // mantido para compatibilidade com a rota

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
