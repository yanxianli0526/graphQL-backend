package resolvers

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/resolvers/excelStyle"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type TransferRefundItem struct {
	ItemName string `json:"itemName"`
	Price    int    `json:"price"`
	Type     string `json:"type"`
}

type PrintPatientBillLeaves struct {
	Dates    []string `json:"dates"`
	Subtotal int      `json:"subtotal"`
}

func (r *queryResolver) PrintPatientBill(ctx context.Context, patientBillIdStr string) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("PrintPatientBill uuid.Parse(userIdStr)", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPatientBill uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	patientBillId, err := uuid.Parse(patientBillIdStr)
	if err != nil {
		r.Logger.Warn("PrintPatientBill uuid.Parse(patientBillIdStr)", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("PrintPatientBill orm.GetUserById", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	paitentBill, err := orm.GetPatientBillById(r.ORM.DB, organizationId, patientBillId)
	if err != nil {
		r.Logger.Error("PrintPatientBill orm.GetPatientBillById", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	publicUrl, err := printPatientBill(r, ctx, paitentBill, user)
	if err != nil {
		r.Logger.Error("PrintPatientBill printPatientBill", zap.Error(err), zap.String("originalUrl", "printPatientBill"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPatientBill run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPatientBill"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

//go:embed excelTemplate/patientBill.xlsx
var excelPatientBillBytes []byte

func printPatientBill(r *queryResolver, ctx context.Context, patientBill *models.PatientBill, user *models.User) (string, error) {
	reader := bytes.NewReader(excelPatientBillBytes)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPatientBill excelize.OpenReader", zap.Error(err))
		return "", err
	}

	patientBillYear := strconv.Itoa(patientBill.BillYear)
	patientBillMonth := strconv.Itoa(patientBill.BillMonth)

	if patientBill.BillMonth < 10 {
		patientBillMonth = "0" + patientBillMonth
	}
	f.SetCellValue("Template", "A1", patientBill.Organization.Name)
	f.SetCellValue("Template", "A2", patientBillYear+"年"+patientBillMonth+"月收費通知單")

	f.SetCellValue("Template", "A3", "床號："+patientBill.Patient.Room+patientBill.Patient.Bed)
	f.SetCellValue("Template", "B3", "姓名："+patientBill.Patient.LastName+patientBill.Patient.FirstName)
	var idNumber string
	if patientBill.Organization.Privacy == "unmask" {
		idNumber = patientBill.Patient.IdNumber
	} else {
		if len(patientBill.Patient.IdNumber) >= 10 {
			idNumber = patientBill.Patient.IdNumber[0:3] + "****" + patientBill.Patient.IdNumber[7:10]
		} else {
			idNumberCount := 0
			for idNumberCount < len(patientBill.Patient.IdNumber) {
				if idNumberCount >= 4 && idNumberCount <= 8 {
					idNumber += "*"
				} else {
					yo := string([]rune(patientBill.Patient.IdNumber)[idNumberCount])
					idNumber += yo
				}
				idNumberCount += 1
			}
		}
	}
	f.SetCellValue("Template", "C3", "身分證字號："+idNumber)

	// 副標題底色
	subTitleStyle, err := excelStyle.GetFillColorStyle(f)
	if err != nil {
		r.Logger.Error("PrintPatientBill excelStyle.GetFillColorStyle", zap.Error(err))
		return "", err
	}
	// 內容框線
	contentBorderStyle, err := excelStyle.GetAllBorderStyle(f)
	if err != nil {
		r.Logger.Error("PrintPatientBill excelStyle.GetAllBorderStyle", zap.Error(err))
		return "", err
	}
	dataCount := 0

	taipeiZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return "", err
	}
	// 基本月費
	if len(patientBill.BasicCharges) > 0 {
		// 寫excel副標題
		f.SetCellValue("Template", "A5", "基本月費")
		// 寫入excel內容
		for i := range patientBill.BasicCharges {
			cellValue := 6 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", cellValue)
			f.SetCellValue("Template", "A"+cellValueStr, patientBill.BasicCharges[i].ItemName)
			f.SetCellValue("Template", "B"+cellValueStr, patientBill.BasicCharges[i].StartDate.In(taipeiZone).Format("01-02")+" 到 "+patientBill.BasicCharges[i].EndDate.In(taipeiZone).Format("01-02"))
			var price int
			if patientBill.BasicCharges[i].Type == "charge" {
				price = patientBill.BasicCharges[i].Price
			} else {
				price = -patientBill.BasicCharges[i].Price
			}
			f.SetCellValue("Template", "E"+cellValueStr, price)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}
	// 補助款
	if len(patientBill.Subsidies) > 0 {
		// 寫excel副標題
		subTitle := strconv.Itoa(6 + dataCount)
		f.InsertRow("Template", 6+dataCount)
		f.SetCellValue("Template", "A"+subTitle, "補助款")
		f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
		dataCount++
		// 寫入excel內容
		for i := range patientBill.Subsidies {
			cellValue := 6 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", cellValue)
			f.SetCellValue("Template", "A"+cellValueStr, patientBill.Subsidies[i].ItemName)
			f.SetCellValue("Template", "B"+cellValueStr, patientBill.Subsidies[i].StartDate.Format("01-02")+" 到 "+patientBill.Subsidies[i].EndDate.Format("01-02"))
			var price int
			if patientBill.Subsidies[i].Type == "charge" {
				price = patientBill.Subsidies[i].Price
			} else {
				price = -patientBill.Subsidies[i].Price
			}
			f.SetCellValue("Template", "E"+cellValueStr, price)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}

	// 異動退費(請假)
	if len(patientBill.TransferRefundLeaves) > 0 {
		// 寫excel副標題
		subTitle := strconv.Itoa(6 + dataCount)
		f.InsertRow("Template", 6+dataCount)
		f.SetCellValue("Template", "A"+subTitle, "請假退費")
		f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
		dataCount++

		leaveElements := make(map[string]*PrintPatientBillLeaves)
		// 這邊是在整理資料
		for i := range patientBill.TransferRefundLeaves {
			// 先把請假裡面的item解成json
			transferRefundItems := []TransferRefundItem{}
			json.Unmarshal(patientBill.TransferRefundLeaves[i].Items, &transferRefundItems)
			// 用item去跑迴圈
			for j := range transferRefundItems {
				date := patientBill.TransferRefundLeaves[i].StartDate.Format("01-02") + " 到 " + patientBill.TransferRefundLeaves[i].EndDate.Format("01-02")
				var price int
				// 看看是收費還是退費
				if transferRefundItems[j].Type == "refund" {
					price = -transferRefundItems[j].Price
				} else {
					price = transferRefundItems[j].Price
				}
				// 用itemName當作key 沒有值的話 就給錢跟日期
				if leaveElements[transferRefundItems[j].ItemName] == nil {
					leaveElements[transferRefundItems[j].ItemName] = &PrintPatientBillLeaves{
						Subtotal: price,
					}
					leaveElements[transferRefundItems[j].ItemName].Dates = append(leaveElements[transferRefundItems[j].ItemName].Dates, date)
				} else {
					// 用itemName當作key 有的值的話就加上去
					leaveElements[transferRefundItems[j].ItemName].Subtotal += price
					leaveElements[transferRefundItems[j].ItemName].Dates = append(leaveElements[transferRefundItems[j].ItemName].Dates, date)
				}
			}
		}
		// 寫入excel內容
		for i, v := range leaveElements {
			cellValue := 6 + dataCount
			cellValueStr := strconv.Itoa(cellValue)
			f.InsertRow("Template", cellValue)
			f.SetCellValue("Template", "A"+cellValueStr, i)
			f.SetCellValue("Template", "B"+cellValueStr, strings.Join(v.Dates[:], "，"))
			f.SetCellValue("Template", "E"+cellValueStr, v.Subtotal)
			f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
			dataCount++
		}
	}
	// 非固定
	if len(patientBill.NonFixedChargeRecords) > 0 {
		// 寫excel副標題
		nonfixedChargeElements := make(map[string]*PatientBillNonFixedChargeRocordData)
		// 這邊是在整理資料
		for _, d := range patientBill.NonFixedChargeRecords {
			dateAndQuantity := make(map[string]*PatientBillNonFixedChargeRocordDateAndQuantity)
			itemCategoryKey := d.ItemCategory
			key := d.ItemName + d.Unit + strconv.Itoa(d.Price) + d.TaxType
			// 先看這個類別有沒有資料了
			if nonfixedChargeElements[itemCategoryKey] == nil {
				// 沒資料就把資料組一組 做成一個struct
				nonfixedChargeElements[itemCategoryKey] = &PatientBillNonFixedChargeRocordData{
					ItemCategory: d.ItemCategory,
				}
				dateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PatientBillNonFixedChargeRocordDateAndQuantity{
					Date:     d.NonFixedChargeDate,
					Quantity: d.Quantity,
				}
				nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas = map[string]*PatientBillNonFixedChargeRocord{}
				nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key] = &PatientBillNonFixedChargeRocord{
					Quantity:        d.Quantity,
					Price:           d.Price,
					ItemName:        d.ItemName,
					TaxType:         d.TaxType,
					Type:            d.Type,
					EarliestDate:    d.NonFixedChargeDate,
					DateAndQuantity: dateAndQuantity,
				}
			} else {
				// 確定這個類別有資料了
				// 檢查是不是一個新的品項+價格+單位的組合
				if nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key] == nil {
					dateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PatientBillNonFixedChargeRocordDateAndQuantity{
						Date:     d.NonFixedChargeDate,
						Quantity: d.Quantity,
					}

					nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key] = &PatientBillNonFixedChargeRocord{
						Quantity:        d.Quantity,
						Price:           d.Price,
						ItemName:        d.ItemName,
						EarliestDate:    d.NonFixedChargeDate,
						DateAndQuantity: dateAndQuantity,
					}
				} else {
					// 表示這個品項+價格+單位的組合 已經有了 那就是要新增個數 還有新增日期
					nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].Quantity += d.Quantity
					// 為了最後列印的排序 可以把時間最早的排在上面
					if d.NonFixedChargeDate.Unix() < nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].EarliestDate.Unix() {
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].EarliestDate = d.NonFixedChargeDate
					}
					//  表示這個品項已經有了(有人耍白痴同一個品項 故意要分開新增)
					if nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")] == nil {
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")] = &PatientBillNonFixedChargeRocordDateAndQuantity{}
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Date = d.NonFixedChargeDate
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Quantity = d.Quantity
					} else {
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Date = d.NonFixedChargeDate
						nonfixedChargeElements[itemCategoryKey].ItemCategoryDatas[key].DateAndQuantity[d.NonFixedChargeDate.Format("01-02")].Quantity += d.Quantity
					}
				}
			}
		}

		// 確保key的順序一致(不做這段 用nonfixedChargeElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(nonfixedChargeElements))
		for key := range nonfixedChargeElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(i, j int) bool {
			return nonfixedChargeElements[keys[i]].ItemCategory < nonfixedChargeElements[keys[j]].ItemCategory
		})

		for i := range keys {
			// 寫excel副標題
			subTitle := strconv.Itoa(6 + dataCount)
			f.InsertRow("Template", 6+dataCount)
			f.SetCellValue("Template", "A"+subTitle, nonfixedChargeElements[keys[i]].ItemCategory)
			f.SetCellStyle("Template", "A"+subTitle, "E"+subTitle, subTitleStyle)
			dataCount++

			// 確保key2的順序一致(不做這段 用nonfixedChargeElements[keys[i]].ItemCategoryDatas跑迴圈順序會亂跳)
			keys2 := make([]string, 0, len(nonfixedChargeElements[keys[i]].ItemCategory))
			for key := range nonfixedChargeElements[keys[i]].ItemCategoryDatas {
				keys2 = append(keys2, key)
			}
			// 把時間做一下排序
			sort.SliceStable(keys2, func(j, k int) bool {
				return nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].EarliestDate.Unix() < nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[k]].EarliestDate.Unix()
			})
			for j := range keys2 {
				cellValue := 6 + dataCount
				cellValueStr := strconv.Itoa(cellValue)
				f.InsertRow("Template", cellValue)
				f.SetCellValue("Template", "A"+cellValueStr, nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].ItemName)
				var nonfixedChargeDate []string
				for _, d2 := range nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].DateAndQuantity {
					nonfixedChargeDate = append(nonfixedChargeDate, d2.Date.Format("01-02")+"("+strconv.Itoa(d2.Quantity)+")")
				}
				f.SetCellValue("Template", "B"+cellValueStr, strings.Join(nonfixedChargeDate[:], "，"))
				var price int
				if nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].Type == "refund" {
					price = -nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].Price
				} else {
					price = nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].Price
				}
				f.SetCellValue("Template", "C"+cellValueStr, price)
				f.SetCellValue("Template", "D"+cellValueStr, nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].Quantity)
				f.SetCellValue("Template", "E"+cellValueStr, price*nonfixedChargeElements[keys[i]].ItemCategoryDatas[keys2[j]].Quantity)
				f.SetCellStyle("Template", "A"+cellValueStr, "E"+cellValueStr, contentBorderStyle)
				dataCount++
			}
		}
	}
	// 塞一些機構資訊
	dataCount += 6
	cellValueStr := strconv.Itoa(dataCount)
	f.RemoveRow("Template", dataCount)
	var remittanceBank string
	if patientBill.Organization.RemittanceBank == nil {
		remittanceBank = ""
	} else {
		remittanceBank = *patientBill.Organization.RemittanceBank
	}
	f.SetCellValue("Template", "A"+cellValueStr, "銀行："+remittanceBank)
	f.SetCellValue("Template", "B"+cellValueStr, "備註："+patientBill.Note)
	f.SetCellFormula("Template", "E"+cellValueStr, "=SUM(E6:E"+strconv.Itoa(dataCount-1)+")")
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	var remittanceIdNumber string
	if patientBill.Organization.RemittanceIdNumber == nil {
		remittanceIdNumber = ""
	} else {
		remittanceIdNumber = *patientBill.Organization.RemittanceIdNumber
	}
	f.SetCellValue("Template", "A"+cellValueStr, "帳號："+remittanceIdNumber)
	f.SetCellFormula("Template", "E"+cellValueStr, "=E"+strconv.Itoa(dataCount-1))
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	var remittanceUserName string
	if patientBill.Organization.RemittanceUserName == nil {
		remittanceUserName = ""
	} else {
		remittanceUserName = *patientBill.Organization.RemittanceUserName
	}
	f.SetCellValue("Template", "A"+cellValueStr, "戶名："+remittanceUserName)
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	f.SetCellValue("Template", "C"+cellValueStr, "製表人："+user.LastName+user.FirstName)
	dataCount++
	cellValueStr = strconv.Itoa(dataCount)
	f.SetCellValue("Template", "C"+cellValueStr, "製表日期："+time.Now().Format("2006-01-02"))

	// 組合sheetName
	sheetName := patientBillYear + "年"
	patientBillMonthInt, _ := strconv.Atoi(patientBillMonth)
	if patientBillMonthInt < 10 {
		sheetName += "0" + patientBillMonth + "月繳費通知單"
	} else {
		sheetName += patientBillMonth + "月繳費通知單"
	}
	// 改sheet名稱
	f.SetSheetName("Template", sheetName)

	fileName := patientBill.Patient.Branch + patientBill.Patient.Room + patientBill.Patient.Bed + patientBill.Patient.LastName + patientBill.Patient.FirstName + "-" + sheetName + ".xlsx"
	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPatientBill f.SaveAs", zap.Error(err))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "patientBill")
	if err != nil {
		r.Logger.Error("PrintPatientBill r.store.UploadFile", zap.Error(err))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "patientBill")
	if err != nil {
		r.Logger.Error("PrintPatientBill r.store.SetMetadata", zap.Error(err))
		return "", err
	}
	publicUrl := r.store.GenPublicLink("patientBill/" + fileName)

	return publicUrl, nil
}

