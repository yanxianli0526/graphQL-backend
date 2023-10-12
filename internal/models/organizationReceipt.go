package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationReceipt struct {
	ID                 uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt          time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          *time.Time
	DeletedAt          *time.Time
	FirstText          string
	Year               string
	YearText           string
	Month              string
	MonthText          string
	LastText           string
	IsResetInNextCycle bool

	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}
