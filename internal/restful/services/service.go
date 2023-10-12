package services

import (
	"context"

	orm "graphql-go-template/internal/database"
)

// Service 可以使用 dao
type Service struct {
	ctx context.Context
	db  *orm.GormDatabase
}

// New 會回傳 Service 的 instance
func New(ctx context.Context, db *orm.GormDatabase) Service {
	service := Service{
		ctx: ctx,
		db:  db,
	}
	return service
}
