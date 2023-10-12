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
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PatientBillNonFixedChargeRocordDateAndQuantity struct {
	Date     time.Time `json:"date"`
	Quantity int       `json:"quantity"`
}

type PatientBillNonFixedChargeRocord struct {
	ItemCategory    string                                                     `json:"itemCategory"`
	ItemName        string                                                     `json:"itemName"`
	Price           int                                                        `json:"price"`
	Quantity        int                                                        `json:"quantity"`
	Unit            string                                                     `json:"unit"`
	TaxType         string                                                     `json:"TaxType"`
	Type            string                                                     `json:"type"`
	Subtotal        int                                                        `json:"subtotal"`
	EarliestDate    time.Time                                                  `json:"earliestDate"`
	DateAndQuantity map[string]*PatientBillNonFixedChargeRocordDateAndQuantity `json:"newTest"`
}

type PatientBillNonFixedChargeRocordData struct {
	ItemCategory      string                                      `json:"itemCategory"`
	ItemCategoryDatas map[string]*PatientBillNonFixedChargeRocord `json:"itemCategoryDatas"`
}

type PatientBillData struct {
	BasicCharges            []*models.BasicCharge          `json:"basicCharges"`
	Subsidise               []*models.Subsidy              `json:"subsidise"`
	TransferRefundLeaves    []*models.TransferRefundLeave  `json:"transferRefundLeaves"`
	NonFixedChargeRecords   []*models.NonFixedChargeRecord `json:"nonFixedChargeRecords"`
	AmountReceived          int                            `json:"amountReceived"`
	FixedChargeStartDate    time.Time                      `json:"fixedChargeStartDate"`
	FixedChargeEndDate      time.Time                      `json:"fixedChargeEndDate"`
	TransferRefundStartDate time.Time                      `json:"transferRefundStartDate"`
	TransferRefundEndDate   time.Time                      `json:"transferRefundEndDate"`
	NonFixedChargeStartDate time.Time                      `json:"nonFixedChargeStartDate"`
	NonFixedChargeEndDate   time.Time                      `json:"nonFixedChargeEndDate"`
	Note                    string                         `json:"note"`
	BillYear                int                            `json:"billYear"`
	BillMonth               int                            `json:"billMonth"`
	OrganizationId          uuid.UUID                      `json:"organizationId"`
	PatientId               uuid.UUID                      `json:"patientId"`
	UserId                  uuid.UUID                      `json:"userId"`
}

