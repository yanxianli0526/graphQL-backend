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
func (r *mutationResolver) CreateOrganizationBasicChargeSetting(ctx context.Context, input *gqlmodels.OrganizationBasicChargeSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateOrganizationBasicChargeSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	basicChargeSetting := models.OrganizationBasicChargeSetting{
		ItemName:       input.ItemName,
		Type:           input.Type,
		Unit:           input.Unit,
		Price:          input.Price,
		TaxType:        input.TaxType,
		OrganizationId: organizationId,
	}
	err = orm.CreateOrganizationBasicChargeSetting(r.ORM.DB, &basicChargeSetting)
	if err != nil {
		r.Logger.Error("CreateOrganizationBasicChargeSetting orm.CreateOrganizationBasicChargeSetting", zap.Error(err), zap.String("fieldName", "createOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createOrganizationBasicChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createOrganizationBasicChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateOrganizationBasicChargeSetting(ctx context.Context, id string, input *gqlmodels.OrganizationBasicChargeSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationBasicChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationBasicChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "updateOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	basicChargeSetting := models.OrganizationBasicChargeSetting{
		ID:       organizationBasicChargeSettingId,
		ItemName: input.ItemName,
		Type:     input.Type,
		Unit:     input.Unit,
		Price:    input.Price,
		TaxType:  input.TaxType,
	}
	err = orm.UpdateOrganizationBasicChargeSettingById(r.ORM.DB, &basicChargeSetting)
	if err != nil {
		r.Logger.Error("UpdateOrganizationBasicChargeSetting orm.UpdateOrganizationBasicChargeSettingById", zap.Error(err), zap.String("fieldName", "updateOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationBasicChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationBasicChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeleteOrganizationBasicChargeSetting(ctx context.Context, id string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationBasicChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("DeleteOrganizationBasicChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "deleteOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	// 除了刪掉機構的設定 還要檢查固定費用模組 有沒有住民用到這個 有的話也一起砍了
	err = orm.DeleteOrganizationBasicChargeSettingAndBasicChargeSetting(r.ORM.DB, organizationBasicChargeSettingId)
	if err != nil {
		r.Logger.Error("DeleteOrganizationBasicChargeSetting orm.DeleteOrganizationBasicChargeSettingAndBasicChargeSetting", zap.Error(err), zap.String("fieldName", "deleteOrganizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteOrganizationBasicChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteOrganizationBasicChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) OrganizationBasicChargeSetting(ctx context.Context, id string) (*models.OrganizationBasicChargeSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationBasicChargeSettingId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("OrganizationBasicChargeSetting uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "organizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicChargeSetting, err := orm.GetOrganizationBasicChargeSettingById(r.ORM.DB, organizationBasicChargeSettingId)
	if err != nil {
		r.Logger.Error("OrganizationBasicChargeSetting orm.GetOrganizationBasicChargeSettingById", zap.Error(err), zap.String("fieldName", "organizationBasicChargeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationBasicChargeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationBasicChargeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return basicChargeSetting, nil
}

func (r *queryResolver) OrganizationBasicChargeSettings(ctx context.Context) ([]*models.OrganizationBasicChargeSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationBasicChargeSettings uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "organizationBasicChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicChargeSettings, err := orm.GetOrganizationBasicChargeSettings(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("OrganizationBasicChargeSettings orm.GetOrganizationBasicChargeSettings", zap.Error(err), zap.String("fieldName", "organizationBasicChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationBasicChargeSettings run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationBasicChargeSettings"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return basicChargeSettings, nil
}

// organizationBasicChargeSetting resolvers
type organizationBasicChargeSettingResolver struct{ *Resolver }

func (r *organizationBasicChargeSettingResolver) ID(ctx context.Context, obj *models.OrganizationBasicChargeSetting) (string, error) {
	return obj.ID.String(), nil
}
