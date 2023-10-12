package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AutoTextField struct {
	ID         uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time
	DeletedAt  gorm.DeletedAt
	ModuleName string `gorm:"not null"`
	ItemName   string `gorm:"not null"`
	Text       string `gorm:"not null"`

	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}
