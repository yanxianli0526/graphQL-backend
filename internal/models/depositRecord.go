package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepositRecord struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
	IdNumber  string `gorm:"not null"`
	Date      time.Time
	Type      string `gorm:"not null"`
	Price     int    `gorm:"not null"`
	Drawee    string
	Note      string
	Invalid   bool `gorm:"not null"`
	// belongs to User
	UserId uuid.UUID
	User   User
	// belongs to Patient
	PatientId uuid.UUID
	Patient   Patient
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}
