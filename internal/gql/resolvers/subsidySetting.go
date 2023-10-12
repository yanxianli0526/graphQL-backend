package resolvers

import (
	"context"
	"fmt"
	orm "graphql-go-template/internal/database"
	"strconv"
	"time"

	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Mutations
func (r *mutationResolver) CreateSubsidiesSetting(ctx context.Context, patientIdStr string, input []*gqlmodels.SubsidySettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateSubsidiesSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createSubsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("CreateSubsidiesSetting uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "createSubsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	if len(input) > 0 {
		var subsidiesSetting []*models.SubsidySetting
		for i := range input {
			subsidySetting := models.SubsidySetting{
				ItemName:       input[i].ItemName,
				Type:           input[i].Type,
				Price:          input[i].Price,
				Unit:           *input[i].Unit,
				IdNumber:       *input[i].IDNumber,
				Note:           *input[i].Note,
				OrganizationId: organizationId,
				PatientId:      patientId,
				SortIndex:      i,
			}
			subsidiesSetting = append(subsidiesSetting, &subsidySetting)
		}

		err = orm.CreateSubsidiesSetting(r.ORM.DB, subsidiesSetting)
		if err != nil {
			r.Logger.Error("CreateSubsidiesSetting orm.CreateSubsidiesSetting", zap.Error(err), zap.String("fieldName", "createSubsidiesSetting"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createSubsidiesSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createSubsidiesSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
	return true, nil
}

func (r *mutationResolver) UpdateSubsidiesSetting(ctx context.Context, patientIdStr string, input []*gqlmodels.SubsidySettingUpdateInput, needDeleteSubsidiesSettingIdStr []string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateSubsidiesSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updateSubsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("UpdateSubsidiesSetting uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "updateSubsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var needDeleteSubsidiesSettingId []uuid.UUID
	for i := range needDeleteSubsidiesSettingIdStr {
		needDeleteSubsidySettingId, err := uuid.Parse(needDeleteSubsidiesSettingIdStr[i])
		if err != nil {
			r.Logger.Warn("UpdateSubsidiesSetting uuid.Parse(needDeleteSubsidiesSettingIdStr)", zap.Error(err), zap.String("fieldName", "updateSubsidiesSetting"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
		needDeleteSubsidiesSettingId = append(needDeleteSubsidiesSettingId, needDeleteSubsidySettingId)
	}

	var updateSubsidiesSetting []*models.SubsidySetting
	var createSubsidiesSetting []*models.SubsidySetting
	for i := range input {
		if input[i].ID != nil {
			id, err := uuid.Parse(*input[i].ID)
			if err != nil {
				r.Logger.Warn("UpdateSubsidiesSetting uuid.Parse(*input[i].ID)", zap.Error(err), zap.String("fieldName", "updateSubsidiesSetting"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return false, err
			}
			updateSubsidySetting := models.SubsidySetting{
				ID:             id,
				ItemName:       input[i].ItemName,
				Type:           input[i].Type,
				Price:          input[i].Price,
				Unit:           *input[i].Unit,
				IdNumber:       *input[i].IDNumber,
				Note:           *input[i].Note,
				SortIndex:      i,
				OrganizationId: organizationId,
				PatientId:      patientId,
			}
			updateSubsidiesSetting = append(updateSubsidiesSetting, &updateSubsidySetting)
		} else {
			createSubsidySetting := models.SubsidySetting{
				ItemName:       input[i].ItemName,
				Type:           input[i].Type,
				Price:          input[i].Price,
				Unit:           *input[i].Unit,
				IdNumber:       *input[i].IDNumber,
				Note:           *input[i].Note,
				SortIndex:      i,
				OrganizationId: organizationId,
				PatientId:      patientId,
			}
			createSubsidiesSetting = append(createSubsidiesSetting, &createSubsidySetting)
		}
	}

	err = orm.UpdateSubsidiesSetting(r.ORM.DB, createSubsidiesSetting, updateSubsidiesSetting, needDeleteSubsidiesSettingId, patientId)
	if err != nil {
		r.Logger.Error("UpdateSubsidiesSetting orm.UpdateSubsidiesSetting", zap.Error(err), zap.String("fieldName", "updateSubsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateSubsidiesSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateSubsidiesSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) SubsidiesSetting(ctx context.Context, patientIdStr string) ([]*models.SubsidySetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("SubsidiesSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "subsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("SubsidiesSetting uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "subsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	subsidiesSetting, err := orm.GetSubsidiesSettingByPatientId(r.ORM.DB, organizationId, patientId)
	if err != nil {
		r.Logger.Error("SubsidiesSetting orm.GetSubsidiesSettingByPatientId", zap.Error(err), zap.String("fieldName", "subsidiesSetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("subsidiesSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "subsidiesSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return subsidiesSetting, nil
}

func (r *queryResolver) SubsidySetting(ctx context.Context, subsidySettingIdStr string) (*models.SubsidySetting, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("SubsidySetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "subsidySetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	subsidySettingId, err := uuid.Parse(subsidySettingIdStr)
	if err != nil {
		r.Logger.Warn("SubsidySetting uuid.Parse(subsidySettingIdStr)", zap.Error(err), zap.String("fieldName", "subsidySetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	subsidySetting, err := orm.GetSubsidySetting(r.ORM.DB, organizationId, subsidySettingId)
	if err != nil {
		r.Logger.Error("SubsidySetting orm.GetSubsidySetting", zap.Error(err), zap.String("fieldName", "subsidySetting"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("subsidySetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "subsidySetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return subsidySetting, nil
}

// subsidySetting resolvers
type subsidySettingResolver struct{ *Resolver }

func (r *subsidySettingResolver) ID(ctx context.Context, obj *models.SubsidySetting) (string, error) {
	return obj.ID.String(), nil
}
