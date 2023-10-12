package resolvers

import (
	"context"

	"graphql-go-template/internal/models"
)

// basicCharge resolvers
type basicChargeResolver struct{ *Resolver }

func (r *basicChargeResolver) ID(ctx context.Context, obj *models.BasicCharge) (string, error) {
	return obj.ID.String(), nil
}
