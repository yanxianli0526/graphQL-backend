package resolvers

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	orm "graphql-go-template/internal/database"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"graphql-go-template/internal/gql/resolvers/excelStyle"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type PayRecordTransferRefundLeaveForPrint struct {
	ItemName   string      `json:"itemName"`
	Type       string      `json:"type"`
	Price      int         `json:"price"`
	StartDates []time.Time `json:"startDate"`
	EndDates   []time.Time `json:"endDate"`
}

type PayRecordNonFixedChargeRecordDataForPrint struct {
	ItemCategory      string                                            `json:"itemCategory"`
	ItemCategoryDatas map[string]*PayRecordNonFixedChargeRecordForPrint `json:"itemCategoryDatas"`
}

type PayRecordNonFixedChargeRecordDateAndQuantityForPrint struct {
	Date     time.Time `json:"date"`
	Quantity int       `json:"quantity"`
}

type PayRecordNonFixedChargeRecordForPrint struct {
	ItemCategory       string                                                           `json:"itemCategory"`
	ItemName           string                                                           `json:"itemName"`
	Type               string                                                           `json:"type"`
	TaxType            string                                                           `json:"taxType"`
	Unit               string                                                           `json:"unit"`
	NonFixedChargeDate time.Time                                                        `json:"nonFixedChargeDate"`
	Quantity           int                                                              `json:"quantity"`
	Price              int                                                              `json:"price"`
	Subtotal           int                                                              `json:"subtotal"`
	DateAndQuantity    map[string]*PayRecordNonFixedChargeRecordDateAndQuantityForPrint `json:"newTest"`
}

// 給列印用的非固定(聯單)
type PayRecordNonFixedChargeRecordDataForPartPrint struct {
	ItemName     string      `json:"itemName"`
	ItemCategory string      `json:"itemCategory"`
	Subtotal     int         `json:"subtotal"`
	Quantities   []int       `json:"quantities"`
	Date         []time.Time `json:"date"`
}

type PayRecordNonFixedChargeRecordDataForMoreTaxTypePrint struct {
	ItemCategory                             string `json:"itemCategory"`
	Subtotal                                 int    `json:"subtotal"`
	StampTaxSubdNonFixedChargeRecordPrice    int    `json:"stampTaxSubdNonFixedChargeRecordPrice"`
	BusinessTaxSubdNonFixedChargeRecordPrice int    `json:"businessTaxSubdNonFixedChargeRecordPrice"`
	NoTaxSubdNonFixedChargeRecordPrice       int    `json:"noTaxSubdNonFixedChargeRecordPrice"`
	OtherSubdNonFixedChargeRecordPrice       int    `json:"otherSubdNonFixedChargeRecordPrice"`
}

//go:embed excelTemplate/receiptDetail.xlsx
var excelReceiptDetail []byte

