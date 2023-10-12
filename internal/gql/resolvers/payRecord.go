package resolvers

import (
	"context"
	_ "embed"
	"encoding/json"
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

type PayRecordBasicCharge struct {
	ItemName  string    `json:"itemName"`
	Type      string    `json:"type"`
	TaxType   string    `json:"taxType"`
	Unit      string    `json:"unit"`
	Price     int       `json:"price"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type PayRecordSubsidy struct {
	ItemName  string    `json:"itemName"`
	Type      string    `json:"type"`
	Unit      string    `json:"unit"`
	Price     int       `json:"price"`
	IdNumber  string    `json:"idNumber"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type PayRecordTransferRefundLeave struct {
	ItemName  string    `json:"itemName"`
	Type      string    `json:"type"`
	Price     int       `json:"price"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type PayRecordNonFixedChargeRecord struct {
	ItemCategory       string    `json:"itemCategory"`
	ItemName           string    `json:"itemName"`
	Type               string    `json:"type"`
	TaxType            string    `json:"taxType"`
	Unit               string    `json:"unit"`
	NonFixedChargeDate time.Time `json:"nonFixedChargeDate"`
	Quantity           int       `json:"quantity"`
	Price              int       `json:"price"`
	Subtotal           int       `json:"subtotal"`
}

type PayRecordsData struct {
	PayRecord                           []*models.PayRecord `json:"payRecord"`
	PayRecordReceiptNumberCount         int                 `json:"payRecordReceiptNumberCount"`         // 用來記錄資料的流水號(因為有多個住民和稅別 所以每次確定有資料時要累加)
	NeedUpdateBasicChargesIdStr         []string            `json:"needUpdateBasicChargesIdStr"`         // 用來記錄需要關帳(開帳)的基本月費id
	NeedUpdateSubsidiesIdStr            []string            `json:"needUpdateSubsidiesIdStr"`            // 用來記錄需要關帳(開帳)的補助款id
	NeedUpdateTransferRefundLeavesIdStr []string            `json:"needUpdateTransferRefundLeavesIdStr"` // 用來記錄需要關帳(開帳)的請假id
	NeedUpdateNonFixedChargesIdStr      []string            `json:"needUpdateNonFixedChargesIdStr"`      // 用來記錄需要關帳(開帳)的非固定id
}

// Mutations
func (r *mutationResolver) CreatePayRecords(ctx context.Context, input gqlmodels.PayRecordInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreatePayRecords uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "createPayRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreatePayRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "createPayRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var patientsId []uuid.UUID
	for i := range input.PatientsID {
		patientId, err := uuid.Parse(input.PatientsID[i])
		if err != nil {
			r.Logger.Warn("CreatePayRecords uuid.Parse(input.PatientsID[i])", zap.Error(err), zap.String("fieldName", "createPayRecords"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, err
		}
		patientsId = append(patientsId, patientId)
	}

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		r.Logger.Error("CreatePayRecords time.LoadLocation", zap.Error(err), zap.String("fieldName", "createPayRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	transferDate := input.PayDate.In(taipeiZone)
	payYear := transferDate.Year()
	payMonth := int(transferDate.Month())

	payRecordCount, err := orm.GetPayRecordCount(r.ORM.DB, organizationId, payYear, payMonth)
	if err != nil {
		r.Logger.Error("CreatePayRecords orm.GetPayRecordCount", zap.Error(err), zap.String("fieldName", "createPayRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	organizationReceipt, err := orm.GetOrganizationReceiptById(r.ORM.DB, organizationId)
	if err != nil {
		return false, err
	}

	var receiptNumber string // 取機構的流水號規則
	// 最先檢查是不是要重啟流水號
	if !organizationReceipt.IsResetInNextCycle {
		var yearText string
		// 西元年規則
		if organizationReceipt.Year == "Christian" {
			yearText = strconv.Itoa(payYear)
		} else if organizationReceipt.Year == "Republican" {
			// 民國年規則
			yearText = strconv.Itoa(payYear - 1911)
		} else {
			yearText = ""
		}
		var monthText string
		if organizationReceipt.Month == "MM" {
			// MM(小於10 要補0)
			if payMonth < 10 {
				monthText = "0" + strconv.Itoa(payMonth)
			} else {
				monthText = strconv.Itoa(payMonth)
			}
		} else if organizationReceipt.Month == "M" {
			monthText = strconv.Itoa(payMonth)
		} else {
			monthText = ""
		}
		var receiptNumberStr string
		payRecordCount, err = orm.GetPayRecordCountByOrganizationId(r.ORM.DB, organizationId)
		if err != nil {
			return false, err
		}
		if payRecordCount == 0 {
			receiptNumberStr = "000000"
		} else {
			lastestPayRecord, err := orm.GetLastestPayRecordByOrganizationId(r.ORM.DB, organizationId)
			if err != nil {
				return false, err
			}
			receiptNumberStr = lastestPayRecord.ReceiptNumber[len(lastestPayRecord.ReceiptNumber)-6:]
		}
		receiptNumber = organizationReceipt.FirstText + yearText + organizationReceipt.YearText + monthText + organizationReceipt.MonthText + organizationReceipt.LastText + receiptNumberStr
	} else if payRecordCount == 0 && organizationReceipt.IsResetInNextCycle {
		// 先看這個月有沒有資料 沒有的話就套用機構的設定
		// 取機構的流水號規則
		if err != nil {
			return false, err
		}
		var yearText string
		// 西元年規則
		if organizationReceipt.Year == "Christian" {
			yearText = strconv.Itoa(payYear)
		} else if organizationReceipt.Year == "Republican" {
			// 民國年規則
			yearText = strconv.Itoa(payYear - 1911)
		} else {
			yearText = ""
		}
		var monthText string
		if organizationReceipt.Month == "MM" {
			// MM(小於10 要補0)
			if payMonth < 10 {
				monthText = "0" + strconv.Itoa(payMonth)
			} else {
				monthText = strconv.Itoa(payMonth)
			}
		} else if organizationReceipt.Month == "M" {
			monthText = strconv.Itoa(payMonth)
		} else {
			monthText = ""
		}
		receiptNumber = organizationReceipt.FirstText + yearText + organizationReceipt.YearText + monthText + organizationReceipt.MonthText + organizationReceipt.LastText + "000000"
	} else {
		// 已經有資料了 就用之前的流水號規則
		lastestPayRecord, err := orm.GetLastestPayRecord(r.ORM.DB, organizationId, payYear, payMonth)
		if err != nil {
			return false, err
		}

		var yearText string
		// 西元年規則
		if organizationReceipt.Year == "Christian" {
			yearText = strconv.Itoa(payYear)
		} else if organizationReceipt.Year == "Republican" {
			// 民國年規則
			yearText = strconv.Itoa(payYear - 1911)
		} else {
			yearText = ""
		}
		var monthText string
		if organizationReceipt.Month == "MM" {
			// MM(小於10 要補0)
			if payMonth < 10 {
				monthText = "0" + strconv.Itoa(payMonth)
			} else {
				monthText = strconv.Itoa(payMonth)
			}
		} else if organizationReceipt.Month == "M" {
			monthText = strconv.Itoa(payMonth)
		} else {
			monthText = ""
		}
		payRecordData := *lastestPayRecord

		receiptNumber = organizationReceipt.FirstText + yearText + organizationReceipt.YearText + monthText + organizationReceipt.MonthText + organizationReceipt.LastText + payRecordData.ReceiptNumber[len(payRecordData.ReceiptNumber)-6:]
	}

	// 用來記錄資料的流水號(因為有多個住民 所以每次確定有資料時要累加)
	var payRecordsData *PayRecordsData
	// 把用到的稅別加進去taxTypes
	var taxTypes []string
	var haveStampTax bool
	// 所有帳單費用結算
	if input.OpenMethod.String() == "allTax" {
		for i := range patientsId {
			// 先看這個人有沒有做過patientBill
			patientBill, havePatientBill := orm.GetPatientBillByPatientIdAndYearMonthHaveData(r.ORM.DB, organizationId, patientsId[i], payYear, payMonth)
			// 有住民帳單才能做繳費單
			if havePatientBill {
				payRecord := models.PayRecord{
					PatientBillId:  patientBill.ID,
					ReceiptNumber:  receiptNumber,
					PayYear:        payYear,
					PayMonth:       payMonth,
					OrganizationId: organizationId,
					PatientId:      patientsId[i],
					UserId:         userId,
					CreatedUserId:  userId,
				}
				if payRecordsData == nil {
					payRecordsData = &PayRecordsData{}
				}
				// 所有帳單費用結算
				payRecordsData, err = getPayRecordsDataByAllTax(*payRecordsData, *patientBill, payRecord)
				if err != nil {
					r.Logger.Error("CreatePayRecords getPayRecordsDataByAllTax", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return false, nil
				}
			}
		}
	} else {
		for j := range input.TaxTypes {
			if input.TaxTypes[j].String() == "stampTax" {
				haveStampTax = true
			}
			taxTypes = append(taxTypes, input.TaxTypes[j].String())
		}
		// 分稅別結算
		for i := range patientsId {
			// 先看這個人有沒有做過patientBill
			patientBill, havePatientBill := orm.GetPatientBillByPatientIdAndYearMonthHaveData(r.ORM.DB, organizationId, patientsId[i], payYear, payMonth)
			// 有住民帳單才能做繳費單
			if havePatientBill {
				payRecord := models.PayRecord{
					PatientBillId:  patientBill.ID,
					ReceiptNumber:  receiptNumber,
					PayYear:        payYear,
					PayMonth:       payMonth,
					OrganizationId: organizationId,
					PatientId:      patientsId[i],
					UserId:         userId,
					CreatedUserId:  userId,
				}
				if payRecordsData == nil {
					payRecordsData = &PayRecordsData{}
				}

				// 分稅別結算
				payRecordsData, err = getPayRecordsDataBySelectTax(*payRecordsData, *patientBill, taxTypes, payRecord)
				if err != nil {
					r.Logger.Error("CreatePayRecords getPayRecordsDataBySelectTax", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return false, nil
				}
			}
		}
	}

	if payRecordsData != nil {
		// 把需要關帳的項目Id做做型別轉換(很白痴的步驟)
		var needUpdateBasicChargesId []uuid.UUID
		var needUpdateSubsidiesId []uuid.UUID
		var needUpdateTransferRefundLeavesId []uuid.UUID
		var needUpdateNonFixedChargesId []uuid.UUID
		for i := range payRecordsData.NeedUpdateBasicChargesIdStr {
			needUpdateBasicChargeId, err := uuid.Parse(payRecordsData.NeedUpdateBasicChargesIdStr[i])
			if err != nil {
				r.Logger.Error("CreatePayRecords uuid.Parse(payRecordsData.NeedUpdateBasicChargesIdStr[i])", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return false, nil
			}
			needUpdateBasicChargesId = append(needUpdateBasicChargesId, needUpdateBasicChargeId)
		}
		for i := range payRecordsData.NeedUpdateSubsidiesIdStr {
			needUpdateSubsidyId, err := uuid.Parse(payRecordsData.NeedUpdateSubsidiesIdStr[i])
			if err != nil {
				r.Logger.Error("CreatePayRecords uuid.Parse(payRecordsData.NeedUpdateSubsidiesIdStr[i])", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return false, nil
			}
			needUpdateSubsidiesId = append(needUpdateSubsidiesId, needUpdateSubsidyId)
		}
		for i := range payRecordsData.NeedUpdateTransferRefundLeavesIdStr {
			needUpdateTransferRefundLeaveId, err := uuid.Parse(payRecordsData.NeedUpdateTransferRefundLeavesIdStr[i])
			if err != nil {
				r.Logger.Error("CreatePayRecords uuid.Parse(payRecordsData.NeedUpdateTransferRefundLeavesIdStr[i])", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return false, nil
			}
			needUpdateTransferRefundLeavesId = append(needUpdateTransferRefundLeavesId, needUpdateTransferRefundLeaveId)
		}
		for i := range payRecordsData.NeedUpdateNonFixedChargesIdStr {
			needUpdateNonFixedChargeId, err := uuid.Parse(payRecordsData.NeedUpdateNonFixedChargesIdStr[i])
			if err != nil {
				r.Logger.Error("CreatePayRecords uuid.Parse(payRecordsData.NeedUpdateNonFixedChargesIdStr[i])", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return false, nil
			}
			needUpdateNonFixedChargesId = append(needUpdateNonFixedChargesId, needUpdateNonFixedChargeId)
		}

		// 這邊有很多步驟 要確定都正確才執行
		// 先新增繳費記錄 才去更新帳單狀態(關帳)
		tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
		err = tx.Transaction(func(tx *gorm.DB) error {
			err = orm.CreatePayRecords(tx, payRecordsData.PayRecord)
			if err != nil {
				r.Logger.Error("CreatePayRecords orm.CreatePayRecords", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			taipeiZone, err := time.LoadLocation("Asia/Taipei")
			if err != nil {
				r.Logger.Error("CreatePayRecords time.LoadLocation", zap.Error(err), zap.String("fieldName", "createPayRecords"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
			payRecordReceiptDate := time.Now().In(taipeiZone)

			if input.OpenMethod.String() == "allTax" {
				// 基本月費
				err = orm.UpdateBasicChargesReceiptStatus(tx, needUpdateBasicChargesId, "issued", &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords allTax orm.UpdateBasicChargesReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}

				// 補助款
				err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "issued", &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords allTax orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}

				// 住院
				err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "issued", &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords allTax orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}

				// 非固定
				err = orm.UpdateNonFixedChargesRecordsReceiptStatus(tx, needUpdateNonFixedChargesId, "issued", &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords allTax orm.UpdateNonFixedChargesRecordsReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			} else {
				// 基本月費
				err = orm.UpdateBasicChargesReceiptStatusInTaxType(tx, needUpdateBasicChargesId, "issued", taxTypes, &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords orm.UpdateBasicChargesReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
				// 表示有印花稅
				if haveStampTax {
					// 補助款
					err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "issued", &payRecordReceiptDate)
					if err != nil {
						r.Logger.Error("CreatePayRecords orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
							zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
						return err
					}

					// 住院
					err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "issued", &payRecordReceiptDate)
					if err != nil {
						r.Logger.Error("CreatePayRecords orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "createPayRecords"),
							zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
						return err
					}
				}
				// 非固定
				err = orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType(tx, needUpdateNonFixedChargesId, "issued", taxTypes, &payRecordReceiptDate)
				if err != nil {
					r.Logger.Error("CreatePayRecords orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "createPayRecords"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			}
			return nil
		})
		if err != nil {
			r.Logger.Error("CreatePayRecords tx.Transaction", zap.Error(err), zap.String("fieldName", "createPayRecords"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return false, nil
		}
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createPayRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createPayRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// 實際上除了把原本的繳費記錄作廢之外 還會把和住民帳單相關的資料狀態一併改為作廢
func (r *mutationResolver) InvalidPayRecord(ctx context.Context, payRecordIdStr string, input gqlmodels.InvalidPayRecordInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("InvalidPayRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("InvalidPayRecord uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	// 這邊有很多步驟 要確定都正確才執行
	// 先改變繳費記錄的狀態
	// 在查出繳費記錄關聯的個模組id 最後再去更新狀態
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.InvalidPayRecord(tx, payRecordId, input, userId)
		if err != nil {
			r.Logger.Error("InvalidPayRecord orm.InvalidPayRecord", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		payRecord, err := orm.GetPayRecordById(tx, payRecordId, false, false, false)
		if err != nil {
			r.Logger.Error("InvalidPayRecord orm.GetPayRecordById", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		var needUpdateBasicChargesId []uuid.UUID
		var needUpdateSubsidiesId []uuid.UUID
		var needUpdateTransferRefundLeavesId []uuid.UUID
		var needUpdateNonFixedChargesId []uuid.UUID
		patientBill := orm.GetPatientBillByMonth(r.ORM.DB, payRecord.PatientId, payRecord.PayYear, payRecord.PayMonth, true)
		if err != nil {
			r.Logger.Error("InvalidPayRecord orm.GetPatientBillByMonth", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		for i := range patientBill.BasicCharges {
			basicChargesId := patientBill.BasicCharges[i].ID
			needUpdateBasicChargesId = append(needUpdateBasicChargesId, basicChargesId)
		}
		for i := range patientBill.Subsidies {
			subsidyId := patientBill.Subsidies[i].ID
			needUpdateSubsidiesId = append(needUpdateSubsidiesId, subsidyId)
		}
		for i := range patientBill.TransferRefundLeaves {
			transferRefundLeaveId := patientBill.TransferRefundLeaves[i].ID
			needUpdateTransferRefundLeavesId = append(needUpdateTransferRefundLeavesId, transferRefundLeaveId)
		}

		for i := range patientBill.NonFixedChargeRecords {
			nonFixedChargesId := patientBill.NonFixedChargeRecords[i].ID
			needUpdateNonFixedChargesId = append(needUpdateNonFixedChargesId, nonFixedChargesId)
		}

		if payRecord.TaxType == "allTax" {
			// 基本月費
			err = orm.UpdateBasicChargesReceiptStatus(tx, needUpdateBasicChargesId, "invalid", nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord allTax orm.UpdateBasicChargesReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 補助款
			err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "invalid", nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord allTax orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 住院
			err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "invalid", nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord allTax orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 非固定
			err = orm.UpdateNonFixedChargesRecordsReceiptStatus(tx, needUpdateNonFixedChargesId, "invalid", nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord allTax orm.UpdateNonFixedChargesRecordsReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		} else {
			var taxType []string
			taxType = append(taxType, payRecord.TaxType)
			// 基本月費
			err = orm.UpdateBasicChargesReceiptStatusInTaxType(tx, needUpdateBasicChargesId, "invalid", taxType, nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord orm.UpdateBasicChargesReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
			if payRecord.TaxType == "stampTax" {
				// 補助款
				err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "invalid", nil)
				if err != nil {
					r.Logger.Error("InvalidPayRecord orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}

				// 住院
				err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "invalid", nil)
				if err != nil {
					r.Logger.Error("InvalidPayRecord orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			}
			// 非固定
			err = orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType(tx, needUpdateNonFixedChargesId, "invalid", taxType, nil)
			if err != nil {
				r.Logger.Error("InvalidPayRecord orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("InvalidPayRecord tx.Transaction", zap.Error(err), zap.String("fieldName", "invalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("invalidPayRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "invalidPayRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// 實際上除了把原本的繳費記錄取消作廢之外 還會把和住民帳單相關的資料狀態一併改為取消作廢
func (r *mutationResolver) CancelInvalidPayRecord(ctx context.Context, payRecordIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CancelInvalidPayRecord uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("CancelInvalidPayRecord uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	// 這邊有很多步驟 要確定都正確才執行
	// 先改變繳費記錄的狀態
	// 在查出繳費記錄關聯的個模組id 最後再去更新狀態
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = orm.CancelInvalidPayRecord(tx, payRecordId, userId)
		if err != nil {
			r.Logger.Error("CancelInvalidPayRecord orm.CancelInvalidPayRecord", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		payRecord, err := orm.GetPayRecordById(tx, payRecordId, false, false, false)
		if err != nil {
			r.Logger.Error("CancelInvalidPayRecord orm.GetPayRecordById", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		var needUpdateBasicChargesId []uuid.UUID
		var needUpdateSubsidiesId []uuid.UUID
		var needUpdateTransferRefundLeavesId []uuid.UUID
		var needUpdateNonFixedChargesId []uuid.UUID
		patientBill := orm.GetPatientBillByMonth(r.ORM.DB, payRecord.PatientId, payRecord.PayYear, payRecord.PayMonth, true)
		if err != nil {
			r.Logger.Error("CancelInvalidPayRecord orm.GetPatientBillByMonth", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		for i := range patientBill.BasicCharges {
			basicChargesId := patientBill.BasicCharges[i].ID
			needUpdateBasicChargesId = append(needUpdateBasicChargesId, basicChargesId)
		}
		for i := range patientBill.Subsidies {
			subsidyId := patientBill.Subsidies[i].ID
			needUpdateSubsidiesId = append(needUpdateSubsidiesId, subsidyId)
		}
		for i := range patientBill.TransferRefundLeaves {
			transferRefundLeaveId := patientBill.TransferRefundLeaves[i].ID
			needUpdateTransferRefundLeavesId = append(needUpdateTransferRefundLeavesId, transferRefundLeaveId)
		}
		for i := range patientBill.NonFixedChargeRecords {
			nonFixedChargesId := patientBill.NonFixedChargeRecords[i].ID
			needUpdateNonFixedChargesId = append(needUpdateNonFixedChargesId, nonFixedChargesId)
		}

		if payRecord.TaxType == "allTax" {
			// 基本月費
			err = orm.UpdateBasicChargesReceiptStatus(tx, needUpdateBasicChargesId, "cancelInvalid", nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord allTax orm.UpdateBasicChargesReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 補助款
			err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "cancelInvalid", nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord allTax orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 住院
			err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "cancelInvalid", nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord allTax orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 非固定
			err = orm.UpdateNonFixedChargesRecordsReceiptStatus(tx, needUpdateNonFixedChargesId, "cancelInvalid", nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord allTax orm.UpdateNonFixedChargesRecordsReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		} else {
			var taxType []string
			taxType = append(taxType, payRecord.TaxType)
			// 基本月費
			err = orm.UpdateBasicChargesReceiptStatusInTaxType(tx, needUpdateBasicChargesId, "cancelInvalid", taxType, nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord orm.UpdateBasicChargesReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
			if payRecord.TaxType == "stampTax" {
				// 補助款
				err = orm.UpdateSubsidiesReceiptStatus(tx, needUpdateSubsidiesId, "cancelInvalid", nil)
				if err != nil {
					r.Logger.Error("CancelInvalidPayRecord orm.UpdateSubsidiesReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}

				// 住院
				err = orm.UpdateTransferRefundLeavesReceiptStatus(tx, needUpdateTransferRefundLeavesId, "cancelInvalid", nil)
				if err != nil {
					r.Logger.Error("CancelInvalidPayRecord orm.UpdateTransferRefundLeavesReceiptStatus", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
						zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			}
			// 非固定
			err = orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType(tx, needUpdateNonFixedChargesId, "cancelInvalid", taxType, nil)
			if err != nil {
				r.Logger.Error("CancelInvalidPayRecord orm.UpdateNonFixedChargesRecordsReceiptStatusInTaxType", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CancelInvalidPayRecord tx.Transaction", zap.Error(err), zap.String("fieldName", "cancelInvalidPayRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("cancelInvalidPayRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "cancelInvalidPayRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdatePayRecorNote(ctx context.Context, payRecordIdStr, note string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Error("UpdatePayRecorNote uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "updatePayRecorNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePayRecorNote uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("fieldName", "updatePayRecorNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	err = orm.UpdatePayRecordNote(r.ORM.DB, note, payRecordId, userId)
	if err != nil {
		r.Logger.Error("UpdatePayRecorNote orm.UpdatePayRecordNote", zap.Error(err), zap.String("fieldName", "updatePayRecorNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePayRecorNote run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePayRecorNote"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) PayRecords(ctx context.Context, payDate time.Time) ([]*models.PayRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PayRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "payRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		r.Logger.Warn("PayRecords time.LoadLocation", zap.Error(err), zap.String("fieldName", "payRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	transferDate := payDate.In(taipeiZone)

	payRecords, err := orm.GetPayRecords(r.ORM.DB, organizationId, transferDate.Year(), int(transferDate.Month()))
	if err != nil {
		r.Logger.Error("PayRecords orm.GetPayRecords", zap.Error(err), zap.String("fieldName", "payRecords"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("payRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "payRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return payRecords, nil
}

func (r *queryResolver) PayRecord(ctx context.Context, payRecordIdStr string) (*models.PayRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("PayRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "payRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	payRecord, err := orm.GetPayRecordById(r.ORM.DB, payRecordId, false, true, true)
	if err != nil {
		r.Logger.Warn("PayRecord orm.GetPayRecordById", zap.Error(err), zap.String("fieldName", "payRecord"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("payRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "payRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return payRecord, nil
}

// payRecord resolvers
type payRecordResolver struct{ *Resolver }

func (r *payRecordResolver) ID(ctx context.Context, obj *models.PayRecord) (string, error) {
	return obj.ID.String(), nil
}

// 這個function和getPayRecordsDataBySelectTax在做的事情幾乎一樣
// 只差在一個有分稅別一個沒有分(影響到回傳的是不是陣列)
func getPayRecordsDataByAllTax(payRecordsData PayRecordsData, patientBill models.PatientBill, payRecordData models.PayRecord) (*PayRecordsData, error) {
	var havePayRecord bool
	var amountDue int
	payRecord := models.PayRecord{
		PayDate:       time.Time{},
		PatientBillId: payRecordData.PatientBillId,
		// ReceiptNumber:  receiptNumber,
		TaxType:        "allTax",
		AmountDue:      0,
		Note:           "",
		IsInvalid:      false,
		PayYear:        payRecordData.PayYear,
		PayMonth:       payRecordData.PayMonth,
		OrganizationId: payRecordData.OrganizationId,
		PatientId:      payRecordData.PatientId,
		UserId:         payRecordData.UserId,
		CreatedUserId:  payRecordData.CreatedUserId,
	}
	/* 基本月費 */
	// 建立一個基本月費的struct
	payRecordBasicCharges := []PayRecordBasicCharge{}
	// 檢查每一個的稅別是否符合(符合的話就把資料塞進去payRecordBasicCharges)
	for i := range patientBill.BasicCharges {
		if patientBill.BasicCharges[i].ReceiptStatus != "issued" && patientBill.BasicCharges[i].ReceiptStatus != "cancelInvalid" {
			havePayRecord = true
			if patientBill.BasicCharges[i].Type == "charge" {
				amountDue += patientBill.BasicCharges[i].Price
			} else {
				amountDue -= patientBill.BasicCharges[i].Price
			}
			payRecordBasicCharge := PayRecordBasicCharge{
				ItemName:  patientBill.BasicCharges[i].ItemName,
				Type:      patientBill.BasicCharges[i].Type,
				TaxType:   patientBill.BasicCharges[i].TaxType,
				Unit:      patientBill.BasicCharges[i].Unit,
				Price:     patientBill.BasicCharges[i].Price,
				StartDate: patientBill.BasicCharges[i].StartDate,
				EndDate:   patientBill.BasicCharges[i].EndDate,
			}
			payRecordsData.NeedUpdateBasicChargesIdStr = append(payRecordsData.NeedUpdateBasicChargesIdStr, patientBill.BasicCharges[i].ID.String())
			payRecordBasicCharges = append(payRecordBasicCharges, payRecordBasicCharge)
		}
	}
	// payRecordBasicCharges塞完資料後轉成json 再丟給payRecord
	payRecordBasicJson, err := json.Marshal(payRecordBasicCharges)
	if err != nil {
		return nil, err
	}
	payRecord.BasicCharge = payRecordBasicJson
	/* 基本月費 */

	/* 補助款 */
	// 建立一個補助款的struct
	payRecordSubsidies := []PayRecordSubsidy{}
	for i := range patientBill.Subsidies {
		if patientBill.Subsidies[i].ReceiptStatus != "issued" && patientBill.Subsidies[i].ReceiptStatus != "cancelInvalid" {
			// 補助款一率歸在印花稅
			havePayRecord = true
			if patientBill.Subsidies[i].Type == "charge" {
				amountDue += patientBill.Subsidies[i].Price
			} else {
				amountDue -= patientBill.Subsidies[i].Price
			}
			payRecordSubsidy := PayRecordSubsidy{
				ItemName:  patientBill.Subsidies[i].ItemName,
				Type:      patientBill.Subsidies[i].Type,
				Unit:      patientBill.Subsidies[i].Unit,
				Price:     patientBill.Subsidies[i].Price,
				StartDate: patientBill.Subsidies[i].StartDate,
				EndDate:   patientBill.Subsidies[i].EndDate,
			}
			payRecordsData.NeedUpdateSubsidiesIdStr = append(payRecordsData.NeedUpdateSubsidiesIdStr, patientBill.Subsidies[i].ID.String())
			payRecordSubsidies = append(payRecordSubsidies, payRecordSubsidy)
		}
	}
	// payRecordSubsidies塞完資料後轉成json 再丟給payRecord
	payRecordSubsidyJson, err := json.Marshal(payRecordSubsidies)
	if err != nil {
		return nil, err
	}
	payRecord.Subsidy = payRecordSubsidyJson
	/* 補助款 */

	/* 異動(請假) */
	// 建立一個異動(請假)的struct
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	for i := range patientBill.TransferRefundLeaves {
		if patientBill.TransferRefundLeaves[i].ReceiptStatus != "issued" && patientBill.TransferRefundLeaves[i].ReceiptStatus != "cancelInvalid" {
			// 請假一率歸在印花稅
			havePayRecord = true
			transferRefundItems := []TransferRefundItem{}
			json.Unmarshal(patientBill.TransferRefundLeaves[i].Items, &transferRefundItems)
			for j := range transferRefundItems {
				if transferRefundItems[j].Type == "charge" {
					amountDue += transferRefundItems[j].Price
				} else {
					amountDue -= transferRefundItems[j].Price
				}
				payRecordTransferRefundLeave := PayRecordTransferRefundLeave{
					ItemName:  transferRefundItems[j].ItemName,
					Type:      transferRefundItems[j].Type,
					Price:     transferRefundItems[j].Price,
					StartDate: patientBill.TransferRefundLeaves[i].StartDate,
					EndDate:   patientBill.TransferRefundLeaves[i].EndDate,
				}
				payRecordsData.NeedUpdateTransferRefundLeavesIdStr = append(payRecordsData.NeedUpdateTransferRefundLeavesIdStr, patientBill.TransferRefundLeaves[i].ID.String())
				payRecordTransferRefundLeaves = append(payRecordTransferRefundLeaves, payRecordTransferRefundLeave)
			}
		}
	}
	// payRecordTransferRefundLeaves塞完資料後轉成json
	payRecordTransferRefundLeaveJson, err := json.Marshal(payRecordTransferRefundLeaves)
	if err != nil {
		return nil, err
	}
	payRecord.TransferRefundLeave = payRecordTransferRefundLeaveJson
	/* 異動(請假) */

	/* 非固定 */
	// 建立一個非固定的struct
	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecord{}
	for i := range patientBill.NonFixedChargeRecords {
		if patientBill.NonFixedChargeRecords[i].ReceiptStatus != "issued" && patientBill.NonFixedChargeRecords[i].ReceiptStatus != "cancelInvalid" {
			// 請假一率歸在印花稅
			havePayRecord = true
			if patientBill.NonFixedChargeRecords[i].Type == "charge" {
				amountDue += patientBill.NonFixedChargeRecords[i].Subtotal
			} else {
				amountDue -= patientBill.NonFixedChargeRecords[i].Subtotal
			}
			payRecordNonFixedChargeRecord := PayRecordNonFixedChargeRecord{
				ItemCategory:       patientBill.NonFixedChargeRecords[i].ItemCategory,
				ItemName:           patientBill.NonFixedChargeRecords[i].ItemName,
				Type:               patientBill.NonFixedChargeRecords[i].Type,
				TaxType:            patientBill.NonFixedChargeRecords[i].TaxType,
				Unit:               patientBill.NonFixedChargeRecords[i].Unit,
				NonFixedChargeDate: patientBill.NonFixedChargeRecords[i].NonFixedChargeDate,
				Quantity:           patientBill.NonFixedChargeRecords[i].Quantity,
				Price:              patientBill.NonFixedChargeRecords[i].Price,
				Subtotal:           patientBill.NonFixedChargeRecords[i].Subtotal,
			}
			payRecordsData.NeedUpdateNonFixedChargesIdStr = append(payRecordsData.NeedUpdateNonFixedChargesIdStr, patientBill.NonFixedChargeRecords[i].ID.String())
			payRecordNonFixedChargeRecords = append(payRecordNonFixedChargeRecords, payRecordNonFixedChargeRecord)
		}
	}
	// payRecordNonFixedChargeRecords塞完資料後轉成json
	payRecordNonFixedChargeRecordJson, err := json.Marshal(payRecordNonFixedChargeRecords)
	if err != nil {
		return nil, err
	}

	payRecord.NonFixedCharge = payRecordNonFixedChargeRecordJson
	/* 非固定 */
	payRecord.AmountDue = amountDue
	if havePayRecord {
		// 取最後六碼流水號(轉成int)
		receiptNumberInt, err := strconv.Atoi(payRecordData.ReceiptNumber[len(payRecordData.ReceiptNumber)-6:]) // result: i = -18
		if err != nil {
			return nil, err
		}
		// 計算流水號
		receiptNumberStr := strconv.Itoa(receiptNumberInt + 1 + payRecordsData.PayRecordReceiptNumberCount)
		receiptNumberStrLength := len(receiptNumberStr)
		// 補0(ex:1=>000001)
		for i := 0; i < 6-receiptNumberStrLength; i++ {
			receiptNumberStr = "0" + receiptNumberStr
		}
		// 拿前面的文字規則+計算完的流水號
		receiptNumber := payRecordData.ReceiptNumber[0:len(payRecordData.ReceiptNumber)-6] + receiptNumberStr
		payRecord.ReceiptNumber = receiptNumber
		payRecordsData.PayRecordReceiptNumberCount++
		payRecordsData.PayRecord = append(payRecordsData.PayRecord, &payRecord)
	}
	return &payRecordsData, nil
}

func getPayRecordsDataBySelectTax(payRecordsData PayRecordsData, patientBill models.PatientBill, taxTypes []string, payRecordData models.PayRecord) (*PayRecordsData, error) {
	// 用來記錄流水號要加幾用的
	for i := range taxTypes {
		var havePayRecord bool
		var amountDue int
		payRecord := models.PayRecord{
			PayDate:       time.Time{},
			PatientBillId: payRecordData.PatientBillId,
			// ReceiptNumber:  receiptNumber,
			TaxType:        taxTypes[i],
			AmountDue:      0,
			Note:           "",
			IsInvalid:      false,
			PayYear:        payRecordData.PayYear,
			PayMonth:       payRecordData.PayMonth,
			OrganizationId: payRecordData.OrganizationId,
			PatientId:      payRecordData.PatientId,
			UserId:         payRecordData.UserId,
			CreatedUserId:  payRecordData.CreatedUserId,
		}
		/* 基本月費 */
		// 建立一個基本月費的struct
		payRecordBasicCharges := []PayRecordBasicCharge{}
		// 檢查每一個的稅別是否符合(符合的話就把資料塞進去payRecordBasicCharges)
		for j := range patientBill.BasicCharges {
			if patientBill.BasicCharges[j].TaxType == taxTypes[i] && patientBill.BasicCharges[j].ReceiptStatus != "issued" && patientBill.BasicCharges[j].ReceiptStatus != "cancelInvalid" {
				havePayRecord = true
				if patientBill.BasicCharges[j].Type == "charge" {
					amountDue += patientBill.BasicCharges[j].Price
				} else {
					amountDue -= patientBill.BasicCharges[j].Price
				}
				payRecordBasicCharge := PayRecordBasicCharge{
					ItemName:  patientBill.BasicCharges[j].ItemName,
					Type:      patientBill.BasicCharges[j].Type,
					TaxType:   patientBill.BasicCharges[j].TaxType,
					Unit:      patientBill.BasicCharges[j].Unit,
					Price:     patientBill.BasicCharges[j].Price,
					StartDate: patientBill.BasicCharges[j].StartDate,
					EndDate:   patientBill.BasicCharges[j].EndDate,
				}
				payRecordsData.NeedUpdateBasicChargesIdStr = append(payRecordsData.NeedUpdateBasicChargesIdStr, patientBill.BasicCharges[j].ID.String())
				payRecordBasicCharges = append(payRecordBasicCharges, payRecordBasicCharge)
			}
		}
		// payRecordBasicCharges塞完資料後轉成json 再丟給payRecord
		payRecordBasicJson, err := json.Marshal(payRecordBasicCharges)
		if err != nil {
			return nil, err
		}
		payRecord.BasicCharge = payRecordBasicJson
		/* 基本月費 */

		/* 補助款 */
		// 建立一個補助款的struct
		payRecordSubsidies := []PayRecordSubsidy{}
		for j := range patientBill.Subsidies {
			// 補助款一率歸在印花稅
			if taxTypes[i] == "stampTax" && patientBill.Subsidies[j].ReceiptStatus != "issued" && patientBill.Subsidies[j].ReceiptStatus != "cancelInvalid" {

				havePayRecord = true
				if patientBill.Subsidies[j].Type == "charge" {
					amountDue += patientBill.Subsidies[j].Price
				} else {
					amountDue -= patientBill.Subsidies[j].Price
				}
				payRecordSubsidy := PayRecordSubsidy{
					ItemName:  patientBill.Subsidies[j].ItemName,
					Type:      patientBill.Subsidies[j].Type,
					Unit:      patientBill.Subsidies[j].Unit,
					Price:     patientBill.Subsidies[j].Price,
					StartDate: patientBill.Subsidies[j].StartDate,
					EndDate:   patientBill.Subsidies[j].EndDate,
				}
				payRecordsData.NeedUpdateSubsidiesIdStr = append(payRecordsData.NeedUpdateSubsidiesIdStr, patientBill.Subsidies[j].ID.String())
				payRecordSubsidies = append(payRecordSubsidies, payRecordSubsidy)
			}
		}
		// payRecordSubsidies塞完資料後轉成json 再丟給payRecord
		payRecordSubsidyJson, err := json.Marshal(payRecordSubsidies)
		if err != nil {
			return nil, err
		}
		payRecord.Subsidy = payRecordSubsidyJson
		/* 補助款 */

		/* 異動(請假) */
		// 建立一個異動(請假)的struct
		payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
		for j := range patientBill.TransferRefundLeaves {
			// 請假一率歸在印花稅
			if taxTypes[i] == "stampTax" && patientBill.TransferRefundLeaves[j].ReceiptStatus != "issued" && patientBill.TransferRefundLeaves[j].ReceiptStatus != "cancelInvalid" {
				havePayRecord = true
				transferRefundItems := []TransferRefundItem{}
				json.Unmarshal(patientBill.TransferRefundLeaves[j].Items, &transferRefundItems)
				for k := range transferRefundItems {
					if transferRefundItems[k].Type == "charge" {
						amountDue += transferRefundItems[k].Price
					} else {
						amountDue -= transferRefundItems[k].Price
					}
					payRecordTransferRefundLeave := PayRecordTransferRefundLeave{
						ItemName:  transferRefundItems[k].ItemName,
						Type:      transferRefundItems[k].Type,
						Price:     transferRefundItems[k].Price,
						StartDate: patientBill.TransferRefundLeaves[j].StartDate,
						EndDate:   patientBill.TransferRefundLeaves[j].EndDate,
					}
					payRecordsData.NeedUpdateTransferRefundLeavesIdStr = append(payRecordsData.NeedUpdateTransferRefundLeavesIdStr, patientBill.TransferRefundLeaves[j].ID.String())
					payRecordTransferRefundLeaves = append(payRecordTransferRefundLeaves, payRecordTransferRefundLeave)
				}
			}
		}
		// payRecordTransferRefundLeaves塞完資料後轉成json
		payRecordTransferRefundLeaveJson, err := json.Marshal(payRecordTransferRefundLeaves)
		if err != nil {
			return nil, err
		}
		payRecord.TransferRefundLeave = payRecordTransferRefundLeaveJson
		/* 異動(請假) */

		/* 非固定 */
		// 建立一個非固定的struct
		payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecord{}
		for j := range patientBill.NonFixedChargeRecords {
			// 請假一率歸在印花稅

			if patientBill.NonFixedChargeRecords[j].TaxType == taxTypes[i] && patientBill.NonFixedChargeRecords[j].ReceiptStatus != "issued" && patientBill.NonFixedChargeRecords[j].ReceiptStatus != "cancelInvalid" {
				havePayRecord = true
				if patientBill.NonFixedChargeRecords[j].Type == "charge" {
					amountDue += patientBill.NonFixedChargeRecords[j].Subtotal
				} else {
					amountDue -= patientBill.NonFixedChargeRecords[j].Subtotal
				}
				payRecordNonFixedChargeRecord := PayRecordNonFixedChargeRecord{
					ItemCategory:       patientBill.NonFixedChargeRecords[j].ItemCategory,
					ItemName:           patientBill.NonFixedChargeRecords[j].ItemName,
					Type:               patientBill.NonFixedChargeRecords[j].Type,
					TaxType:            patientBill.NonFixedChargeRecords[j].TaxType,
					Unit:               patientBill.NonFixedChargeRecords[j].Unit,
					NonFixedChargeDate: patientBill.NonFixedChargeRecords[j].NonFixedChargeDate,
					Quantity:           patientBill.NonFixedChargeRecords[j].Quantity,
					Price:              patientBill.NonFixedChargeRecords[j].Price,
					Subtotal:           patientBill.NonFixedChargeRecords[j].Subtotal,
				}
				payRecordsData.NeedUpdateNonFixedChargesIdStr = append(payRecordsData.NeedUpdateNonFixedChargesIdStr, patientBill.NonFixedChargeRecords[j].ID.String())
				payRecordNonFixedChargeRecords = append(payRecordNonFixedChargeRecords, payRecordNonFixedChargeRecord)
			}
		}
		// payRecordNonFixedChargeRecords塞完資料後轉成json
		payRecordNonFixedChargeRecordJson, err := json.Marshal(payRecordNonFixedChargeRecords)
		if err != nil {
			return nil, err
		}
		payRecord.NonFixedCharge = payRecordNonFixedChargeRecordJson
		/* 非固定 */
		payRecord.AmountDue = amountDue
		if havePayRecord {
			// 取最後六碼流水號(轉成int)
			receiptNumberInt, err := strconv.Atoi(payRecordData.ReceiptNumber[len(payRecordData.ReceiptNumber)-6:])
			if err != nil {
				return nil, err

			}
			// 計算流水號
			// receiptNumberInt是固定的 上一次create後最新的流水號(如果沒有就是000001)
			// payRecordReceiptNumberCount 是在計算到底有多少個稅別
			// 上一次是000002的話 ex:2+0+1 => 2+1+1 => 2+1+2 => 2+1+2
			receiptNumberStr := strconv.Itoa(receiptNumberInt + 1 + payRecordsData.PayRecordReceiptNumberCount)
			receiptNumberStrLength := len(receiptNumberStr)
			// 補0(ex:1=>000001)
			for i := 0; i < 6-receiptNumberStrLength; i++ {
				receiptNumberStr = "0" + receiptNumberStr
			}
			// 拿前面的文字規則+計算完的流水號
			receiptNumber := payRecordData.ReceiptNumber[0:len(payRecordData.ReceiptNumber)-6] + receiptNumberStr
			payRecord.ReceiptNumber = receiptNumber
			payRecordsData.PayRecordReceiptNumberCount++
			payRecordsData.PayRecord = append(payRecordsData.PayRecord, &payRecord)
		}
	}
	return &payRecordsData, nil
}
