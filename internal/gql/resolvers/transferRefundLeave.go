package resolvers

import (
	"context"
	"encoding/json"
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
func (r *mutationResolver) CreateTransferRefundLeave(ctx context.Context, patientIdStr string, input gqlmodels.TransferRefundLeaveInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreateTransferRefundLeave uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateTransferRefundLeave uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("CreateTransferRefundLeave uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	// 計算這次新增的總金額(目的是要算應繳金額)
	var itemSubtotal int
	for i := range input.Items {
		if input.Items[i].Type == "charge" {
			itemSubtotal += input.Items[i].Price
		} else {
			itemSubtotal -= input.Items[i].Price
		}
	}

	items, err := json.Marshal(input.Items)
	if err != nil {
		r.Logger.Error("CreateTransferRefundLeave json.Marshal(input.Items)", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	var note string
	if input.Note != nil {
		note = *input.Note
	}

	transferRefundLeave := models.TransferRefundLeave{
		ID:             uuid.New(),
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Reason:         input.Reason,
		IsReserveBed:   input.IsReserveBed,
		Note:           note,
		Items:          items,
		OrganizationId: organizationId,
		PatientId:      patientId,
		UserId:         userId,
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CreateTransferRefundLeave(r.ORM.DB, &transferRefundLeave)
		if err != nil {
			r.Logger.Error("CreateTransferRefundLeave orm.CreateTransferRefundLeave", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 以下都在更新住民帳單

		//	如果沒有住帳單:世界和平
		//  如果有:要檢查住民帳單跟非固定的關聯
		patientBills := orm.GetPatientBillsByTransferRefundLeaveDate(tx, patientId, input.EndDate)
		if len(patientBills) == 0 {
			return nil
		} else {
			// 表示已經開帳了 要檢查異動的日期 是不是在住民帳單的異動的區間
			// 在區間:這次新增的全部加進去
			// 不在區間:邏輯上不會到這邊,這個情況是(有跑到這邊代表彥賢有問題....)
			for i := range patientBills {
				if input.EndDate.Unix() >= patientBills[i].TransferRefundStartDate.Unix() && input.EndDate.Unix() <= patientBills[i].TransferRefundEndDate.Unix() {
					err = orm.AppendAssociationsPatientBillTransferRefundLeave(tx, patientBills[i], transferRefundLeave)
					if err != nil {
						r.Logger.Error("CreateTransferRefundLeave orm.AppendAssociationsPatientBillTransferRefundLeave", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
							zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
						return err
					}
					updatePatientBillAmountDueStruct := UpdatePatientBillAmountDueStruct{
						PatientBill:  patientBills[i],
						NewAmountDue: patientBills[i].AmountDue + itemSubtotal,
						Tx:           tx,
					}
					// 更新應繳金額
					err = UpdatePatientBillAmountDue(updatePatientBillAmountDueStruct)
					if err != nil {
						r.Logger.Error("CreateTransferRefundLeave UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "createTransferRefundLeave"),
							zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
						return err
					}
				} else {
					//	邏輯上不會到這邊(有跑到這邊代表彥賢有問題....)
					r.Logger.Error("CreateTransferRefundLeave logic error", zap.Error(fmt.Errorf("UpdateNonFixedChargeRecord have some error")), zap.String("fieldName", "createTransferRefundLeave"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
					return fmt.Errorf("CreateTransferRefundLeave have some error")
				}
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreateTransferRefundLeave tx.Transaction", zap.Error(err), zap.String("fieldName", "createTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createTransferRefundLeave run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createTransferRefundLeave"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateTransferRefundLeave(ctx context.Context, transferRefundLeaveIdStr string, input gqlmodels.TransferRefundLeaveInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdateTransferRefundLeave uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateTransferRefundLeave uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	transferRefundLeaveId, err := uuid.Parse(transferRefundLeaveIdStr)
	if err != nil {
		r.Logger.Warn("UpdateTransferRefundLeave uuid.Parse(transferRefundLeaveIdStr)", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	items, err := json.Marshal(input.Items)
	if err != nil {
		r.Logger.Error("UpdateTransferRefundLeave json.Marshal(input.Items)", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	var note string
	if input.Note != nil {
		note = *input.Note
	}

	transferRefundLeave := models.TransferRefundLeave{
		ID:             transferRefundLeaveId,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Reason:         input.Reason,
		IsReserveBed:   input.IsReserveBed,
		Note:           note,
		Items:          items,
		OrganizationId: organizationId,
		UserId:         userId,
	}

	// 要先抓更新前的時間,後面檢查住民帳單需要
	beforeUpdatedTransferRefundLeave, err := orm.GetTransferRefundLeaveById(r.ORM.DB, organizationId, transferRefundLeaveId)
	if err != nil {
		r.Logger.Error("UpdateTransferRefundLeave orm.GetTransferRefundLeaveById", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		_, err = orm.UpdateTransferRefundLeave(tx, &transferRefundLeave)
		if err != nil {
			r.Logger.Error("UpdateTransferRefundLeave orm.UpdateTransferRefundLeave", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}

		// 價差的規則 刪除只要把舊的扣掉就好 如果是更新的話 要把舊的扣掉之外 還要再加上新的
		var priceSpread int
		// 這個是用來算住民帳單的應繳金額用的 會等住民帳單的關聯都拉完後 最後再跑一個回圈 算有影響到的住民帳單
		needUpdateAmountDuePatientBills := make(map[uuid.UUID]*models.PatientBill)
		// 以下都在更新住民帳單
		// 先看有沒有需要刪除的資料(用更新前的結束時間去看全部的住民帳單)
		needDeletePatientBills := orm.GetPatientBillsByTransferRefundLeaveDate(tx, beforeUpdatedTransferRefundLeave.PatientId, beforeUpdatedTransferRefundLeave.EndDate)
		if len(needDeletePatientBills) != 0 {
			// 先把舊的扣掉
			transferRefundItems := []TransferRefundItem{}
			err = json.Unmarshal(beforeUpdatedTransferRefundLeave.Items, &transferRefundItems)
			if err != nil {
				r.Logger.Error("UpdateTransferRefundLeave json.Unmarshal", zap.Error(err), zap.String("originalUrl", "updateTransferRefundLeave"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return err
			}
			for j := range transferRefundItems {
				// 計算這次新增的總金額(目的是要算應繳金額)
				if transferRefundItems[j].Type == "charge" {
					priceSpread -= transferRefundItems[j].Price
				} else {
					priceSpread += transferRefundItems[j].Price
				}
			}
			for i := range needDeletePatientBills {
				for j := range needDeletePatientBills[i].TransferRefundLeaves {
					// 找到id一樣的後 把他remove後 更新資料庫
					if transferRefundLeaveId == needDeletePatientBills[i].TransferRefundLeaves[j].ID {
						err = orm.DeleteAssociationsPatientBillTransferRefundLeave(tx, needDeletePatientBills[i], transferRefundLeave)
						if err != nil {
							r.Logger.Error("UpdateTransferRefundLeave orm.DeleteAssociationsPatientBillTransferRefundLeave", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
								zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
							return err
						}
						needUpdateAmountDuePatientBills[needDeletePatientBills[i].ID] = needDeletePatientBills[i]
						needUpdateAmountDuePatientBills[needDeletePatientBills[i].ID].AmountDue = needDeletePatientBills[i].AmountDue + priceSpread
					}
				}
			}
		}
		// 先看有沒有需要新增的資料(用新的結束時間去看全部的住民帳單)
		// 在區間:如果沒資料也不用比對了 代表一定是新的資料直接新增
		// 如果有數筆就要檢查是不是原本就有了 沒有的話要新增
		// 不在區間:邏輯上不會到這邊,這個情況是(有跑到這邊代表彥賢有問題....)
		patientBills := orm.GetPatientBillsByTransferRefundLeaveDate(tx, beforeUpdatedTransferRefundLeave.PatientId, input.EndDate)
		// 計算這次新增的總金額(目的是要算應繳金額)
		for i := range input.Items {
			if input.Items[i].Type == "charge" {
				priceSpread += input.Items[i].Price
			} else {
				priceSpread -= input.Items[i].Price
			}
		}
		for i := range patientBills {
			if input.EndDate.Unix() >= patientBills[i].TransferRefundStartDate.Unix() && input.EndDate.Unix() <= patientBills[i].TransferRefundEndDate.Unix() {
				if len(patientBills[i].TransferRefundLeaves) == 0 {
					err = orm.AppendAssociationsPatientBillTransferRefundLeave(tx, patientBills[i], transferRefundLeave)
					if err != nil {
						r.Logger.Error("UpdateTransferRefundLeave orm.AppendAssociationsPatientBillTransferRefundLeave", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
							zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
						return err
					}
					needUpdateAmountDuePatientBills[patientBills[i].ID] = patientBills[i]
					needUpdateAmountDuePatientBills[patientBills[i].ID].AmountDue = patientBills[i].AmountDue + priceSpread
				} else {
					for j := range patientBills[i].TransferRefundLeaves {
						// 表示這次更新的要新加進來
						if transferRefundLeaveId != patientBills[i].TransferRefundLeaves[j].ID {
							err = orm.AppendAssociationsPatientBillTransferRefundLeave(tx, patientBills[i], transferRefundLeave)
							if err != nil {
								r.Logger.Error("UpdateTransferRefundLeave orm.AppendAssociationsPatientBillTransferRefundLeave", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
									zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
								return err
							}
							needUpdateAmountDuePatientBills[patientBills[i].ID] = patientBills[i]
							needUpdateAmountDuePatientBills[patientBills[i].ID].AmountDue = patientBills[i].AmountDue + priceSpread
						}
					}
				}
			} else {
				// 不在區間:邏輯上不會到這邊,這個情況是(有跑到這邊代表彥賢有問題....)
				r.Logger.Error("UpdateTransferRefundLeave logic error", zap.Error(fmt.Errorf("UpdateTransferRefundLeave have some error")), zap.String("fieldName", "updateTransferRefundLeave"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return fmt.Errorf("UpdateTransferRefundLeave have some error")
			}
		}
		for i := range needUpdateAmountDuePatientBills {
			updatePatientBillAmountDueStruct := UpdatePatientBillAmountDueStruct{
				PatientBill:  needUpdateAmountDuePatientBills[i],
				NewAmountDue: needUpdateAmountDuePatientBills[i].AmountDue,
				Tx:           tx,
			}
			err = UpdatePatientBillAmountDue(updatePatientBillAmountDueStruct)
			if err != nil {
				r.Logger.Error("UpdateTransferRefundLeave UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "updateTransferRefundLeave"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdateTransferRefundLeave tx.Transaction", zap.Error(err), zap.String("fieldName", "updateTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateTransferRefundLeave run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateTransferRefundLeave"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeleteTransferRefundLeave(ctx context.Context, transferRefundLeaveIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeleteTransferRefundLeave uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "deleteTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	transferRefundLeaveId, err := uuid.Parse(transferRefundLeaveIdStr)
	if err != nil {
		r.Logger.Warn("DeleteTransferRefundLeave uuid.Parse(transferRefundLeaveIdStr)", zap.Error(err), zap.String("fieldName", "deleteTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	transferRefundLeave := models.TransferRefundLeave{
		ID:             transferRefundLeaveId,
		OrganizationId: organizationId,
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdatedTransferRefundLeave, err := orm.GetTransferRefundLeaveById(tx, organizationId, transferRefundLeaveId)
		if err != nil {
			r.Logger.Error("DeleteTransferRefundLeave orm.GetTransferRefundLeaveById", zap.Error(err), zap.String("fieldName", "deleteTransferRefundLeave"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needDeletePatientBills := orm.GetPatientBillsByTransferRefundLeaveDate(tx, beforeUpdatedTransferRefundLeave.PatientId, beforeUpdatedTransferRefundLeave.EndDate)
		// 這邊是在計算 該住民帳單底下的所有退費金額
		var priceSpread int
		transferRefundItems := []TransferRefundItem{}
		err = json.Unmarshal(beforeUpdatedTransferRefundLeave.Items, &transferRefundItems)
		if err != nil {
			r.Logger.Error("DeleteTransferRefundLeave json.Unmarshal", zap.Error(err), zap.String("originalUrl", "deleteTransferRefundLeave"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		for j := range transferRefundItems {
			// 計算這次新增的總金額(目的是要算應繳金額)
			if transferRefundItems[j].Type == "charge" {
				priceSpread -= transferRefundItems[j].Price
			} else {
				priceSpread += transferRefundItems[j].Price
			}
		}
		// 更新住民帳單狀態(紀錄成應繳金額不是最新)
		for i := range needDeletePatientBills {
			updatePatientBillAmountDue := UpdatePatientBillAmountDueStruct{
				PatientBill:  needDeletePatientBills[i],
				NewAmountDue: needDeletePatientBills[i].AmountDue + priceSpread,
				Tx:           tx,
			}

			// 更新應繳金額
			err = UpdatePatientBillAmountDue(updatePatientBillAmountDue)
			if err != nil {
				r.Logger.Error("DeleteTransferRefundLeave UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "deleteTransferRefundLeave"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return err
			}
		}
		err = orm.DeleteTransferRefundLeave(tx, &transferRefundLeave)
		if err != nil {
			r.Logger.Error("DeleteTransferRefundLeave orm.DeleteTransferRefundLeave", zap.Error(err), zap.String("fieldName", "deleteTransferRefundLeave"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("DeleteTransferRefundLeave tx.Transaction", zap.Error(err), zap.String("originalUrl", "deleteTransferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteTransferRefundLeave run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteTransferRefundLeave"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) TransferRefundLeave(ctx context.Context, transferRefundLeaveIdStr string) (*models.TransferRefundLeave, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("TransferRefundLeave uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "transferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	transferRefundLeaveId, err := uuid.Parse(transferRefundLeaveIdStr)
	if err != nil {
		r.Logger.Warn("TransferRefundLeave uuid.Parse(transferRefundLeaveIdStr)", zap.Error(err), zap.String("originalUrl", "transferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	transferRefundLeave, err := orm.GetTransferRefundLeaveById(r.ORM.DB, organizationId, transferRefundLeaveId)
	if err != nil {
		r.Logger.Error("TransferRefundLeave orm.GetTransferRefundLeaveById", zap.Error(err), zap.String("originalUrl", "transferRefundLeave"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("transferRefundLeave run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "transferRefundLeave"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return transferRefundLeave, nil
}

func (r *queryResolver) TransferRefundLeaves(ctx context.Context, patientIdStr string, startDate, endDate time.Time) ([]*models.TransferRefundLeave, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("TransferRefundLeaves uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "transferRefundLeaves"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("TransferRefundLeaves uuid.Parse(patientIdStr)", zap.Error(err), zap.String("originalUrl", "transferRefundLeaves"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdAndDate(r.ORM.DB, organizationId, patientId, startDate, endDate)
	if err != nil {
		r.Logger.Error("TransferRefundLeaves orm.GetTransferRefundLeavesByPatientIdAndDate", zap.Error(err), zap.String("originalUrl", "transferRefundLeaves"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("transferRefundLeaves run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "transferRefundLeaves"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return transferRefundLeaves, nil
}

// transferRefundLeave resolvers
type transferRefundLeaveResolver struct{ *Resolver }

func (r *transferRefundLeaveResolver) ID(ctx context.Context, obj *models.TransferRefundLeave) (string, error) {
	return obj.ID.String(), nil
}
