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
func (r *mutationResolver) CreateOrganizationNonFixedChargeSetting(ctx context.Context, input *gqlmodels.OrganizationNonFixedChargeSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateOrganizationNonFixedChargeSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	fixedChargeSetting := models.OrganizationNonFixedChargeSetting{
		ItemCategory:   input.ItemCategory,
		ItemName:       input.ItemName,
		Type:           input.Type,
		Unit:           input.Unit,
		Price:          input.Price,
		TaxType:        input.TaxType,
		OrganizationId: organizationId,
	}
	err = orm.CreateOrganizationNonFixedChargeSetting(r.ORM.DB, &fixedChargeSetting)
	if err != nil {
		r.Logger.Error("CreateOrganizationNonFixedChargeSetting orm.CreateOrganizationNonFixedChargeSetting", zap.Error(err), zap.String("fieldName", "createOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createOrganizationNonFixedChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createOrganizationNonFixedChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateOrganizationNonFixedChargeSetting(ctx context.Context, id string, input *gqlmodels.OrganizationNonFixedChargeSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationNonFixedChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationNonFixedChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "updateOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	fixedChargeSetting := models.OrganizationNonFixedChargeSetting{
		ID:           organizationNonFixedChargeSettingId,
		ItemCategory: input.ItemCategory,
		ItemName:     input.ItemName,
		Type:         input.Type,
		Unit:         input.Unit,
		Price:        input.Price,
		TaxType:      input.TaxType,
	}
	err = orm.UpdateOrganizationNonFixedChargeSettingById(r.ORM.DB, &fixedChargeSetting)
	if err != nil {
		r.Logger.Error("UpdateOrganizationNonFixedChargeSetting orm.UpdateOrganizationNonFixedChargeSettingById", zap.Error(err), zap.String("fieldName", "updateOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationNonFixedChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationNonFixedChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeleteOrganizationNonFixedChargeSetting(ctx context.Context, id string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationFixedChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("DeleteOrganizationNonFixedChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "deleteOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	// 除了刪掉機構的設定 (還沒做
	// 還沒做
	// 還沒做
	// 還沒做 （很重要講三次)
	// 要檢查非固定費用模組 有沒有住民用到這個 有的話也一起砍了
	err = orm.DeleteOrganizationNonFixedChargeSettingAndFixedChargeSetting(r.ORM.DB, organizationFixedChargeSettingId)
	if err != nil {
		r.Logger.Error("DeleteOrganizationNonFixedChargeSetting orm.DeleteOrganizationNonFixedChargeSettingAndFixedChargeSetting", zap.Error(err), zap.String("fieldName", "deleteOrganizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteOrganizationNonFixedChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteOrganizationNonFixedChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) OrganizationNonFixedChargeSetting(ctx context.Context, id string) (*models.OrganizationNonFixedChargeSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationNonFixedChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("OrganizationNonFixedChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "organizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	nonFixedChargeSetting, err := orm.GetOrganizationNonFixedChargeSettingById(r.ORM.DB, organizationNonFixedChargeSettingId)
	if err != nil {
		r.Logger.Error("OrganizationNonFixedChargeSetting orm.GetOrganizationNonFixedChargeSettingById", zap.Error(err), zap.String("fieldName", "organizationNonFixedChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationNonFixedChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationNonFixedChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return nonFixedChargeSetting, nil
}

func (r *queryResolver) OrganizationNonFixedChargeSettings(ctx context.Context) ([]*models.OrganizationNonFixedChargeSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationNonFixedChargeSettings uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "organizationNonFixedChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	nonFixedChargeSetting, err := orm.GetOrganizationNonFixedChargeSettings(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("OrganizationNonFixedChargeSettings orm.GetOrganizationNonFixedChargeSettings", zap.Error(err), zap.String("fieldName", "organizationNonFixedChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationNonFixedChargeSettings run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationNonFixedChargeSettings"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return nonFixedChargeSetting, nil
}

// organizationNonFixedChargeSetting resolvers
type organizationNonFixedChargeSettingResolver struct{ *Resolver }

func (r *organizationNonFixedChargeSettingResolver) ID(ctx context.Context, obj *models.OrganizationNonFixedChargeSetting) (string, error) {
	return obj.ID.String(), nil
}
