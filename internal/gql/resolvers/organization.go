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
func (r *mutationResolver) UpdateOrganization(ctx context.Context, input *gqlmodels.OrganizationSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganization uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updateOrganization"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	organization := models.Organization{
		ID:                  organizationId,
		Name:                input.Name,
		AddressCity:         input.AddressCity,
		AddressDistrict:     input.AddressDistrict,
		Address:             input.Address,
		Phone:               input.Phone,
		Fax:                 input.Fax,
		Owner:               input.Owner,
		Email:               input.Email,
		TaxIdNumber:         input.TaxIDNumber,
		RemittanceBank:      input.RemittanceBank,
		RemittanceIdNumber:  input.RemittanceIDNumber,
		RemittanceUserName:  input.RemittanceUserName,
		EstablishmentNumber: input.EstablishmentNumber,
	}
	err = orm.UpdateOrganzationById(r.ORM.DB, &organization)
	if err != nil {
		r.Logger.Error("UpdateOrganization orm.UpdateOrganzationById", zap.Error(err), zap.String("originalUrl", "updateOrganization"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganization run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganization"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateOrganizationPrivacy(ctx context.Context, input *gqlmodels.OrganizationPrivacyInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationPrivacy uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updateOrganizationPrivacy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	organization := models.Organization{
		ID:      organizationId,
		Privacy: input.Privacy.String(),
	}
	err = orm.UpdateOrganizationPrivacy(r.ORM.DB, &organization)
	if err != nil {
		r.Logger.Error("UpdateOrganizationPrivacy orm.UpdateOrganizationPrivacy", zap.Error(err), zap.String("originalUrl", "updateOrganizationPrivacy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationPrivacy run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationPrivacy"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateOrganizationBillDateRangeSetting(ctx context.Context, input *gqlmodels.OrganizationBillDateRangeSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationBillDateRangeSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updateOrganizationBillDateRangeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	organization := models.Organization{
		ID:                       organizationId,
		FixedChargeStartMonth:    input.FixedChargeStartMonth,
		FixedChargeStartDate:     input.FixedChargeStartDate,
		FixedChargeEndMonth:      input.FixedChargeEndMonth,
		FixedChargeEndDate:       input.FixedChargeEndDate,
		NonFixedChargeStartMonth: input.NonFixedChargeStartMonth,
		NonFixedChargeStartDate:  input.NonFixedChargeStartDate,
		NonFixedChargeEndMonth:   input.NonFixedChargeEndMonth,
		NonFixedChargeEndDate:    input.NonFixedChargeEndDate,
		TransferRefundStartMonth: input.TransferRefundStartMonth,
		TransferRefundStartDate:  input.TransferRefundStartDate,
		TransferRefundEndMonth:   input.TransferRefundEndMonth,
		TransferRefundEndDate:    input.TransferRefundEndDate,
	}
	err = orm.UpdateOrganizationBillDateRangeSetting(r.ORM.DB, &organization)
	if err != nil {
		r.Logger.Error("UpdateOrganizationBillDateRangeSetting orm.UpdateOrganizationBillDateRangeSetting", zap.Error(err), zap.String("originalUrl", "updateOrganizationBillDateRangeSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationBillDateRangeSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationBillDateRangeSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) Organization(ctx context.Context) (*models.Organization, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("Organization uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "organization"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	organization, err := orm.GetOrganizationById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("Organization orm.GetOrganizationById", zap.Error(err), zap.String("originalUrl", "organization"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organization run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organization"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return organization, nil
}

// organization resolvers
type organizationResolver struct{ *Resolver }

func (r *organizationResolver) ID(ctx context.Context, obj *models.Organization) (string, error) {
	return obj.ID.String(), nil
}

func (r *organizationResolver) Privacy(ctx context.Context, obj *models.Organization) (gqlmodels.PrivacyType, error) {
	privacy := gqlmodels.PrivacyType(obj.Privacy)
	isValid := gqlmodels.PrivacyType.IsValid(gqlmodels.PrivacyType(privacy))
	if !isValid {
		r.Logger.Error("Organization Privacy is inValid", zap.String("fieldName", "organization"), zap.Int64("timestamp", time.Now().Unix()))
		return "", fmt.Errorf("Organization Privacy is inValid ")
	} else {
		return privacy, nil
	}
}
