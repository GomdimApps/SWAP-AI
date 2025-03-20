package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Clients      []Client  `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// Novo campo para a data de expiração
	ExpiresAt time.Time `gorm:"not null"`
	// Status do usuário: "ativo" ou "inativo"
	Status string `gorm:"default:ativo"`
}

func (User) TableName() string {
	return "system.users"
}

type Folder struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
}

func (Folder) TableName() string {
	return "files.folders"
}

type Client struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"not null"`
	FolderID     uuid.UUID `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Status       string    `gorm:"type:varchar(10);check:status IN ('ativo', 'inativo')"`
	Email        string
	WhatsApp     string `gorm:"column:whatsapp"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Client) TableName() string {
	return "clients.clients_swap"
}

// BeforeCreate sets the default status before record creation.
func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.Status == "" {
		c.Status = "ativo"
	}
	return nil
}

type FileDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID     uuid.UUID `gorm:"not null"`
	FolderID     uuid.UUID `gorm:"not null"`
	OriginalName string    `gorm:"not null"`
	FilePath     string    `gorm:"not null"`
	FileSize     int64     `gorm:"not null"`
	FileType     string
	ProcessedBy  uuid.UUID `gorm:"type:uuid"`
	IsProcessed  bool      `gorm:"default:false"`
	ProjectID    uuid.UUID // Novo campo para associar o projeto
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (FileDocument) TableName() string {
	return "files.documents"
}

// Nova struct para o Projeto (ajustada)
type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID  uuid.UUID `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Project) TableName() string {
	return "files.project"
}
