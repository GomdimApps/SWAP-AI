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
	filePath := c.Param("filename")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Caminho do arquivo é obrigatório"})
		return
	}

	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	fullPath := filepath.Join(uploadDir, filePath)

	if !strings.HasPrefix(filepath.Clean(fullPath), uploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao arquivo solicitado"})
		return
	}

	sendFileToClient(c, fullPath)
}

func SendProcessedFile(c *gin.Context) {
	filePath := c.Param("filename")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Caminho do arquivo é obrigatório"})
		return
	}

	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	adminUploadDir := os.Getenv("ADMIN_UPLOAD_DIR")
	fullPath := filepath.Join(adminUploadDir, filePath)

	if !strings.HasPrefix(filepath.Clean(fullPath), adminUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao arquivo solicitado"})
		return
	}

	sendFileToClient(c, fullPath)
}

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
