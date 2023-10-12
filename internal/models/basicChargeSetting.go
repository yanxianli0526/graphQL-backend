package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BasicChargeSetting struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
	SortIndex int
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
	// belongs to OrganizationBasicChargeSetting
	OrganizationBasicChargeSettingId uuid.UUID
	OrganizationBasicChargeSetting   OrganizationBasicChargeSetting
	// belongs to Patient
	PatientId uuid.UUID
	Patient   Patient
	// belongs to  User
	UserId uuid.UUID
	User   User
}
