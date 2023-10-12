package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PatientBill struct {
	ID        uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt

	AmountReceived          int
	AmountDue               int
	Note                    string
	EditNoteDate            *time.Time
	FixedChargeStartDate    time.Time
	FixedChargeEndDate      time.Time
	TransferRefundStartDate time.Time
	TransferRefundEndDate   time.Time
	NonFixedChargeStartDate time.Time
	NonFixedChargeEndDate   time.Time
	BasicChargesSortIds     pq.StringArray `gorm:"type:uuid[]"`
	SubsidiesSortIds        pq.StringArray `gorm:"type:uuid[]"`

	BillYear  int
	BillMonth int
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
	EditNoteUserId *uuid.UUID
	EditNoteUser   User

	// many to many
	BasicCharges          []*BasicCharge          `gorm:"many2many:patientBill_basic_charges"`
	Subsidies             []*Subsidy              `gorm:"many2many:patientBill_subsidies"`
	TransferRefundLeaves  []*TransferRefundLeave  `gorm:"many2many:patientBill_transfer_refund_leaves"`
	NonFixedChargeRecords []*NonFixedChargeRecord `gorm:"many2many:patientBill_non_fixed_charge_records"`
}
