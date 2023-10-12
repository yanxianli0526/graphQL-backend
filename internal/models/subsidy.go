package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 未來要改回來
type Subsidy struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	ItemName      string `gorm:"not null"`
	Type          string `gorm:"not null"`
	Price         int    `gorm:"not null"`
	Unit          string
	IdNumber      string
	Note          string
	StartDate     time.Time
	EndDate       time.Time
	ReceiptStatus string
	ReceiptDate   time.Time
	SortIndex     int

	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
	// belongs to Patient
	PatientId uuid.UUID
	Patient   Patient
	// belongs to  User
	UserId uuid.UUID
	User   User
}
