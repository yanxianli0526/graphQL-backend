package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PayRecord struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	BasicCharge         datatypes.JSON
	Subsidy             datatypes.JSON
	TransferRefundLeave datatypes.JSON
	NonFixedCharge      datatypes.JSON

	PayDate        time.Time
	ReceiptNumber  string
	TaxType        string
	AmountDue      int
	PaidAmount     int
	Note           string
	IsInvalid      bool
	InvalidCaption string
	InvalidDate    time.Time
	PayYear        int
	PayMonth       int

	// belongs to PatientBill
	PatientBillId uuid.UUID
	PatientBill   PatientBill
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
	// belongs to Patient
	PatientId uuid.UUID
	Patient   Patient
	// belongs to  User
	UserId uuid.UUID
	User   User
	// belongs to  User
	CreatedUserId uuid.UUID
	CreatedUser   User
	// belongs to  User
	InvalidUserId *uuid.UUID
	InvalidUser   User

	// many to many
	PayRecordDetails []*PayRecordDetail `gorm:"many2many:payRecord_payRecord_detail"`
}
