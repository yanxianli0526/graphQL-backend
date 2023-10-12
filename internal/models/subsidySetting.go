package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubsidySetting struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
	SortIndex int

	ItemName string `gorm:"not null"`
	Type     string `gorm:"not null"`
	Price    int    `gorm:"not null"`
	Unit     string
	IdNumber string
	Note     string

	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
	// belongs to Patient
	PatientId uuid.UUID
	Patient   Patient
}