func (r *queryResolver) PrintPayRecordDetail(ctx context.Context, payRecordIdStr string) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordDetail uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecord, err := orm.GetPayRecordById(r.ORM.DB, payRecordId, true, true, false)
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", nil
	}

	reader := bytes.NewReader(excelReceiptDetail)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail excelize.OpenReader", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecordYear := strconv.Itoa(payRecord.PayYear)
	payRecordMonth := strconv.Itoa(payRecord.PayMonth)

	if payRecord.PayMonth < 10 {
		payRecordMonth = "0" + payRecordMonth
	}

	f.SetCellValue("Template", "A1", payRecord.Organization.Name)
	f.SetCellValue("Template", "A2", payRecordYear+"年"+payRecordMonth+"月收據明細")
	f.SetCellValue("Template", "E3", "收據編號："+payRecord.ReceiptNumber)
	f.SetCellValue("Template", "A4", "床號："+payRecord.Patient.Room+payRecord.Patient.Bed)
	f.SetCellValue("Template", "B4", "姓名："+payRecord.Patient.LastName+payRecord.Patient.FirstName)
	var idNumber string
	if payRecord.Organization.Privacy == "unmask" {
		idNumber = payRecord.Patient.IdNumber
	} else {
		if len(payRecord.Patient.IdNumber) >= 10 {
			idNumber = payRecord.Patient.IdNumber[0:3] + "****" + payRecord.Patient.IdNumber[7:10]
		} else {
			idNumberCount := 0
			for idNumberCount < len(payRecord.Patient.IdNumber) {
				if idNumberCount >= 4 && idNumberCount <= 8 {
					idNumber += "*"
				} else {
					yo := string([]rune(payRecord.Patient.IdNumber)[idNumberCount])
					idNumber += yo
				}
				idNumberCount += 1
			}
		}
	}
	f.SetCellValue("Template", "C4", "身分證字號："+idNumber)
	// 副標題底色
	subTitleStyle, err := excelStyle.GetFillColorStyle(f)
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail excelStyle.GetFillColorStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 內容框線
	contentBorderStyle, err := excelStyle.GetAllBorderStyle(f)
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail excelStyle.GetAllBorderStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 基本月費
	payRecordBasicCharges := []PayRecordBasicCharge{}
	json.Unmarshal(payRecord.BasicCharge, &payRecordBasicCharges)
	dataCount := 0

	sort.Slice(payRecordBasicCharges, func(i, j int) bool {
		return payRecordBasicCharges[i].ItemName < payRecordBasicCharges[j].ItemName
	})

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail time.LoadLocation", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	if len(payRecordBasicCharges) > 0 {
		// 寫excel副標題
		f.SetCellValue("Template", "A6", "基本月費")
		// 寫入excel內容
		for i := range payRecordBasicCharges {
			cellValue := 7 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", cellValue)
			f.SetCellValue("Template", "A"+cellValueStr, payRecordBasicCharges[i].ItemName)
			f.SetCellValue("Template", "B"+cellValueStr, payRecordBasicCharges[i].StartDate.In(taipeiZone).Format("01-02")+" 到 "+payRecordBasicCharges[i].EndDate.In(taipeiZone).Format("01-02"))
			var price int
			if payRecordBasicCharges[i].Type == "charge" {
				price = payRecordBasicCharges[i].Price
			} else {
				price = -payRecordBasicCharges[i].Price
			}
			f.SetCellValue("Template", "E"+cellValueStr, price)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}

	// 補助款
	payRecordSubsidies := []PayRecordSubsidy{}
	json.Unmarshal(payRecord.Subsidy, &payRecordSubsidies)
	if len(payRecordSubsidies) > 0 {
		// 寫excel副標題
		subTitle := strconv.Itoa(7 + dataCount)
		f.InsertRow("Template", 7+dataCount)
		f.SetCellValue("Template", "A"+subTitle, "補助款")
		f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
		dataCount++
		// 寫入excel內容
		for i := range payRecordSubsidies {
			cellValue := 7 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", cellValue)
			f.SetCellValue("Template", "A"+cellValueStr, payRecordSubsidies[i].ItemName)
			f.SetCellValue("Template", "B"+cellValueStr, payRecordSubsidies[i].StartDate.In(taipeiZone).Format("01-02")+" 到 "+payRecordSubsidies[i].EndDate.In(taipeiZone).Format("01-02"))
			var price int
			if payRecordSubsidies[i].Type == "charge" {
				price = payRecordSubsidies[i].Price
			} else {
				price = -payRecordSubsidies[i].Price
			}
			f.SetCellValue("Template", "E"+cellValueStr, price)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}

	//異動(請假)
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	json.Unmarshal(payRecord.TransferRefundLeave, &payRecordTransferRefundLeaves)
	if len(payRecordTransferRefundLeaves) > 0 {
		subTitle := strconv.Itoa(7 + dataCount)
		f.InsertRow("Template", 7+dataCount)
		f.SetCellValue("Template", "A"+subTitle, "請假退費")
		f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
		dataCount++
		payRecordTransferRefundLeaveElements := make(map[string]*PayRecordTransferRefundLeaveForPrint)
		for _, d := range payRecordTransferRefundLeaves {
			if payRecordTransferRefundLeaveElements[d.ItemName] == nil {
				payRecordTransferRefundLeaveElements[d.ItemName] = &PayRecordTransferRefundLeaveForPrint{
					ItemName: d.ItemName,
					Type:     d.Type,
				}
				if d.Type == "charge" {
					payRecordTransferRefundLeaveElements[d.ItemName].Price = d.Price
				} else {
					payRecordTransferRefundLeaveElements[d.ItemName].Price = -d.Price
				}
				payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate)
				payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate)
			} else {
				payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate)
				payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate)
				if d.Type == "charge" {
					payRecordTransferRefundLeaveElements[d.ItemName].Price += d.Price
				} else {
					payRecordTransferRefundLeaveElements[d.ItemName].Price -= d.Price
				}
			}
		}

		for i := range payRecordTransferRefundLeaveElements {
			cellValue := 7 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", 7+dataCount)
			// 組時間內容(ex:02-10 到 02-20。02-21 到 02-22)
			var dateData []string
			// 由小到大排序
			sort.Slice(payRecordTransferRefundLeaveElements[i].StartDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[i].StartDates[j].Unix() < payRecordTransferRefundLeaveElements[i].StartDates[k].Unix()
			})
			sort.Slice(payRecordTransferRefundLeaveElements[i].EndDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[i].EndDates[j].Unix() < payRecordTransferRefundLeaveElements[i].EndDates[k].Unix()
			})
			for j := range payRecordTransferRefundLeaveElements[i].StartDates {
				dateData = append(dateData, payRecordTransferRefundLeaveElements[i].StartDates[j].Format("01-02")+"到"+payRecordTransferRefundLeaveElements[i].EndDates[j].Format("01-02"))
			}
			f.SetCellValue("Template", "A"+cellValueStr, i)
			f.SetCellValue("Template", "B"+cellValueStr, strings.Join(dateData, "。"))
			f.SetCellValue("Template", "E"+cellValueStr, payRecordTransferRefundLeaveElements[i].Price)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}
	// 非固定
	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
	json.Unmarshal(payRecord.NonFixedCharge, &payRecordNonFixedChargeRecords)
	if len(payRecordNonFixedChargeRecords) > 0 {
		// 寫excel副標題
		payRecordNonFixedChargeRecordElements := make(map[string]*PayRecordNonFixedChargeRecordDataForPrint)
		for _, d := range payRecordNonFixedChargeRecords {
			dateAndQuantity := make(map[string]*PayRecordNonFixedChargeRecordDateAndQuantityForPrint)
			itemCategoryKey := d.ItemCategory
			key := d.ItemName + d.Unit + strconv.Itoa(d.Price) + d.TaxType
			if payRecordNonFixedChargeRecordElements[itemCategoryKey] == nil {
				// 沒資料就把資料組一組 做成一個struct
				payRecordNonFixedChargeRecordElements[itemCategoryKey] = &PayRecordNonFixedChargeRecordDataForPrint{
					ItemCategory: d.ItemCategory,
				}
				dateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PayRecordNonFixedChargeRecordDateAndQuantityForPrint{
					Date:     d.NonFixedChargeDate,
					Quantity: d.Quantity,
				}
				payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas = map[string]*PayRecordNonFixedChargeRecordForPrint{}
				payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key] = &PayRecordNonFixedChargeRecordForPrint{
					Quantity: d.Quantity,
					Price:    d.Price,
					ItemName: d.ItemName,
					TaxType:  d.TaxType,
					Type:     d.Type,
					// Subtotal: d.Quantity * d.Price,
					DateAndQuantity: dateAndQuantity,
				}
			} else {
				// 確定這個類別有資料了
				// 檢查是不是一個新的品項+價格+單位的組合
				if payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key] == nil {
					dateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PayRecordNonFixedChargeRecordDateAndQuantityForPrint{
						Date:     d.NonFixedChargeDate,
						Quantity: d.Quantity,
					}

					payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key] = &PayRecordNonFixedChargeRecordForPrint{
						Quantity: d.Quantity,
						Price:    d.Price,
						ItemName: d.ItemName,
						TaxType:  d.TaxType,
						Type:     d.Type,
						// Subtotal: d.Quantity * d.Price,
						DateAndQuantity: dateAndQuantity,
					}
				} else {
					// 表示這個品項+價格+單位的組合 已經有了 那就是要新增個數 還有新增日期
					payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].Quantity += d.Quantity
					//  表示這個品項已經有了(有人耍白痴同一個品項 故意要分開新增)
					if payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")] == nil {
						payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PayRecordNonFixedChargeRecordDateAndQuantityForPrint{}
						payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Date = d.NonFixedChargeDate
						payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Quantity = d.Quantity
					} else {
						payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Date = d.NonFixedChargeDate
						payRecordNonFixedChargeRecordElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Quantity += d.Quantity
					}
				}
			}
		}

		// 確保key的順序一致(不做這段 用nonfixedChargeElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(payRecordNonFixedChargeRecordElements))
		for key := range payRecordNonFixedChargeRecordElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(i, j int) bool {
			return payRecordNonFixedChargeRecordElements[keys[i]].ItemCategory < payRecordNonFixedChargeRecordElements[keys[j]].ItemCategory
		})

		for i := range keys {
			// 寫excel副標題
			subTitle := strconv.Itoa(7 + dataCount)
			f.InsertRow("Template", 7+dataCount)
			f.SetCellValue("Template", "A"+subTitle, payRecordNonFixedChargeRecordElements[keys[i]].ItemCategory)
			f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
			dataCount++
			for _, d := range payRecordNonFixedChargeRecordElements[keys[i]].ItemCategoryDatas {
				cellValue := 7 + dataCount
				cellValueStr := strconv.Itoa(cellValue)
				f.InsertRow("Template", cellValue)
				f.SetCellValue("Template", "A"+cellValueStr, d.ItemName)
				var nonfixedChargeDate []string
				for _, d2 := range d.DateAndQuantity {
					nonfixedChargeDate = append(nonfixedChargeDate, d2.Date.Format("01-02")+"("+strconv.Itoa(d2.Quantity)+")")
				}
				f.SetCellValue("Template", "B"+cellValueStr, strings.Join(nonfixedChargeDate[:], "，"))
				var price int

				if d.Type == "charge" {
					price = d.Price
				} else {
					price = -d.Price
				}
				f.SetCellValue("Template", "C"+cellValueStr, price)
				f.SetCellValue("Template", "D"+cellValueStr, d.Quantity)
				f.SetCellValue("Template", "E"+cellValueStr, price*d.Quantity)
				f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
				dataCount++
			}
		}
	}

	// 塞一些機構資訊
	dataCount += 7
	cellValueStr := strconv.Itoa(dataCount)
	f.RemoveRow("Template", dataCount)
	var remittanceBank string
	if payRecord.Organization.RemittanceBank == nil {
		remittanceBank = ""
	} else {
		remittanceBank = *payRecord.Organization.RemittanceBank
	}
	f.SetCellValue("Template", "A"+cellValueStr, "銀行："+remittanceBank)
	f.SetCellValue("Template", "B"+cellValueStr, "備註："+payRecord.Note)
	f.SetCellFormula("Template", "E"+cellValueStr, "=SUM(E6:E"+strconv.Itoa(dataCount-1)+")")
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	var remittanceIdNumber string
	if payRecord.Organization.RemittanceIdNumber == nil {
		remittanceIdNumber = ""
	} else {
		remittanceIdNumber = *payRecord.Organization.RemittanceIdNumber
	}
	f.SetCellValue("Template", "A"+cellValueStr, "帳號："+remittanceIdNumber)
	f.SetCellFormula("Template", "E"+cellValueStr, "=E"+strconv.Itoa(dataCount-1))
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	var remittanceUserName string
	if payRecord.Organization.RemittanceUserName == nil {
		remittanceUserName = ""
	} else {
		remittanceUserName = *payRecord.Organization.RemittanceUserName
	}
	f.SetCellValue("Template", "A"+cellValueStr, "戶名："+remittanceUserName)
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	f.SetCellValue("Template", "C"+cellValueStr, "製表人："+payRecord.User.LastName+payRecord.User.FirstName)
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	f.SetCellValue("Template", "C"+cellValueStr, "製表日期："+time.Now().Format("2006-01-02"))

	// 組合sheetName
	sheetName := payRecordYear + "年" + payRecordMonth + "月收據明細"
	// 改sheet名稱
	f.SetSheetName("Template", sheetName)

	fileName := payRecord.Patient.Branch + payRecord.Patient.Room + payRecord.Patient.Bed + payRecord.Patient.LastName + payRecord.Patient.FirstName + " " + sheetName + payRecord.ReceiptNumber + ".xlsx"

	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPayRecordDetail f.SaveAs", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "payRecordDetail")
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail r.store.UploadFile", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "payRecordDetail")
	if err != nil {
		r.Logger.Error("PrintPayRecordDetail r.store.SetMetadata", zap.Error(err), zap.String("originalUrl", "printPayRecordDetail"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	publicUrl := r.store.GenPublicLink("payRecordDetail/" + fileName)
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPayRecordDetail run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPayRecordDetail"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

//go:embed excelTemplate/receiptPart.xlsx
var excelReceiptPart []byte

func (r *queryResolver) PrintPayRecordPart(ctx context.Context, payRecordIdStr string) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordPart uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecordId, err := uuid.Parse(payRecordIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordPart uuid.Parse(payRecordIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	payRecord, err := orm.GetPayRecordById(r.ORM.DB, payRecordId, true, true, false)
	if err != nil {
		r.Logger.Error("PrintPayRecordPart orm.GetPayRecordById", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	organizationReceiptTemplateSetting, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(r.ORM.DB, organizationId, payRecord.TaxType)
	if err != nil {
		r.Logger.Error("PrintPayRecordPart orm.GetOrganizationReceiptTemplateSettingInTaxType", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	reader := bytes.NewReader(excelReceiptPart)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPayRecordPart excelize.OpenReader", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

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

	organizationReceiptTemplateSettingPatientInfo := OrganizationReceiptTemplateSettingPatientInfo{
		ShowAreaAndClass:  organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "areaAndClass"),
		ShowBedAndRoom:    organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "bedAndRoom"),
		ShowSex:           organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "sex"),
		ShowBirthday:      organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "birthday"),
		ShowAge:           organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "age"),
		ShowCheckInDate:   organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "checkInDate"),
		ShowPatientNumber: organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "patientNumber"),
		ShowRecordNumber:  organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "recordNumber"),
		ShowIdNumber:      organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting, "idNumber"),
	}

	// 組合sheetName
	sheetName := payRecordYear + "年" + payRecordMonth + "月收據聯單"
	// 代表有作廢(要增加作廢說明和改sheetName)
	if isInvalid {
		sheetName += " -作廢"
	}

	f.NewSheet(sheetName)
	// 組一個塞住民資訊的struct
	excelPatientDataStruct := ExcelPatientDataStruct{
		f:                                  f,
		PayRecord:                          payRecord,
		SheetName:                          sheetName,
		PayRecordYear:                      payRecordYear,
		PayRecordMonth:                     payRecordMonth,
		InvalidText:                        invalidText,
		OrganizationReceiptTemplateSetting: organizationReceiptTemplateSetting,
		OrganizationReceiptTemplateSettingPatientInfo: organizationReceiptTemplateSettingPatientInfo,
	}
	// 塞住民資訊
	SetExcelPatientData(excelPatientDataStruct)

	// 金額框線
	priceStyle, err := excelStyle.GetPriceFormatAndTwoBorderAndFontStyle(f, 11, []string{"right", "bottom"}, "Calibri (本文)", []int{2, 1})
	if err != nil {
		r.Logger.Error("PrintPayRecordPart excelStyle.GetPriceFormatAndTwoBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	var insertRowCount int
	// 這邊才是在塞費用的內容
	if organizationReceiptTemplateSetting.PriceShowType == "classAddUp" {
		// 內容框線
		classStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 11)
		if err != nil {
			r.Logger.Error("PrintPayRecordPart excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
		insertRowCount = printClassAddUp(payRecord, f, classStyle, priceStyle, sheetName)
	} else {
		// 內容框線
		contentStyle, err := excelStyle.GetRightAndBottomBorderAndLeftAlignmentAndFontStyle(f)
		if err != nil {
			r.Logger.Error("PrintPayRecordPart excelStyle.GetRightAndBottomBorderAndLeftAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
		// 內容框線
		dateStyle, err := excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle(f, 10)
		if err != nil {
			r.Logger.Error("PrintPayRecordPart excelStyle.GetTopAndRightAndBottomBorderAndCenterAlignmentAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return "", err
		}
		insertRowCount = printItem(payRecord, f, contentStyle, dateStyle, priceStyle, sheetName)
	}

	// 塞機構資訊
	organizationReceiptTemplateSettingOrganizationInfo := OrganizationReceiptTemplateSettingOrganizationInfo{
		ShowTaxIdNumber:         organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSetting, "taxIdNumber"),
		ShowPhone:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSetting, "phone"),
		ShowFax:                 organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSetting, "fax"),
		ShowOwner:               organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSetting, "owner"),
		ShowEstablishmentNumber: organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "establishmentNumber"),
		ShowAddress:             organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "address"),
		ShowEmail:               organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "email"),
		ShowRemittanceBank:      organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "remittanceBank"),
		ShowRemittanceIdNumber:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "remittanceIdNumber"),
		ShowRemittanceUserName:  organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting, "remittanceUserName"),
	}
	excelOrganizationDataStruct := ExcelOrganizationDataStruct{
		f:                                  f,
		PayRecord:                          payRecord,
		OrganizationReceiptTemplateSetting: organizationReceiptTemplateSetting,
		OrganizationReceiptTemplateSettingOrganizationInfo: organizationReceiptTemplateSettingOrganizationInfo,
		SheetName:      sheetName,
		InsertRowCount: insertRowCount,
		IsInvalid:      isInvalid,
	}

	SetExcelOrganizationData(excelOrganizationDataStruct)

	// 下載印章圖片
	downloadExcelOrganizationSealStruct := DownloadExcelOrganizationSealStruct{
		OrganizationReceiptTemplateSetting: organizationReceiptTemplateSetting,
		r:                                  r,
		ctx:                                ctx,
	}

	err = DownloadExcelOrganizationSeal(downloadExcelOrganizationSealStruct)
	if err != nil {
		r.Logger.Error("PrintPayRecordPart DownloadExcelOrganizationSeal", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 把印章圖片塞到excel
	organizationSealDataStruct := OrganizationSealDataStruct{
		f:                                  f,
		OrganizationReceiptTemplateSetting: organizationReceiptTemplateSetting,
		SheetName:                          sheetName,
		InsertRowCount:                     insertRowCount,
	}

	err = SetExcelOrganizationSealData(organizationSealDataStruct)
	if err != nil {
		r.Logger.Error("PrintPayRecordPart SetExcelOrganizationSealData", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 改sheet名稱
	f.DeleteSheet("Template")
	// 組合sheetName
	fileName := payRecord.Patient.Branch + payRecord.Patient.Room + payRecord.Patient.Bed + payRecord.Patient.LastName + payRecord.Patient.FirstName + " " + payRecordYear + "年" + payRecordMonth + "月收據聯單" + payRecord.ReceiptNumber + ".xlsx"
	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPayRecordPart f.SaveAs", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordPart r.store.UploadFile", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordPart r.store.SetMetadata", zap.Error(err), zap.String("originalUrl", "printPayRecordPart"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	publicUrl := r.store.GenPublicLink("payRecordPart/" + fileName)

	// 刪掉下載的印章檔案
	if organizationReceiptTemplateSetting.OrganizationPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.OrganizationPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		os.Remove(fileName[1] + ".jpg")
	}
	if organizationReceiptTemplateSetting.SealOnePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealOnePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		os.Remove(fileName[1] + ".jpg")
	}
	if organizationReceiptTemplateSetting.SealTwoPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealTwoPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		os.Remove(fileName[1] + ".jpg")
	}
	if organizationReceiptTemplateSetting.SealThreePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealThreePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		os.Remove(fileName[1] + ".jpg")
	}
	if organizationReceiptTemplateSetting.SealFourPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealFourPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		os.Remove(fileName[1] + ".jpg")
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPayRecordPart run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPayRecordPart"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

