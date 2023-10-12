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
func (r *mutationResolver) AddPatientBillBasicCharge(ctx context.Context, input *gqlmodels.CreatePatientBillBasicChargeInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("AddPatientBillBasicCharge uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("AddPatientBillBasicCharge uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	patientBillId, err := uuid.Parse(input.ID)
	if err != nil {
		r.Logger.Warn("AddPatientBillBasicCharge uuid.Parse(input.ID)", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	// 先找出該資料
	patientBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, patientBillId)
	if err != nil {
		r.Logger.Error("AddPatientBillBasicCharge orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	basicCharge := models.BasicCharge{
		ID:             uuid.New(),
		ItemName:       input.ItemName,
		Type:           input.Type,
		Unit:           input.Unit,
		Price:          input.Price,
		TaxType:        input.TaxType,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Note:           *input.Note,
		SortIndex:      len(patientBill.BasicCharges),
		UserId:         userId,
		PatientId:      patientBill.PatientId,
		OrganizationId: patientBill.OrganizationId,
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CreateBasicCharge(tx, &basicCharge)
		if err != nil {
			r.Logger.Error("AddPatientBillBasicCharge orm.CreateBasicCharge", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}

		err = orm.AppendAssociationsPatientBillBasicCharge(tx, patientBill, basicCharge)
		if err != nil {
			r.Logger.Error("AddPatientBillBasicCharge orm.AppendAssociationsPatientBillBasicCharge", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
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
		r.Logger.Error("AddPatientBillBasicCharge tx.Transaction", zap.Error(err), zap.String("originalUrl", "addPatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("addPatientBillBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "addPatientBillBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdatePatientBillBasicCharge(ctx context.Context, input *gqlmodels.UpdatePatientBillBasicChargeInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillBasicCharge uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillBasicCharge uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	basicChargeId, err := uuid.Parse(input.BasicChargeID)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillBasicCharge uuid.Parse(input.BasicChargeID)", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	basicCharge := models.BasicCharge{
		ID:             basicChargeId,
		ItemName:       input.ItemName,
		Type:           input.Type,
		Unit:           input.Unit,
		Price:          input.Price,
		TaxType:        input.TaxType,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Note:           *input.Note,
		OrganizationId: organizationId,
		UserId:         userId,
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdateBasicCharge, err := orm.GetBasicCharge(tx, organizationId, basicChargeId)
		if err != nil {
			r.Logger.Error("UpdatePatientBillBasicCharge orm.GetBasicCharge", zap.Error(err), zap.String("fieldName", "updatePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needUpdatePatientBill := orm.GetPatientBillHaveBasicChargeId(tx, beforeUpdateBasicCharge.PatientId, organizationId, basicChargeId)
		if needUpdatePatientBill == nil {
			r.Logger.Error("DeletePatientBillBasicCharge orm.GetPatientBillHaveBasicChargeId", zap.Error(err), zap.String("fieldName", "deletePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return fmt.Errorf("orm.GetPatientBillHaveBasicChargeId have error")
		}
		// 更新住民帳單狀態(紀錄成應繳金額不是最新)
		var priceSpread int
		if beforeUpdateBasicCharge.Type == "charge" {
			priceSpread = -beforeUpdateBasicCharge.Price
		} else {
			priceSpread = beforeUpdateBasicCharge.Price
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
			r.Logger.Error("UpdatePatientBillBasicCharge UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		err = orm.UpdatePatientBillBasicCharge(tx, &basicCharge)
		if err != nil {
			r.Logger.Error("UpdatePatientBillBasicCharge orm.UpdatePatientBillBasicCharge", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdatePatientBillBasicCharge tx.Transaction", zap.Error(err), zap.String("originalUrl", "updatePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePatientBillBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePatientBillBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// 這裡的刪除並不是真的把資料庫的fixed_charge刪除 而且刪除裡面json中的其中一項
func (r *mutationResolver) DeletePatientBillBasicCharge(ctx context.Context, basicChargeIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeletePatientBillBasicCharge uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "deletePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	basicChargeId, err := uuid.Parse(basicChargeIdStr)
	if err != nil {
		r.Logger.Warn("DeletePatientBillBasicCharge uuid.Parse(basicChargeIdStr)", zap.Error(err), zap.String("originalUrl", "deletePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdateBasicCharge, err := orm.GetBasicCharge(tx, organizationId, basicChargeId)
		if err != nil {
			r.Logger.Error("DeletePatientBillBasicCharge orm.GetBasicCharge", zap.Error(err), zap.String("fieldName", "deletePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needUpdatePatientBill := orm.GetPatientBillHaveBasicChargeId(tx, beforeUpdateBasicCharge.PatientId, organizationId, basicChargeId)
		if needUpdatePatientBill == nil {
			r.Logger.Error("DeletePatientBillBasicCharge orm.GetPatientBillHaveBasicChargeId", zap.Error(err), zap.String("fieldName", "deletePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return fmt.Errorf("orm.GetPatientBillHaveBasicChargeId have error")
		}
		// 更新住民帳單狀態(紀錄成應繳金額不是最新)
		var priceSpread int
		if beforeUpdateBasicCharge.Type == "charge" {
			priceSpread = -beforeUpdateBasicCharge.Price
		} else {
			priceSpread = beforeUpdateBasicCharge.Price
		}
		updatePatientBillAmountDue := UpdatePatientBillAmountDueStruct{
			PatientBill:  needUpdatePatientBill,
			NewAmountDue: needUpdatePatientBill.AmountDue + priceSpread,
			Tx:           tx,
		}
		// 更新應繳金額
		err = UpdatePatientBillAmountDue(updatePatientBillAmountDue)
		if err != nil {
			r.Logger.Error("DeletePatientBillBasicCharge UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "deletePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}

		err = orm.DeleteBasicCharge(tx, organizationId, basicChargeId)
		if err != nil {
			r.Logger.Error("DeletePatientBillBasicCharge orm.DeleteBasicCharge", zap.Error(err), zap.String("originalUrl", "deletePatientBillBasicCharge"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("DeletePatientBillBasicCharge tx.Transaction", zap.Error(err), zap.String("originalUrl", "deletePatientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deletePatientBillBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deletePatientBillBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
// 前端的邏輯異動和非固定都是關聯 固定月費建立完就脫勾了,所以會有自己的查詢(編輯自己的form))
func (r *queryResolver) PatientBillBasicCharge(ctx context.Context, basicChargeIdStr string) (*models.BasicCharge, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientBillBasicCharge uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "patientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicChargeId, err := uuid.Parse(basicChargeIdStr)
	if err != nil {
		r.Logger.Warn("PatientBillBasicCharge uuid.Parse(basicChargeIdStr)", zap.Error(err), zap.String("originalUrl", "patientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicCharge, err := orm.GetBasicCharge(r.ORM.DB, organizationId, basicChargeId)
	if err != nil {
		r.Logger.Error("PatientBillBasicCharge orm.GetBasicCharge", zap.Error(err), zap.String("originalUrl", "patientBillBasicCharge"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientBillBasicCharge run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientBillBasicCharge"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return basicCharge, nil
}
