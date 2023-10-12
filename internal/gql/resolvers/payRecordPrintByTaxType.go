package resolvers

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/resolvers/excelStyle"

	"graphql-go-template/internal/models"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type OrganizationReceiptTemplateSettingPatientInfo struct {
	ShowAreaAndClass  bool
	ShowBedAndRoom    bool
	ShowSex           bool
	ShowBirthday      bool
	ShowAge           bool
	ShowCheckInDate   bool
	ShowPatientNumber bool
	ShowRecordNumber  bool
	ShowIdNumber      bool
}

type OrganizationReceiptTemplateSettingOrganizationInfo struct {
	ShowTaxIdNumber         bool
	ShowPhone               bool
	ShowFax                 bool
	ShowOwner               bool
	ShowEstablishmentNumber bool
	ShowAddress             bool
	ShowEmail               bool
	ShowRemittanceBank      bool
	ShowRemittanceIdNumber  bool
	ShowRemittanceUserName  bool
}

type PrintPayRecordPartByTaxTypeStruct struct {
	f                                   *excelize.File
	PayRecord                           *models.PayRecord
	OrganizationReceiptTemplateSettings []*models.OrganizationReceiptTemplateSetting
	r                                   *queryResolver
	Ctx                                 context.Context
	TaxTypes                            []string
}

type PrintPayRecordPartsByTaxTypeStruct struct {
	F                                              *excelize.File
	PayRecord                                      *models.PayRecord
	OrganizationReceiptTemplateSettings            []*models.OrganizationReceiptTemplateSetting
	OrganizationReceiptTemplateSettingsPatientInfo []OrganizationReceiptTemplateSettingPatientInfo
	R                                              *queryResolver
	Ctx                                            context.Context
	PayRecordYear                                  string
	PayRecordMonth                                 string
	InvalidText                                    string
	IsInvalid                                      bool
	TaxTypes                                       []string
	ClassStyle                                     int
	ContentStyle                                   int
	DateStyle                                      int
	PriceStyle                                     int
}

//go:embed excelTemplate/receiptPartSelectTaxType.xlsx
var receiptPartSelectTaxType []byte

