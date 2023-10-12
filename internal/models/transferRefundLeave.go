package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TransferRefundLeave struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
	Reason        string    `gorm:"not null"`
	IsReserveBed  string    `gorm:"not null"`
	Note          string
	Items         datatypes.JSON
	ReceiptStatus string
	ReceiptDate   time.Time

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