//go:embed excelTemplate/patientBillGeneralTable.xlsx
var excelPatientBillGeneralTableBytes []byte

func (r *queryResolver) PrintPatientBillGeneralTable(ctx context.Context, BillDate time.Time) (string, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("PrintPatientBillGeneralTable uuid.Parse(userIdStr)", zap.Error(err))
		return "", err
	}
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintPatientBillGeneralTable uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "printPatientBillGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable orm.GetUserById", zap.Error(err), zap.String("originalUrl", "printPatientBillGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	organization, err := orm.GetOrganizationById(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable orm.GetOrganizationById", zap.Error(err), zap.String("originalUrl", "printPatientBillGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	taipeiZone, _ := time.LoadLocation("Asia/Taipei")
	billDataTimeZoneInTaipei := BillDate.UTC().In(taipeiZone)

	paitentBills, err := orm.GetPatientBillsByDate(r.ORM.DB, organizationId, billDataTimeZoneInTaipei.Year(), int(billDataTimeZoneInTaipei.Month()))
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable orm.GetPatientBillsByDate", zap.Error(err), zap.String("originalUrl", "printPatientBillGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}

	sort.SliceStable(paitentBills, func(i, j int) bool {
		if paitentBills[i].Patient.Branch > paitentBills[j].Patient.Branch {
			return false
		}
		if paitentBills[i].Patient.Branch < paitentBills[j].Patient.Branch {
			return true
		}
		if paitentBills[i].Patient.Room > paitentBills[j].Patient.Room {
			return false
		}
		if paitentBills[i].Patient.Room < paitentBills[j].Patient.Room {
			return true
		}
		if paitentBills[i].Patient.Bed > paitentBills[j].Patient.Bed {
			return false
		}
		if paitentBills[i].Patient.Bed < paitentBills[j].Patient.Bed {
			return true
		}
		if paitentBills[i].Patient.Status > paitentBills[j].Patient.Status {
			return false
		}
		if paitentBills[i].Patient.Status < paitentBills[j].Patient.Status {
			return true
		}
		if paitentBills[i].Patient.LastName > paitentBills[j].Patient.LastName {
			return false
		}
		if paitentBills[i].Patient.LastName < paitentBills[j].Patient.LastName {
			return true
		}
		if paitentBills[i].Patient.FirstName > paitentBills[j].Patient.FirstName {
			return false
		}
		if paitentBills[i].Patient.FirstName < paitentBills[j].Patient.FirstName {
			return true
		}

		return false
	})

	publicUrl, err := printPatientBillGeneralTable(r, ctx, paitentBills, user, organization, billDataTimeZoneInTaipei.Year(), int(billDataTimeZoneInTaipei.Month()))
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable printPatientBillGeneralTable", zap.Error(err), zap.String("originalUrl", "printPatientBillGeneralTable"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return "", err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printPatientBillGeneralTable run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printPatientBillGeneralTable"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return publicUrl, nil
}

