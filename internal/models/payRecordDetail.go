package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 未來要改回來
type PayRecordDetail struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	RecordDate time.Time
	Type       string `gorm:"not null"`
	Price      int    `gorm:"not null"`
	Method     string `gorm:"not null"`
	Payer      string
	Handler    string
	Note       string

	// belongs to PayRecord
	PayRecordId uuid.UUID
	PayRecord   PayRecord
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
