package controllers

import (
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	FilePath string
}

func GetUnprocessedFileInfo(clientID, projectID, fileID string) (FileInfo, error) {
	baseDir := "./uploads/unprocessed"
	filePath := path.Join(baseDir, clientID, projectID, fileID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return FileInfo{}, errors.New("arquivo não encontrado")
	}
	return FileInfo{FilePath: filePath}, nil
}

// DownloadUnprocessedFile serve o arquivo não processado com base nos parâmetros.
func DownloadUnprocessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileID := c.Param("fileid")
	// Realiza consulta para obter informações do arquivo usando clientID, projectID e fileID.
	// Supõe-se que files.GetUnprocessedFileInfo retorne uma struct com FilePath.
	fileInfo, err := GetUnprocessedFileInfo(clientID, projectID, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	// Disponibiliza o arquivo via download streamer usando a localização completa.
	c.File(fileInfo.FilePath)
}

// DownloadProcessedFile serve o arquivo processado com base nos parâmetros.
func DownloadProcessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileID := c.Param("fileid")
	// Diretório base para arquivos processados (ajuste conforme necessário)
	baseDir := "./uploads/processed"
	filePath := path.Join(baseDir, clientID, projectID, fileID)
	// Opcional: verificar se o arquivo existe antes de enviar
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	c.File(filePath)
}
