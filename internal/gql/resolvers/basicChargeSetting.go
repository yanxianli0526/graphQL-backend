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
func (r *mutationResolver) CreateBasicCharge(ctx context.Context, input *gqlmodels.BasicChargeInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreateBasicCharge uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "createBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	if err != nil {
		r.Logger.Warn("CreateBasicCharge uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "createBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	patientId, err := uuid.Parse(input.PatientID)
	if err != nil {
		r.Logger.Warn("CreateBasicCharge uuid.Parse(input.PatientID)", zap.Error(err), zap.String("fieldName", "createBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var basicChargeSettingsUUID []uuid.UUID
	for i := range input.OrganizationBasicChargeID {
		basicChargeSettingUUID, err := uuid.Parse(input.OrganizationBasicChargeID[i])
		if err != nil {
			r.Logger.Warn("CreateBasicCharge uuid.Parse(input.OrganizationBasicChargeID[i])", zap.Error(err), zap.String("fieldName", "createBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
		basicChargeSettingsUUID = append(basicChargeSettingsUUID, basicChargeSettingUUID)
	}

	if len(input.OrganizationBasicChargeID) > 0 {
		var basicChargeSettings []*models.BasicChargeSetting

		for i := range input.OrganizationBasicChargeID {
			basicChargeSetting := models.BasicChargeSetting{
				OrganizationId:                   organizationId,
				OrganizationBasicChargeSettingId: basicChargeSettingsUUID[i],
				PatientId:                        patientId,
				UserId:                           userId,
				SortIndex:                        i,
			}
			basicChargeSettings = append(basicChargeSettings, &basicChargeSetting)
		}

		err = orm.CreateBasicChargeSettings(r.ORM.DB, basicChargeSettings)
		if err != nil {
			r.Logger.Error("CreateBasicCharge orm.CreateBasicChargeSettings", zap.Error(err), zap.String("fieldName", "createBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateBasicCharge(ctx context.Context, input *gqlmodels.BasicChargeInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdateBasicCharge uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updateBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	if err != nil {
		r.Logger.Warn("UpdateBasicCharge uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "updateBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientId, err := uuid.Parse(input.PatientID)
	if err != nil {
		r.Logger.Warn("UpdateBasicCharge uuid.Parse(input.PatientID)", zap.Error(err), zap.String("fieldName", "updateBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var basicChargeSettingsUUID []uuid.UUID
	for i := range input.OrganizationBasicChargeID {
		basicChargeSettingUUID, err := uuid.Parse(input.OrganizationBasicChargeID[i])
		if err != nil {
			r.Logger.Warn("UpdateBasicCharge uuid.Parse(input.OrganizationBasicChargeID[i])", zap.Error(err), zap.String("fieldName", "updateBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
		basicChargeSettingsUUID = append(basicChargeSettingsUUID, basicChargeSettingUUID)
	}

	var basicChargeSettings []*models.BasicChargeSetting

	for i := range input.OrganizationBasicChargeID {
		basicChargeSetting := models.BasicChargeSetting{
			OrganizationId:                   organizationId,
			OrganizationBasicChargeSettingId: basicChargeSettingsUUID[i],
			PatientId:                        patientId,
			UserId:                           userId,
			SortIndex:                        i,
		}
		basicChargeSettings = append(basicChargeSettings, &basicChargeSetting)
	}

	err = orm.UpdateBasicChargeSetting(r.ORM.DB, basicChargeSettings, patientId, basicChargeSettingsUUID)
	if err != nil {
		r.Logger.Error("UpdateBasicCharge orm.UpdateBasicChargeSetting", zap.Error(err), zap.String("fieldName", "updateBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
	return true, nil
}

// Queries
func (r *queryResolver) BasicChargeSettings(ctx context.Context, patientIdStr string) ([]*models.BasicChargeSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("BasicChargeSettings uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "basicChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("BasicChargeSettings uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "basicChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicSettings, err := orm.GetBasicChargeSettingsByPatientId(r.ORM.DB, organizationId, patientId)
	if err != nil {
		r.Logger.Error("BasicChargeSettings orm.GetBasicChargeSettingsByPatientId", zap.Error(err), zap.String("fieldName", "basicChargeSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("basicChargeSettings run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "basicChargeSettings"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return basicSettings, nil
}

// basicChargeSetting resolvers
type basicChargeSettingResolver struct{ *Resolver }

func (r *basicChargeSettingResolver) ID(ctx context.Context, obj *models.BasicChargeSetting) (string, error) {
	return obj.ID.String(), nil
}
