package files

import (
	"fmt"
	"log"
	"time"

	"os"
	"path/filepath"

	"github.com/NightGhost404/SWAP_AI/App/database"
	"github.com/NightGhost404/SWAP_AI/App/models"
)

const (
	cleanupInterval = 1 * time.Minute
	maxFileAge      = 5 * 24 * time.Hour
)

func StartCleanupRoutine() {
	log.Println("Iniciando rotina de limpeza automática de arquivos")

	go func() {
		for {
			time.Sleep(cleanupInterval)

			err := cleanupOldFiles()
			if err != nil {
				log.Printf("Erro durante a limpeza de arquivos: %v", err)
			}
		}
	}()
}

func cleanupOldFiles() error {
	cutoffDate := time.Now().Add(-maxFileAge)
	log.Printf("Verificando arquivos não processados criados antes de %v", cutoffDate.Format(time.RFC3339))

	var oldFiles []models.FileDocument

	if err := database.DB.Where("is_processed = false AND created_at < ?", cutoffDate).Find(&oldFiles).Error; err != nil {
		return fmt.Errorf("erro ao buscar arquivos antigos: %w", err)
	}

	if len(oldFiles) == 0 {
		log.Println("Nenhum arquivo antigo encontrado para exclusão")
		return nil
	}

	log.Printf("Encontrados %d arquivos não processados para exclusão", len(oldFiles))

	config := GetDefaultConfig()
	deletedCount := 0

	for _, file := range oldFiles {
		fullPath := filepath.Join(config.BaseDir, file.FilePath)

		if err := os.Remove(fullPath); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Aviso: Não foi possível excluir o arquivo %s: %v", fullPath, err)
			}
		}

		if err := database.DB.Delete(&file).Error; err != nil {
			log.Printf("Aviso: Não foi possível excluir o registro do arquivo %s: %v", file.ID, err)
		} else {
			deletedCount++
		}
	}

	log.Printf("Limpeza concluída: %d/%d arquivos removidos com sucesso", deletedCount, len(oldFiles))
	return nil
}
