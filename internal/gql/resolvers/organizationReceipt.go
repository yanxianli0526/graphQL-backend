package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Mutations
func (r *mutationResolver) UpdateOrganizationReceipt(ctx context.Context, input *gqlmodels.OrganizationReceiptInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationReceipt uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updateOrganizationReceipt"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	var year string
	if input.Year != nil {
		year = input.Year.String()
	}

	var month string
	if input.Month != nil {
		month = input.Month.String()
	}

	var isResetInNextCycle bool
	if input.IsResetInNextCycle != nil {
		isResetInNextCycle = *input.IsResetInNextCycle
	}

	organizationReceipt := models.OrganizationReceipt{
		FirstText:          *input.FirstText,
		Year:               year,
		YearText:           *input.YearText,
		Month:              month,
		MonthText:          *input.MonthText,
		LastText:           *input.LastText,
		IsResetInNextCycle: isResetInNextCycle,
		OrganizationId:     organizationId,
	}

	err = orm.UpdateOrganizationReceiptById(r.ORM.DB, &organizationReceipt)
	if err != nil {
		r.Logger.Error("UpdateOrganizationReceipt orm.UpdateOrganizationReceiptById", zap.Error(err), zap.String("fieldName", "updateOrganizationReceipt"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationReceipt run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationReceipt"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) OrganizationReceipt(ctx context.Context) (*models.OrganizationReceipt, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationReceipt uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "organizationReceipt"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	organizationReceipt, err := orm.GetOrganizationReceiptById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("OrganizationReceipt orm.GetOrganizationReceiptById", zap.Error(err), zap.String("fieldName", "organizationReceipt"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationReceipt run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationReceipt"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return organizationReceipt, nil
}

// organizationReceipt resolvers
type organizationReceiptResolver struct{ *Resolver }

func (r *organizationReceiptResolver) ID(ctx context.Context, obj *models.OrganizationReceipt) (string, error) {
	return obj.ID.String(), nil
}

func (r *organizationReceiptResolver) Year(ctx context.Context, obj *models.OrganizationReceipt) (*gqlmodels.YearType, error) {
	year := gqlmodels.YearType(obj.Year)
	isValid := gqlmodels.YearType.IsValid(gqlmodels.YearType(year))
	if !isValid && year != "" {
		r.Logger.Error("OrganizationReceipt Year is inValid", zap.String("fieldName", "organizationReceipt"), zap.Int64("timestamp", time.Now().Unix()))
		return nil, fmt.Errorf("OrganizationReceipt Year is inValid ")
	} else {
		return &year, nil
	}
}

func (r *organizationReceiptResolver) Month(ctx context.Context, obj *models.OrganizationReceipt) (*gqlmodels.MonthType, error) {
	month := gqlmodels.MonthType(obj.Month)
	isValid := gqlmodels.MonthType.IsValid(gqlmodels.MonthType(month))
	if !isValid && month != "" {
		r.Logger.Error("OrganizationReceipt Month is inValid", zap.String("fieldName", "organizationReceipt"), zap.Int64("timestamp", time.Now().Unix()))
		return nil, fmt.Errorf("OrganizationReceipt Month is inValid ")
	} else {
		return &month, nil
	}
}
