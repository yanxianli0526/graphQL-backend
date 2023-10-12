package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NonFixedChargeRecord struct {
	ID                 uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt          time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          *time.Time
	DeletedAt          gorm.DeletedAt
	NonFixedChargeDate time.Time `gorm:"not null"`
	ItemCategory       string    `gorm:"not null"`
	ItemName           string    `gorm:"not null"`
	Type               string    `gorm:"not null"`
	Unit               string    `gorm:"not null"`
	Price              int       `gorm:"not null"`
	Quantity           int       `gorm:"not null"`
	Subtotal           int       `gorm:"not null"`
	Note               string
	IsTax              string `gorm:"not null"` // 下一次要砍掉
	TaxType            string
	ReceiptStatus      string
	ReceiptDate        time.Time

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
