package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Organization struct {
	ID                       uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt                time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
	Name                     string  `gorm:"not null"`
	AddressCity              *string // 地址(縣市)
	AddressDistrict          *string // 地址(區鄉鎮)
	Address                  *string // 地址
	Phone                    *string // 電話
	Fax                      *string // 傳真
	Owner                    *string // 負責人
	Email                    *string // 機構信箱
	TaxIdNumber              *string // 統一編號
	RemittanceBank           *string // 匯款銀行
	RemittanceIdNumber       *string // 匯款帳號
	RemittanceUserName       *string // 匯款戶名
	EstablishmentNumber      *string // 設立許可文號
	Solution                 string
	FixedChargeStartMonth    string `gorm:"not null"`
	FixedChargeStartDate     int    `gorm:"not null"`
	FixedChargeEndMonth      string `gorm:"not null"`
	FixedChargeEndDate       int    `gorm:"not null"`
	NonFixedChargeStartMonth string `gorm:"not null"`
	NonFixedChargeStartDate  int    `gorm:"not null"`
	NonFixedChargeEndMonth   string `gorm:"not null"`
	NonFixedChargeEndDate    int    `gorm:"not null"`
	TransferRefundStartMonth string `gorm:"not null"`
	TransferRefundStartDate  int    `gorm:"not null"`
	TransferRefundEndMonth   string `gorm:"not null"`
	TransferRefundEndDate    int    `gorm:"not null"`

	Branchs       pq.StringArray `gorm:"type:text[]"`
	ProviderOrgId string         // NIS的Id
	TestTime      time.Time
	Privacy       string
}

func MarshalPgStringArray(a pq.StringArray) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		data, _ := json.Marshal(a)
		io.WriteString(w, string(data))
	})
}

func UnmarshalPgStringArray(v interface{}) (pq.StringArray, error) {
	a, ok := v.(pq.StringArray)
	if !ok {
		return nil, errors.New("failed to cast to pq.StringArray")
	}
	return a, nil
}
