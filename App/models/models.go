package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Clients      []Client  `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName overrides the table name
func (User) TableName() string {
	return "system.users"
}

// Folder represents a folder for client documents
type Folder struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
}

// TableName overrides the table name
func (Folder) TableName() string {
	return "files.folders"
}

// Client represents a client belonging to a user
type Client struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"not null"`
	FolderID  uuid.UUID `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Status    string    `gorm:"type:varchar(10);check:status IN ('ativo', 'inativo')"`
	Email     string
	WhatsApp  string `gorm:"column:whatsapp"` // Match column name in database
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the table name
func (Client) TableName() string {
	return "clients.clients_swap"
}

// BeforeCreate sets default values before creating a record
func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.Status == "" {
		c.Status = "ativo"
	}
	return nil
}

// FileDocument represents an uploaded document associated with a client
type FileDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID     uuid.UUID `gorm:"not null"`
	FolderID     uuid.UUID `gorm:"not null"`
	OriginalName string    `gorm:"not null"`
	FilePath     string    `gorm:"not null"`
	FileSize     int64     `gorm:"not null"`
	FileType     string
	ProcessedBy  uuid.UUID `gorm:"type:uuid"`     // User ID that processed this file (if admin uploaded)
	IsProcessed  bool      `gorm:"default:false"` // Flag indicating if this is a processed file
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName overrides the table name
func (FileDocument) TableName() string {
	return "files.documents"
}
