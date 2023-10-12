package resolvers

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

// Mutations
func (r *mutationResolver) CreateDepositRecord(ctx context.Context, input gqlmodels.DepositRecordInput) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Info("CreateDepositRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateDepositRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	patientId, err := uuid.Parse(input.PatientID)
	if err != nil {
		r.Logger.Warn("CreateDepositRecord uuid.Parse(input.PatientID)", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	idNumber, err := orm.GetDepositRecordLastIdNumber(r.ORM.DB, organizationId, input.Date)
	if err != nil {
		r.Logger.Warn("CreateDepositRecord orm.GetDepositRecordLastIdNumber", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 如果是no record代表他今天是第一次新增
	if idNumber == "no record" {
		idNumber = input.Date.Format("20060102") + "001"
	} else {
		idNumberInt, err := strconv.Atoi(idNumber) // result: i = -18
		if err != nil {
			r.Logger.Error("CreateDepositRecord strconv.Atoi(idNumber)", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
		idNumber = strconv.Itoa(idNumberInt + 1)
	}

	depositRecord := models.DepositRecord{
		IdNumber:       idNumber,
		Date:           input.Date,
		Type:           input.Type,
		Price:          input.Price,
		Drawee:         *input.Drawee,
		Note:           *input.Note,
		PatientId:      patientId,
		OrganizationId: organizationId,
		UserId:         userId,
		Invalid:        false,
	}
	createdDepositRecordId, err := orm.CreateDepositRecord(r.ORM.DB, &depositRecord)
	if err != nil {
		r.Logger.Error("CreateDepositRecord orm.CreateDepositRecord", zap.Error(err), zap.String("originalUrl", "createDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createDepositRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createDepositRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return createdDepositRecordId.String(), nil
}

func (r *mutationResolver) UpdateDepositRecord(ctx context.Context, id string, input gqlmodels.DepositRecordUpdateInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdateDepositRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "updateDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	depositRecordId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("UpdateDepositRecord uuid.Parse(id)", zap.Error(err), zap.String("originalUrl", "updateDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	depositRecord := models.DepositRecord{
		ID:     depositRecordId,
		Note:   *input.Note,
		UserId: userId,
	}
	err = orm.UpdateDepositRecordById(r.ORM.DB, &depositRecord)
	if err != nil {
		r.Logger.Error("UpdateDepositRecord orm.UpdateDepositRecordById", zap.Error(err), zap.String("originalUrl", "updateDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateDepositRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateDepositRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) InvalidDepositRecord(ctx context.Context, id string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	depositRecordId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("InvalidDepositRecord uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "invalidDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	err = orm.UpdateInvalidDepositRecordById(r.ORM.DB, depositRecordId)
	if err != nil {
		r.Logger.Error("InvalidDepositRecord orm.UpdateInvalidDepositRecordById", zap.Error(err), zap.String("fieldName", "invalidDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("invalidDepositRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "invalidDepositRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) DepositRecords(ctx context.Context, patientID string) ([]*models.DepositRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DepositRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "depositRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientId, err := uuid.Parse(patientID)
	if err != nil {
		r.Logger.Warn("DepositRecords uuid.Parse(patientID)", zap.Error(err), zap.String("fieldName", "depositRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	collectFieldsCtx := graphql.CollectFieldsCtx(ctx, nil)

	var preloadUser bool
	var preloadPatient bool
	// 可以用這個減低db的負擔 只針對需要的欄位去進行preload
	for i := range collectFieldsCtx {
		if collectFieldsCtx[i].Name == "user" {
			preloadUser = true
		}
		if collectFieldsCtx[i].Name == "patient" {
			preloadPatient = true
		}
	}

	depositRecords, err := orm.GetDepositRecords(r.ORM.DB, organizationId, patientId, preloadUser, preloadPatient)
	if err != nil {
		r.Logger.Error("DepositRecords orm.GetDepositRecords", zap.Error(err), zap.String("fieldName", "depositRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("depositRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "depositRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return depositRecords, nil
}

func (r *queryResolver) DepositRecord(ctx context.Context, id string) (*models.DepositRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DepositRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "depositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	depositRecordId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("DepositRecord uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "depositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	collectFieldsCtx := graphql.CollectFieldsCtx(ctx, nil)

	var preloadUser bool
	var preloadPatient bool
	// 可以用這個減低db的負擔 只針對需要的欄位去進行preload
	for i := range collectFieldsCtx {
		if collectFieldsCtx[i].Name == "user" {
			preloadUser = true
		}
		if collectFieldsCtx[i].Name == "patient" {
			preloadPatient = true
		}
	}

	depositRecord, err := orm.GetDepositRecordById(r.ORM.DB, depositRecordId, organizationId, preloadUser, preloadPatient)
	if err != nil {
		r.Logger.Error("DepositRecord orm.GetDepositRecordById", zap.Error(err), zap.String("fieldName", "depositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("depositRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "depositRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return depositRecord, nil
}

func (r *queryResolver) PatientLatestDepositRecords(ctx context.Context) (*gqlmodels.PatientLatestDepositRecords, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientLatestDepositRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "patientLatestDepositRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	depositRecordsDescByDate, depositRecordsDescByUpdatedAt, err := orm.PatientLatestDepositRecords(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PatientLatestDepositRecords orm.PatientLatestDepositRecords", zap.Error(err), zap.String("fieldName", "patientLatestDepositRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientLatestDepositRecord := gqlmodels.PatientLatestDepositRecords{
		LatestRecord:        depositRecordsDescByDate,
		LatestUpdatedRecord: depositRecordsDescByUpdatedAt,
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientLatestDepositRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientLatestDepositRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return &patientLatestDepositRecord, nil
}

// depositRecord resolvers
type depositRecordResolver struct{ *Resolver }

func (r *depositRecordResolver) ID(ctx context.Context, obj *models.DepositRecord) (string, error) {
	return obj.ID.String(), nil
}