// Mutations
func (r *mutationResolver) CreatePatientBill(ctx context.Context, input *gqlmodels.PatientBillInput) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreatePatientBill uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreatePatientBill uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	patientId, err := uuid.Parse(input.PatientID)
	if err != nil {
		r.Logger.Warn("CreatePatientBill uuid.Parse(input.PatientID)", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	organization, err := orm.GetOrganizationById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("CreatePatientBill orm.GetOrganizationById", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		r.Logger.Error("CreatePatientBill time.LoadLocation", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 帳開的起迄日期要依照機構的設定呈現
	// 找出機構設定固定費用(基本月費及補助款)的時間區間
	fixedChargeStartDate, fixedChargeEndDate, err := getPatientBillDateRange(input.BillDate.In(taipeiZone), organization.FixedChargeStartMonth, organization.FixedChargeEndMonth, organization.FixedChargeStartDate, organization.FixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBill FixedCharge getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 基本月費
	basicChargeSettings, err := orm.GetBasicChargeSettingsByPatientId(r.ORM.DB, organizationId, patientId)
	if err != nil {
		r.Logger.Error("CreatePatientBill orm.GetBasicChargeSettingsByPatientId", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	var basicChargesSortId pq.StringArray
	patientBillBasicCharge := getBasicCharge(basicChargeSettings, *fixedChargeStartDate, *fixedChargeEndDate, patientId, userId)
	basicCharges := make([]*models.BasicCharge, len(patientBillBasicCharge))
	for i := range patientBillBasicCharge {
		basicChargesSortId = append(basicChargesSortId, patientBillBasicCharge[i].ID.String())
		basicCharges[i] = &models.BasicCharge{
			ID: patientBillBasicCharge[i].ID,
		}
	}

	// 補助款
	subsidiesSetting, err := orm.GetSubsidiesSettingByPatientId(r.ORM.DB, organizationId, patientId)
	if err != nil {
		r.Logger.Error("CreatePatientBill orm.GetSubsidiesSettingByPatientId", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	patientBillSubsidies := getPatientBillSubsidies(subsidiesSetting, *fixedChargeStartDate, *fixedChargeEndDate, userId, patientId, organizationId)
	subsidies := make([]*models.Subsidy, len(patientBillSubsidies))
	for i := range patientBillSubsidies {
		subsidies[i] = &models.Subsidy{
			ID: patientBillSubsidies[i].ID,
		}
	}

	// 異動退費
	// 找出機構設定異動退費的時間區間
	transferRefundStartDate, transferRefundEndDate, err := getPatientBillDateRange(input.BillDate.In(taipeiZone), organization.TransferRefundStartMonth, organization.TransferRefundEndMonth, organization.TransferRefundStartDate, organization.TransferRefundEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBill TransferRefund getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 請假
	transferRefundLeavesData, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(r.ORM.DB, organizationId, patientId, *transferRefundStartDate, *transferRefundEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBill orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	transferRefundLeaves := make([]*models.TransferRefundLeave, len(transferRefundLeavesData))
	for i := range transferRefundLeavesData {
		transferRefundLeaves[i] = &models.TransferRefundLeave{
			ID: transferRefundLeavesData[i].ID,
		}
	}

	// 非固定
	// 找出機構設定非固定費用的時間區間
	nonFixedChargeStartDate, nonFixedChargeEndDate, err := getPatientBillDateRange(input.BillDate.In(taipeiZone), organization.NonFixedChargeStartMonth, organization.NonFixedChargeEndMonth, organization.NonFixedChargeStartDate, organization.NonFixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBill NonFixedCharge getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 用時間區間找非固定費用
	nonFixedCharges, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(r.ORM.DB, patientId, *nonFixedChargeStartDate, *nonFixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBill orm.GetNonFixedChargeRecordsByPatientIdAndDate", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	nonFixedChargeRecords := make([]*models.NonFixedChargeRecord, len(nonFixedCharges))
	for i := range nonFixedCharges {
		nonFixedChargeRecords[i] = &models.NonFixedChargeRecord{
			ID: nonFixedCharges[i].ID,
		}
	}

	// 住民帳單的id
	var createdPatientBillId *uuid.UUID
	// 這邊要確保 基本月費 補助款 新增住民(刪除)帳單都成功
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		if len(patientBillBasicCharge) > 0 {
			// 先新增基本月費 住民帳單才能做關聯
			err = orm.CreateBasicCharges(tx, patientBillBasicCharge)
			if err != nil {
				r.Logger.Error("CreatePatientBill orm.CreateBasicCharges", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		if len(patientBillSubsidies) > 0 {
			// 先新增補助款 住民帳單才能做關聯
			err = orm.CreateSubsidies(tx, patientBillSubsidies)
			if err != nil {
				r.Logger.Error("CreatePatientBill orm.CreateSubsidies", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}

		patientBill := models.PatientBill{
			AmountReceived:          0,
			FixedChargeStartDate:    *fixedChargeStartDate,
			FixedChargeEndDate:      *fixedChargeEndDate,
			NonFixedChargeStartDate: *nonFixedChargeStartDate,
			NonFixedChargeEndDate:   *nonFixedChargeEndDate,
			TransferRefundStartDate: *transferRefundStartDate,
			TransferRefundEndDate:   *transferRefundEndDate,
			Note:                    "",
			BillYear:                input.BillDate.In(taipeiZone).Year(),
			BillMonth:               int(input.BillDate.In(taipeiZone).Month()),
			OrganizationId:          organizationId,
			PatientId:               patientId,
			UserId:                  userId,
			BasicCharges:            basicCharges,
			BasicChargesSortIds:     basicChargesSortId,
			Subsidies:               subsidies,
			NonFixedChargeRecords:   nonFixedChargeRecords,
			TransferRefundLeaves:    transferRefundLeaves,
		}

		// 新增前要先看一下這個住民在該月有沒有新增過帳單
		// 如果有查到ID 要把舊的這筆刪掉
		oldPatientBill := orm.GetPatientBillByMonth(tx, patientId, input.BillDate.In(taipeiZone).Year(), int(input.BillDate.In(taipeiZone).Month()), true)
		// 新增新的帳單
		createdPatientBillId, err = orm.CreatePatientBill(tx, &patientBill)
		if err != nil {
			r.Logger.Error("CreatePatientBill orm.CreatePatientBill", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		// 刪掉舊的帳單
		if oldPatientBill != nil {
			var needDeleteBasicChargesID []uuid.UUID
			for i := range oldPatientBill.BasicCharges {
				needDeleteBasicChargesID = append(needDeleteBasicChargesID, oldPatientBill.BasicCharges[i].ID)
			}
			var needDeleteSubsidiesID []uuid.UUID
			for i := range oldPatientBill.Subsidies {
				needDeleteSubsidiesID = append(needDeleteSubsidiesID, oldPatientBill.Subsidies[i].ID)
			}
			err = orm.DeletePatientBillById(tx, oldPatientBill.ID, needDeleteBasicChargesID, needDeleteSubsidiesID)
			if err != nil {
				r.Logger.Error("CreatePatientBill orm.DeletePatientBillById", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreatePatientBill tx.Transaction", zap.Error(err), zap.String("originalUrl", "createPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createPatientBill run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createPatientBill"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return createdPatientBillId.String(), nil
}

func (r *mutationResolver) CreatePatientBills(ctx context.Context, input *gqlmodels.PatientBillsInput) ([]*models.PatientBill, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("CreatePatientBills uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreatePatientBills uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	var patientsId []uuid.UUID
	for i := range input.PatientID {
		patientId, err := uuid.Parse(input.PatientID[i])
		if err != nil {
			r.Logger.Warn("CreatePatientBills uuid.Parse(input.PatientID[i])", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return nil, err
		}
		patientsId = append(patientsId, patientId)
	}

	organization, err := orm.GetOrganizationById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetOrganizationById", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	taipeiZone, _ := time.LoadLocation("Asia/Taipei")

	billDataTimeZoneInTaipei := input.BillDate.UTC().In(taipeiZone)
	// 帳開的起迄日期要依照機構的設定呈現
	// 找出機構設定固定費用的時間區間
	fixedChargeStartDate, fixedChargeEndDate, err := getPatientBillDateRange(billDataTimeZoneInTaipei, organization.FixedChargeStartMonth, organization.FixedChargeEndMonth, organization.FixedChargeStartDate, organization.FixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBills FixedCharge getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	// 找出機構設定異動退費的時間區間
	transferRefundStartDate, transferRefundEndDate, err := getPatientBillDateRange(billDataTimeZoneInTaipei, organization.TransferRefundStartMonth, organization.TransferRefundEndMonth, organization.TransferRefundStartDate, organization.TransferRefundEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBills TransferRefund getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	// 找出機構設定非固定費用的時間區間
	nonFixedChargeStartDate, nonFixedChargeEndDate, err := getPatientBillDateRange(billDataTimeZoneInTaipei, organization.NonFixedChargeStartMonth, organization.NonFixedChargeEndMonth, organization.NonFixedChargeStartDate, organization.NonFixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBills NonFixedCharge getPatientBillDateRange", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	elements := make(map[uuid.UUID]*PatientBillData)

	// 全部的住民都要跑一包patient
	for i := range patientsId {
		elements[patientsId[i]] = &PatientBillData{
			FixedChargeStartDate:    *fixedChargeStartDate,
			FixedChargeEndDate:      *fixedChargeEndDate,
			TransferRefundStartDate: *transferRefundStartDate,
			TransferRefundEndDate:   *transferRefundEndDate,
			NonFixedChargeStartDate: *nonFixedChargeStartDate,
			NonFixedChargeEndDate:   *nonFixedChargeEndDate,
			Note:                    "",
			BillYear:                billDataTimeZoneInTaipei.Year(),
			BillMonth:               int(billDataTimeZoneInTaipei.Month()), OrganizationId: organizationId,
			PatientId: patientsId[i],
			UserId:    userId,
		}
	}

	// 基本月費
	var basicCharges []*models.BasicCharge
	fixedChargeSettings, err := orm.GetBasicChargeSettingsInPatientIds(r.ORM.DB, organizationId, patientsId)
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetBasicChargeSettingsInPatientIds", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	for i := range fixedChargeSettings {
		basicCharge := models.BasicCharge{
			ID:             uuid.New(),
			ItemName:       fixedChargeSettings[i].OrganizationBasicChargeSetting.ItemName,
			Type:           fixedChargeSettings[i].OrganizationBasicChargeSetting.Type,
			Unit:           fixedChargeSettings[i].OrganizationBasicChargeSetting.Unit,
			Price:          fixedChargeSettings[i].OrganizationBasicChargeSetting.Price,
			TaxType:        fixedChargeSettings[i].OrganizationBasicChargeSetting.TaxType,
			StartDate:      *fixedChargeStartDate,
			EndDate:        *fixedChargeEndDate,
			Note:           "",
			ReceiptStatus:  "",
			OrganizationId: organizationId,
			PatientId:      fixedChargeSettings[i].PatientId,
			UserId:         userId,
		}
		basicCharges = append(basicCharges, &basicCharge)
		elements[fixedChargeSettings[i].PatientId].BasicCharges = append(elements[fixedChargeSettings[i].PatientId].BasicCharges, &basicCharge)
	}

	// 補助款
	var subsidies []*models.Subsidy
	subsidiesSetting, err := orm.GetSubsidiesSettingInPatientIds(r.ORM.DB, organizationId, patientsId)
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetSubsidiesSettingInPatientIds", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	for i := range subsidiesSetting {
		subsidy := models.Subsidy{
			ID:             uuid.New(),
			ItemName:       subsidiesSetting[i].ItemName,
			Type:           subsidiesSetting[i].Type,
			Price:          subsidiesSetting[i].Price,
			Unit:           subsidiesSetting[i].Unit,
			IdNumber:       subsidiesSetting[i].IdNumber,
			Note:           "",
			StartDate:      *fixedChargeStartDate,
			EndDate:        *fixedChargeEndDate,
			ReceiptStatus:  "",
			OrganizationId: organizationId,
			PatientId:      subsidiesSetting[i].PatientId,
			UserId:         userId,
		}
		subsidies = append(subsidies, &subsidy)
		elements[subsidiesSetting[i].PatientId].Subsidise = append(elements[subsidiesSetting[i].PatientId].Subsidise, &subsidy)
	}

	// 異動退費
	// 請假
	transferRefundLeavesData, err := orm.GetTransferRefundLeavesInPatientIdsBetweenEndDate(r.ORM.DB, organizationId, patientsId, *transferRefundStartDate, *transferRefundEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetTransferRefundLeavesInPatientIdsBetweenEndDate", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	for i := range transferRefundLeavesData {
		transferRefundLeave := models.TransferRefundLeave{
			ID: transferRefundLeavesData[i].ID,
		}
		elements[transferRefundLeavesData[i].PatientId].TransferRefundLeaves = append(elements[transferRefundLeavesData[i].PatientId].TransferRefundLeaves, &transferRefundLeave)
	}

	// 非固定
	// 用時間區間找非固定費用
	nonFixedChargeSettings, err := orm.GetNonFixedChargeRecordsInPatientIdsAndDate(r.ORM.DB, patientsId, *nonFixedChargeStartDate, *nonFixedChargeEndDate)
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetNonFixedChargeRecordsInPatientIdsAndDate", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	for i := range nonFixedChargeSettings {
		nonFixedChargeRecord := models.NonFixedChargeRecord{
			ID: nonFixedChargeSettings[i].ID,
		}
		elements[nonFixedChargeSettings[i].PatientId].NonFixedChargeRecords = append(elements[nonFixedChargeSettings[i].PatientId].NonFixedChargeRecords, &nonFixedChargeRecord)
	}

	var patientBills []*models.PatientBill
	// 整理全部的patientBill
	for i := range elements {
		patientBill := models.PatientBill{
			AmountReceived:          0,
			Note:                    elements[i].Note,
			FixedChargeStartDate:    elements[i].FixedChargeStartDate,
			FixedChargeEndDate:      elements[i].FixedChargeEndDate,
			TransferRefundStartDate: elements[i].TransferRefundStartDate,
			TransferRefundEndDate:   elements[i].TransferRefundEndDate,
			NonFixedChargeStartDate: elements[i].NonFixedChargeStartDate,
			NonFixedChargeEndDate:   elements[i].NonFixedChargeEndDate,
			BillYear:                elements[i].BillYear,
			BillMonth:               elements[i].BillMonth,
			OrganizationId:          organizationId,
			PatientId:               i,
			UserId:                  userId,
			BasicCharges:            elements[i].BasicCharges,
			Subsidies:               elements[i].Subsidise,
			NonFixedChargeRecords:   elements[i].NonFixedChargeRecords,
			TransferRefundLeaves:    elements[i].TransferRefundLeaves,
		}
		patientBills = append(patientBills, &patientBill)
	}

	// // 新增前要先看一下這個住民在該月有沒有新增過帳單
	// // 如果有查到ID 要把舊的這筆刪掉
	oldPatientBills := orm.GetPatientBillsByMonth(r.ORM.DB, patientsId, billDataTimeZoneInTaipei.Year(), int(billDataTimeZoneInTaipei.Month()))
	var createdPatientBillIds []uuid.UUID
	for i := range oldPatientBills {
		createdPatientBillIds = append(createdPatientBillIds, oldPatientBills[i].ID)
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		if len(basicCharges) > 0 {
			// 先新增基本月費 住民帳單才能做關聯
			err = orm.CreateBasicCharges(tx, basicCharges)
			if err != nil {
				r.Logger.Error("CreatePatientBills orm.CreateBasicCharges", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		if len(subsidies) > 0 {
			// 先新增補助款 住民帳單才能做關聯
			err = orm.CreateSubsidies(tx, subsidies)
			if err != nil {
				r.Logger.Error("CreatePatientBills orm.CreateSubsidies", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		// 新增新的帳單
		err = orm.CreatePatientBills(tx, patientBills)
		if err != nil {
			r.Logger.Error("CreatePatientBills orm.CreatePatientBills", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}

		// 刪掉舊的帳單
		if oldPatientBills != nil {
			var needDeleteBasicChargesID []uuid.UUID
			var needDeleteSubsidiesID []uuid.UUID
			for i := range oldPatientBills {
				for j := range oldPatientBills[i].BasicCharges {
					needDeleteBasicChargesID = append(needDeleteBasicChargesID, oldPatientBills[i].BasicCharges[j].ID)
				}
				for j := range oldPatientBills[i].Subsidies {
					needDeleteSubsidiesID = append(needDeleteSubsidiesID, oldPatientBills[i].Subsidies[j].ID)
				}
			}
			err = orm.DeletePatientBillsInId(tx, createdPatientBillIds, needDeleteBasicChargesID, needDeleteSubsidiesID)
			if err != nil {
				r.Logger.Error("CreatePatientBills orm.DeletePatientBillsInId", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
					zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreatePatientBills tx.Transaction", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	createdPatientBills, err := orm.GetPatientBillsByDate(r.ORM.DB, organizationId, billDataTimeZoneInTaipei.Year(), int(billDataTimeZoneInTaipei.Month()))
	if err != nil {
		r.Logger.Error("CreatePatientBills orm.GetPatientBillsByDate", zap.Error(err), zap.String("originalUrl", "createPatientBills"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createPatientBills run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createPatientBills"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return createdPatientBills, nil
}

func (r *mutationResolver) UpdatePatientBillNote(ctx context.Context, input *gqlmodels.UpdatePatientBillNoteInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillNote uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "updatePatientBillNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillNote uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updatePatientBillNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBillId, err := uuid.Parse(input.ID)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillNote uuid.Parse(input.ID)", zap.Error(err), zap.String("fieldName", "updatePatientBillNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBill := models.PatientBill{
		ID:             patientBillId,
		Note:           *input.Note,
		OrganizationId: organizationId,
		UserId:         userId,
	}
	err = orm.UpdatePatientBillNote(r.ORM.DB, &patientBill)
	if err != nil {
		r.Logger.Error("UpdatePatientBillNote orm.UpdatePatientBillNote", zap.Error(err), zap.String("fieldName", "updatePatientBillNote"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePatientBillNote run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePatientBillNote"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdatePatientBillChargeDates(ctx context.Context, patientBillIdStr string, input gqlmodels.UpdatePatientBillChargeDatesInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillChargeDates uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	patientBillId, err := uuid.Parse(patientBillIdStr)
	if err != nil {
		r.Logger.Warn("UpdatePatientBillChargeDates orm.uuid.Parse(patientBillIdStr)", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		patientBill, err := orm.GetPatientBillById(tx, organizationId, patientBillId)
		if err != nil {
			r.Logger.Error("UpdatePatientBillChargeDates orm.GetPatientBillById", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		// 有傳的再檢查否則每次都要重算這兩個
		if input.TransferRefundStartDate != nil || input.TransferRefundEndDate != nil {
			if input.TransferRefundStartDate != nil {
				patientBill.TransferRefundStartDate = *input.TransferRefundStartDate
			}
			if input.TransferRefundEndDate != nil {
				patientBill.TransferRefundEndDate = *input.TransferRefundEndDate
			}

			// 請假
			err = orm.ClearAssociationsPatientBillTransferRefundLeaves(tx, patientBill)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.ClearAssociationsPatientBillTransferRefundLeaves", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 重新找一次 符合時間區間的請假資料
			transferRefundLeavesData, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(r.ORM.DB, organizationId, patientBill.PatientId, patientBill.TransferRefundStartDate, patientBill.TransferRefundEndDate)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			err = orm.AppendAssociationsPatientBillTransferRefundLeaves(tx, patientBill, transferRefundLeavesData)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.AppendAssociationsPatientBillTransferRefundLeaves", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		if input.NonFixedChargeStartDate != nil || input.NonFixedChargeEndDate != nil {
			if input.NonFixedChargeStartDate != nil {
				patientBill.NonFixedChargeStartDate = *input.NonFixedChargeStartDate
			}
			if input.NonFixedChargeEndDate != nil {
				patientBill.NonFixedChargeEndDate = *input.NonFixedChargeEndDate
			}
			// 非固定
			err = orm.ClearAssociationsPatientBillNonFixedChargeRecords(tx, patientBill)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.ClearAssociationsPatientBillNonFixedChargeRecords", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			// 用時間區間找非固定費用
			nonFixedCharges, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(r.ORM.DB, patientBill.PatientId, patientBill.NonFixedChargeStartDate, patientBill.NonFixedChargeEndDate)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.GetNonFixedChargeRecordsByPatientIdAndDate", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}

			err = orm.AppendAssociationsPatientBillNonFixedChargeRecords(tx, patientBill, nonFixedCharges)
			if err != nil {
				r.Logger.Error("UpdatePatientBillChargeDates orm.AppendAssociationsPatientBillNonFixedChargeRecords", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
					zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
				return err
			}
		}
		// 更新住民帳單的顯示區間
		err = orm.UpdatePatientBillChargeDates(tx, patientBill)
		if err != nil {
			r.Logger.Error("UpdatePatientBillChargeDates orm.UpdatePatientBillChargeDates", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdatePatientBillChargeDates tx.Transaction", zap.Error(err), zap.String("fieldName", "updatePatientBillChargeDates"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updatePatientBillChargeDates run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updatePatientBillChargeDates"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) PatientBill(ctx context.Context, patientIdStr string, billYear, billMonth int) (*models.PatientBill, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientBill uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "patientBill"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientId, err := uuid.Parse(patientIdStr)
	if err != nil {
		r.Logger.Warn("PatientBill uuid.Parse(patientIdStr)", zap.Error(err), zap.String("fieldName", "patientBill"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientBill, err := orm.GetPatientBillByPatientIdAndYearMonth(r.ORM.DB, organizationId, patientId, billYear, billMonth)
	if err != nil {
		r.Logger.Error("PatientBill orm.GetPatientBillByPatientIdAndYearMonth", zap.Error(err), zap.String("fieldName", "patientBill"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	if patientBill.EditNoteUserId != nil {
		editNoteUser, err := orm.GetUserById(r.ORM.DB, *patientBill.EditNoteUserId)
		if err != nil {
			r.Logger.Error("PatientBill orm.GetUserById", zap.Error(err), zap.String("fieldName", "patientBill"),
				zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
			return nil, err
		}
		patientBill.EditNoteUser = *editNoteUser
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientBill run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientBill"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return patientBill, nil
}

func (r *queryResolver) PatientBills(ctx context.Context, BillDate time.Time) ([]*models.PatientBill, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientBills uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "patientBills"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	taipeiZone, _ := time.LoadLocation("Asia/Taipei")

	patientBills, err := orm.GetPatientBillsByDate(r.ORM.DB, organizationId, BillDate.In(taipeiZone).Year(), int(BillDate.In(taipeiZone).Month()))
	if err != nil {
		r.Logger.Error("PatientBills orm.GetPatientBillsByDate", zap.Error(err), zap.String("fieldName", "patientBills"),
			zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientBills run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientBills"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return patientBills, nil
}

// patientBill resolvers
type patientBillResolver struct{ *Resolver }

func (r *patientBillResolver) ID(ctx context.Context, obj *models.PatientBill) (string, error) {
	return obj.ID.String(), nil
}

// 捏出固定費用的基本月費資料
func getBasicCharge(fixedChargeSettings []*models.BasicChargeSetting, fixedChargeStartDate, fixedChargeEndDate time.Time, patientId, userId uuid.UUID) []*models.BasicCharge {
	patientBillBasicCharges := []*models.BasicCharge{}
	for i := range fixedChargeSettings {
		patientBillBasicCharge := models.BasicCharge{
			ID:             uuid.New(),
			ItemName:       fixedChargeSettings[i].OrganizationBasicChargeSetting.ItemName,
			Type:           fixedChargeSettings[i].OrganizationBasicChargeSetting.Type,
			Unit:           fixedChargeSettings[i].OrganizationBasicChargeSetting.Unit,
			Price:          fixedChargeSettings[i].OrganizationBasicChargeSetting.Price,
			TaxType:        fixedChargeSettings[i].OrganizationBasicChargeSetting.TaxType,
			StartDate:      fixedChargeStartDate,
			EndDate:        fixedChargeEndDate,
			Note:           "",
			SortIndex:      i,
			OrganizationId: fixedChargeSettings[i].OrganizationId,
			PatientId:      patientId,
			UserId:         userId,
		}
		patientBillBasicCharges = append(patientBillBasicCharges, &patientBillBasicCharge)
	}
	return patientBillBasicCharges
}

// 捏出補助款資料
func getPatientBillSubsidies(subsidiesSetting []*models.SubsidySetting, fixedChargeStartDate, fixedChargeEndDate time.Time, userId, patientId, organizationId uuid.UUID) []*models.Subsidy {
	// 補助款(時間和固定費用共用)
	patientBillSubsidies := []*models.Subsidy{}
	for i := range subsidiesSetting {
		patientBillSubsidy := models.Subsidy{
			ID:             uuid.New(),
			ItemName:       subsidiesSetting[i].ItemName,
			Type:           subsidiesSetting[i].Type,
			Price:          subsidiesSetting[i].Price,
			Unit:           subsidiesSetting[i].Unit,
			IdNumber:       subsidiesSetting[i].IdNumber,
			Note:           subsidiesSetting[i].Note,
			StartDate:      fixedChargeStartDate,
			EndDate:        fixedChargeEndDate,
			ReceiptStatus:  "",
			SortIndex:      i,
			OrganizationId: organizationId,
			PatientId:      patientId,
			UserId:         userId,
		}
		patientBillSubsidies = append(patientBillSubsidies, &patientBillSubsidy)
	}
	return patientBillSubsidies
}

// 取得固定or非固定or異動的時間區間
func getPatientBillDateRange(billDate time.Time, startStatus, endStatus string, startDay, endDay int) (*time.Time, *time.Time, error) {
	// 先算起的月份
	var fixedChargeMonthStartStatus int
	if startStatus == "thisMonth" {
		fixedChargeMonthStartStatus = 0
	} else if startStatus == "lastMonth" {
		fixedChargeMonthStartStatus = -1
	} else if startStatus == "twoMonthsAgo" {
		fixedChargeMonthStartStatus = -2
	} else {
		return nil, nil, fmt.Errorf("startStatus error")
	}
	// 算出住民帳單上起的月份
	addTime := billDate.AddDate(0, fixedChargeMonthStartStatus, 0)

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, nil, fmt.Errorf("taipeiZone error")
	}
	// 算出住民帳單上起的日期
	fixedChargeStartDate := time.Date(addTime.Year(), addTime.Month(), startDay, 0, 0, 0, 0, taipeiZone)
	// 算迄的月份
	var fixedChargeMonthEndStatus int
	if endStatus == "thisMonth" {
		fixedChargeMonthEndStatus = 0
	} else if endStatus == "lastMonth" {
		fixedChargeMonthEndStatus = -1
	} else if endStatus == "twoMonthsAgo" {
		fixedChargeMonthEndStatus = -2
	} else {
		return nil, nil, fmt.Errorf("endStatus error")
	}

	// 算出住民帳單上迄的月份
	addTime = billDate.AddDate(0, fixedChargeMonthEndStatus, 0)
	var fixedChargeEndDate time.Time
	// 算出住民帳單上迄的日期
	if endDay == 31 {

		fixedChargeEndDate = addTime.AddDate(0, 0, -addTime.Day()+1).AddDate(0, 1, -1)

		fixedChargeEndDate = time.Date(fixedChargeEndDate.Year(), fixedChargeEndDate.Month(), fixedChargeEndDate.Day(), 23, 59, 59, 0, taipeiZone)

	} else {
		fixedChargeEndDate = time.Date(addTime.Year(), addTime.Month(), endDay, 23, 59, 59, 0, taipeiZone)

	}
	return &fixedChargeStartDate, &fixedChargeEndDate, nil
}
