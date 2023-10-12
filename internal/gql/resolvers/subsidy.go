package resolvers

import (
	"context"

	"graphql-go-template/internal/models"
)

// basicCharge resolvers
type subsidyResolver struct{ *Resolver }

func (r *subsidyResolver) ID(ctx context.Context, obj *models.Subsidy) (string, error) {
	return obj.ID.String(), nil
}