func printPatientBillGeneralTable(r *queryResolver, ctx context.Context, patientBills []*models.PatientBill, user *models.User, organization *models.Organization, year, month int) (string, error) {
	reader := bytes.NewReader(excelPatientBillGeneralTableBytes)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelize.OpenReader", zap.Error(err))
		return "", err
	}

	alignmentStyle, err := excelStyle.GetAlignmentAndFontStyle(f)
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelStyle.GetAlignmentAndFontStyle", zap.Error(err))
		return "", err
	}

	fontStyle, err := excelStyle.GetFontStyle(f, "儷宋 Pro")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelStyle.GetFontStyle", zap.Error(err))
		return "", err
	}
	rightBorderAndFontStyle, err := excelStyle.GetRightBorderAndFontStyle(f, "儷宋 Pro")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelStyle.GetRightBorderAndFontStyle", zap.Error(err))
		return "", err
	}

	priceFormatAndRightBorderAndFontStyle, err := excelStyle.GetPriceFormatAndRightBorderAndFontStyle(f, "儷宋 Pro")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelStyle.GetPriceFormatAndRightBorderAndFontStyle", zap.Error(err))
		return "", err
	}

	priceFormatAndFontStyle, err := excelStyle.GetPriceFormatAndFontStyle(f, "儷宋 Pro")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable excelStyle.GetPriceFormatAndFontStyle", zap.Error(err))
		return "", err
	}

	f.SetCellValue("Template", "A1", organization.Name)
	f.SetCellValue("Template", "F2", strconv.Itoa(year)+"年"+strconv.Itoa(month)+"月")
	var rowCount int
	for i := range patientBills {
		rowNumberStr := strconv.Itoa(i + 4)
		if i != 0 {
			f.InsertRow("Template", i+4)
		}
		patient := patientBills[i].Patient
		f.SetCellValue("Template", "A"+rowNumberStr, patient.Branch)
		f.SetCellValue("Template", "B"+rowNumberStr, patient.Room+patient.Bed)
		f.SetCellValue("Template", "C"+rowNumberStr, patient.LastName+patient.FirstName)
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
		f.SetCellValue("Template", "D"+rowNumberStr, status)
		f.SetCellValue("Template", "E"+rowNumberStr, patientBills[i].CreatedAt.Format("2006-01-02"))
		var subTotal int
		// 固定費用
		var basicChargePrice int
		for j := 0; j < len(patientBills[i].BasicCharges); j++ {
			if patientBills[i].BasicCharges[j].Type == "charge" {
				basicChargePrice += patientBills[i].BasicCharges[j].Price
			} else {
				basicChargePrice -= patientBills[i].BasicCharges[j].Price
			}
		}
		for j := 0; j < len(patientBills[i].Subsidies); j++ {
			if patientBills[i].Subsidies[j].Type == "charge" {
				basicChargePrice += patientBills[i].Subsidies[j].Price
			} else {
				basicChargePrice -= patientBills[i].Subsidies[j].Price
			}
		}
		f.SetCellValue("Template", "H"+rowNumberStr, basicChargePrice)
		// 異動費用
		var transferRefundPrice int
		for j := 0; j < len(patientBills[i].TransferRefundLeaves); j++ {
			transferRefundItems := []TransferRefundItem{}
			json.Unmarshal(patientBills[i].TransferRefundLeaves[j].Items, &transferRefundItems)
			for k := 0; k < len(transferRefundItems); k++ {
				if transferRefundItems[k].Type == "charge" {
					transferRefundPrice += transferRefundItems[k].Price
				} else {
					transferRefundPrice -= transferRefundItems[k].Price
				}
			}
		}
		f.SetCellValue("Template", "J"+rowNumberStr, transferRefundPrice)
		// 非固定
		var nonFixedChargePrice int
		for j := 0; j < len(patientBills[i].NonFixedChargeRecords); j++ {
			if patientBills[i].NonFixedChargeRecords[j].Type == "charge" {
				nonFixedChargePrice += patientBills[i].NonFixedChargeRecords[j].Subtotal
			} else {
				nonFixedChargePrice -= patientBills[i].NonFixedChargeRecords[j].Subtotal
			}
		}
		f.SetCellValue("Template", "I"+rowNumberStr, nonFixedChargePrice)
		// 已收金額
		f.SetCellValue("Template", "G"+rowNumberStr, patientBills[i].AmountReceived)
		// 應收金額
		subTotal += basicChargePrice + transferRefundPrice + nonFixedChargePrice
		f.SetCellValue("Template", "F"+rowNumberStr, subTotal)
		// 備註
		f.SetCellValue("Template", "K"+rowNumberStr, patientBills[i].Note)
		// 文字style
		f.SetCellStyle("Template", "A"+rowNumberStr, "E"+rowNumberStr, fontStyle)
		f.SetCellStyle("Template", "F"+rowNumberStr, "F"+rowNumberStr, priceFormatAndFontStyle)
		f.SetCellStyle("Template", "H"+rowNumberStr, "I"+rowNumberStr, priceFormatAndFontStyle)
		// 框線style
		f.SetCellStyle("Template", "K"+rowNumberStr, "K"+rowNumberStr, rightBorderAndFontStyle)
		// 金額框線style
		f.SetCellStyle("Template", "G"+rowNumberStr, "G"+rowNumberStr, priceFormatAndRightBorderAndFontStyle)
		f.SetCellStyle("Template", "J"+rowNumberStr, "J"+rowNumberStr, priceFormatAndRightBorderAndFontStyle)
		rowCount = i
	}

	rowNumberStr := strconv.Itoa(rowCount + 6)

	// 共計
	f.SetCellValue("Template", "A"+rowNumberStr, "共計 "+strconv.Itoa(rowCount+1)+" 筆")
	f.SetCellFormula("Template", "F"+rowNumberStr, "=SUM(F4:F"+strconv.Itoa(rowCount+4)+")")
	f.SetCellFormula("Template", "G"+rowNumberStr, "=SUM(G4:G"+strconv.Itoa(rowCount+4)+")")
	f.SetCellFormula("Template", "H"+rowNumberStr, "=SUM(H4:H"+strconv.Itoa(rowCount+4)+")")
	f.SetCellFormula("Template", "I"+rowNumberStr, "=SUM(I4:I"+strconv.Itoa(rowCount+4)+")")
	f.SetCellFormula("Template", "J"+rowNumberStr, "=SUM(J4:J"+strconv.Itoa(rowCount+4)+")")

	// 製表人 製表日期
	f.SetCellValue("Template", "K"+strconv.Itoa(rowCount+7), "製表人："+user.DisplayName)
	f.SetCellValue("Template", "K"+strconv.Itoa(rowCount+8), "製表日期："+time.Now().Format("2006-01-02"))
	// 文字style

	f.SetCellStyle("Template", "K"+strconv.Itoa(rowCount+7), "K"+strconv.Itoa(rowCount+7), alignmentStyle)
	f.SetCellStyle("Template", "K"+strconv.Itoa(rowCount+8), "K"+strconv.Itoa(rowCount+8), alignmentStyle)
	f.SetCellStyle("Template", "K"+strconv.Itoa(rowCount+7), "K"+strconv.Itoa(rowCount+7), alignmentStyle)
	f.SetCellStyle("Template", "K"+strconv.Itoa(rowCount+8), "K"+strconv.Itoa(rowCount+8), alignmentStyle)

	f.RemoveRow("Template", rowCount+5)
	f.SetSheetName("Template", strconv.Itoa(year)+"-"+strconv.Itoa(month))
	fileName := strconv.Itoa(year) + "年" + strconv.Itoa(month) + "月住民帳單總表.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable f.SaveAs", zap.Error(err))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "patientBill")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable r.store.UploadFile", zap.Error(err))
		return "", err
	}

	err = r.store.SetMetadata(ctx, fileName, "patientBill")
	if err != nil {
		r.Logger.Error("PrintPatientBillGeneralTable r.store.SetMetadata", zap.Error(err))
		return "", err
	}

	publicUrl := r.store.GenPublicLink("patientBill/" + fileName)
	return publicUrl, nil
}
