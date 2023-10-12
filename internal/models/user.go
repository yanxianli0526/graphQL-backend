package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User defines a user for the app
type User struct {
	ID             uuid.UUID `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
	FirstName      string
	LastName       string
	DisplayName    string
	IdNumber       string
	Preference     datatypes.JSON
	Token          string
	TokenExpiredAt *time.Time
	ProviderToken  datatypes.JSON
	ProviderId     string `gorm:"uniqueIndex"`
	Username       string
	Password       string
	// belongs to Organization
	OrganizationId uuid.UUID
	Organization   Organization
}

func MarshalTransferRefundItemsJsonType(a datatypes.JSON) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		data, _ := json.Marshal(a)
		io.WriteString(w, string(data))
	})
}

func UnmarshalTransferRefundItemsJsonType(v interface{}) (datatypes.JSON, error) {
	a, ok := v.(datatypes.JSON)
	if !ok {
		return nil, errors.New("failed to cast to datatypes.JSON")
	}
	return a, nil
}
