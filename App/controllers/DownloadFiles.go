package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
)

// DownloadUnprocessedFile serve o arquivo não processado com base nos parâmetros.
func DownloadUnprocessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileName := c.Param("filename")
	// Diretório base para arquivos não processados (ajuste conforme necessário)
	baseDir := "./uploads/unprocessed"
	filePath := path.Join(baseDir, clientID, projectID, fileName)
	c.File(filePath)
}

// DownloadProcessedFile serve o arquivo processado com base nos parâmetros.
func DownloadProcessedFile(c *gin.Context) {
	clientID := c.Param("clientid")
	projectID := c.Param("projectid")
	fileName := c.Param("filename")
	// Diretório base para arquivos processados (ajuste conforme necessário)
	baseDir := "./uploads/processed"
	filePath := path.Join(baseDir, clientID, projectID, fileName)
	c.File(filePath)
}
