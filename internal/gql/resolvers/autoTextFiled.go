package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Mutations
func (r *mutationResolver) CreateAutoTextField(ctx context.Context, input gqlmodels.AutoTextFieldInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateAutoTextField uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	formatIsCorrect := strings.Contains(input.Field, ".")
	if !formatIsCorrect {
		r.Logger.Error("CreateAutoTextField formatIsCorrect", zap.Error(err), zap.String("fieldName", "createAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, fmt.Errorf("name format error")
	}

	moduleNameAndItemName := strings.Split(input.Field, ".")

	autoTextField := models.AutoTextField{
		ID:             uuid.New(),
		ModuleName:     moduleNameAndItemName[0],
		ItemName:       moduleNameAndItemName[1],
		Text:           input.Value,
		OrganizationId: organizationId,
	}

	err = orm.CreateAutoTextField(r.ORM.DB, &autoTextField)
	if err != nil {
		r.Logger.Error("CreateAutoTextField orm.CreateAutoTextField", zap.Error(err), zap.String("fieldName", "createAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createAutoTextField run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createAutoTextField"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}
func (r *mutationResolver) DeleteAutoTextField(ctx context.Context, input gqlmodels.AutoTextFieldInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeleteAutoTextField uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "deleteAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	formatIsCorrect := strings.Contains(input.Field, ".")
	if !formatIsCorrect {
		r.Logger.Error("DeleteAutoTextField formatIsCorrect", zap.Error(err), zap.String("fieldName", "deleteAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, fmt.Errorf("name format error")
	}

	moduleNameAndItemName := strings.Split(input.Field, ".")

	deleteAutoTextField := models.AutoTextField{
		ModuleName:     moduleNameAndItemName[0],
		ItemName:       moduleNameAndItemName[1],
		Text:           input.Value,
		OrganizationId: organizationId,
	}

	err = orm.DeleteAutoTextField(r.ORM.DB, &deleteAutoTextField)
	if err != nil {
		r.Logger.Error("DeleteAutoTextField orm.DeleteAutoTextField", zap.Error(err), zap.String("fieldName", "deleteAutoTextField"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteAutoTextField run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteAutoTextField"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) AutoTextFields(ctx context.Context, field string) ([]*models.AutoTextField, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("AutoTextFields uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "autoTextFields"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	formatIsCorrect := strings.Contains(field, ".")
	if !formatIsCorrect {
		r.Logger.Error("AutoTextFields formatIsCorrect", zap.Error(err), zap.String("fieldName", "autoTextFields"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, fmt.Errorf("name format error")
	}

	moduleNameAndItemName := strings.Split(field, ".")

	autoTextField := models.AutoTextField{
		ModuleName:     moduleNameAndItemName[0],
		ItemName:       moduleNameAndItemName[1],
		OrganizationId: organizationId,
	}

	autoTextFields, err := orm.GetAutoTextFields(r.ORM.DB, &autoTextField)
	if err != nil {
		r.Logger.Error("AutoTextFields orm.GetAutoTextFields", zap.Error(err), zap.String("fieldName", "autoTextFields"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("autoTextFields run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "autoTextFields"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return autoTextFields, nil
}

// autoTextField resolvers
type autoTextFieldResolver struct{ *Resolver }

func (r *autoTextFieldResolver) ID(ctx context.Context, obj *models.AutoTextField) (string, error) {
	return obj.ID.String(), nil
}
