package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel defines the common columns that all db structs should hold, usually
// db structs based on this have no soft delete
type BaseModel struct {
	// ID should use uuid_generate_v4() for the pk's
	ID        uuid.UUID  `gorm:"primaryKey;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `gorm:"index;not null;default:CURRENT_TIMESTAMP"` // (My|Postgre)SQL
	UpdatedAt *time.Time `gorm:"index"`
}

// BaseModelSoftDelete defines the common columns that all db structs should
// hold, usually. This struct also defines the fields for GORM triggers to
// detect the entity should soft delete
type BaseModelSoftDelete struct {
	BaseModel
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
// func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
// 	id := uuid.New()
// 	if err != nil {
// 		return err
// 	}
// 	return scope.SetColumn("ID", id)
// }
