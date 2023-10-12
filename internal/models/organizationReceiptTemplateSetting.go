package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type OrganizationReceiptTemplateSetting struct {
	ID                  uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt           *time.Time
	DeletedAt           gorm.DeletedAt
	Name                string         `gorm:"not null"`
	TaxTypes            pq.StringArray `gorm:"type:text[]"`
	OrganizationPicture string
	TitleName           string `gorm:"not null"`
	// 住民資料

	PatientInfo pq.StringArray `gorm:"type:text[]"`
	// 金額顯示方式
	PriceShowType string `gorm:"not null"`
	// 機構資料 1
	OrganizationInfoOne pq.StringArray `gorm:"type:text[]"`
	// 機構資料 2
	OrganizationInfoTwo pq.StringArray `gorm:"type:text[]"`
	// 備註固定文字
	NoteText string
	// 白痴印章
	SealOneName      string
	SealOnePicture   string
	SealTwoName      string
	SealTwoPicture   string
	SealThreeName    string
	SealThreePicture string
	SealFourName     string
	SealFourPicture  string
	// 第一聯名稱
	PartOneName string `gorm:"not null"`
	// 第二聯名稱
	PartTwoName string `gorm:"not null"`
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}
