package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	FilePath string
}

func SendPendingFile(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome do arquivo é obrigatório"})
		return
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	filePath := filepath.Join(uploadDir, fileName)

	if !strings.HasPrefix(filepath.Clean(filePath), uploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao arquivo solicitado"})
		return
	}

	sendFileToClient(c, filePath)
}

func SendProcessedFile(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome do arquivo é obrigatório"})
		return
	}

	adminUploadDir := os.Getenv("ADMIN_UPLOAD_DIR")
	filePath := filepath.Join(adminUploadDir, fileName)

	if !strings.HasPrefix(filepath.Clean(filePath), adminUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao arquivo solicitado"})
		return
	}

	sendFileToClient(c, filePath)
}

// sendFileToClient envia o arquivo ao cliente para download
func sendFileToClient(c *gin.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo"})
			fmt.Println("Erro ao abrir arquivo:", err)
		}
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do arquivo"})
		fmt.Println("Erro ao obter informações do arquivo:", err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar o arquivo"})
		fmt.Println("Erro ao enviar o arquivo:", err)
	}
}
