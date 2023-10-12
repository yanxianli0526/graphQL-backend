package resolvers

import (
	"context"
	_ "embed"
	"fmt"
	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"strconv"
	"time"

	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Mutations
func (r *mutationResolver) CreatePayRecordDetail(ctx context.Context, payRecrodIdStr string, input gqlmodels.PayRecordDetailInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreatePayRecordDetail uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreatePayRecordDetail uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrodId, err := uuid.Parse(payRecrodIdStr)
	if err != nil {
		r.Logger.Warn("CreatePayRecordDetail uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrod, err := orm.GetPayRecordById(r.ORM.DB, payRecrodId, false, false, false)
	if err != nil {
		r.Logger.Error("CreatePayRecordDetail orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, payRecrod.PatientBillId)
	if err != nil {
		r.Logger.Error("CreatePayRecordDetail orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecordDetail := models.PayRecordDetail{
		ID:             uuid.New(),
		RecordDate:     input.RecordDate,
		Type:           string(input.Type),
		Price:          input.Price,
		Method:         input.Method,
		Payer:          *input.Payer,
		Handler:        *input.Handler,
		Note:           *input.Note,
		OrganizationId: organizationId,
		PatientId:      payRecrod.PatientId,
		UserId:         userId,
		PayRecordId:    payRecrodId,
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CreatePayRecordDetail(tx, &payRecordDetail)
		if err != nil {
			r.Logger.Error("CreatePayRecordDetail orm.CreatePayRecordDetail", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		err = orm.AppendAssociationsPayRecordPayRecordDetail(tx, payRecrod, payRecordDetail)
		if err != nil {
			r.Logger.Error("CreatePayRecordDetail orm.AppendAssociationsPayRecordPayRecordDetail", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		paidAmount := payRecrod.PaidAmount
		if string(input.Type) == "charge" {
			paidAmount += input.Price
		} else {
			paidAmount -= input.Price
		}

		// 更新完繳費記錄後再更新住民帳單
		err = orm.UpdatePayRecordPaidAmount(tx, paidAmount, payRecrod.ID)
		if err != nil {
			r.Logger.Error("CreatePayRecordDetail orm.UpdatePayRecordPaidAmount", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		patientBill.AmountReceived += paidAmount

		// 更新住民帳單
		err = orm.UpdatePatientBillAmountReceived(tx, patientBill)
		if err != nil {
			r.Logger.Error("CreatePayRecordDetail orm.UpdatePatientBillAmountReceived", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreatePayRecordDetail tx.Transaction", zap.Error(err), zap.String("originalUrl", "createPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createPayRecordDetail run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createPayRecordDetail"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdatePayRecordDetail(ctx context.Context, payRecrodDetailIdStr string, input gqlmodels.PayRecordDetailInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrodDetailId, err := uuid.Parse(payRecrodDetailIdStr)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail uuid.Parse(payRecrodDetailIdStr)", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	getPayRecordDetail, err := orm.GetPayRecordDetail(r.ORM.DB, payRecrodDetailId, organizationId)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail orm.GetPayRecordDetail", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrod, err := orm.GetPayRecordById(r.ORM.DB, getPayRecordDetail.PayRecordId, false, false, true)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, payRecrod.PatientBillId)
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecordDetail := models.PayRecordDetail{
		ID:             payRecrodDetailId,
		RecordDate:     input.RecordDate,
		Type:           string(input.Type),
		Price:          input.Price,
		Method:         input.Method,
		Payer:          *input.Payer,
		Handler:        *input.Handler,
		Note:           *input.Note,
		OrganizationId: organizationId,
		PatientId:      payRecrod.PatientId,
		UserId:         userId,
	}

	var paidAmount int
	for i := range payRecrod.PayRecordDetails {
		if payRecrod.PayRecordDetails[i].ID == payRecrodDetailId {
			// 表示是現在正在更新的這筆
			if string(input.Type) == "charge" {
				paidAmount += input.Price
			} else {
				paidAmount -= input.Price
			}
		} else {
			if payRecrod.PayRecordDetails[i].Type == "charge" {
				paidAmount += payRecrod.PayRecordDetails[i].Price
			} else {
				paidAmount -= payRecrod.PayRecordDetails[i].Price
			}
		}
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.UpdatePayRecordDetail(tx, &payRecordDetail, organizationId)
		if err != nil {
			r.Logger.Error("UpdatePayRecordDetail orm.UpdatePayRecordDetail", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		// 更新完繳費記錄後再更新住民帳單
		err = orm.UpdatePayRecordPaidAmount(tx, paidAmount, payRecrod.ID)
		if err != nil {
			r.Logger.Error("UpdatePayRecordDetail orm.UpdatePayRecordDetail", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		patientBill.AmountReceived += paidAmount
		// 更新住民帳單
		err = orm.UpdatePatientBillAmountReceived(tx, patientBill)
		if err != nil {
			r.Logger.Error("UpdatePayRecordDetail orm.UpdatePatientBillAmountReceived", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdatePayRecordDetail tx.Transaction", zap.Error(err), zap.String("originalUrl", "updatePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePayRecordDetail run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePayRecordDetail"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeletePayRecordDetail(ctx context.Context, payRecrodDetailIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("DeletePayRecordDetail uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeletePayRecordDetail uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrodDetailId, err := uuid.Parse(payRecrodDetailIdStr)
	if err != nil {
		r.Logger.Warn("DeletePayRecordDetail uuid.Parse(payRecrodDetailIdStr)", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrodDetail, err := orm.GetPayRecordDetail(r.ORM.DB, payRecrodDetailId, organizationId)
	if err != nil {
		r.Logger.Error("DeletePayRecordDetail orm.GetPayRecordDetail", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	payRecrod, err := orm.GetPayRecordById(r.ORM.DB, payRecrodDetail.PayRecordId, false, false, false)
	if err != nil {
		r.Logger.Error("DeletePayRecordDetail orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, payRecrod.PatientBillId)
	if err != nil {
		r.Logger.Error("DeletePayRecordDetail orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	paidAmount := payRecrod.PaidAmount

	if payRecrodDetail.Type == "charge" {
		paidAmount -= payRecrodDetail.Price
	} else {
		paidAmount += payRecrodDetail.Price
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		updatePayRecordDetail := models.PayRecordDetail{
			ID:     payRecrodDetailId,
			UserId: userId,
		}
		// 記錄最後一個刪除的人
		err = orm.UpdatePayRecordDetailUser(tx, &updatePayRecordDetail, organizationId)
		// 把資料刪掉
		err = orm.DeletePayRecordDetail(tx, payRecrodDetailId, organizationId)
		if err != nil {
			r.Logger.Error("DeletePayRecordDetail orm.DeletePayRecordDetail", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		// 更新完繳費記錄後再更新住民帳單
		err = orm.UpdatePayRecordPaidAmount(tx, paidAmount, payRecrod.ID)
		if err != nil {
			r.Logger.Error("DeletePayRecordDetail orm.UpdatePayRecordPaidAmount", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		patientBill.AmountReceived = paidAmount
		// 更新住民帳單
		err = orm.UpdatePatientBillAmountReceived(tx, patientBill)
		if err != nil {
			r.Logger.Error("DeletePayRecordDetail orm.UpdatePatientBillAmountReceived", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("DeletePayRecordDetail tx.Transaction", zap.Error(err), zap.String("originalUrl", "deletePayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deletePayRecordDetail run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deletePayRecordDetail"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) PayRecordDetail(ctx context.Context, payRecrodDetailIdStr string) (*models.PayRecordDetail, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PayRecordDetail uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "payRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	payRecrodDetailId, err := uuid.Parse(payRecrodDetailIdStr)
	if err != nil {
		r.Logger.Warn("PayRecordDetail uuid.Parse(payRecrodDetailIdStr)", zap.Error(err), zap.String("originalUrl", "payRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	payRecrodDetail, err := orm.GetPayRecordDetail(r.ORM.DB, payRecrodDetailId, organizationId)
	if err != nil {
		r.Logger.Error("PayRecordDetail orm.GetPayRecordDetail", zap.Error(err), zap.String("originalUrl", "payRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("payRecordDetail run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "payRecordDetail"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return payRecrodDetail, nil
}

// payRecordDetail resolvers
type payRecordDetailResolver struct{ *Resolver }

func (r *payRecordDetailResolver) ID(ctx context.Context, obj *models.PayRecordDetail) (string, error) {
	return obj.ID.String(), nil
}

func (r *payRecordDetailResolver) Type(ctx context.Context, obj *models.PayRecordDetail) (gqlmodels.PayRecordDetailType, error) {
	typeStr := gqlmodels.PayRecordDetailType(obj.Type)
	isValid := gqlmodels.PayRecordDetailType.IsValid(gqlmodels.PayRecordDetailType(typeStr))
	if !isValid {
		r.Logger.Error("PayRecordDetail Type is inValid", zap.String("fieldName", "payRecordDetail"), zap.Int64("timestamp", time.Now().Unix()))
		return "", fmt.Errorf("PettyCash Type is inValid ")
	} else {
		return typeStr, nil
	}
}