func (r *queryResolver) PrintPayRecordPartByTaxType(ctx context.Context, payRecordIdStr string) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordPartByTaxType uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordPartByTaxType uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecord, err := orm.GetPayRecordById(r.ORM.DB, payRecordId, true, true, false)
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", nil
	}

	reader := bytes.NewReader(receiptPartSelectTaxType)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType excelize.OpenReader", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	var organizationReceiptTemplateSettings []*models.OrganizationReceiptTemplateSetting
	stampTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "stampTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is stampTax", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, stampTaxSetting)
	businessTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "businessTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is businessTax", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, businessTaxSetting)
	noTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "noTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is noTax", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, noTaxSetting)
	otherTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "other")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is other", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, otherTaxSetting)

	taxTypes := []string{"stampTax", "businessTax", "noTax", "other"}

	printPayRecordPartByTaxTypeStruct := PrintPayRecordPartByTaxTypeStruct{
		f:                                   f,
		PayRecord:                           payRecord,
		OrganizationReceiptTemplateSettings: organizationReceiptTemplateSettings,
		r:                                   r,
		Ctx:                                 ctx,
		TaxTypes:                            taxTypes,
	}

	// 在這邊先下載印章(不然底下才下載的話 會重複下載)
	for i := range organizationReceiptTemplateSettings {
		downloadExcelOrganizationSealStruct := DownloadExcelOrganizationSealStruct{
			OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
			r:                                  r,
			ctx:                                ctx,
		}
		err = DownloadExcelOrganizationSeal(downloadExcelOrganizationSealStruct)
		if err != nil {
			r.Logger.Error("PrintPayRecordPartByTaxType DownloadExcelOrganizationSeal", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
	}

	fileName, err := PrintPayRecordPartByTaxType(printPayRecordPartByTaxTypeStruct)
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType PrintPayRecordPartByTaxType", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType f.SaveAs", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType r.store.UploadFile", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordPartByTaxType r.store.SetMetadata", zap.Error(err), zap.String("originalUrl", "printPayRecordPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	publicUrl := r.store.GenPublicLink("payRecordPart/" + fileName)

	// 刪掉下載的印章檔案
	for i := range organizationReceiptTemplateSettings {
		if organizationReceiptTemplateSettings[i].OrganizationPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].OrganizationPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealOnePicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealOnePicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealTwoPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealTwoPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealThreePicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealThreePicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealFourPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealFourPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPayRecordPartByTaxType run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPayRecordPartByTaxType"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

func PrintPayRecordPartByTaxType(printPayRecordPartByTaxTypeStruct PrintPayRecordPartByTaxTypeStruct) (string, error) {
	f := printPayRecordPartByTaxTypeStruct.f
	organizationReceiptTemplateSettings := printPayRecordPartByTaxTypeStruct.OrganizationReceiptTemplateSettings
	taxTypes := printPayRecordPartByTaxTypeStruct.TaxTypes
	payRecord := printPayRecordPartByTaxTypeStruct.PayRecord
	haveTaxBool := []bool{false, false, false, false}

	payRecordYear := strconv.Itoa(payRecord.PayYear)
	payRecordMonth := strconv.Itoa(payRecord.PayMonth)

	if payRecord.PayMonth < 10 {
		payRecordMonth = "0" + payRecordMonth
	}
	var invalidText string
	var isInvalid bool
	if payRecord.IsInvalid {
		isInvalid = true
		invalidText = "（作廢）"
	}

	// 內容框線
	classStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 11)
	if err != nil {
		return "", err
	}

	// 內容框線(日期)
	dateStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 10)
	if err != nil {
		return "", err
	}

	// 金額框線
	priceFormatAndRightBottomBorderAndFontStyle, err := excelStyle.GetPriceFormatAndTwoBorderAndFontStyle(f, 11, []string{"right", "bottom"}, "Calibri (本文)", []int{2, 2})
	if err != nil {
		return "", err
	}

	for i := range taxTypes {
		organizationReceiptTemplateSettingPatientInfo := OrganizationReceiptTemplateSettingPatientInfo{
			ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "areaAndClass"),
			ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "bedAndRoom"),
			ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "sex"),
			ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "birthday"),
			ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "age"),
			ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "checkInDate"),
			ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "patientNumber"),
			ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "recordNumber"),
			ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSettings[i], "idNumber"),
		}
		// 組一個塞住民資訊的struct
		excelPatientDataStruct := ExcelPatientDataStruct{
			f:                                  f,
			PayRecord:                          payRecord,
			SheetName:                          taxTypes[i],
			PayRecordYear:                      payRecordYear,
			PayRecordMonth:                     payRecordMonth,
			InvalidText:                        invalidText,
			OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
			OrganizationReceiptTemplateSettingPatientInfo: organizationReceiptTemplateSettingPatientInfo,
		}
		// 要先塞住民資訊再用費用的原因是
		// 因為聯單的問題 費用內容可能會新增列 影響到第二聯(其實也可以寫成 檢查InsertRowCount的形式 但好像沒什麼必要)
		SetExcelPatientData(excelPatientDataStruct)
		if organizationReceiptTemplateSettings[i].PriceShowType == "classAddUp" {
			// 塞費用的內容(回傳有沒有資料和總共新增幾列)
			taxHaveData, insertRowCount := printClassAddUpByTaxType(payRecord, f, classStyle, priceFormatAndRightBottomBorderAndFontStyle, i, taxTypes[i], taxTypes[i])
			haveTaxBool[i] = taxHaveData
			// 檢查這個稅別有沒有資料 有的話才塞機構資料
			if taxHaveData {
				organizationReceiptTemplateSettingOrganizationInfo := OrganizationReceiptTemplateSettingOrganizationInfo{
					ShowTaxIdNumber:         organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "taxIdNumber"),
					ShowPhone:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "phone"),
					ShowFax:                 organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "fax"),
					ShowOwner:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "owner"),
					ShowEstablishmentNumber: organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "establishmentNumber"),
					ShowAddress:             organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "address"),
					ShowEmail:               organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "email"),
					ShowRemittanceBank:      organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceBank"),
					ShowRemittanceIdNumber:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceIdNumber"),
					ShowRemittanceUserName:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceUserName"),
				}
				// 組一個塞機構資訊的struct
				excelOrganizationDataStruct := ExcelOrganizationDataStruct{
					f:                                  f,
					PayRecord:                          payRecord,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					OrganizationReceiptTemplateSettingOrganizationInfo: organizationReceiptTemplateSettingOrganizationInfo,
					SheetName:      taxTypes[i],
					InsertRowCount: insertRowCount,
					IsInvalid:      isInvalid,
				}
				// 塞機構資訊
				SetExcelOrganizationData(excelOrganizationDataStruct)
				// 組一個塞印章資訊的struct
				organizationSealDataStruct := OrganizationSealDataStruct{
					f:                                  f,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					SheetName:                          taxTypes[i],
					InsertRowCount:                     insertRowCount,
				}
				// 塞印章資訊
				err = SetExcelOrganizationSealData(organizationSealDataStruct)
				if err != nil {
					return "", err
				}
			}
		} else {
			// 塞費用的內容(回傳有沒有資料和總共新增幾列)
			printItemByTaxTypeStruct := &PrintItemByTaxTypeStruct{
				PayRecord:  payRecord,
				F:          f,
				ClassStyle: classStyle,
				DateStyle:  dateStyle,
				PriceStyle: priceFormatAndRightBottomBorderAndFontStyle,
				TaxCount:   i,
				TaxType:    taxTypes[i],
				SheetName:  taxTypes[i],
			}
			taxHaveData, insertRowCount := printItemByTaxType(printItemByTaxTypeStruct)
			haveTaxBool[i] = taxHaveData
			// 檢查這個稅別有沒有資料 有的話才塞機構資料
			if taxHaveData {
				organizationReceiptTemplateSettingOrganizationInfo := OrganizationReceiptTemplateSettingOrganizationInfo{
					ShowTaxIdNumber:         organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "taxIdNumber"),
					ShowPhone:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "phone"),
					ShowFax:                 organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "fax"),
					ShowOwner:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "owner"),
					ShowEstablishmentNumber: organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "establishmentNumber"),
					ShowAddress:             organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "address"),
					ShowEmail:               organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "email"),
					ShowRemittanceBank:      organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceBank"),
					ShowRemittanceIdNumber:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceIdNumber"),
					ShowRemittanceUserName:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceUserName"),
				}
				// 組一個塞機構資訊的struct
				excelOrganizationDataStruct := ExcelOrganizationDataStruct{
					f:                                  f,
					PayRecord:                          payRecord,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					OrganizationReceiptTemplateSettingOrganizationInfo: organizationReceiptTemplateSettingOrganizationInfo,
					SheetName:      taxTypes[i],
					InsertRowCount: insertRowCount,
					IsInvalid:      isInvalid,
				}
				// 塞機構資訊
				SetExcelOrganizationData(excelOrganizationDataStruct)
				// 組一個塞印章資訊的struct
				organizationSealDataStruct := OrganizationSealDataStruct{
					f:                                  f,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					SheetName:                          taxTypes[i],
					InsertRowCount:                     insertRowCount,
				}
				// 塞印章資訊
				err = SetExcelOrganizationSealData(organizationSealDataStruct)
				if err != nil {
					return "", err
				}
			}
		}
	}
	// 檢查該稅別是不是有資料
	// 有=>改名稱
	// 沒有=>刪掉
	changeSheetNames := []string{"印花稅", "營業稅", "免稅", "其他"}
	for i := range haveTaxBool {
		if haveTaxBool[i] {
			f.SetSheetName(taxTypes[i], changeSheetNames[i])
		} else {
			f.DeleteSheet(taxTypes[i])
		}
	}
	// 組合sheetName
	fileName := payRecord.Patient.Branch + payRecord.Patient.Room + payRecord.Patient.Bed + payRecord.Patient.LastName + payRecord.Patient.FirstName + " " + payRecordYear + "年" + payRecordMonth + "月收據聯單" + payRecord.ReceiptNumber + "(按稅別分頁).xlsx"
	if err := printPayRecordPartByTaxTypeStruct.f.SaveAs(fileName); err != nil {
		return "", err
	}
	return fileName, nil
}

