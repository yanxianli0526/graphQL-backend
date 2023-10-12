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
	"gorm.io/gorm"
)

// Mutations
func (r *mutationResolver) CreateNonFixedChargeRecord(ctx context.Context, patientIdStr string, input []*gqlmodels.NonFixedChargeRecordInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreateNonFixedChargeRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateNonFixedChargeRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("CreateNonFixedChargeRecord uuid.Parse(patientIdStr)", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	// 計算這次新增的總金額(目的是要算應繳金額)
	var itemSubtotal int
	var nonFixedChargeRecords []*models.NonFixedChargeRecord
	for i := range input {
		nonFixedChargeRecord := models.NonFixedChargeRecord{
			ID:                 uuid.New(),
			NonFixedChargeDate: input[i].NonFixedChargeDate,
			ItemCategory:       input[i].ItemCategory,
			ItemName:           input[i].ItemName,
			Type:               input[i].Type,
			Unit:               input[i].Unit,
			Price:              input[i].Price,
			Quantity:           input[i].Quantity,
			Subtotal:           input[i].Subtotal,
			Note:               *input[i].Note,
			TaxType:            input[i].TaxType,
			OrganizationId:     organizationId,
			PatientId:          patientId,
			UserId:             userId,
		}
		if input[i].Type == "charge" {
			itemSubtotal += input[i].Subtotal
		} else {
			itemSubtotal -= input[i].Subtotal
		}
		nonFixedChargeRecords = append(nonFixedChargeRecords, &nonFixedChargeRecord)
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CreateNonFixedChargeRecords(tx, nonFixedChargeRecords)
		if err != nil {
			r.Logger.Error("CreateNonFixedChargeRecord CreateNonFixedChargeRecords", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 以下都在更新住民帳單

		//	如果沒有住帳單:世界和平
		//  如果有:要檢查住民帳單跟非固定的關聯
		patientBills := orm.GetPatientBillsByNonFixedChargeDate(tx, patientId, input[0].NonFixedChargeDate)
		if len(patientBills) == 0 {
			return nil
		} else {
			// 表示已經開帳了 要檢查非固定非用的日期 是不是在住民帳單的非固定的區間
			// 在區間:這次新增的全部加進去
			// 不在區間:邏輯上不會到這邊,這個情況是(有跑到這邊代表彥賢有問題....)
			for i := range patientBills {
				if input[0].NonFixedChargeDate.Unix() >= patientBills[i].NonFixedChargeStartDate.Unix() && input[0].NonFixedChargeDate.Unix() <= patientBills[i].NonFixedChargeEndDate.Unix() {
					err = orm.AppendAssociationsPatientBillNonFixedChargeRecords(tx, patientBills[i], nonFixedChargeRecords)
					if err != nil {
						r.Logger.Error("CreateNonFixedChargeRecord orm.AppendAssociationsPatientBillNonFixedChargeRecords", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
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
						r.Logger.Error("CreateNonFixedChargeRecord UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
							zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
						return err
					}
				} else {
					r.Logger.Error("CreateNonFixedChargeRecord logic error", zap.Error(fmt.Errorf("CreateNonFixedChargeRecord have some error")), zap.String("originalUrl", "createNonFixedChargeRecord"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
					//	邏輯上不會到這邊(有跑到這邊代表彥賢有問題....)
					return fmt.Errorf("CreateNonFixedChargeRecord have some error")
				}
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreateNonFixedChargeRecord tx.Transaction", zap.Error(err), zap.String("originalUrl", "createNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createNonFixedChargeRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createNonFixedChargeRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateNonFixedChargeRecord(ctx context.Context, nonFixedChargeRecordIdStr string, input *gqlmodels.NonFixedChargeRecordInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdateNonFixedChargeRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateNonFixedChargeRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	nonFixedChargeRecordId, err := uuid.Parse(nonFixedChargeRecordIdStr)
	if err != nil {
		r.Logger.Warn("UpdateNonFixedChargeRecord uuid.Parse(nonFixedChargeRecordIdStr)", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	updateNonFixedChargeRecord := models.NonFixedChargeRecord{
		ID:                 nonFixedChargeRecordId,
		NonFixedChargeDate: input.NonFixedChargeDate,
		ItemCategory:       input.ItemCategory,
		ItemName:           input.ItemName,
		Type:               input.Type,
		Unit:               input.Unit,
		Price:              input.Price,
		Quantity:           input.Quantity,
		Subtotal:           input.Subtotal,
		Note:               *input.Note,
		TaxType:            input.TaxType,
		OrganizationId:     organizationId,
		UserId:             userId,
	}

	// 要先抓更新前的時間,後面檢查住民帳單需要
	beforeUpdatedNonFixedChargeRecord, err := orm.GetNonFixedChargeRecord(r.ORM.DB, nonFixedChargeRecordId, false, false)
	if err != nil {
		r.Logger.Error("UpdateNonFixedChargeRecord orm.GetNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		updatedNonFixedChargeRecord, err := orm.UpdateAndGetNonFixedChargeRecord(tx, &updateNonFixedChargeRecord)
		if err != nil {
			r.Logger.Error("UpdateNonFixedChargeRecord orm.UpdateAndGetNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 價差的規則 刪除只要把舊的扣掉就好 如果是更新的話 要把舊的扣掉之外 還要再加上新的
		var priceSpread int
		// 這個是用來算住民帳單的應繳金額用的 會等住民帳單的關聯都拉完後 最後再跑一個回圈 算有影響到的住民帳單
		needUpdateAmountDuePatientBills := make(map[uuid.UUID]*models.PatientBill)
		// 以下都在更新住民帳單
		// 先看有沒有需要刪除的資料(用更新前的結束時間去看全部的住民帳單)
		needDeletePatientBills := orm.GetPatientBillsByNonFixedChargeDate(tx, beforeUpdatedNonFixedChargeRecord.PatientId, beforeUpdatedNonFixedChargeRecord.NonFixedChargeDate)
		if len(needDeletePatientBills) != 0 {
			// 先把舊的扣掉
			if beforeUpdatedNonFixedChargeRecord.Type == "charge" {
				priceSpread = -beforeUpdatedNonFixedChargeRecord.Subtotal
			} else {
				priceSpread = beforeUpdatedNonFixedChargeRecord.Subtotal
			}
			for i := range needDeletePatientBills {
				for j := range needDeletePatientBills[i].NonFixedChargeRecords {
					// 除了id一樣之外
					if nonFixedChargeRecordId == needDeletePatientBills[i].NonFixedChargeRecords[j].ID {
						err = orm.DeleteAssociationsPatientBillNonFixedChargeRecord(tx, needDeletePatientBills[i], updateNonFixedChargeRecord)
						if err != nil {
							r.Logger.Error("UpdateNonFixedChargeRecord orm.DeleteAssociationsPatientBillNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
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
		patientBills := orm.GetPatientBillsByNonFixedChargeDate(tx, beforeUpdatedNonFixedChargeRecord.PatientId, input.NonFixedChargeDate)
		// 上面已經把舊的扣掉了 這邊還要再加上新的
		if input.Type == "charge" {
			priceSpread += input.Subtotal
		} else {
			priceSpread -= input.Subtotal
		}
		for i := range patientBills {
			if input.NonFixedChargeDate.Unix() >= patientBills[i].NonFixedChargeStartDate.Unix() && input.NonFixedChargeDate.Unix() <= patientBills[i].NonFixedChargeEndDate.Unix() {
				// 這裡原本一直有一個bug就是 如果住民帳單沒有非固定的話 更新會沒辦法被append
				if len(patientBills[i].NonFixedChargeRecords) == 0 {
					err = orm.AppendAssociationsPatientBillNonFixedChargeRecord(tx, patientBills[i], *updatedNonFixedChargeRecord)
					if err != nil {
						r.Logger.Error("UpdateNonFixedChargeRecord orm.AppendAssociationsPatientBillNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
							zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
						return err
					}
					needUpdateAmountDuePatientBills[patientBills[i].ID] = patientBills[i]
					fmt.Println("patientBills[i].AmountDue", patientBills[i].AmountDue)
					fmt.Println("priceSpread", priceSpread)

					needUpdateAmountDuePatientBills[patientBills[i].ID].AmountDue = patientBills[i].AmountDue + priceSpread
				} else {
					for j := range patientBills[i].NonFixedChargeRecords {
						// 表示這次更新的要新加進來
						if nonFixedChargeRecordId != patientBills[i].NonFixedChargeRecords[j].ID {
							err = orm.AppendAssociationsPatientBillNonFixedChargeRecord(tx, patientBills[i], *updatedNonFixedChargeRecord)
							if err != nil {
								r.Logger.Error("UpdateNonFixedChargeRecord orm.AppendAssociationsPatientBillNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
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
				r.Logger.Error("UpdateNonFixedChargeRecord logic error", zap.Error(fmt.Errorf("UpdateNonFixedChargeRecord have some error")), zap.String("originalUrl", "updateNonFixedChargeRecord"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return fmt.Errorf("UpdateNonFixedChargeRecord have some error")
			}
		}
		for i := range needUpdateAmountDuePatientBills {
			updatePatientBillAmountDueStruct := UpdatePatientBillAmountDueStruct{
				PatientBill:  needUpdateAmountDuePatientBills[i],
				NewAmountDue: needUpdateAmountDuePatientBills[i].AmountDue,
				Tx:           tx,
			}
			fmt.Println("updatePatientBillAmountDueStruct", updatePatientBillAmountDueStruct)
			err = UpdatePatientBillAmountDue(updatePatientBillAmountDueStruct)
			if err != nil {
				r.Logger.Error("UpdateNonFixedChargeRecord UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdateNonFixedChargeRecord tx.Transaction", zap.Error(err), zap.String("originalUrl", "updateNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, nil
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateNonFixedChargeRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateNonFixedChargeRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeleteNonFixedChargeRecord(ctx context.Context, nonFixedChargeRecordIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeleteNonFixedChargeRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "deleteNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	nonFixedChargeRecordId, err := uuid.Parse(nonFixedChargeRecordIdStr)
	if err != nil {
		r.Logger.Warn("DeleteNonFixedChargeRecord uuid.Parse(nonFixedChargeRecordIdStr)", zap.Error(err), zap.String("originalUrl", "deleteNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 要先抓更新前的時間,後面檢查住民帳單需要
		beforeUpdatedNonFixedCharge, err := orm.GetNonFixedChargeRecord(tx, nonFixedChargeRecordId, false, false)
		if err != nil {
			r.Logger.Error("DeleteNonFixedChargeRecord orm.GetNonFixedChargeRecord", zap.Error(err), zap.String("fieldName", "deleteNonFixedChargeRecord"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}
		// 找出有使用到那個區間的住民帳單
		needDeletePatientBills := orm.GetPatientBillsByNonFixedChargeDate(tx, beforeUpdatedNonFixedCharge.PatientId, beforeUpdatedNonFixedCharge.NonFixedChargeDate)
		var priceSpread int
		if beforeUpdatedNonFixedCharge.Type == "charge" {
			priceSpread = -beforeUpdatedNonFixedCharge.Subtotal
		} else {
			priceSpread = beforeUpdatedNonFixedCharge.Subtotal
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
				r.Logger.Error("DeleteNonFixedChargeRecord UpdatePatientBillAmountDue", zap.Error(err), zap.String("originalUrl", "deleteNonFixedChargeRecord"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
				return err
			}
		}
		err = orm.DeleteNonFixedChargeRecord(tx, organizationId, nonFixedChargeRecordId)
		if err != nil {
			r.Logger.Error("DeleteNonFixedChargeRecord orm.DeleteNonFixedChargeRecord", zap.Error(err), zap.String("originalUrl", "deleteNonFixedChargeRecord"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return err
		}

		return nil
	})
	if err != nil {
		r.Logger.Error("DeleteNonFixedChargeRecord tx.Transaction", zap.Error(err), zap.String("originalUrl", "deleteNonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteNonFixedChargeRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteNonFixedChargeRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) NonFixedChargeRecord(ctx context.Context, nonFixedChargeRecordIdStr string) (*models.NonFixedChargeRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	nonFixedChargeRecordId, err := uuid.Parse(nonFixedChargeRecordIdStr)
	if err != nil {
		r.Logger.Warn("NonFixedChargeRecord uuid.Parse(nonFixedChargeRecordIdStr)", zap.Error(err), zap.String("fieldName", "nonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	nonFixedChargeRecord, err := orm.GetNonFixedChargeRecord(r.ORM.DB, nonFixedChargeRecordId, true, true)
	if err != nil {
		r.Logger.Error("NonFixedChargeRecord orm.GetNonFixedChargeRecord", zap.Error(err), zap.String("fieldName", "nonFixedChargeRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("nonFixedChargeRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "nonFixedChargeRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return nonFixedChargeRecord, nil
}

func (r *queryResolver) NonFixedChargeRecords(ctx context.Context, patientIdStr string, startDate, endDate time.Time) ([]*models.NonFixedChargeRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("NonFixedChargeRecords uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "nonFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(r.ORM.DB, patientId, startDate, endDate)
	if err != nil {
		r.Logger.Error("NonFixedChargeRecords orm.GetNonFixedChargeRecordsByPatientIdAndDate", zap.Error(err), zap.String("fieldName", "nonFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("nonFixedChargeRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "nonFixedChargeRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return nonFixedChargeRecords, nil
}

func (r *queryResolver) PatientLatestNonFixedChargeRecords(ctx context.Context) (*gqlmodels.PatientLatestNonFixedChargeRecords, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientLatestNonFixedChargeRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "patientLatestNonFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	// 這邊要撈本月總金額跟上的月的資料 還有最後一次更新 所以要回傳一堆
	thisMonthNonFixedChargeRecords, lastMonthNonFixedChargeRecords, nonFixedChargeRecordsDescByUpdatedAt, err := orm.PatientLatestNonFixedChargeRecords(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PatientLatestNonFixedChargeRecords orm.PatientLatestNonFixedChargeRecords", zap.Error(err), zap.String("fieldName", "patientLatestNonFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	elements := make(map[uuid.UUID]*gqlmodels.PatientLatestNonFixedChargeRecords)
	// 塞本月
	for i := range thisMonthNonFixedChargeRecords {
		if elements[thisMonthNonFixedChargeRecords[i].PatientId] == nil {
			elements[thisMonthNonFixedChargeRecords[i].PatientId] = &gqlmodels.PatientLatestNonFixedChargeRecords{}
			elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth = []*models.NonFixedChargeRecord{}
			elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth = append(elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth, &models.NonFixedChargeRecord{ID: uuid.New()})
			elements[thisMonthNonFixedChargeRecords[i].PatientId].LastMonth = []*models.NonFixedChargeRecord{}
			elements[thisMonthNonFixedChargeRecords[i].PatientId].LastMonth = append(elements[thisMonthNonFixedChargeRecords[i].PatientId].LastMonth, &models.NonFixedChargeRecord{ID: uuid.New()})
		}
		if thisMonthNonFixedChargeRecords[i].Type == "charge" {
			elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth[0].Subtotal += thisMonthNonFixedChargeRecords[i].Subtotal
		} else {
			elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth[0].Subtotal += -thisMonthNonFixedChargeRecords[i].Subtotal
		}
		elements[thisMonthNonFixedChargeRecords[i].PatientId].ThisMonth[0].Patient = thisMonthNonFixedChargeRecords[i].Patient
	}

	// 塞上月
	for i := range lastMonthNonFixedChargeRecords {
		if elements[lastMonthNonFixedChargeRecords[i].PatientId] == nil {
			elements[lastMonthNonFixedChargeRecords[i].PatientId] = &gqlmodels.PatientLatestNonFixedChargeRecords{}
			elements[lastMonthNonFixedChargeRecords[i].PatientId].ThisMonth = []*models.NonFixedChargeRecord{}
			elements[lastMonthNonFixedChargeRecords[i].PatientId].ThisMonth = append(elements[lastMonthNonFixedChargeRecords[i].PatientId].ThisMonth, &models.NonFixedChargeRecord{ID: uuid.New()})
			elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth = []*models.NonFixedChargeRecord{}
			elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth = append(elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth, &models.NonFixedChargeRecord{ID: uuid.New()})
		}
		if lastMonthNonFixedChargeRecords[i].Type == "charge" {
			elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth[0].Subtotal += lastMonthNonFixedChargeRecords[i].Subtotal
		} else {
			elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth[0].Subtotal += -lastMonthNonFixedChargeRecords[i].Subtotal
		}
		elements[lastMonthNonFixedChargeRecords[i].PatientId].LastMonth[0].Patient = lastMonthNonFixedChargeRecords[i].Patient

	}

	// 把本月和上月的資料整成合法的type
	var newThisMonthNonFixedChargeRecords []*models.NonFixedChargeRecord
	var newLastMonthNonFixedChargeRecords []*models.NonFixedChargeRecord
	for i := range elements {
		if elements[i].ThisMonth[0].Subtotal != 0 {
			newThisMonthNonFixedChargeRecord := models.NonFixedChargeRecord{
				ID:       elements[i].ThisMonth[0].ID,
				Subtotal: elements[i].ThisMonth[0].Subtotal,
				Patient:  elements[i].ThisMonth[0].Patient,
			}
			newThisMonthNonFixedChargeRecords = append(newThisMonthNonFixedChargeRecords, &newThisMonthNonFixedChargeRecord)
		}
		if elements[i].LastMonth[0].Subtotal != 0 {
			newLastMonthNonFixedChargeRecord := models.NonFixedChargeRecord{
				ID:       elements[i].LastMonth[0].ID,
				Subtotal: elements[i].LastMonth[0].Subtotal,
				Patient:  elements[i].LastMonth[0].Patient,
			}
			newLastMonthNonFixedChargeRecords = append(newLastMonthNonFixedChargeRecords, &newLastMonthNonFixedChargeRecord)
		}
	}

	patientLatestNonFixedChargeRecords := gqlmodels.PatientLatestNonFixedChargeRecords{
		ThisMonth:           newThisMonthNonFixedChargeRecords,
		LastMonth:           newLastMonthNonFixedChargeRecords,
		LatestUpdatedRecord: nonFixedChargeRecordsDescByUpdatedAt,
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientLatestNonFixedChargeRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientLatestNonFixedChargeRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return &patientLatestNonFixedChargeRecords, nil
}

// nonFixedCharge resolvers
type nonFixedChargeRecordResolver struct{ *Resolver }

func (r *nonFixedChargeRecordResolver) ID(ctx context.Context, obj *models.NonFixedChargeRecord) (string, error) {

	return obj.ID.String(), nil
}