//go:embed excelTemplate/receiptGeneralTabler.xlsx
var excelReceiptGeneralTablerBytes []byte

func (r *queryResolver) PrintPayRecordGeneralTable(ctx context.Context, BillDate time.Time) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordGeneralTable uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPayRecordGeneralTable uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	organization, err := orm.GetOrganizationById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable orm.GetOrganizationById", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable orm.GetUserById", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	reader := bytes.NewReader(excelReceiptGeneralTablerBytes)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelize.OpenReader", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	taipeiZone, _ := time.LoadLocation("Asia/Taipei")
	billDataTimeZoneInTaipei := BillDate.UTC().In(taipeiZone)

	year := billDataTimeZoneInTaipei.Year()
	month := int(billDataTimeZoneInTaipei.Month())
	f.SetCellValue("有效", "A1", organization.Name)
	f.SetCellValue("作廢", "A1", organization.Name)
	f.SetCellValue("有效", "G2", strconv.Itoa(year)+"年"+strconv.Itoa(month)+"月")
	f.SetCellValue("作廢", "G2", strconv.Itoa(year)+"年"+strconv.Itoa(month)+"月")

	validPayRecords, err := orm.GetPayRecordsByInvalidStatus(r.ORM.DB, organizationId, year, month, false)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable orm.GetPayRecordsByInvalidStatus is false", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	invalidPayRecords, err := orm.GetPayRecordsByInvalidStatus(r.ORM.DB, organizationId, year, month, true)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable orm.GetPayRecordsByInvalidStatus is true", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 純文字
	fontStyle, err := excelStyle.GetFontStyle(f, "Calibri (本文)")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 右框線
	rightBorderStyle, err := excelStyle.GetRightBorderAndFontStyle(f, "Calibri (本文)")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetRightBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	// 金額
	priceFormatAndFontStyle, err := excelStyle.GetPriceFormatAndFontStyle(f, "Calibri (本文)")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetPriceFormatAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 金額加左框線
	priceLeftBorderStyle, err := excelStyle.GetPriceFormatAndLeftBorderAndFontStyle(f)
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetPriceFormatAndLeftBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 金額加右框線
	priceRightBorderStyle, err := excelStyle.GetPriceFormatAndRightBorderAndFontStyle(f, "Calibri (本文)")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetPriceFormatAndLeftBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	// 金額加左右框線
	priceLeftAndRightBorderStyle, err := excelStyle.GetPriceFormatAndTwoBorderAndFontStyle(f, 11, []string{"right", "left"}, "儷宋 Pro", []int{2, 2})
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable excelStyle.GetPriceFormatAndTwoBorderAndFontStyle", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	var validRowCount int
	var validSubTotal int
	var validFixChargeSubTotal int
	var validNonFixChargeSubTotal int
	var validTransferRefundChargeSubTotal int
	var validPaidAmountSubTotal int
	var validUnpaidSubTotal int
	for i := range validPayRecords {
		rowNumberStr := strconv.Itoa(validRowCount + 4)
		f.InsertRow("有效", validRowCount+4)
		patient := validPayRecords[i].Patient
		f.SetCellValue("有效", "A"+rowNumberStr, validPayRecords[i].ReceiptNumber)
		f.SetCellValue("有效", "B"+rowNumberStr, validPayRecords[i].CreatedAt.Format("2006-01-02"))
		var taxType string
		if validPayRecords[i].TaxType == "allTax" {
			taxType = "所有項目"
		} else if validPayRecords[i].TaxType == "businessTax" {
			taxType = "營業稅"
		} else if validPayRecords[i].TaxType == "noTax" {
			taxType = "免稅"
		} else if validPayRecords[i].TaxType == "other" {
			taxType = "其他"
		} else if validPayRecords[i].TaxType == "stampTax" {
			taxType = "印花稅"
		}
		f.SetCellValue("有效", "C"+rowNumberStr, taxType)
		f.SetCellValue("有效", "D"+rowNumberStr, patient.Branch)
		f.SetCellValue("有效", "E"+rowNumberStr, patient.Room+patient.Bed)
		f.SetCellValue("有效", "F"+rowNumberStr, patient.LastName+patient.FirstName)
		var status string
		if patient.Status == "present" {
			status = "在院"
		} else if patient.Status == "hospital" {
			status = "住院"
		} else if patient.Status == "away" {
			status = "請假"
		} else if patient.Status == "withdraw" {
			status = "在院"
		} else if patient.Status == "reservation" {
			status = "退住"
		} else if patient.Status == "unpresented" {
			status = "取消預約"
		} else {
			status = ""
		}
		f.SetCellValue("有效", "G"+rowNumberStr, status)
		payRecordBasicCharges := []PayRecordBasicCharge{}
		json.Unmarshal(validPayRecords[i].BasicCharge, &payRecordBasicCharges)

		var subTotal int
		// 固定費用
		var basicChargePrice int
		for j := range payRecordBasicCharges {
			if payRecordBasicCharges[j].Type == "charge" {
				basicChargePrice += payRecordBasicCharges[j].Price
			} else {
				basicChargePrice -= payRecordBasicCharges[j].Price
			}
		}

		// 補助款
		payRecordSubsidies := []PayRecordSubsidy{}
		json.Unmarshal(validPayRecords[i].Subsidy, &payRecordSubsidies)
		for j := range payRecordSubsidies {
			if payRecordSubsidies[j].Type == "charge" {
				basicChargePrice += payRecordSubsidies[j].Price
			} else {
				basicChargePrice -= payRecordSubsidies[j].Price
			}
		}
		f.SetCellValue("有效", "I"+rowNumberStr, basicChargePrice)
		validFixChargeSubTotal += basicChargePrice
		// 異動(住院)
		var transferRefundPrice int
		//異動(請假)
		payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
		json.Unmarshal(validPayRecords[i].TransferRefundLeave, &payRecordTransferRefundLeaves)
		for j := 0; j < len(payRecordTransferRefundLeaves); j++ {
			if payRecordTransferRefundLeaves[j].Type == "charge" {
				transferRefundPrice += payRecordTransferRefundLeaves[j].Price
			} else {
				transferRefundPrice -= payRecordTransferRefundLeaves[j].Price
			}
		}
		f.SetCellValue("有效", "K"+rowNumberStr, transferRefundPrice)
		validTransferRefundChargeSubTotal += transferRefundPrice
		// 非固定
		var nonFixedChargePrice int
		payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
		json.Unmarshal(validPayRecords[i].NonFixedCharge, &payRecordNonFixedChargeRecords)
		for j := 0; j < len(payRecordNonFixedChargeRecords); j++ {
			if payRecordNonFixedChargeRecords[j].Type == "charge" {
				nonFixedChargePrice += payRecordNonFixedChargeRecords[j].Subtotal
			} else {
				nonFixedChargePrice -= payRecordNonFixedChargeRecords[j].Subtotal
			}
		}
		f.SetCellValue("有效", "J"+rowNumberStr, nonFixedChargePrice)
		validNonFixChargeSubTotal += nonFixedChargePrice
		subTotal += basicChargePrice + transferRefundPrice + nonFixedChargePrice
		// 應繳金額
		f.SetCellValue("有效", "H"+rowNumberStr, subTotal)
		validSubTotal += subTotal
		// 備註
		f.SetCellValue("有效", "L"+rowNumberStr, validPayRecords[i].Note)
		// 已繳金額
		f.SetCellValue("有效", "M"+rowNumberStr, validPayRecords[i].PaidAmount)
		validPaidAmountSubTotal += validPayRecords[i].PaidAmount
		// 未繳金額
		unpaid := subTotal - validPayRecords[i].PaidAmount
		f.SetCellValue("有效", "N"+rowNumberStr, unpaid)
		validUnpaidSubTotal += unpaid
		// style
		f.SetCellStyle("有效", "A"+rowNumberStr, "G"+rowNumberStr, fontStyle)
		f.SetCellStyle("有效", "H"+rowNumberStr, "H"+rowNumberStr, priceLeftAndRightBorderStyle)
		f.SetCellStyle("有效", "K"+rowNumberStr, "K"+rowNumberStr, priceRightBorderStyle)
		f.SetCellStyle("有效", "M"+rowNumberStr, "M"+rowNumberStr, priceLeftBorderStyle)
		f.SetCellStyle("有效", "N"+rowNumberStr, "N"+rowNumberStr, priceRightBorderStyle)

		f.SetCellStyle("有效", "I"+rowNumberStr, "J"+rowNumberStr, priceFormatAndFontStyle)
		f.SetCellStyle("有效", "R"+rowNumberStr, "R"+rowNumberStr, priceFormatAndFontStyle)
		f.SetCellStyle("有效", "T"+rowNumberStr, "T"+rowNumberStr, rightBorderStyle)

		// 繳費記錄
		for j := range validPayRecords[i].PayRecordDetails {
			if j != 0 {
				validRowCount++
				rowNumberStr = strconv.Itoa(validRowCount + 4)
				f.InsertRow("有效", validRowCount+4)
				// 這邊只是在補框線而已
				f.SetCellStyle("有效", "T"+rowNumberStr, "T"+rowNumberStr, rightBorderStyle)
				f.SetCellStyle("有效", "H"+rowNumberStr, "H"+rowNumberStr, priceLeftAndRightBorderStyle)
				f.SetCellStyle("有效", "K"+rowNumberStr, "K"+rowNumberStr, priceRightBorderStyle)
				f.SetCellStyle("有效", "M"+rowNumberStr, "M"+rowNumberStr, priceLeftBorderStyle)
				f.SetCellStyle("有效", "N"+rowNumberStr, "N"+rowNumberStr, priceRightBorderStyle)
			}
			f.SetCellValue("有效", "O"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].RecordDate.Format("2006-01-02"))
			f.SetCellValue("有效", "P"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].Method)
			f.SetCellValue("有效", "Q"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].Payer)
			f.SetCellValue("有效", "R"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].Price)
			f.SetCellValue("有效", "S"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].Note)
			f.SetCellValue("有效", "T"+rowNumberStr, validPayRecords[i].PayRecordDetails[j].Handler)
			f.SetCellStyle("有效", "O"+rowNumberStr, "Q"+rowNumberStr, fontStyle)
			f.SetCellStyle("有效", "R"+rowNumberStr, "R"+rowNumberStr, priceFormatAndFontStyle)
			f.SetCellStyle("有效", "S"+rowNumberStr, "S"+rowNumberStr, fontStyle)
		}
		validRowCount++
	}

	f.SetCellValue("有效", "H"+strconv.Itoa(validRowCount+5), validSubTotal)
	f.SetCellValue("有效", "I"+strconv.Itoa(validRowCount+5), validFixChargeSubTotal)
	f.SetCellValue("有效", "J"+strconv.Itoa(validRowCount+5), validNonFixChargeSubTotal)
	f.SetCellValue("有效", "K"+strconv.Itoa(validRowCount+5), validTransferRefundChargeSubTotal)
	f.SetCellValue("有效", "M"+strconv.Itoa(validRowCount+5), validPaidAmountSubTotal)
	f.SetCellValue("有效", "N"+strconv.Itoa(validRowCount+5), validUnpaidSubTotal)
	f.SetCellValue("有效", "T"+strconv.Itoa(validRowCount+6), "製表人："+user.DisplayName)
	f.SetCellValue("有效", "T"+strconv.Itoa(validRowCount+7), "製表日期："+time.Now().Format("2006-01-02"))

	f.RemoveRow("有效", validRowCount+4)
	f.SetSheetName("有效", strconv.Itoa(year)+"-"+strconv.Itoa(month)+"有效")

	// 作廢
	var invalidRowCount int
	var invalidSubTotal int
	var invalidFixChargeSubTotal int
	var invalidNonFixChargeSubTotal int
	var invalidTransferRefundChargeSubTotal int
	var invalidPaidAmountSubTotal int
	var invalidUnpaidSubTotal int
	for i := range invalidPayRecords {
		rowNumberStr := strconv.Itoa(invalidRowCount + 4)
		f.InsertRow("作廢", invalidRowCount+4)
		patient := invalidPayRecords[i].Patient
		f.SetCellValue("作廢", "A"+rowNumberStr, invalidPayRecords[i].ReceiptNumber)
		f.SetCellValue("作廢", "B"+rowNumberStr, invalidPayRecords[i].CreatedAt.Format("2006-01-02"))
		var taxType string
		if invalidPayRecords[i].TaxType == "allTax" {
			taxType = "所有項目"
		} else if invalidPayRecords[i].TaxType == "businessTax" {
			taxType = "營業稅"
		} else if invalidPayRecords[i].TaxType == "noTax" {
			taxType = "免稅"
		} else if invalidPayRecords[i].TaxType == "other" {
			taxType = "其他"
		} else if invalidPayRecords[i].TaxType == "stampTax" {
			taxType = "印花稅"
		}
		f.SetCellValue("作廢", "C"+rowNumberStr, taxType)
		f.SetCellValue("作廢", "D"+rowNumberStr, patient.Branch)
		f.SetCellValue("作廢", "E"+rowNumberStr, patient.Room+patient.Bed)
		f.SetCellValue("作廢", "F"+rowNumberStr, patient.LastName+patient.FirstName)
		var status string
		if patient.Status == "present" {
			status = "在院"
		} else if patient.Status == "hospital" {
			status = "住院"
		} else if patient.Status == "away" {
			status = "請假"
		} else if patient.Status == "withdraw" {
			status = "在院"
		} else if patient.Status == "reservation" {
			status = "退住"
		} else if patient.Status == "unpresented" {
			status = "取消預約"
		} else {
			status = ""
		}
		f.SetCellValue("作廢", "G"+rowNumberStr, status)
		payRecordBasicCharges := []PayRecordBasicCharge{}
		json.Unmarshal(invalidPayRecords[i].BasicCharge, &payRecordBasicCharges)

		var subTotal int
		// 固定費用
		var basicChargePrice int
		for j := range payRecordBasicCharges {
			if payRecordBasicCharges[j].Type == "charge" {
				basicChargePrice += payRecordBasicCharges[j].Price
			} else {
				basicChargePrice -= payRecordBasicCharges[j].Price
			}
		}

		// 補助款
		payRecordSubsidies := []PayRecordSubsidy{}
		json.Unmarshal(invalidPayRecords[i].Subsidy, &payRecordSubsidies)
		for j := range payRecordSubsidies {
			if payRecordSubsidies[j].Type == "charge" {
				basicChargePrice += payRecordSubsidies[j].Price
			} else {
				basicChargePrice -= payRecordSubsidies[j].Price
			}
		}
		f.SetCellValue("作廢", "I"+rowNumberStr, basicChargePrice)
		invalidFixChargeSubTotal += basicChargePrice
		// 異動(住院)
		var transferRefundPrice int
		//異動(請假)
		payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
		json.Unmarshal(invalidPayRecords[i].TransferRefundLeave, &payRecordTransferRefundLeaves)
		for j := 0; j < len(payRecordTransferRefundLeaves); j++ {
			if payRecordTransferRefundLeaves[j].Type == "charge" {
				transferRefundPrice += payRecordTransferRefundLeaves[j].Price
			} else {
				transferRefundPrice -= payRecordTransferRefundLeaves[j].Price
			}
		}
		f.SetCellValue("有效", "K"+rowNumberStr, transferRefundPrice)
		f.SetCellValue("作廢", "K"+rowNumberStr, transferRefundPrice)
		invalidTransferRefundChargeSubTotal += transferRefundPrice
		// 非固定
		var nonFixedChargePrice int
		payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
		json.Unmarshal(invalidPayRecords[i].NonFixedCharge, &payRecordNonFixedChargeRecords)
		for j := 0; j < len(payRecordNonFixedChargeRecords); j++ {
			if payRecordNonFixedChargeRecords[j].Type == "charge" {
				nonFixedChargePrice += payRecordNonFixedChargeRecords[j].Subtotal
			} else {
				nonFixedChargePrice -= payRecordNonFixedChargeRecords[j].Subtotal
			}
		}
		f.SetCellValue("作廢", "J"+rowNumberStr, nonFixedChargePrice)
		invalidNonFixChargeSubTotal += nonFixedChargePrice

		subTotal += basicChargePrice + transferRefundPrice + nonFixedChargePrice
		// 應繳金額
		f.SetCellValue("作廢", "H"+rowNumberStr, subTotal)
		invalidSubTotal += subTotal
		// 備註
		f.SetCellValue("作廢", "L"+rowNumberStr, invalidPayRecords[i].Note)
		// 已繳金額
		f.SetCellValue("作廢", "M"+rowNumberStr, invalidPayRecords[i].PaidAmount)
		invalidPaidAmountSubTotal += invalidPayRecords[i].PaidAmount
		// 未繳金額
		unpaid := subTotal - invalidPayRecords[i].PaidAmount
		f.SetCellValue("作廢", "N"+rowNumberStr, unpaid)
		invalidUnpaidSubTotal += unpaid
		// style
		f.SetCellStyle("作廢", "A"+rowNumberStr, "G"+rowNumberStr, fontStyle)
		f.SetCellStyle("作廢", "H"+rowNumberStr, "H"+rowNumberStr, priceLeftAndRightBorderStyle)
		f.SetCellStyle("作廢", "K"+rowNumberStr, "K"+rowNumberStr, priceRightBorderStyle)
		f.SetCellStyle("作廢", "M"+rowNumberStr, "M"+rowNumberStr, priceLeftBorderStyle)
		f.SetCellStyle("作廢", "N"+rowNumberStr, "N"+rowNumberStr, priceRightBorderStyle)

		f.SetCellStyle("作廢", "I"+rowNumberStr, "J"+rowNumberStr, priceFormatAndFontStyle)
		f.SetCellStyle("作廢", "R"+rowNumberStr, "R"+rowNumberStr, priceFormatAndFontStyle)
		f.SetCellStyle("作廢", "T"+rowNumberStr, "T"+rowNumberStr, rightBorderStyle)

		// 繳費記錄
		for j := range invalidPayRecords[i].PayRecordDetails {
			if j != 0 {
				invalidRowCount++
				rowNumberStr = strconv.Itoa(invalidRowCount + 4)
				f.InsertRow("作廢", invalidRowCount+4)
				// 這邊只是在補框線而已
				f.SetCellStyle("作廢", "T"+rowNumberStr, "T"+rowNumberStr, rightBorderStyle)
				f.SetCellStyle("作廢", "H"+rowNumberStr, "H"+rowNumberStr, priceLeftAndRightBorderStyle)
				f.SetCellStyle("作廢", "K"+rowNumberStr, "K"+rowNumberStr, priceRightBorderStyle)
				f.SetCellStyle("作廢", "M"+rowNumberStr, "M"+rowNumberStr, priceLeftBorderStyle)
				f.SetCellStyle("作廢", "N"+rowNumberStr, "N"+rowNumberStr, priceRightBorderStyle)
			}
			f.SetCellValue("作廢", "O"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].RecordDate.Format("2006-01-02"))
			f.SetCellValue("作廢", "P"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].Method)
			f.SetCellValue("作廢", "Q"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].Payer)
			f.SetCellValue("作廢", "R"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].Price)
			f.SetCellValue("作廢", "S"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].Note)
			f.SetCellValue("作廢", "T"+rowNumberStr, invalidPayRecords[i].PayRecordDetails[j].Handler)
			f.SetCellStyle("作廢", "O"+rowNumberStr, "Q"+rowNumberStr, fontStyle)
			f.SetCellStyle("作廢", "R"+rowNumberStr, "R"+rowNumberStr, priceFormatAndFontStyle)
			f.SetCellStyle("作廢", "S"+rowNumberStr, "S"+rowNumberStr, fontStyle)
		}
		invalidRowCount++
	}

	f.SetCellValue("作廢", "H"+strconv.Itoa(invalidRowCount+5), invalidSubTotal)
	f.SetCellValue("作廢", "I"+strconv.Itoa(invalidRowCount+5), invalidFixChargeSubTotal)
	f.SetCellValue("作廢", "J"+strconv.Itoa(invalidRowCount+5), invalidNonFixChargeSubTotal)
	f.SetCellValue("作廢", "K"+strconv.Itoa(invalidRowCount+5), invalidTransferRefundChargeSubTotal)
	f.SetCellValue("作廢", "M"+strconv.Itoa(invalidRowCount+5), invalidPaidAmountSubTotal)
	f.SetCellValue("作廢", "N"+strconv.Itoa(invalidRowCount+5), invalidUnpaidSubTotal)
	f.SetCellValue("作廢", "T"+strconv.Itoa(invalidRowCount+6), "製表人："+user.DisplayName)
	f.SetCellValue("作廢", "T"+strconv.Itoa(invalidRowCount+7), "製表日期："+time.Now().Format("2006-01-02"))

	f.RemoveRow("作廢", invalidRowCount+4)
	f.SetSheetName("作廢", strconv.Itoa(year)+"-"+strconv.Itoa(month)+"作廢")

	fileName := strconv.Itoa(year) + "年" + strconv.Itoa(month) + "月繳費紀錄總表.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable f.SaveAs", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable r.store.UploadFile", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "payRecordPart")
	if err != nil {
		r.Logger.Error("PrintPayRecordGeneralTable r.store.SetMetadata", zap.Error(err), zap.String("originalUrl", "printPayRecordGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	publicUrl := r.store.GenPublicLink("payRecordPart/" + fileName)
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPayRecordGeneralTable run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPayRecordGeneralTable"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}
