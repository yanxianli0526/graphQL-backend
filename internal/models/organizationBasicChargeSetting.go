package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationBasicChargeSetting struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
	ItemName  string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Unit      string `gorm:"not null"`
	Price     int    `gorm:"not null"`
	IsTax     string `gorm:"not null"` // 下一次要砍掉
	TaxType   string

	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}
