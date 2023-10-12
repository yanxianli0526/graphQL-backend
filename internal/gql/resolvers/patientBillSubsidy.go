package resolvers

import (
	"context"
	"fmt"
	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Mutations
func (r *mutationResolver) AddPatientBillSubsidy(ctx context.Context, input *gqlmodels.CreatePatientBillSubsidyInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("AddPatientBillSubsidy uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("AddPatientBillSubsidy uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	patientBillId, err := uuid.Parse(input.ID)
	if err != nil {
		r.Logger.Warn("AddPatientBillSubsidy uuid.Parse(input.ID)", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	// 先找出該資料
	patientBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, patientBillId)
	if err != nil {
		r.Logger.Error("AddPatientBillSubsidy orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	subsidy := models.Subsidy{
		ID:             uuid.New(),
		ItemName:       input.ItemName,
		Type:           input.Type,
		Price:          input.Price,
		Unit:           *input.Unit,
		IdNumber:       *input.IDNumber,
		Note:           *input.Note,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		ReceiptStatus:  "",
		SortIndex:      len(patientBill.Subsidies),
		PatientId:      patientBill.PatientId,
		OrganizationId: organizationId,
		UserId:         userId,
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CreateSubsidy(tx, &subsidy)
		if err != nil {
			r.Logger.Error("AddPatientBillSubsidy orm.CreateSubsidy", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}

		err = orm.AppendAssociationsPatientBillSubSidy(tx, patientBill, subsidy)
		if err != nil {
			r.Logger.Error("AddPatientBillSubsidy orm.AppendAssociationsPatientBillSubSidy", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		var priceSpread int
		if input.Type == "charge" {
			priceSpread = input.Price
		} else {
			priceSpread = -input.Price
		}
		updatePatientBillAmountDueStruct := UpdatePatientBillAmountDueStruct{
			PatientBill:  patientBill,
			NewAmountDue: patientBill.AmountDue + priceSpread,
			Tx:           tx,
		}
		// 更新應繳金額
		err = UpdatePatientBillAmountDue(updatePatientBillAmountDueStruct)
		if err != nil {
			r.Logger.Error("AddPatientBillBasicCharge UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("AddPatientBillSubsidy tx.Transaction", zap.Error(err), zap.String("originalUrl", "addPatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("addPatientBillSubsidy run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "addPatientBillSubsidy"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdatePatientBillSubsidy(ctx context.Context, input *gqlmodels.UpdatePatientBillSubsidyInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillSubsidy uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "updatePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillSubsidy uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updatePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	subsidyId, err := uuid.Parse(input.SubsidyID)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillSubsidy uuid.Parse(input.SubsidyID)", zap.Error(err), zap.String("originalUrl", "updatePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	subsidy := models.Subsidy{
		ID:             subsidyId,
		ItemName:       *input.ItemName,
		Type:           input.Type,
		Price:          input.Price,
		Unit:           input.Unit,
		IdNumber:       *input.IDNumber,
		Note:           *input.Note,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		ReceiptStatus:  "",
		OrganizationId: organizationId,
		UserId:         userId,
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdateSubsidy, err := orm.GetSubsidy(tx, organizationId, subsidyId)
		if err != nil {
			r.Logger.Error("UpdatePatientBillSubsidy orm.GetSubsidy", zap.Error(err), zap.String("fieldName", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needUpdatePatientBill := orm.GetPatientBillHaveSubsidyId(tx, beforeUpdateSubsidy.PatientId, organizationId, subsidyId)
		if needUpdatePatientBill == nil {
			r.Logger.Error("UpdatePatientBillSubsidy orm.GetPatientBillHaveSubsidyId", zap.Error(err), zap.String("fieldName", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return fmt.Errorf("orm.GetPatientBillHaveSubsidyId have error")
		}
		// 更新住民帳單狀態(紀錄成應繳金額不是最新)
		var priceSpread int
		if beforeUpdateSubsidy.Type == "charge" {
			priceSpread = -beforeUpdateSubsidy.Price
		} else {
			priceSpread = beforeUpdateSubsidy.Price
		}
		if input.Type == "charge" {
			priceSpread += input.Price
		} else {
			priceSpread -= input.Price
		}
		updatePatientBillAmountDue := UpdatePatientBillAmountDueStruct{
			PatientBill:  needUpdatePatientBill,
			NewAmountDue: needUpdatePatientBill.AmountDue + priceSpread,
			Tx:           tx,
		}
		// 更新應繳金額
		err = UpdatePatientBillAmountDue(updatePatientBillAmountDue)
		if err != nil {
			r.Logger.Error("UpdatePatientBillSubsidy UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		err = orm.UpdatePatientBillSubsidy(tx, &subsidy)
		if err != nil {
			r.Logger.Error("UpdatePatientBillSubsidy orm.UpdatePatientBillSubsidy", zap.Error(err), zap.String("originalUrl", "updatePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdatePatientBillSubsidy tx.Transaction", zap.Error(err), zap.String("originalUrl", "updatePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePatientBillSubsidy run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePatientBillSubsidy"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeletePatientBillSubsidy(ctx context.Context, subsidyIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeletePatientBillSubsidy uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	subsidyId, err := uuid.Parse(subsidyIdStr)
	if err != nil {
		r.Logger.Warn("DeletePatientBillSubsidy uuid.Parse(subsidyIdStr)", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdateSubsidy, err := orm.GetSubsidy(tx, organizationId, subsidyId)
		if err != nil {
			r.Logger.Error("DeletePatientBillSubsidy orm.GetSubsidy", zap.Error(err), zap.String("fieldName", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needUpdatePatientBill := orm.GetPatientBillHaveSubsidyId(tx, beforeUpdateSubsidy.PatientId, organizationId, subsidyId)
		if needUpdatePatientBill == nil {
			r.Logger.Error("DeletePatientBillSubsidy orm.GetPatientBillHaveSubsidyId", zap.Error(err), zap.String("fieldName", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return fmt.Errorf("orm.GetPatientBillHaveSubsidyId have error")
		}
		// 更新住民帳單狀態(紀錄成應繳金額不是最新)
		var priceSpread int
		if beforeUpdateSubsidy.Type == "charge" {
			priceSpread = -beforeUpdateSubsidy.Price
		} else {
			priceSpread = beforeUpdateSubsidy.Price
		}
		updatePatientBillAmountDue := UpdatePatientBillAmountDueStruct{
			PatientBill:  needUpdatePatientBill,
			NewAmountDue: needUpdatePatientBill.AmountDue + priceSpread,
			Tx:           tx,
		}
		// 更新應繳金額
		err = UpdatePatientBillAmountDue(updatePatientBillAmountDue)
		if err != nil {
			r.Logger.Error("DeletePatientBillSubsidy UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		err = orm.DeleteSubsidy(tx, organizationId, subsidyId)
		if err != nil {
			r.Logger.Error("DeletePatientBillSubsidy orm.DeleteSubsidy", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("DeletePatientBillSubsidy tx.Transaction", zap.Error(err), zap.String("originalUrl", "deletePatientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deletePatientBillSubsidy run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deletePatientBillSubsidy"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
// 前端的邏輯異動和非固定都是關聯 固定月費建立完就脫勾了,所以會有自己的查詢(編輯自己的form))
func (r *queryResolver) PatientBillSubsidy(ctx context.Context, subsidyIdStr string) (*models.Subsidy, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientBillSubsidy uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "patientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	subsidyId, err := uuid.Parse(subsidyIdStr)
	if err != nil {
		r.Logger.Warn("PatientBillSubsidy uuid.Parse(subsidyIdStr)", zap.Error(err), zap.String("originalUrl", "patientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	subsidy, err := orm.GetSubsidy(r.ORM.DB, organizationId, subsidyId)
	if err != nil {
		r.Logger.Error("PatientBillSubsidy orm.GetSubsidy", zap.Error(err), zap.String("originalUrl", "patientBillSubsidy"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientBillSubsidy run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientBillSubsidy"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return subsidy, nil
}
