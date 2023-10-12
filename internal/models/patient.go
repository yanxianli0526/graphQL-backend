package models

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID             uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
	FirstName      string
	LastName       string
	IdNumber       string
	PhotoUrl       string
	PhotoXPosition int
	PhotoYPosition int
	ProviderId     string `gorm:"uniqueIndex"`
	Status         string
	Branch         string
	Room           string
	Bed            string
	Sex            string
	Birthday       time.Time
	CheckInDate    time.Time
	PatientNumber  string
	RecordNumber   string
	Numbering      string
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization

	// many to many
	Users []*User `gorm:"many2many:patient_user_relation"`
}
