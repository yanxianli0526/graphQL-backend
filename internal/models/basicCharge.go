package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BasicCharge struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	ItemName      string
	Type          string
	Unit          string
	Price         int
	TaxType       string
	StartDate     time.Time
	EndDate       time.Time
	Note          string
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
