package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"graphql-go-template/internal/models"
	"image"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelPatientDataStruct struct {
	f                                             *excelize.File
	PayRecord                                     *models.PayRecord
	SheetName                                     string
	PayRecordYear                                 string
	PayRecordMonth                                string
	InvalidText                                   string
	OrganizationReceiptTemplateSetting            *models.OrganizationReceiptTemplateSetting
	OrganizationReceiptTemplateSettingPatientInfo OrganizationReceiptTemplateSettingPatientInfo
}

type ExcelOrganizationDataStruct struct {
	f                                                  *excelize.File
	PayRecord                                          *models.PayRecord
	OrganizationReceiptTemplateSetting                 *models.OrganizationReceiptTemplateSetting
	OrganizationReceiptTemplateSettingOrganizationInfo OrganizationReceiptTemplateSettingOrganizationInfo
	SheetName                                          string
	InsertRowCount                                     int
	IsInvalid                                          bool
}

type DownloadExcelOrganizationSealStruct struct {
	OrganizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting
	r                                  *queryResolver
	ctx                                context.Context
}

type OrganizationSealDataStruct struct {
	f                                  *excelize.File
	OrganizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting
	SheetName                          string
	InsertRowCount                     int
}

type CheckoutCellStruct struct {
	F              *excelize.File
	DataCount      int
	CellValue      int
	InsertRowCount int
	ContentStyle   int
	DateStyle      int
	PriceStyle     int
	Price          int
	SheetName      string
	ItemName       string
	DateStr        *string
}

type PrintItemByTaxTypeStruct struct {
	PayRecord  *models.PayRecord
	F          *excelize.File
	ClassStyle int
	DateStyle  int
	PriceStyle int
	TaxCount   int
	TaxType    string
	SheetName  string
}