func (r *queryResolver) PrintPayRecordsPartByTaxType(ctx context.Context, payRecordsIdStr []string) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordsPartByTaxType uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	var payRecordsId []uuid.UUID

	for i := range payRecordsIdStr {
		payRecordId, err := uuid.Parse(payRecordsIdStr[i])
		if err != nil {
			r.Logger.Warn("PrintPayRecordsPartByTaxType uuid.Parse(payRecordsIdStr[i])", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
		payRecordsId = append(payRecordsId, payRecordId)
	}

	var organizationReceiptTemplateSettings []*models.OrganizationReceiptTemplateSetting
	// 這邊先把 各個稅別的列印設定存起來 因為是批次列印 如果在後面(GetPayRecordsForPrint)做的話 會跑無謂的迴圈
	var organizationReceiptTemplateSettingsPatientInfo []OrganizationReceiptTemplateSettingPatientInfo
	stampTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "stampTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is stampTax", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, stampTaxSetting)
	organizationReceiptTemplateSettingsPatientInfo = append(organizationReceiptTemplateSettingsPatientInfo, OrganizationReceiptTemplateSettingPatientInfo{
		ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "areaAndClass"),
		ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "bedAndRoom"),
		ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "sex"),
		ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "birthday"),
		ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "age"),
		ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "checkInDate"),
		ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "patientNumber"),
		ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "recordNumber"),
		ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(stampTaxSetting, "idNumber"),
	})
	businessTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "businessTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is businessTax", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, businessTaxSetting)
	organizationReceiptTemplateSettingsPatientInfo = append(organizationReceiptTemplateSettingsPatientInfo, OrganizationReceiptTemplateSettingPatientInfo{
		ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "areaAndClass"),
		ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "bedAndRoom"),
		ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "sex"),
		ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "birthday"),
		ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "age"),
		ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "checkInDate"),
		ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "patientNumber"),
		ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "recordNumber"),
		ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(businessTaxSetting, "idNumber"),
	})
	noTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "noTax")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is noTax", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, noTaxSetting)
	organizationReceiptTemplateSettingsPatientInfo = append(organizationReceiptTemplateSettingsPatientInfo, OrganizationReceiptTemplateSettingPatientInfo{
		ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "areaAndClass"),
		ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "bedAndRoom"),
		ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "sex"),
		ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "birthday"),
		ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "age"),
		ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "checkInDate"),
		ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "patientNumber"),
		ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "recordNumber"),
		ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(noTaxSetting, "idNumber"),
	})
	otherTaxSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, "other")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType orm.GetOrganizationReceiptTemplateSettingInTaxType is other", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationReceiptTemplateSettingsPatientInfo = append(organizationReceiptTemplateSettingsPatientInfo, OrganizationReceiptTemplateSettingPatientInfo{
		ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "areaAndClass"),
		ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "bedAndRoom"),
		ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "sex"),
		ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "birthday"),
		ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "age"),
		ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "checkInDate"),
		ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "patientNumber"),
		ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "recordNumber"),
		ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(otherTaxSetting, "idNumber"),
	})
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, otherTaxSetting)

	payRecords, err := orm.GetPayRecordsForPrint(r.ORM.DB, payRecordsId, true, true)
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType orm.GetPayRecordsForPrint", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", nil
	}

	reader := bytes.NewReader(receiptPartSelectTaxType)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType excelize.OpenReader", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecordYear := strconv.Itoa(payRecords[0].PayYear)
	payRecordMonth := strconv.Itoa(payRecords[0].PayMonth)

	if payRecords[0].PayMonth < 10 {
		payRecordMonth = "0" + payRecordMonth
	}

	taxTypes := []string{"stampTax", "businessTax", "noTax", "other"}
	fileName := time.Now().Format("01-02") + "收據聯單(按稅別分頁).xlsx"
	// 內容框線
	classStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 11)
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 內容框線
	contentStyle, err := excelStyle.GetRightAndBottomBorderAndLeftAlignmentAndFontStyle(f)
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType excelStyle.GetRightAndBottomBorderAndLeftAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 內容框線(日期)
	dateStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 10)
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	priceStyle, err := excelStyle.GetPriceFormatAndTwoBorderAndFontStyle(f, 11, []string{"right", "bottom"}, "Calibri (本文)", []int{2, 2})
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType excelStyle.GetPriceFormatAndTwoBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	if err != nil {
		return "", err
	}

	// 在這邊先下載印章(不然底下才下載的話 會重複下載)
	for i := range organizationReceiptTemplateSettings {
		downloadExcelOrganizationSealStruct := DownloadExcelOrganizationSealStruct{
			OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
			r:                                  r,
			ctx:                                ctx,
		}
		err = DownloadExcelOrganizationSeal(downloadExcelOrganizationSealStruct)
		if err != nil {
			r.Logger.Error("PrintPayRecordsPartByTaxType DownloadExcelOrganizationSeal", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
	}

	// r.Logger.Info("start PrintPayRecordPartsByTaxType")

	for i := range payRecords {
		var isInvalid bool
		var invalidText string
		if payRecords[i].IsInvalid {
			isInvalid = true
			invalidText = "（作廢）"
		}
		printPayRecordPartsByTaxTypeStruct := PrintPayRecordPartsByTaxTypeStruct{
			F:                                   f,
			PayRecord:                           payRecords[i],
			OrganizationReceiptTemplateSettings: organizationReceiptTemplateSettings,
			OrganizationReceiptTemplateSettingsPatientInfo: organizationReceiptTemplateSettingsPatientInfo,
			R:              r,
			Ctx:            ctx,
			PayRecordYear:  payRecordYear,
			PayRecordMonth: payRecordMonth,
			InvalidText:    invalidText,
			IsInvalid:      isInvalid,
			TaxTypes:       taxTypes,
			ClassStyle:     classStyle,
			ContentStyle:   contentStyle,
			DateStyle:      dateStyle,
			PriceStyle:     priceStyle,
		}
		err = PrintPayRecordPartsByTaxType(printPayRecordPartsByTaxTypeStruct)
		if err != nil {
			r.Logger.Error("PrintPayRecordsPartByTaxType PrintPayRecordPartsByTaxType", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
	}

	// 把預設的四個sheet清掉
	for i := range taxTypes {
		f.DeleteSheet(taxTypes[i])
	}

	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType f.SaveAs", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType r.store.UploadFile", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordsPartByTaxType r.store.SetMetadata", zap.Error(err), zap.String("originalUrl", "printPayRecordsPartByTaxType"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	publicUrl := r.store.GenPublicLink("payRecordPart/" + fileName)
	// 刪掉下載的印章檔案
	for i := range organizationReceiptTemplateSettings {
		if organizationReceiptTemplateSettings[i].OrganizationPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].OrganizationPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealOnePicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealOnePicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealTwoPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealTwoPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealThreePicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealThreePicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
		if organizationReceiptTemplateSettings[i].SealFourPicture != "" {
			orgAndFileId := strings.Split(organizationReceiptTemplateSettings[i].SealFourPicture, "inventory-tool/")
			fileName := strings.Split(orgAndFileId[1], "/")
			os.Remove(fileName[1] + ".jpg")
		}
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPayRecordsPartByTaxType run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPayRecordsPartByTaxType"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

func PrintPayRecordPartsByTaxType(printPayRecordPartsByTaxTypeStruct PrintPayRecordPartsByTaxTypeStruct) error {
	f := printPayRecordPartsByTaxTypeStruct.F
	payRecord := printPayRecordPartsByTaxTypeStruct.PayRecord
	taxTypes := printPayRecordPartsByTaxTypeStruct.TaxTypes
	payRecordYear := printPayRecordPartsByTaxTypeStruct.PayRecordYear
	payRecordMonth := printPayRecordPartsByTaxTypeStruct.PayRecordMonth
	invalidText := printPayRecordPartsByTaxTypeStruct.InvalidText
	organizationReceiptTemplateSettings := printPayRecordPartsByTaxTypeStruct.OrganizationReceiptTemplateSettings
	organizationReceiptTemplateSettingsPatientInfo := printPayRecordPartsByTaxTypeStruct.OrganizationReceiptTemplateSettingsPatientInfo
	classStyle := printPayRecordPartsByTaxTypeStruct.ClassStyle
	contentStyle := printPayRecordPartsByTaxTypeStruct.ContentStyle
	dateStyle := printPayRecordPartsByTaxTypeStruct.DateStyle
	priceStyle := printPayRecordPartsByTaxTypeStruct.PriceStyle
	isInvalid := printPayRecordPartsByTaxTypeStruct.IsInvalid
	r := printPayRecordPartsByTaxTypeStruct.R
	haveTaxBool := []bool{false, false, false, false}
	sheetsName := []string{payRecord.ReceiptNumber + "-印花稅", payRecord.ReceiptNumber + "-營業稅", payRecord.ReceiptNumber + "-免稅", payRecord.ReceiptNumber + "-其他"}
	f.NewSheet(payRecord.ReceiptNumber + "-印花稅")
	f.NewSheet(payRecord.ReceiptNumber + "-營業稅")
	f.NewSheet(payRecord.ReceiptNumber + "-免稅")
	f.NewSheet(payRecord.ReceiptNumber + "-其他")
	for i := range taxTypes {
		// 組一個塞住民資訊的struct
		excelPatientDataStruct := ExcelPatientDataStruct{
			f:                                  f,
			PayRecord:                          payRecord,
			SheetName:                          sheetsName[i],
			PayRecordYear:                      payRecordYear,
			PayRecordMonth:                     payRecordMonth,
			InvalidText:                        invalidText,
			OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
			OrganizationReceiptTemplateSettingPatientInfo: organizationReceiptTemplateSettingsPatientInfo[i],
		}
		// 要先塞住民資訊再用費用的原因是
		// 因為聯單的問題 費用內容可能會新增列 影響到第二聯(其實也可以寫成 檢查InsertRowCount的形式 但好像沒什麼必要)
		SetExcelPatientData(excelPatientDataStruct)
		if organizationReceiptTemplateSettings[i].PriceShowType == "classAddUp" {
			// 塞費用的內容(回傳有沒有資料和總共新增幾列)
			taxHaveData, insertRowCount := printClassAddUpByTaxType(payRecord, f, classStyle, priceStyle, i, taxTypes[i], sheetsName[i])
			haveTaxBool[i] = taxHaveData
			// 檢查這個稅別有沒有資料 有的話改塞機構資料
			// 沒有的話把sheet刪掉
			if taxHaveData {
				organizationReceiptTemplateSettingOrganizationInfo := OrganizationReceiptTemplateSettingOrganizationInfo{
					ShowTaxIdNumber:         organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "taxIdNumber"),
					ShowPhone:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "phone"),
					ShowFax:                 organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "fax"),
					ShowOwner:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "owner"),
					ShowEstablishmentNumber: organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "establishmentNumber"),
					ShowAddress:             organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "address"),
					ShowEmail:               organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "email"),
					ShowRemittanceBank:      organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceBank"),
					ShowRemittanceIdNumber:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceIdNumber"),
					ShowRemittanceUserName:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceUserName"),
				}
				// 組一個塞機構資訊的struct
				excelOrganizationDataStruct := ExcelOrganizationDataStruct{
					f:                                  f,
					PayRecord:                          payRecord,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					OrganizationReceiptTemplateSettingOrganizationInfo: organizationReceiptTemplateSettingOrganizationInfo,
					SheetName:      sheetsName[i],
					InsertRowCount: insertRowCount,
					IsInvalid:      isInvalid,
				}
				// 塞機構資訊
				SetExcelOrganizationData(excelOrganizationDataStruct)
				// 組一個塞印章資訊的struct
				organizationSealDataStruct := OrganizationSealDataStruct{
					f:                                  f,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					SheetName:                          sheetsName[i],
					InsertRowCount:                     insertRowCount,
				}
				// 塞印章資訊
				err := SetExcelOrganizationSealData(organizationSealDataStruct)
				if err != nil {
					r.Logger.Error("PrintPayRecordsPartByTaxType classAddUp SetExcelOrganizationSealData", zap.Error(err))
					return err
				}
			}
		} else {
			// 塞費用的內容(回傳有沒有資料和總共新增幾列)
			printItemByTaxTypeStruct := &PrintItemByTaxTypeStruct{
				PayRecord:  payRecord,
				F:          f,
				ClassStyle: contentStyle,
				DateStyle:  dateStyle,
				PriceStyle: priceStyle,
				TaxCount:   i,
				TaxType:    taxTypes[i],
				SheetName:  sheetsName[i],
			}
			taxHaveData, insertRowCount := printItemByTaxType(printItemByTaxTypeStruct)
			haveTaxBool[i] = taxHaveData
			// 檢查這個稅別有沒有資料 有的話改塞機構資料
			// 沒有的話把sheet刪掉
			if taxHaveData {
				organizationReceiptTemplateSettingOrganizationInfo := OrganizationReceiptTemplateSettingOrganizationInfo{
					ShowTaxIdNumber:         organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "taxIdNumber"),
					ShowPhone:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "phone"),
					ShowFax:                 organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "fax"),
					ShowOwner:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSettings[i], "owner"),
					ShowEstablishmentNumber: organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "establishmentNumber"),
					ShowAddress:             organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "address"),
					ShowEmail:               organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "email"),
					ShowRemittanceBank:      organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceBank"),
					ShowRemittanceIdNumber:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceIdNumber"),
					ShowRemittanceUserName:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSettings[i], "remittanceUserName"),
				}
				// 組一個塞機構資訊的struct
				excelOrganizationDataStruct := ExcelOrganizationDataStruct{
					f:                                  f,
					PayRecord:                          payRecord,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					OrganizationReceiptTemplateSettingOrganizationInfo: organizationReceiptTemplateSettingOrganizationInfo,
					SheetName:      sheetsName[i],
					InsertRowCount: insertRowCount,
					IsInvalid:      isInvalid,
				}
				// 塞機構資訊
				SetExcelOrganizationData(excelOrganizationDataStruct)
				// 組一個塞印章資訊的struct
				organizationSealDataStruct := OrganizationSealDataStruct{
					f:                                  f,
					OrganizationReceiptTemplateSetting: organizationReceiptTemplateSettings[i],
					SheetName:                          sheetsName[i],
					InsertRowCount:                     insertRowCount,
				}
				// 塞印章資訊
				err := SetExcelOrganizationSealData(organizationSealDataStruct)
				if err != nil {
					r.Logger.Error("PrintPayRecordsPartByTaxType item SetExcelOrganizationSealData error", zap.Error(err))
					return err
				}
			}
		}
	}
	// 刪除 	// 檢查該稅別是不是有資料
	// 沒有=>刪掉
	for i := range haveTaxBool {
		if !haveTaxBool[i] {
			f.DeleteSheet(sheetsName[i])
		}
	}
	return nil
}