// 塞一些住民資料
func SetExcelPatientData(excelPatientDataStruct ExcelPatientDataStruct) {
	sheetName := excelPatientDataStruct.SheetName
	f := excelPatientDataStruct.f
	payRecord := excelPatientDataStruct.PayRecord
	payRecordYear := excelPatientDataStruct.PayRecordYear
	payRecordMonth := excelPatientDataStruct.PayRecordMonth
	invalidText := excelPatientDataStruct.InvalidText
	organizationReceiptTemplateSetting := excelPatientDataStruct.OrganizationReceiptTemplateSetting
	organizationReceiptTemplateSettingPatientInfo := excelPatientDataStruct.OrganizationReceiptTemplateSettingPatientInfo
	// 設定版面列印的縮放比例改為最適大小
	f.SetSheetPrOptions(sheetName, excelize.FitToPage(true))
	// 設定版面列印的縮放比例的參數
	f.SetPageLayout(
		sheetName,
		excelize.FitToHeight(1),
		excelize.FitToWidth(1),
		excelize.PageLayoutPaperSize(9), // A4大小
	)
	// 設定版面列印的邊界值
	f.SetPageMargins(sheetName,
		excelize.PageMarginBottom(0.25),
		excelize.PageMarginFooter(0),
		excelize.PageMarginHeader(0),
		excelize.PageMarginLeft(0.25),
		excelize.PageMarginRight(0.25),
		excelize.PageMarginTop(0.25),
	)
	// 把該merge的欄位做一做
	f.MergeCell(sheetName, "A1", "A2")
	f.MergeCell(sheetName, "A19", "A20") // 聯單複寫
	f.MergeCell(sheetName, "B1", "F1")
	f.MergeCell(sheetName, "B19", "F19") // 聯單複寫
	f.MergeCell(sheetName, "B2", "F2")
	f.MergeCell(sheetName, "B20", "F20") // 聯單複寫
	f.MergeCell(sheetName, "A3", "B3")
	f.MergeCell(sheetName, "A21", "B21") // 聯單複寫
	f.MergeCell(sheetName, "A4", "F4")
	f.MergeCell(sheetName, "A22", "F22") // 聯單複寫
	f.MergeCell(sheetName, "A5", "B5")
	f.MergeCell(sheetName, "A23", "B23") // 聯單複寫
	f.MergeCell(sheetName, "D5", "E5")
	f.MergeCell(sheetName, "D23", "E23") // 聯單複寫
	f.MergeCell(sheetName, "A11", "B11")
	f.MergeCell(sheetName, "A29", "B29") // 聯單複寫
	f.MergeCell(sheetName, "D11", "E11")
	f.MergeCell(sheetName, "D29", "E29") // 聯單複寫
	f.MergeCell(sheetName, "A12", "F12")
	f.MergeCell(sheetName, "A30", "F30") // 聯單複寫
	f.MergeCell(sheetName, "A13", "F13")
	f.MergeCell(sheetName, "A31", "F31") // 聯單複寫
	f.MergeCell(sheetName, "A14", "F14")
	f.MergeCell(sheetName, "A32", "F32") // 聯單複寫

	//  調整每一欄位的寬度(這邊設定的值 最後再看的時候會差0.83=>excel為16)
	f.SetColWidth(sheetName, "A", "F", 16.83)

	//  調整每一列的高度
	f.SetRowHeight(sheetName, 1, 32)
	f.SetRowHeight(sheetName, 19, 32) // 聯單複寫
	f.SetRowHeight(sheetName, 2, 16)
	f.SetRowHeight(sheetName, 20, 16) // 聯單複寫
	f.SetRowHeight(sheetName, 3, 19)
	f.SetRowHeight(sheetName, 21, 19) // 聯單複寫
	f.SetRowHeight(sheetName, 4, 21)
	f.SetRowHeight(sheetName, 22, 21) // 聯單複寫
	f.SetRowHeight(sheetName, 5, 15.75)
	f.SetRowHeight(sheetName, 23, 15.75) // 聯單複寫
	f.SetRowHeight(sheetName, 6, 16)
	f.SetRowHeight(sheetName, 24, 16) // 聯單複寫
	f.SetRowHeight(sheetName, 7, 17)
	f.SetRowHeight(sheetName, 25, 17) // 聯單複寫
	f.SetRowHeight(sheetName, 8, 17)
	f.SetRowHeight(sheetName, 26, 17) // 聯單複寫
	f.SetRowHeight(sheetName, 9, 16)
	f.SetRowHeight(sheetName, 27, 16) // 聯單複寫
	f.SetRowHeight(sheetName, 10, 17)
	f.SetRowHeight(sheetName, 28, 17) // 聯單複寫
	f.SetRowHeight(sheetName, 11, 16.5)
	f.SetRowHeight(sheetName, 29, 16.5) // 聯單複寫
	f.SetRowHeight(sheetName, 12, 21)
	f.SetRowHeight(sheetName, 30, 21) // 聯單複寫
	f.SetRowHeight(sheetName, 13, 68)
	f.SetRowHeight(sheetName, 31, 68) // 聯單複寫
	f.SetRowHeight(sheetName, 14, 72)
	f.SetRowHeight(sheetName, 32, 72) // 聯單複寫
	f.SetRowHeight(sheetName, 15, 16)
	f.SetRowHeight(sheetName, 33, 16) // 聯單複寫
	f.SetRowHeight(sheetName, 16, 60)
	f.SetRowHeight(sheetName, 34, 60) // 聯單複寫
	f.SetRowHeight(sheetName, 17, 19.5)
	f.SetRowHeight(sheetName, 35, 19.5) // 聯單複寫
	f.SetRowHeight(sheetName, 18, 19.5)
	f.SetRowHeight(sheetName, 36, 19.5) // 聯單複寫
	f.SetRowHeight(sheetName, 19, 32)
	f.SetRowHeight(sheetName, 37, 32) // 聯單複寫
	f.SetRowHeight(sheetName, 20, 16)
	f.SetRowHeight(sheetName, 38, 16) // 聯單複寫
	f.SetRowHeight(sheetName, 21, 19)
	f.SetRowHeight(sheetName, 39, 19) // 聯單複寫
	f.SetRowHeight(sheetName, 22, 21)
	f.SetRowHeight(sheetName, 40, 21) // 聯單複寫

	// 開始宣告一狗票的style
	numberFormat := "#,##0 "
	rightStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", WrapText: false},
	})

	fillStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "left",
				Color: "#000000",
				Style: 2,
			},
		},
		Fill:      excelize.Fill{Color: []string{"#D3D3D4"}, Type: "pattern", Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	fillStyle2, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 2,
			},
		},
		Fill:      excelize.Fill{Color: []string{"#D3D3D4"}, Type: "pattern", Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	fillStyle3, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "bottom",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "left",
				Color: "#000000",
				Style: 2,
			},
		},
		Fill:      excelize.Fill{Color: []string{"#D3D3D4"}, Type: "pattern", Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	fillStyle4, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "bottom",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 2,
			},
		},
		Fill:         excelize.Fill{Color: []string{"#D3D3D4"}, Type: "pattern", Pattern: 1},
		CustomNumFmt: &numberFormat,
	})

	fillStyle5, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "bottom",
				Color: "#000000",
				Style: 2,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			},
		},
		Fill:      excelize.Fill{Color: []string{"#D3D3D4"}, Type: "pattern", Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "left",
				Color: "#000000",
				Style: 2,
			},
		},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
	})
	borderStyle1, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 1,
			},
		},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	})
	borderStyle2, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			}, {
				Type:  "right",
				Color: "#000000",
				Style: 2,
			},
		},
		CustomNumFmt: &numberFormat,
		Alignment:    &excelize.Alignment{Vertical: "center"},
	})

	borderStyle3, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 4,
			},
		},
		CustomNumFmt: &numberFormat,
		Alignment:    &excelize.Alignment{Vertical: "center"},
	})

	fontStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 4,
			},
		},
		Alignment: &excelize.Alignment{Horizontal: "right", WrapText: false, Vertical: "center"},
		Font:      &excelize.Font{Color: "#A5A5A5", Size: 8},
	})

	fontStyle2, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", WrapText: false, Vertical: "center"},
		Font:      &excelize.Font{Color: "#A5A5A5", Size: 8},
	})

	fontStyle3, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 14},
	})

	// 收據編號跟 第x聯
	f.SetCellStyle(sheetName, "F3", "F3", rightStyle)
	f.SetCellStyle(sheetName, "F21", "F21", rightStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "F17", "F17", rightStyle)
	f.SetCellStyle(sheetName, "F35", "F35", rightStyle) // 聯單複寫
	// (項目,金額)
	f.SetCellStyle(sheetName, "A5", "B5", fillStyle)
	f.SetCellStyle(sheetName, "A23", "B23", fillStyle) // 聯單複寫
	f.SetCellValue(sheetName, "C5", "金額")
	f.SetCellValue(sheetName, "C23", "金額") // 聯單複寫
	f.SetCellStyle(sheetName, "C5", "C5", fillStyle2)
	f.SetCellStyle(sheetName, "C23", "C23", fillStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D5", "E5", fillStyle)
	f.SetCellStyle(sheetName, "D23", "E23", fillStyle) // 聯單複寫
	f.SetCellValue(sheetName, "F5", "金額")
	f.SetCellValue(sheetName, "F23", "金額") // 聯單複寫
	f.SetCellStyle(sheetName, "F5", "F5", fillStyle2)
	f.SetCellStyle(sheetName, "F23", "F23", fillStyle2)

	// 項目內容跟金額內容
	f.SetCellStyle(sheetName, "A6", "A6", borderStyle)
	f.SetCellStyle(sheetName, "A24", "A24", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "B6", "B6", borderStyle1)
	f.SetCellStyle(sheetName, "B24", "B24", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "C6", "C6", borderStyle2)
	f.SetCellStyle(sheetName, "C24", "C24", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D6", "D6", borderStyle)
	f.SetCellStyle(sheetName, "D24", "D24", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "E6", "E6", borderStyle1)
	f.SetCellStyle(sheetName, "E24", "E24", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "F6", "F6", borderStyle2)
	f.SetCellStyle(sheetName, "F24", "F24", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "A7", "A7", borderStyle)
	f.SetCellStyle(sheetName, "A25", "A25", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "B7", "B7", borderStyle1)
	f.SetCellStyle(sheetName, "B25", "B25", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "C7", "C7", borderStyle2)
	f.SetCellStyle(sheetName, "C25", "C25", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D7", "D7", borderStyle)
	f.SetCellStyle(sheetName, "D25", "D25", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "E7", "E7", borderStyle1)
	f.SetCellStyle(sheetName, "E25", "E25", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "F7", "F7", borderStyle2)
	f.SetCellStyle(sheetName, "F25", "F25", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "A8", "A8", borderStyle)
	f.SetCellStyle(sheetName, "A26", "A26", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "B8", "B8", borderStyle1)
	f.SetCellStyle(sheetName, "B26", "B26", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "C8", "C8", borderStyle2)
	f.SetCellStyle(sheetName, "C26", "C26", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D8", "D8", borderStyle)
	f.SetCellStyle(sheetName, "D26", "D26", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "E8", "E8", borderStyle1)
	f.SetCellStyle(sheetName, "E26", "E26", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "F8", "F8", borderStyle2)
	f.SetCellStyle(sheetName, "F26", "F26", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "A9", "A9", borderStyle)
	f.SetCellStyle(sheetName, "A27", "A27", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "B9", "B9", borderStyle1)
	f.SetCellStyle(sheetName, "B27", "B27", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "C9", "C9", borderStyle2)
	f.SetCellStyle(sheetName, "C27", "C27", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D9", "D9", borderStyle)
	f.SetCellStyle(sheetName, "D27", "D27", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "E9", "E9", borderStyle1)
	f.SetCellStyle(sheetName, "E27", "E27", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "F9", "F9", borderStyle2)
	f.SetCellStyle(sheetName, "F27", "F27", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "A10", "A10", borderStyle)
	f.SetCellStyle(sheetName, "A28", "A28", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "B10", "B10", borderStyle1)
	f.SetCellStyle(sheetName, "B28", "B28", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "C10", "C10", borderStyle2)
	f.SetCellStyle(sheetName, "C28", "C28", borderStyle2) // 聯單複寫
	f.SetCellStyle(sheetName, "D10", "D10", borderStyle)
	f.SetCellStyle(sheetName, "D28", "D28", borderStyle) // 聯單複寫
	f.SetCellStyle(sheetName, "E10", "E10", borderStyle1)
	f.SetCellStyle(sheetName, "E28", "E28", borderStyle1) // 聯單複寫
	f.SetCellStyle(sheetName, "F10", "F10", borderStyle2)
	f.SetCellStyle(sheetName, "F28", "F28", borderStyle2) // 聯單複寫

	// 本期費用總計
	f.SetCellStyle(sheetName, "A11", "B11", fillStyle3)
	f.SetCellStyle(sheetName, "A29", "B29", fillStyle3) // 聯單複寫
	f.SetCellValue(sheetName, "A11", "本期費用總計")
	f.SetCellValue(sheetName, "A29", "本期費用總計") // 聯單複寫
	f.SetCellStyle(sheetName, "C11", "C11", fillStyle4)
	f.SetCellStyle(sheetName, "C29", "C29", fillStyle4) // 聯單複寫
	// 應繳（退）金額
	f.SetCellValue(sheetName, "D11", "應繳（退）金額")
	f.SetCellValue(sheetName, "D29", "應繳（退）金額") // 聯單複寫
	f.SetCellStyle(sheetName, "D11", "E11", fillStyle5)
	f.SetCellStyle(sheetName, "D29", "E29", fillStyle5) // 聯單複寫
	f.SetCellStyle(sheetName, "F11", "F11", fillStyle4)
	f.SetCellStyle(sheetName, "F29", "F29", fillStyle4) // 聯單複寫

	// 程式改版日期：2022
	f.SetCellValue(sheetName, "F18", "程式改版日期：2022-09-15 | Jubo智慧照護平台")
	f.SetCellValue(sheetName, "F36", "程式改版日期：2022-09-15 | Jubo智慧照護平台") // 聯單複寫
	f.SetCellStyle(sheetName, "A18", "F18", borderStyle3)              //  這個不用複寫
	f.SetCellStyle(sheetName, "F18", "F18", fontStyle)
	f.SetCellStyle(sheetName, "F36", "F36", fontStyle2) // 聯單複寫

	// 機構名稱 style
	f.SetCellStyle(sheetName, "B1", "B1", fontStyle3)
	f.SetCellStyle(sheetName, "B19", "B19", fontStyle3) // 聯單複寫

	//	寫機構(住民)檔案
	f.SetCellValue(sheetName, "B1", payRecord.Organization.Name)
	f.SetCellValue(sheetName, "B19", payRecord.Organization.Name) // 聯單複寫
	f.SetCellValue(sheetName, "B2", payRecordYear+"年"+payRecordMonth+"月"+organizationReceiptTemplateSetting.TitleName+invalidText)
	f.SetCellValue(sheetName, "B20", payRecordYear+"年"+payRecordMonth+"月"+organizationReceiptTemplateSetting.TitleName+invalidText) // 聯單複寫
	f.SetCellValue(sheetName, "A3", "姓名："+payRecord.Patient.LastName+payRecord.Patient.FirstName)
	f.SetCellValue(sheetName, "A21", "姓名："+payRecord.Patient.LastName+payRecord.Patient.FirstName) // 聯單複寫
	f.SetCellValue(sheetName, "F3", "收據編號："+payRecord.ReceiptNumber)
	f.SetCellValue(sheetName, "F21", "收據編號："+payRecord.ReceiptNumber) // 聯單複寫
	var patientTextArray []string
	var patientText string
	if organizationReceiptTemplateSettingPatientInfo.ShowAreaAndClass {
		if payRecord.Patient.Branch != "" {
			patientTextArray = append(patientTextArray, "區域/組別："+payRecord.Patient.Branch+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowBedAndRoom {
		if payRecord.Patient.Room != "" || payRecord.Patient.Bed != "" {
			patientTextArray = append(patientTextArray, "房號/床號："+payRecord.Patient.Room+payRecord.Patient.Bed+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowSex {
		if payRecord.Patient.Sex != "" {
			var sex string
			if payRecord.Patient.Sex == "male" {
				sex = "男"
			} else if payRecord.Patient.Sex == "female" {
				sex = "女"
			}
			patientTextArray = append(patientTextArray, "性別："+sex+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowBirthday {
		if payRecord.Patient.Birthday.Format("2006-01-02") != "0001-01-01" {
			patientTextArray = append(patientTextArray, "生日："+payRecord.Patient.Birthday.Format("2006-01-02")+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowAge {
		if payRecord.Patient.Birthday.Format("2006-01-02") != "0001-01-01" {
			nowTime := time.Now()
			age := nowTime.Year() - payRecord.Patient.Birthday.Year()
			if payRecord.Patient.Birthday.Month() >= nowTime.Month() && payRecord.Patient.Birthday.Day() >= nowTime.Day() {
				age += 1
			}

			patientTextArray = append(patientTextArray, "年齡："+strconv.Itoa(age)+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowCheckInDate {
		if payRecord.Patient.CheckInDate.Format("2006-01-02") != "0001-01-01" {
			patientTextArray = append(patientTextArray, "入住日期："+payRecord.Patient.CheckInDate.Format("2006-01-02")+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowPatientNumber {
		if payRecord.Patient.PatientNumber != "" {
			patientTextArray = append(patientTextArray, "住民編號："+payRecord.Patient.PatientNumber+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowRecordNumber {
		if payRecord.Patient.RecordNumber != "" {
			patientTextArray = append(patientTextArray, "病歷號："+payRecord.Patient.RecordNumber+"    ")
		}
	}
	if organizationReceiptTemplateSettingPatientInfo.ShowIdNumber {
		if payRecord.Patient.IdNumber != "" {
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
			patientTextArray = append(patientTextArray, "身分證字號："+idNumber+"    ")
		}
	}

	for i := range patientTextArray {
		if i == 4 {
			f.SetRowHeight(sheetName, 4, 35)
			f.SetRowHeight(sheetName, 22, 35)
			patientText += "\n"
		} else if i == 8 {
			f.SetRowHeight(sheetName, 4, 50)
			f.SetRowHeight(sheetName, 22, 50)
			patientText += "\n"
		}
		patientText += patientTextArray[i]
	}
	f.SetCellValue(sheetName, "A4", patientText)
	f.SetCellValue(sheetName, "A22", patientText) // 聯單複寫
}

// 塞一些機構資訊
func SetExcelOrganizationData(excelOrganizationDataStruct ExcelOrganizationDataStruct) {
	organizationReceiptTemplateSettingOrganizationInfo := excelOrganizationDataStruct.OrganizationReceiptTemplateSettingOrganizationInfo
	payRecord := excelOrganizationDataStruct.PayRecord
	f := excelOrganizationDataStruct.f
	sheetName := excelOrganizationDataStruct.SheetName
	insertRowCount := excelOrganizationDataStruct.InsertRowCount
	organizationReceiptTemplateSetting := excelOrganizationDataStruct.OrganizationReceiptTemplateSetting
	// 塞一些機構資訊
	var organizationText string
	if organizationReceiptTemplateSettingOrganizationInfo.ShowTaxIdNumber {
		var taxIdNumber string
		if payRecord.Organization.TaxIdNumber != nil {
			taxIdNumber = *payRecord.Organization.TaxIdNumber
			organizationText += "機構統一編號：" + taxIdNumber + "    "
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowPhone {
		var phone string
		if payRecord.Organization.Phone != nil {
			phone = *payRecord.Organization.Phone
			organizationText += "機構電話：" + phone + "    "
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowFax {
		var fax string
		if payRecord.Organization.Fax != nil {
			fax = *payRecord.Organization.Fax
			organizationText += "機構傳真：" + fax + "    "
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowOwner {
		var owner string
		if payRecord.Organization.Owner != nil {
			owner = *payRecord.Organization.Owner
			organizationText += "負責人：" + owner + "    "
		}
	}
	f.SetCellValue(sheetName, "A"+strconv.Itoa(12+insertRowCount), organizationText)
	f.SetCellValue(sheetName, "A"+strconv.Itoa(30+insertRowCount*2), organizationText)

	// 以下的要塞到A13
	var organizationTextArray []string
	var organizationTextCount int
	if organizationReceiptTemplateSettingOrganizationInfo.ShowEstablishmentNumber {
		var establishmentNumber string
		if payRecord.Organization.EstablishmentNumber != nil {
			establishmentNumber = *payRecord.Organization.EstablishmentNumber
			if establishmentNumber != "" {
				organizationTextArray = append(organizationTextArray, "設立許可文號："+establishmentNumber+"    ")
				organizationTextCount++
			}
		}
	}
	// 機構地址
	if organizationReceiptTemplateSettingOrganizationInfo.ShowAddress {
		var address string
		if payRecord.Organization.AddressCity != nil {
			address = *payRecord.Organization.AddressCity
		}
		if payRecord.Organization.AddressDistrict != nil {
			address += *payRecord.Organization.AddressDistrict
		}
		if payRecord.Organization.Address != nil {
			address += *payRecord.Organization.Address
		}
		if address != "" {
			organizationTextArray = append(organizationTextArray, "機構地址："+address+"    ")
			organizationTextCount++
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowEmail {
		var email string
		if payRecord.Organization.Email != nil {
			email = *payRecord.Organization.Email
			if email != "" {
				organizationTextArray = append(organizationTextArray, "機構信箱："+email+"    ")
				organizationTextCount++
			}
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowRemittanceBank {
		var remittanceBank string
		if payRecord.Organization.RemittanceBank != nil {
			remittanceBank = *payRecord.Organization.RemittanceBank
			if remittanceBank != "" {
				organizationTextArray = append(organizationTextArray, "匯款銀行："+remittanceBank+"    ")
				organizationTextCount++
			}
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowRemittanceIdNumber {
		var remittanceIdNumber string
		if payRecord.Organization.RemittanceIdNumber != nil {
			remittanceIdNumber = *payRecord.Organization.RemittanceIdNumber
			if remittanceIdNumber != "" {
				organizationTextArray = append(organizationTextArray, "匯款帳號："+remittanceIdNumber+"    ")
				organizationTextCount++
			}
		}
	}
	if organizationReceiptTemplateSettingOrganizationInfo.ShowRemittanceUserName {
		var remittanceUserName string
		if payRecord.Organization.RemittanceUserName != nil {
			remittanceUserName = *payRecord.Organization.RemittanceUserName
			if remittanceUserName != "" {
				organizationTextArray = append(organizationTextArray, "匯款戶名："+remittanceUserName+"    ")
				organizationTextCount++
			}
		}
	}
	organizationText = ""
	for i := range organizationTextArray {
		if i == 0 {
			organizationText += organizationTextArray[i]
		} else {
			organizationText += "\n" + organizationTextArray[i]
		}
	}
	f.SetCellValue(sheetName, "A"+strconv.Itoa(13+insertRowCount), organizationText)
	f.SetCellValue(sheetName, "A"+strconv.Itoa(31+insertRowCount*2), organizationText) // 聯單複寫
	// 代表不只一筆 要調整高度
	if organizationTextCount > 4 {
		f.SetRowHeight(sheetName, 13+insertRowCount, float64(organizationTextCount*17))
		f.SetRowHeight(sheetName, 31+insertRowCount*2, float64(organizationTextCount*17)) // 聯單複寫
	} else {
		f.SetRowHeight(sheetName, 13+insertRowCount, 68)
		f.SetRowHeight(sheetName, 31+insertRowCount*2, 68) // 聯單複寫
	}

	// 備註
	var noteText string
	if organizationReceiptTemplateSetting.NoteText != "" {
		noteText += organizationReceiptTemplateSetting.NoteText
	}
	if payRecord.Note != "" {
		noteText += "\n" + payRecord.Note
	}
	if excelOrganizationDataStruct.IsInvalid {
		noteText += "\n" + payRecord.InvalidDate.Format("2006-01-02 15:04") + "作廢，作廢說明：" + payRecord.InvalidCaption
	}
	f.SetCellValue(sheetName, "A"+strconv.Itoa(14+insertRowCount), noteText)
	f.SetCellValue(sheetName, "A"+strconv.Itoa(32+insertRowCount*2), noteText) // 聯單複寫
	f.SetCellValue(sheetName, "F"+strconv.Itoa(17+insertRowCount), organizationReceiptTemplateSetting.PartOneName)
	f.SetCellValue(sheetName, "F"+strconv.Itoa(35+insertRowCount*2), organizationReceiptTemplateSetting.PartTwoName)
}

func DownloadExcelOrganizationSeal(downloadExcelOrganizationSealStruct DownloadExcelOrganizationSealStruct) error {
	r := downloadExcelOrganizationSealStruct.r
	ctx := downloadExcelOrganizationSealStruct.ctx
	organizationReceiptTemplateSetting := downloadExcelOrganizationSealStruct.OrganizationReceiptTemplateSetting

	if organizationReceiptTemplateSetting.OrganizationPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.OrganizationPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		urlAndBucketName := strings.Split(organizationReceiptTemplateSetting.OrganizationPicture, "inventory-tool")
		getBucketName := strings.Split(urlAndBucketName[0], "/")
		err := r.store.DownloadPictureFile(ctx, getBucketName[3], fileName[1], "inventory-tool/"+fileName[0])
		if err != nil {
			return err
		}
	}

	if organizationReceiptTemplateSetting.SealOnePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealOnePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		urlAndBucketName := strings.Split(organizationReceiptTemplateSetting.SealOnePicture, "inventory-tool")
		getBucketName := strings.Split(urlAndBucketName[0], "/")
		err := r.store.DownloadPictureFile(ctx, getBucketName[3], fileName[1], "inventory-tool/"+fileName[0])
		if err != nil {
			return err
		}
	}

	if organizationReceiptTemplateSetting.SealTwoPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealTwoPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		urlAndBucketName := strings.Split(organizationReceiptTemplateSetting.SealTwoPicture, "inventory-tool")
		getBucketName := strings.Split(urlAndBucketName[0], "/")
		err := r.store.DownloadPictureFile(ctx, getBucketName[3], fileName[1], "inventory-tool/"+fileName[0])
		if err != nil {
			return err
		}
	}
	if organizationReceiptTemplateSetting.SealThreePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealThreePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		urlAndBucketName := strings.Split(organizationReceiptTemplateSetting.SealThreePicture, "inventory-tool")
		getBucketName := strings.Split(urlAndBucketName[0], "/")
		err := r.store.DownloadPictureFile(ctx, getBucketName[3], fileName[1], "inventory-tool/"+fileName[0])
		if err != nil {
			return err
		}
	}
	if organizationReceiptTemplateSetting.SealFourPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealFourPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		urlAndBucketName := strings.Split(organizationReceiptTemplateSetting.SealFourPicture, "inventory-tool")
		getBucketName := strings.Split(urlAndBucketName[0], "/")
		err := r.store.DownloadPictureFile(ctx, getBucketName[3], fileName[1], "inventory-tool/"+fileName[0])
		if err != nil {
			return err
		}
	}

	return nil
}

// 塞一些印章資訊
func SetExcelOrganizationSealData(organizationSealDataStruct OrganizationSealDataStruct) error {
	f := organizationSealDataStruct.f
	sheetName := organizationSealDataStruct.SheetName
	insertRowCount := organizationSealDataStruct.InsertRowCount
	organizationReceiptTemplateSetting := organizationSealDataStruct.OrganizationReceiptTemplateSetting

	// 印章
	f.SetCellValue(sheetName, "A"+strconv.Itoa(15+insertRowCount), organizationReceiptTemplateSetting.SealOneName)
	f.SetCellValue(sheetName, "A"+strconv.Itoa(33+insertRowCount*2), organizationReceiptTemplateSetting.SealOneName) // 聯單複寫
	f.SetCellValue(sheetName, "B"+strconv.Itoa(15+insertRowCount), organizationReceiptTemplateSetting.SealTwoName)
	f.SetCellValue(sheetName, "B"+strconv.Itoa(33+insertRowCount*2), organizationReceiptTemplateSetting.SealTwoName) // 聯單複寫
	f.SetCellValue(sheetName, "C"+strconv.Itoa(15+insertRowCount), organizationReceiptTemplateSetting.SealThreeName)
	f.SetCellValue(sheetName, "C"+strconv.Itoa(33+insertRowCount*2), organizationReceiptTemplateSetting.SealThreeName) // 聯單複寫
	f.SetCellValue(sheetName, "D"+strconv.Itoa(15+insertRowCount), organizationReceiptTemplateSetting.SealFourName)
	f.SetCellValue(sheetName, "D"+strconv.Itoa(33+insertRowCount*2), organizationReceiptTemplateSetting.SealFourName) // 聯單複寫

	if organizationReceiptTemplateSetting.OrganizationPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.OrganizationPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		organizationPicture, _ := os.Open(fileName[1] + ".jpg")
		organizationPictureImage, _, err := image.Decode(organizationPicture)
		if err != nil {
			return err
		}

		organizationPictureImageBounds := organizationPictureImage.Bounds()
		organizationWidth := organizationPictureImageBounds.Max.X
		organizationHeight := organizationPictureImageBounds.Max.Y
		var organizationPictureWidth float64
		var organizationPictureWidthStr string

		widthPercent := float64(organizationWidth) / 96 / 89
		heightPercent := float64(organizationHeight) / 96 / 48
		if widthPercent > heightPercent {
			// 這邊是在做單位的換算 還有找出適合的大小(比例)
			organizationPictureWidth = 110 / (float64(organizationWidth) / 96) * 0.01
			organizationPictureWidthStr = fmt.Sprintf("%f", organizationPictureWidth)
		} else {
			organizationPictureWidth = 63 / (float64(organizationHeight) / 96) * 0.01
			organizationPictureWidthStr = fmt.Sprintf("%f", organizationPictureWidth)
		}

		if err := f.AddPicture(sheetName, "A"+strconv.Itoa(1), fileName[1]+".jpg", `{
		"x_offset": 3,
		"y_offset": 3,
		"x_scale": `+organizationPictureWidthStr+`,
		"y_scale": `+organizationPictureWidthStr+`,
		"positioning":"oneCell"
}`); err != nil {
			return err
		}

		if err := f.AddPicture(sheetName, "A"+strconv.Itoa(19+insertRowCount), fileName[1]+".jpg", `{
		"x_offset": 3,
		"y_offset": 3,
		"x_scale": `+organizationPictureWidthStr+`,
		"y_scale": `+organizationPictureWidthStr+`,
		"positioning":"oneCell"
		}`); err != nil {
			return err
		}
		err = organizationPicture.Close()
		if err != nil {
			return err
		}
	}

	if organizationReceiptTemplateSetting.SealOnePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealOnePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")

		sealOnePictureFile, _ := os.Open(fileName[1] + ".jpg")
		sealOnePictureImage, _, err := image.Decode(sealOnePictureFile)
		if err != nil {
			return err
		}

		sealOnePictureImageBounds := sealOnePictureImage.Bounds()
		sealOnePictureWidth := sealOnePictureImageBounds.Max.X
		// 這邊是在做單位的換算 還有找出適合的大小(比例)
		sealOneWidth := 79 / (float64(sealOnePictureWidth) / 96) * 0.01
		sealOneWidthStr := fmt.Sprintf("%f", sealOneWidth)

		if err := f.AddPicture(sheetName, "A"+strconv.Itoa(16+insertRowCount), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealOneWidthStr+`,
		"y_scale": `+sealOneWidthStr+`,
		"positioning":"oneCell"
}`); err != nil {
			return err
		}

		if err := f.AddPicture(sheetName, "A"+strconv.Itoa(34+insertRowCount*2), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealOneWidthStr+`,
		"y_scale": `+sealOneWidthStr+`,
		"positioning":"oneCell"
		}`); err != nil {
			return err
		}
		err = sealOnePictureFile.Close()
		if err != nil {
			return err
		}
	}
	if organizationReceiptTemplateSetting.SealTwoPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealTwoPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		sealTwoPictureFile, _ := os.Open(fileName[1] + ".jpg")
		sealTwoPictureImage, _, err := image.Decode(sealTwoPictureFile)
		if err != nil {
			return err
		}

		sealTwoPictureImageBounds := sealTwoPictureImage.Bounds()
		sealTwoPictureWidth := sealTwoPictureImageBounds.Max.X
		// 這邊是在做單位的換算 還有找出適合的大小(比例)
		sealTwoWidth := 79 / (float64(sealTwoPictureWidth) / 96) * 0.01
		sealTwoWidthStr := fmt.Sprintf("%f", sealTwoWidth)
		if err := f.AddPicture(sheetName, "B"+strconv.Itoa(16+insertRowCount), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealTwoWidthStr+`,
		"y_scale": `+sealTwoWidthStr+`,
		"positioning":"oneCell"
}`); err != nil {
			return err
		}

		if err := f.AddPicture(sheetName, "B"+strconv.Itoa(34+insertRowCount*2), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealTwoWidthStr+`,
		"y_scale": `+sealTwoWidthStr+`,
		"positioning":"oneCell"
		}`); err != nil {
			return err
		}
		err = sealTwoPictureFile.Close()
		if err != nil {
			return err
		}
	}
	if organizationReceiptTemplateSetting.SealThreePicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealThreePicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		sealThreePictureFile, _ := os.Open(fileName[1] + ".jpg")
		sealThreePictureImage, _, err := image.Decode(sealThreePictureFile)
		if err != nil {
			return err
		}

		sealThreePictureImageBounds := sealThreePictureImage.Bounds()
		sealThreePictureWidth := sealThreePictureImageBounds.Max.X
		// 這邊是在做單位的換算 還有找出適合的大小(比例)
		sealThreeWidth := 79 / (float64(sealThreePictureWidth) / 96) * 0.01
		sealThreeWidthStr := fmt.Sprintf("%f", sealThreeWidth)
		if err := f.AddPicture(sheetName, "C"+strconv.Itoa(16+insertRowCount), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealThreeWidthStr+`,
		"y_scale": `+sealThreeWidthStr+`,
		"positioning":"oneCell"
}`); err != nil {
			return err
		}

		if err := f.AddPicture(sheetName, "C"+strconv.Itoa(34+insertRowCount*2), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealThreeWidthStr+`,
		"y_scale": `+sealThreeWidthStr+`,
		"positioning":"oneCell"
		}`); err != nil {
			return err
		}
		err = sealThreePictureFile.Close()
		if err != nil {
			return err
		}
	}
	if organizationReceiptTemplateSetting.SealFourPicture != "" {
		orgAndFileId := strings.Split(organizationReceiptTemplateSetting.SealFourPicture, "inventory-tool/")
		fileName := strings.Split(orgAndFileId[1], "/")
		sealFourPictureFile, _ := os.Open(fileName[1] + ".jpg")
		sealFourPictureImage, _, err := image.Decode(sealFourPictureFile)
		if err != nil {
			return err
		}

		sealFourPictureImageBounds := sealFourPictureImage.Bounds()
		sealFourPictureWidth := sealFourPictureImageBounds.Max.X
		// 這邊是在做單位的換算 還有找出適合的大小(比例)
		sealFourWidth := 79 / (float64(sealFourPictureWidth) / 96) * 0.01
		sealFourWidthStr := fmt.Sprintf("%f", sealFourWidth)
		if err := f.AddPicture(sheetName, "D"+strconv.Itoa(16+insertRowCount), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealFourWidthStr+`,
		"y_scale": `+sealFourWidthStr+`,
		"positioning":"oneCell"
}`); err != nil {
			return err
		}

		if err := f.AddPicture(sheetName, "D"+strconv.Itoa(34+insertRowCount*2), fileName[1]+".jpg", `{
		"x_offset": 3,
		"x_scale": `+sealFourWidthStr+`,
		"y_scale": `+sealFourWidthStr+`,
		"positioning":"oneCell"
		}`); err != nil {
			return err
		}
		err = sealFourPictureFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// 機構設定的金額顯示方式-類別加總
func printClassAddUp(payRecord *models.PayRecord, f *excelize.File, classStyle, priceStyle int, sheetName string) int {
	f.SetCellValue(sheetName, "A5", "類別")
	f.SetCellValue(sheetName, "A23", "類別")
	f.SetCellValue(sheetName, "D5", "類別")
	f.SetCellValue(sheetName, "D23", "類別")
	f.MergeCell(sheetName, "A6", "B6")
	f.MergeCell(sheetName, "A24", "B24") // 聯單複寫
	f.MergeCell(sheetName, "A7", "B7")
	f.MergeCell(sheetName, "A25", "B25") // 聯單複寫
	f.MergeCell(sheetName, "A8", "B8")
	f.MergeCell(sheetName, "A26", "B26") // 聯單複寫
	f.MergeCell(sheetName, "A9", "B9")
	f.MergeCell(sheetName, "A27", "B27") // 聯單複寫
	f.MergeCell(sheetName, "A10", "B10")
	f.MergeCell(sheetName, "A28", "B28") // 聯單複寫
	f.MergeCell(sheetName, "D6", "E6")
	f.MergeCell(sheetName, "D24", "E24") // 聯單複寫
	f.MergeCell(sheetName, "D7", "E7")
	f.MergeCell(sheetName, "D25", "E25") // 聯單複寫
	f.MergeCell(sheetName, "D8", "E8")
	f.MergeCell(sheetName, "D26", "E26") // 聯單複寫
	f.MergeCell(sheetName, "D9", "E9")
	f.MergeCell(sheetName, "D27", "E27") // 聯單複寫
	f.MergeCell(sheetName, "D10", "E10")
	f.MergeCell(sheetName, "D28", "E28") // 聯單複寫
	var subtotal int
	// 基本月費
	payRecordBasicCharges := []PayRecordBasicCharge{}
	json.Unmarshal(payRecord.BasicCharge, &payRecordBasicCharges)
	var dataCount int

	if len(payRecordBasicCharges) > 0 {
		// 寫excel副標題
		f.SetCellValue(sheetName, "A6", "基本月費")
		f.SetCellValue(sheetName, "A24", "基本月費")
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		// 寫入excel內容
		for i := range payRecordBasicCharges {
			if payRecordBasicCharges[i].Type == "charge" {
				price += payRecordBasicCharges[i].Price
			} else {
				price += -payRecordBasicCharges[i].Price
			}
		}
		subtotal += price
		f.SetCellValue(sheetName, "C"+cellValueStr, price)
		f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
		dataCount++
	}

	// 補助款
	payRecordSubsidies := []PayRecordSubsidy{}
	json.Unmarshal(payRecord.Subsidy, &payRecordSubsidies)
	if len(payRecordSubsidies) > 0 {
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		// 寫excel副標題
		f.SetCellValue(sheetName, "A"+cellValueStr, "補助款")
		f.SetCellValue(sheetName, "A"+partCellValueStr, "補助款")

		// 寫入excel內容
		for i := range payRecordSubsidies {
			if payRecordSubsidies[i].Type == "charge" {
				price += payRecordSubsidies[i].Price
			} else {
				price += -payRecordSubsidies[i].Price
			}
		}
		subtotal += price
		f.SetCellValue(sheetName, "C"+cellValueStr, price)
		f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
		dataCount++
	}

	// 異動(請假)
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	json.Unmarshal(payRecord.TransferRefundLeave, &payRecordTransferRefundLeaves)
	if len(payRecordTransferRefundLeaves) > 0 {
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		// 寫excel副標題
		f.SetCellValue(sheetName, "A"+cellValueStr, "請假退費")
		f.SetCellValue(sheetName, "A"+partCellValueStr, "請假退費")

		// 寫入excel內容
		for i := range payRecordTransferRefundLeaves {
			if payRecordTransferRefundLeaves[i].Type == "charge" {
				price += payRecordTransferRefundLeaves[i].Price
			} else {
				price += -payRecordTransferRefundLeaves[i].Price
			}
		}
		subtotal += price
		f.SetCellValue(sheetName, "C"+cellValueStr, price)
		f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
		dataCount++
	}

	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
	json.Unmarshal(payRecord.NonFixedCharge, &payRecordNonFixedChargeRecords)
	// 排序要一樣(否則可能每次印出來的不太一樣)
	sort.Slice(payRecordNonFixedChargeRecords, func(j, k int) bool {
		return payRecordNonFixedChargeRecords[j].ItemCategory < payRecordNonFixedChargeRecords[k].ItemCategory
	})
	// 動態新增的行數
	var insertRowCount int
	if len(payRecordNonFixedChargeRecords) > 0 {
		payRecordNonFixedChargeRecordElements := make(map[string]*PayRecordNonFixedChargeRecordDataForPartPrint)
		for _, d := range payRecordNonFixedChargeRecords {
			if payRecordNonFixedChargeRecordElements[d.ItemCategory] == nil {
				var price int
				if d.Type == "charge" {
					price = d.Subtotal
				} else {
					price = -d.Subtotal
				}
				subtotal += price
				payRecordNonFixedChargeRecordElements[d.ItemCategory] = &PayRecordNonFixedChargeRecordDataForPartPrint{
					ItemCategory: d.ItemCategory,
					Subtotal:     price,
				}
			} else {
				if d.Type == "charge" {
					subtotal += d.Subtotal
					payRecordNonFixedChargeRecordElements[d.ItemCategory].Subtotal += d.Subtotal
				} else {
					subtotal -= d.Subtotal
					payRecordNonFixedChargeRecordElements[d.ItemCategory].Subtotal -= d.Subtotal
				}
			}
		}

		// 確保key的順序一致(不做這段 用payRecordNonFixedChargeRecordElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(payRecordNonFixedChargeRecordElements))
		for key := range payRecordNonFixedChargeRecordElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(j, k int) bool {
			return payRecordNonFixedChargeRecordElements[keys[j]].ItemCategory < payRecordNonFixedChargeRecordElements[keys[k]].ItemCategory
		})
		for i := range keys {
			cellValue := 6 + dataCount
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   classStyle,
				DateStyle:      classStyle,
				PriceStyle:     priceStyle,
				Price:          payRecordNonFixedChargeRecordElements[keys[i]].Subtotal,
				SheetName:      sheetName,
				ItemName:       payRecordNonFixedChargeRecordElements[keys[i]].ItemCategory,
				DateStr:        nil,
			}

			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 如果有超過10 表示有新增欄位 要重算總和
	if dataCount >= 11 {
		subTotalcellValueStr := strconv.Itoa(dataCount)
		partSubTotalcellValueStr := strconv.Itoa(29 + insertRowCount*2)
		f.SetCellFormula(sheetName, "C"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
	} else {
		f.SetCellFormula(sheetName, "C11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C29", strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F29", strconv.Itoa(subtotal)) // 聯單複寫
	}
	return insertRowCount
}

// 機構設定的金額顯示方式-類別加總(分稅別)
func printClassAddUpByTaxType(payRecord *models.PayRecord, f *excelize.File, classStyle, priceFormatAndRightBottomBorderAndFontStyle, taxCount int, taxType, sheetName string) (bool, int) {
	var taxHaveData bool
	f.SetCellValue(sheetName, "A5", "類別")
	f.SetCellValue(sheetName, "A23", "類別")
	f.SetCellValue(sheetName, "D5", "類別")
	f.SetCellValue(sheetName, "D23", "類別")
	f.MergeCell(sheetName, "A6", "B6")
	f.MergeCell(sheetName, "A24", "B24") // 聯單複寫
	f.MergeCell(sheetName, "A7", "B7")
	f.MergeCell(sheetName, "A25", "B25") // 聯單複寫
	f.MergeCell(sheetName, "A8", "B8")
	f.MergeCell(sheetName, "A26", "B26") // 聯單複寫
	f.MergeCell(sheetName, "A9", "B9")
	f.MergeCell(sheetName, "A27", "B27") // 聯單複寫
	f.MergeCell(sheetName, "A10", "B10")
	f.MergeCell(sheetName, "A28", "B28") // 聯單複寫
	f.MergeCell(sheetName, "D6", "E6")
	f.MergeCell(sheetName, "D24", "E24") // 聯單複寫
	f.MergeCell(sheetName, "D7", "E7")
	f.MergeCell(sheetName, "D25", "E25") // 聯單複寫
	f.MergeCell(sheetName, "D8", "E8")
	f.MergeCell(sheetName, "D26", "E26") // 聯單複寫
	f.MergeCell(sheetName, "D9", "E9")
	f.MergeCell(sheetName, "D27", "E27") // 聯單複寫
	f.MergeCell(sheetName, "D10", "E10")
	f.MergeCell(sheetName, "D28", "E28") // 聯單複寫
	var subtotal int
	// 基本月費
	payRecordBasicCharges := []PayRecordBasicCharge{}
	json.Unmarshal(payRecord.BasicCharge, &payRecordBasicCharges)
	var dataCount int
	var haveBasicCharge bool
	if len(payRecordBasicCharges) > 0 {
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		// 寫入excel內容
		for i := range payRecordBasicCharges {
			// 稅別要一樣才新增
			if taxType == payRecordBasicCharges[i].TaxType {
				haveBasicCharge = true
				taxHaveData = true
				if payRecordBasicCharges[i].Type == "charge" {
					price += payRecordBasicCharges[i].Price
				} else {
					price += -payRecordBasicCharges[i].Price
				}
			}
		}
		if haveBasicCharge {
			// 寫excel副標題
			f.SetCellValue(sheetName, "A6", "基本月費")
			f.SetCellValue(sheetName, "A24", "基本月費")
			subtotal += price
			f.SetCellValue(sheetName, "C"+cellValueStr, price)
			f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
			dataCount++
		}
	}

	// 補助款
	payRecordSubsidies := []PayRecordSubsidy{}
	json.Unmarshal(payRecord.Subsidy, &payRecordSubsidies)
	var haveSubsidy bool
	if len(payRecordSubsidies) > 0 {
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		// 寫入excel內容
		for i := range payRecordSubsidies {
			if taxType == "stampTax" {
				haveSubsidy = true
				taxHaveData = true
				if payRecordSubsidies[i].Type == "charge" {
					price += payRecordSubsidies[i].Price
				} else {
					price += -payRecordSubsidies[i].Price
				}
			}
		}
		if haveSubsidy {
			// 寫excel副標題
			f.SetCellValue(sheetName, "A"+cellValueStr, "補助款")
			f.SetCellValue(sheetName, "A"+partCellValueStr, "補助款")
			subtotal += price
			f.SetCellValue(sheetName, "C"+cellValueStr, price)
			f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
			dataCount++
		}
	}

	// 異動(請假)
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	json.Unmarshal(payRecord.TransferRefundLeave, &payRecordTransferRefundLeaves)
	var haveTransferRefundLeave bool
	if len(payRecordTransferRefundLeaves) > 0 {
		var price int
		cellValue := 6 + dataCount
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)

		// 寫入excel內容
		for i := range payRecordTransferRefundLeaves {
			if taxType == "stampTax" {
				taxHaveData = true
				haveTransferRefundLeave = true
				if payRecordTransferRefundLeaves[i].Type == "charge" {
					price += payRecordTransferRefundLeaves[i].Price
				} else {
					price += -payRecordTransferRefundLeaves[i].Price
				}
			}
		}
		if haveTransferRefundLeave {
			// 寫excel副標題
			f.SetCellValue(sheetName, "A"+cellValueStr, "請假退費")
			f.SetCellValue(sheetName, "A"+partCellValueStr, "請假退費")
			subtotal += price
			f.SetCellValue(sheetName, "C"+cellValueStr, price)
			f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
			dataCount++
		}
	}

	// 非固定
	var insertRowCount int
	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
	json.Unmarshal(payRecord.NonFixedCharge, &payRecordNonFixedChargeRecords)
	// 排序要一樣(否則可能每次印出來的不太一樣)
	sort.Slice(payRecordNonFixedChargeRecords, func(j, k int) bool {
		return payRecordNonFixedChargeRecords[j].ItemCategory < payRecordNonFixedChargeRecords[k].ItemCategory
	})
	// 動態新增的行數
	if len(payRecordNonFixedChargeRecords) > 0 {
		payRecordNonFixedChargeRecordElements := make(map[string]*PayRecordNonFixedChargeRecordDataForPartPrint)
		for _, d := range payRecordNonFixedChargeRecords {
			if d.TaxType == taxType {
				if payRecordNonFixedChargeRecordElements[d.ItemCategory] == nil && d.TaxType == taxType {
					taxHaveData = true
					var price int
					if d.Type == "charge" {
						price = d.Subtotal
					} else {
						price = -d.Subtotal
					}
					subtotal += price
					payRecordNonFixedChargeRecordElements[d.ItemCategory] = &PayRecordNonFixedChargeRecordDataForPartPrint{
						ItemCategory: d.ItemCategory,
						Subtotal:     price,
					}
				} else {
					if d.Type == "charge" {
						subtotal += d.Subtotal
						payRecordNonFixedChargeRecordElements[d.ItemCategory].Subtotal += d.Subtotal
					} else {
						subtotal -= d.Subtotal
						payRecordNonFixedChargeRecordElements[d.ItemCategory].Subtotal -= d.Subtotal
					}
				}
			}
		}

		// 確保key的順序一致(不做這段 用payRecordNonFixedChargeRecordElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(payRecordNonFixedChargeRecordElements))
		for key := range payRecordNonFixedChargeRecordElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(j, k int) bool {
			return payRecordNonFixedChargeRecordElements[keys[j]].ItemCategory < payRecordNonFixedChargeRecordElements[keys[k]].ItemCategory
		})
		for i := range keys {
			cellValue := 6 + dataCount
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   classStyle,
				DateStyle:      classStyle,
				PriceStyle:     priceFormatAndRightBottomBorderAndFontStyle,
				Price:          payRecordNonFixedChargeRecordElements[keys[i]].Subtotal,
				SheetName:      sheetName,
				ItemName:       payRecordNonFixedChargeRecordElements[keys[i]].ItemCategory,
				DateStr:        nil,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 如果有超過10 表示有新增欄位 要重算總和
	if dataCount >= 11 {
		subTotalcellValueStr := strconv.Itoa(dataCount)
		partSubTotalcellValueStr := strconv.Itoa(29 + insertRowCount*2)
		f.SetCellFormula(sheetName, "C"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
	} else {
		f.SetCellFormula(sheetName, "C11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C29", strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F29", strconv.Itoa(subtotal)) // 聯單複寫
	}

	return taxHaveData, insertRowCount
}

// 機構設定的金額顯示方式-項目
func printItem(payRecord *models.PayRecord, f *excelize.File, contentStyle, dateStyle, priceStyle int, sheetName string) int {
	f.SetCellValue(sheetName, "A5", "項目")
	f.SetCellValue(sheetName, "A23", "項目")
	f.SetCellValue(sheetName, "D5", "項目")
	f.SetCellValue(sheetName, "D23", "項目")
	loc, _ := time.LoadLocation("Asia/Taipei")

	var subtotal int
	// 基本月費
	payRecordBasicCharges := []PayRecordBasicCharge{}
	json.Unmarshal(payRecord.BasicCharge, &payRecordBasicCharges)
	// 資料總筆數
	var dataCount int
	// 動態新增的行數
	var insertRowCount int
	if len(payRecordBasicCharges) > 0 {
		// 寫入excel內容
		for i := range payRecordBasicCharges {
			cellValue := 6 + dataCount
			var price int
			if payRecordBasicCharges[i].Type == "charge" {
				price += payRecordBasicCharges[i].Price
			} else {
				price += -payRecordBasicCharges[i].Price
			}
			subtotal += price
			dateStr := payRecordBasicCharges[i].StartDate.In(loc).Format("01-02") + "到" + payRecordBasicCharges[i].EndDate.In(loc).Format("01-02")
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   contentStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          price,
				SheetName:      sheetName,
				ItemName:       payRecordBasicCharges[i].ItemName,
				DateStr:        &dateStr,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 補助款
	payRecordSubsidies := []PayRecordSubsidy{}
	json.Unmarshal(payRecord.Subsidy, &payRecordSubsidies)
	if len(payRecordSubsidies) > 0 {
		// 寫入excel內容
		for i := range payRecordSubsidies {
			cellValue := 6 + dataCount
			var price int
			if payRecordSubsidies[i].Type == "charge" {
				price += payRecordSubsidies[i].Price
			} else {
				price += -payRecordSubsidies[i].Price
			}
			subtotal += price
			dateStr := payRecordSubsidies[i].StartDate.In(loc).Format("01-02") + "到" + payRecordSubsidies[i].EndDate.In(loc).Format("01-02")
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   contentStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          price,
				SheetName:      sheetName,
				ItemName:       payRecordSubsidies[i].ItemName,
				DateStr:        &dateStr,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 異動(請假)
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	json.Unmarshal(payRecord.TransferRefundLeave, &payRecordTransferRefundLeaves)
	if len(payRecordTransferRefundLeaves) > 0 {
		payRecordTransferRefundLeaveElements := make(map[string]*PayRecordTransferRefundLeaveForPrint)
		// 寫excel副標題
		// 寫入excel內容
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
				payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate.In(loc))
				payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate.In(loc))
			} else {
				payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate.In(loc))
				payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate.In(loc))
				if d.Type == "charge" {
					payRecordTransferRefundLeaveElements[d.ItemName].Price += d.Price
				} else {
					payRecordTransferRefundLeaveElements[d.ItemName].Price -= d.Price
				}
			}
		}

		keys := make([]string, 0, len(payRecordTransferRefundLeaveElements))
		for key := range payRecordTransferRefundLeaveElements {
			keys = append(keys, key)
		}

		for i := range keys {
			cellValue := 6 + dataCount
			var dateData []string
			// 由小到大排序
			sort.Slice(payRecordTransferRefundLeaveElements[keys[i]].StartDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[keys[i]].StartDates[j].Unix() < payRecordTransferRefundLeaveElements[keys[i]].StartDates[k].Unix()
			})
			sort.Slice(payRecordTransferRefundLeaveElements[keys[i]].EndDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[keys[i]].EndDates[j].Unix() < payRecordTransferRefundLeaveElements[keys[i]].EndDates[k].Unix()
			})

			for j := range payRecordTransferRefundLeaveElements[keys[i]].StartDates {
				dateData = append(dateData, payRecordTransferRefundLeaveElements[keys[i]].StartDates[j].In(loc).Format("01-02")+"到"+payRecordTransferRefundLeaveElements[keys[i]].EndDates[j].In(loc).Format("01-02"))
			}
			subtotal += payRecordTransferRefundLeaveElements[keys[i]].Price
			dateStr := strings.Join(dateData, "，")
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   contentStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          payRecordTransferRefundLeaveElements[keys[i]].Price,
				SheetName:      sheetName,
				ItemName:       "請假 " + payRecordTransferRefundLeaveElements[keys[i]].ItemName,
				DateStr:        &dateStr,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 非固定
	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
	json.Unmarshal(payRecord.NonFixedCharge, &payRecordNonFixedChargeRecords)
	// 排序要一樣(否則可能每次印出來的不太一樣)
	sort.Slice(payRecordNonFixedChargeRecords, func(j, k int) bool {
		return payRecordNonFixedChargeRecords[j].ItemName < payRecordNonFixedChargeRecords[k].ItemName
	})
	if len(payRecordNonFixedChargeRecords) > 0 {
		payRecordNonFixedChargeRecordElements := make(map[string]*PayRecordNonFixedChargeRecordDataForPartPrint)
		for _, d := range payRecordNonFixedChargeRecords {
			if payRecordNonFixedChargeRecordElements[d.ItemName] == nil {
				var price int
				if d.Type == "charge" {
					price = d.Subtotal
				} else {
					price = -d.Subtotal
				}
				subtotal += price
				payRecordNonFixedChargeRecordElements[d.ItemName] = &PayRecordNonFixedChargeRecordDataForPartPrint{
					ItemName:   d.ItemName,
					Subtotal:   price,
					Quantities: []int{d.Quantity},
					Date:       []time.Time{d.NonFixedChargeDate},
				}
			} else {
				if d.Type == "charge" {
					subtotal += d.Subtotal
					payRecordNonFixedChargeRecordElements[d.ItemName].Subtotal += d.Subtotal
				} else {
					subtotal -= d.Subtotal
					payRecordNonFixedChargeRecordElements[d.ItemName].Subtotal -= d.Subtotal
				}
				payRecordNonFixedChargeRecordElements[d.ItemName].Quantities = append(payRecordNonFixedChargeRecordElements[d.ItemName].Quantities, d.Quantity)
				payRecordNonFixedChargeRecordElements[d.ItemName].Date = append(payRecordNonFixedChargeRecordElements[d.ItemName].Date, d.NonFixedChargeDate)
			}
		}

		// 確保key的順序一致(不做這段 用payRecordNonFixedChargeRecordElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(payRecordNonFixedChargeRecordElements))
		for key := range payRecordNonFixedChargeRecordElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(j, k int) bool {
			return payRecordNonFixedChargeRecordElements[keys[j]].ItemName < payRecordNonFixedChargeRecordElements[keys[k]].ItemName
		})
		for i := range keys {
			var dateStrArray []string
			// 把時間做一下排序
			sort.SliceStable(payRecordNonFixedChargeRecordElements[keys[i]].Date, func(j, k int) bool {
				return payRecordNonFixedChargeRecordElements[keys[i]].Date[j].Unix() < payRecordNonFixedChargeRecordElements[keys[i]].Date[k].Unix()
			})
			for j := range payRecordNonFixedChargeRecordElements[keys[i]].Date {
				dateStrArray = append(dateStrArray, payRecordNonFixedChargeRecordElements[keys[i]].Date[j].In(loc).Format("01-02")+"("+strconv.Itoa(payRecordNonFixedChargeRecordElements[keys[i]].Quantities[j])+")")
			}
			dateStr := strings.Join(dateStrArray, "，")
			cellValue := 6 + dataCount
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   contentStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          payRecordNonFixedChargeRecordElements[keys[i]].Subtotal,
				SheetName:      sheetName,
				ItemName:       payRecordNonFixedChargeRecordElements[keys[i]].ItemName,
				DateStr:        &dateStr,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}
	// 如果有超過10 表示有新增欄位 要重算總和
	if dataCount >= 11 {
		subTotalcellValueStr := strconv.Itoa(insertRowCount + 11)
		partSubTotalcellValueStr := strconv.Itoa(29 + insertRowCount*2)
		f.SetCellFormula(sheetName, "C"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
	} else {
		f.SetCellFormula(sheetName, "C11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C29", strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F29", strconv.Itoa(subtotal)) // 聯單複寫
	}
	return insertRowCount
}

// 機構設定的金額顯示方式-項目(分稅別)
func printItemByTaxType(printItemByTaxTypeStruct *PrintItemByTaxTypeStruct) (bool, int) {
	f := printItemByTaxTypeStruct.F
	sheetName := printItemByTaxTypeStruct.SheetName
	payRecord := printItemByTaxTypeStruct.PayRecord
	taxType := printItemByTaxTypeStruct.TaxType
	classStyle := printItemByTaxTypeStruct.ClassStyle
	dateStyle := printItemByTaxTypeStruct.DateStyle
	priceStyle := printItemByTaxTypeStruct.PriceStyle

	loc, _ := time.LoadLocation("Asia/Taipei")
	var taxHaveData bool
	f.SetCellValue(sheetName, "A5", "項目")
	f.SetCellValue(sheetName, "A23", "項目")
	f.SetCellValue(sheetName, "D5", "項目")
	f.SetCellValue(sheetName, "D23", "項目")
	var subtotal int
	// 基本月費
	payRecordBasicCharges := []PayRecordBasicCharge{}
	json.Unmarshal(payRecord.BasicCharge, &payRecordBasicCharges)
	// 資料總筆數
	var dataCount int
	// 動態新增的行數
	var insertRowCount int
	if len(payRecordBasicCharges) > 0 {
		// 寫入excel內容
		for i := range payRecordBasicCharges {
			if taxType == payRecordBasicCharges[i].TaxType {
				taxHaveData = true
				cellValue := 6 + dataCount
				var price int
				if payRecordBasicCharges[i].Type == "charge" {
					price += payRecordBasicCharges[i].Price
				} else {
					price += -payRecordBasicCharges[i].Price
				}
				subtotal += price
				dateStr := payRecordBasicCharges[i].StartDate.In(loc).Format("01-02") + "到" + payRecordBasicCharges[i].EndDate.In(loc).Format("01-02")
				checkoutCellStruct := &CheckoutCellStruct{
					F:              f,
					DataCount:      dataCount,
					CellValue:      cellValue,
					InsertRowCount: insertRowCount,
					ContentStyle:   classStyle,
					DateStyle:      dateStyle,
					PriceStyle:     priceStyle,
					Price:          price,
					SheetName:      sheetName,
					ItemName:       payRecordBasicCharges[i].ItemName,
					DateStr:        &dateStr,
				}
				insertRowCount = checkCell(checkoutCellStruct)
				dataCount++
			}
		}
	}

	// 補助款
	payRecordSubsidies := []PayRecordSubsidy{}
	json.Unmarshal(payRecord.Subsidy, &payRecordSubsidies)
	if len(payRecordSubsidies) > 0 {
		// 寫入excel內容
		for i := range payRecordSubsidies {
			if taxType == "stampTax" {
				taxHaveData = true
				cellValue := 6 + dataCount
				var price int
				if payRecordSubsidies[i].Type == "charge" {
					price += payRecordSubsidies[i].Price
				} else {
					price += -payRecordSubsidies[i].Price
				}
				subtotal += price
				dateStr := payRecordSubsidies[i].StartDate.In(loc).Format("01-02") + "到" + payRecordSubsidies[i].EndDate.In(loc).Format("01-02")
				checkoutCellStruct := &CheckoutCellStruct{
					F:              f,
					DataCount:      dataCount,
					CellValue:      cellValue,
					InsertRowCount: insertRowCount,
					ContentStyle:   classStyle,
					DateStyle:      dateStyle,
					PriceStyle:     priceStyle,
					Price:          price,
					SheetName:      sheetName,
					ItemName:       payRecordSubsidies[i].ItemName,
					DateStr:        &dateStr,
				}
				insertRowCount = checkCell(checkoutCellStruct)
				dataCount++
			}
		}
	}
	// 異動(請假)
	payRecordTransferRefundLeaves := []PayRecordTransferRefundLeave{}
	json.Unmarshal(payRecord.TransferRefundLeave, &payRecordTransferRefundLeaves)
	if len(payRecordTransferRefundLeaves) > 0 {
		payRecordTransferRefundLeaveElements := make(map[string]*PayRecordTransferRefundLeaveForPrint)
		// 寫excel副標題
		// 寫入excel內容
		for _, d := range payRecordTransferRefundLeaves {
			if taxType == "stampTax" {
				taxHaveData = true
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
					payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate.In(loc))
					payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate.In(loc))
				} else {
					payRecordTransferRefundLeaveElements[d.ItemName].StartDates = append(payRecordTransferRefundLeaveElements[d.ItemName].StartDates, d.StartDate.In(loc))
					payRecordTransferRefundLeaveElements[d.ItemName].EndDates = append(payRecordTransferRefundLeaveElements[d.ItemName].EndDates, d.EndDate.In(loc))
					if d.Type == "charge" {
						payRecordTransferRefundLeaveElements[d.ItemName].Price += d.Price
					} else {
						payRecordTransferRefundLeaveElements[d.ItemName].Price -= d.Price
					}
				}
			}
		}

		for i := range payRecordTransferRefundLeaveElements {
			cellValue := 6 + dataCount
			var dateData []string
			// 由小到大排序
			sort.Slice(payRecordTransferRefundLeaveElements[i].StartDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[i].StartDates[j].Unix() < payRecordTransferRefundLeaveElements[i].StartDates[k].Unix()
			})
			sort.Slice(payRecordTransferRefundLeaveElements[i].EndDates, func(j, k int) bool {
				return payRecordTransferRefundLeaveElements[i].EndDates[j].Unix() < payRecordTransferRefundLeaveElements[i].EndDates[k].Unix()
			})

			for j := range payRecordTransferRefundLeaveElements[i].StartDates {
				dateData = append(dateData, payRecordTransferRefundLeaveElements[i].StartDates[j].In(loc).Format("01-02")+"到"+payRecordTransferRefundLeaveElements[i].EndDates[j].In(loc).Format("01-02"))
			}
			subtotal += payRecordTransferRefundLeaveElements[i].Price
			dateStr := strings.Join(dateData, "，")
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   classStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          payRecordTransferRefundLeaveElements[i].Price,
				SheetName:      sheetName,
				ItemName:       "請假 " + payRecordTransferRefundLeaveElements[i].ItemName,
				DateStr:        &dateStr,
			}
			insertRowCount = checkCell(checkoutCellStruct)
			dataCount++
		}
	}

	// 非固定
	payRecordNonFixedChargeRecords := []PayRecordNonFixedChargeRecordForPrint{}
	json.Unmarshal(payRecord.NonFixedCharge, &payRecordNonFixedChargeRecords)
	// 排序要一樣(否則可能每次印出來的不太一樣)
	sort.Slice(payRecordNonFixedChargeRecords, func(j, k int) bool {
		return payRecordNonFixedChargeRecords[j].ItemName < payRecordNonFixedChargeRecords[k].ItemName
	})
	if len(payRecordNonFixedChargeRecords) > 0 {
		payRecordNonFixedChargeRecordElements := make(map[string]*PayRecordNonFixedChargeRecordDataForPartPrint)
		for _, d := range payRecordNonFixedChargeRecords {
			if d.TaxType == taxType {
				taxHaveData = true
				if payRecordNonFixedChargeRecordElements[d.ItemName] == nil {
					var price int
					if d.Type == "charge" {
						price = d.Subtotal
					} else {
						price = -d.Subtotal
					}
					subtotal += price
					payRecordNonFixedChargeRecordElements[d.ItemName] = &PayRecordNonFixedChargeRecordDataForPartPrint{
						ItemName:   d.ItemName,
						Subtotal:   price,
						Quantities: []int{d.Quantity},
						Date:       []time.Time{d.NonFixedChargeDate},
					}
				} else {
					if d.Type == "charge" {
						subtotal += d.Subtotal
						payRecordNonFixedChargeRecordElements[d.ItemName].Subtotal += d.Subtotal
					} else {
						subtotal -= d.Subtotal
						payRecordNonFixedChargeRecordElements[d.ItemName].Subtotal -= d.Subtotal
					}
					payRecordNonFixedChargeRecordElements[d.ItemName].Quantities = append(payRecordNonFixedChargeRecordElements[d.ItemName].Quantities, d.Quantity)
					payRecordNonFixedChargeRecordElements[d.ItemName].Date = append(payRecordNonFixedChargeRecordElements[d.ItemName].Date, d.NonFixedChargeDate)
				}
			}
		}

		// 確保key的順序一致(不做這段 用payRecordNonFixedChargeRecordElements跑迴圈順序會亂跳)
		keys := make([]string, 0, len(payRecordNonFixedChargeRecordElements))
		for key := range payRecordNonFixedChargeRecordElements {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(j, k int) bool {
			return payRecordNonFixedChargeRecordElements[keys[j]].ItemName < payRecordNonFixedChargeRecordElements[keys[k]].ItemName
		})
		for i := range keys {
			var dateStrArray []string
			// 把時間做一下排序
			sort.SliceStable(payRecordNonFixedChargeRecordElements[keys[i]].Date, func(j, k int) bool {
				return payRecordNonFixedChargeRecordElements[keys[i]].Date[j].Unix() < payRecordNonFixedChargeRecordElements[keys[i]].Date[k].Unix()
			})
			for j := range payRecordNonFixedChargeRecordElements[keys[i]].Date {
				dateStrArray = append(dateStrArray, payRecordNonFixedChargeRecordElements[keys[i]].Date[j].In(loc).Format("01-02")+"("+strconv.Itoa(payRecordNonFixedChargeRecordElements[keys[i]].Quantities[j])+")")
			}
			dateStr := strings.Join(dateStrArray, "，")
			cellValue := 6 + dataCount
			checkoutCellStruct := &CheckoutCellStruct{
				F:              f,
				DataCount:      dataCount,
				CellValue:      cellValue,
				InsertRowCount: insertRowCount,
				ContentStyle:   classStyle,
				DateStyle:      dateStyle,
				PriceStyle:     priceStyle,
				Price:          payRecordNonFixedChargeRecordElements[keys[i]].Subtotal,
				SheetName:      sheetName,
				ItemName:       payRecordNonFixedChargeRecordElements[keys[i]].ItemName,
				DateStr:        &dateStr,
			}
			if i < 13 {
				insertRowCount = checkCell(checkoutCellStruct)
			}
			dataCount++
		}
	}
	// 如果有超過10 表示有新增欄位 要重算總和
	if dataCount >= 11 {
		subTotalcellValueStr := strconv.Itoa(insertRowCount + 11)
		partSubTotalcellValueStr := strconv.Itoa(29 + insertRowCount*2)
		f.SetCellFormula(sheetName, "C"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F"+subTotalcellValueStr, strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F"+partSubTotalcellValueStr, strconv.Itoa(subtotal)) // 聯單複寫
	} else {
		f.SetCellFormula(sheetName, "C11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "F11", strconv.Itoa(subtotal))
		f.SetCellFormula(sheetName, "C29", strconv.Itoa(subtotal)) // 聯單複寫
		f.SetCellFormula(sheetName, "F29", strconv.Itoa(subtotal)) // 聯單複寫
	}
	return taxHaveData, insertRowCount
}

// 這邊是在判斷需不需要新增列及塞資料用的function
func checkCell(checkoutCellStruct *CheckoutCellStruct) int {
	itemName := checkoutCellStruct.ItemName
	dataCount := checkoutCellStruct.DataCount
	cellValue := checkoutCellStruct.CellValue
	insertRowCount := checkoutCellStruct.InsertRowCount
	f := checkoutCellStruct.F
	sheetName := checkoutCellStruct.SheetName
	contentStyle := checkoutCellStruct.ContentStyle
	dateStyle := checkoutCellStruct.DateStyle
	priceStyle := checkoutCellStruct.PriceStyle
	dateStr := checkoutCellStruct.DateStr
	price := checkoutCellStruct.Price
	// 上面的聯單
	var fontSize int
	if dateStr == nil {
		fontSize = len(itemName) / 23
	} else {
		if len(*dateStr) > len(itemName) {
			fontSize = len(*dateStr) / 23
		} else {
			fontSize = len(itemName) / 23
		}
	}
	fmt.Println("dataCount", dataCount)
	if dataCount >= 10 {
		fmt.Println("a")
		cellValueStr := strconv.Itoa(cellValue - 5 - insertRowCount)
		var partCellValueStr string
		cellItemStr := "D"
		cellDateStr := "E"
		cellPriceStr := "F"
		if dataCount%2 == 0 {
			// 先跑下面的聯單
			f.InsertRow(sheetName, dataCount+19)
			f.InsertRow(sheetName, dataCount+1-insertRowCount)
			cellItemStr = "A"
			cellDateStr = "B"
			cellPriceStr = "C"
			partCellValueStr = strconv.Itoa(cellValue + 14)
			//	取第幾列
			cellValue, _ := strconv.Atoi(cellValueStr)
			partCellValue, _ := strconv.Atoi(partCellValueStr)
			// 調高度
			f.SetRowHeight(sheetName, cellValue, float64((fontSize+1)*16))
			f.SetRowHeight(sheetName, partCellValue, float64((fontSize+1)*16))
			// 調style
			f.SetCellStyle(sheetName, "A"+cellValueStr, "A"+cellValueStr, contentStyle)
			f.SetCellStyle(sheetName, "B"+cellValueStr, "B"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "C"+cellValueStr, "C"+cellValueStr, priceStyle)
			f.SetCellStyle(sheetName, "D"+cellValueStr, "D"+cellValueStr, contentStyle)
			f.SetCellStyle(sheetName, "E"+cellValueStr, "E"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "F"+cellValueStr, "F"+cellValueStr, priceStyle)
			f.SetCellStyle(sheetName, "A"+partCellValueStr, "A"+partCellValueStr, contentStyle) // 聯單複寫
			f.SetCellStyle(sheetName, "B"+partCellValueStr, "B"+partCellValueStr, dateStyle)    // 聯單複寫
			f.SetCellStyle(sheetName, "C"+partCellValueStr, "C"+partCellValueStr, priceStyle)   // 聯單複寫
			f.SetCellStyle(sheetName, "D"+partCellValueStr, "D"+partCellValueStr, contentStyle) // 聯單複寫
			f.SetCellStyle(sheetName, "E"+partCellValueStr, "E"+partCellValueStr, dateStyle)    // 聯單複寫
			f.SetCellStyle(sheetName, "F"+partCellValueStr, "F"+partCellValueStr, priceStyle)   // 聯單複寫
			insertRowCount++
		} else {
			partCellValueStr = strconv.Itoa(cellValue + 13)
			var rowHieght float64
			rowHieght, _ = f.GetRowHeight(sheetName, cellValue-5-insertRowCount)
			if rowHieght < float64((fontSize+1)*16) {
				rowHieght = float64((fontSize + 1) * 16)
			}
			f.SetRowHeight(sheetName, cellValue-5-insertRowCount, rowHieght)
			f.SetRowHeight(sheetName, cellValue+13, rowHieght) // 聯單複寫
		}
		f.SetCellValue(sheetName, cellItemStr+cellValueStr, itemName)
		f.SetCellValue(sheetName, cellPriceStr+cellValueStr, price)
		if dateStr != nil {
			// 表示需要給項目
			// 給時間
			f.SetCellValue(sheetName, cellDateStr+cellValueStr, *dateStr)
			f.SetCellValue(sheetName, cellDateStr+partCellValueStr, *dateStr) // 聯單複寫
		} else {
			// 表示為類別加總
			// 需要合併儲存格
			f.MergeCell(sheetName, "A"+cellValueStr, "B"+cellValueStr)
			f.MergeCell(sheetName, "A"+partCellValueStr, "B"+partCellValueStr) // 聯單複寫
			f.MergeCell(sheetName, "D"+cellValueStr, "E"+cellValueStr)
			f.MergeCell(sheetName, "D"+partCellValueStr, "E"+partCellValueStr) // 聯單複寫
		}
		f.SetCellValue(sheetName, cellItemStr+partCellValueStr, itemName) // 聯單複寫
		f.SetCellValue(sheetName, cellPriceStr+partCellValueStr, price)   // 聯單複寫
	} else if dataCount >= 5 && dataCount < 10 {
		cellValueStr := strconv.Itoa(cellValue - 5)
		partCellValueStr := strconv.Itoa(cellValue + 13)
		var rowHieght float64
		rowHieght, _ = f.GetRowHeight(sheetName, cellValue-5)
		if rowHieght < float64((fontSize+1)*16) {
			rowHieght = float64((fontSize + 1) * 16)
		}
		f.SetRowHeight(sheetName, cellValue-5, rowHieght)
		f.SetRowHeight(sheetName, cellValue+13, rowHieght) // 聯單複寫

		// 寫excel項目名稱
		f.SetCellValue(sheetName, "D"+cellValueStr, itemName)
		f.SetCellValue(sheetName, "D"+partCellValueStr, itemName) // 聯單複寫
		if dateStr != nil {
			// 表示需要給時間
			f.SetCellStyle(sheetName, "B"+cellValueStr, "B"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "B"+partCellValueStr, "B"+partCellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "E"+cellValueStr, "E"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "E"+partCellValueStr, "E"+partCellValueStr, dateStyle)

			f.SetCellValue(sheetName, "E"+cellValueStr, *dateStr)
			f.SetCellValue(sheetName, "E"+partCellValueStr, *dateStr) // 聯單複寫
		}
		f.SetCellValue(sheetName, "F"+cellValueStr, price)
		f.SetCellValue(sheetName, "F"+partCellValueStr, price) // 聯單複寫
	} else {
		f.SetRowHeight(sheetName, cellValue, float64((fontSize+1)*16))
		f.SetRowHeight(sheetName, cellValue+18, float64((fontSize+1)*16)) // 聯單複寫
		cellValueStr := strconv.Itoa(cellValue)
		partCellValueStr := strconv.Itoa(cellValue + 18)
		f.SetCellValue(sheetName, "A"+cellValueStr, itemName)
		f.SetCellValue(sheetName, "A"+partCellValueStr, itemName) // 聯單複寫
		if dateStr != nil {
			// 表示需要給時間
			f.SetCellStyle(sheetName, "B"+cellValueStr, "B"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "B"+partCellValueStr, "B"+partCellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "E"+cellValueStr, "E"+cellValueStr, dateStyle)
			f.SetCellStyle(sheetName, "E"+partCellValueStr, "E"+partCellValueStr, dateStyle)
			f.SetCellValue(sheetName, "B"+cellValueStr, *dateStr)
			f.SetCellValue(sheetName, "B"+partCellValueStr, *dateStr) // 聯單複寫
		}
		f.SetCellValue(sheetName, "C"+cellValueStr, price)
		f.SetCellValue(sheetName, "C"+partCellValueStr, price) // 聯單複寫
	}
	return insertRowCount
}

// 用來查organizationReceiptTemplateSetting的PatientInfo help function
func organizationReceiptTemplateSettingPatientInfoIncludes(organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting, includeString string) bool {
	for i := range organizationReceiptTemplateSetting.PatientInfo {
		if organizationReceiptTemplateSetting.PatientInfo[i] == includeString {
			return true
		}
	}
	return false
}

// 用來查organizationReceiptTemplateSetting的OrganizationInfoOne help function
func organizationReceiptTemplateSettingOrganizationInfoOneIncludes(organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting, includeString string) bool {
	for i := range organizationReceiptTemplateSetting.OrganizationInfoOne {
		if organizationReceiptTemplateSetting.OrganizationInfoOne[i] == includeString {
			return true
		}
	}
	return false
}

// 用來查organizationReceiptTemplateSetting的OrganizationInfoTwo help function
func organizationReceiptTemplateSettingOrganizationInfoTwoIncludes(organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting, includeString string) bool {
	for i := range organizationReceiptTemplateSetting.OrganizationInfoTwo {
		if organizationReceiptTemplateSetting.OrganizationInfoTwo[i] == includeString {
			return true
		}
	}
	return false
}
