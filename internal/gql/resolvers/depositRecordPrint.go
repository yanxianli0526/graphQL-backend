package resolvers

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/models"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

// Queries
func (r *queryResolver) PrintDepositRecord(ctx context.Context, id string) (*string, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PrintDepositRecord uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "printDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	depositRecordId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("PrintDepositRecord uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "printDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	depositRecord, err := orm.GetDepositRecordForPrintById(r.ORM.DB, depositRecordId, organizationId)
	if err != nil {
		r.Logger.Error("PrintDepositRecord orm.GetDepositRecordForPrintById", zap.Error(err), zap.String("fieldName", "printDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	publicUrl, err := printDepositRecord(r, ctx, depositRecord)
	if err != nil {
		r.Logger.Error("PrintDepositRecord printDepositRecord", zap.Error(err), zap.Error(err), zap.String("fieldName", "printDepositRecord"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("printDepositRecord run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "printDepositRecord"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return &publicUrl, nil
}

//go:embed excelTemplate/depositRecord.xlsx
var excelDepositRecordBytes []byte

func printDepositRecord(r *queryResolver, ctx context.Context, depositRecord *models.DepositRecord) (string, error) {

	reader := bytes.NewReader(excelDepositRecordBytes)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.Logger.Error("printDepositRecord excelize.OpenReader", zap.Error(err))
		return "", err
	}

	f.SetCellValue("保證金收據", "A1", depositRecord.Organization.Name)
	f.SetCellValue("保證金收據", "A3", "床號："+depositRecord.Patient.Room+depositRecord.Patient.Bed)

	f.SetCellValue("保證金收據", "C3", "姓名："+depositRecord.Patient.LastName+depositRecord.Patient.FirstName)

	var idNumber string
	if depositRecord.Organization.Privacy == "unmask" {
		idNumber = depositRecord.Patient.IdNumber
	} else {
		if len(depositRecord.Patient.IdNumber) >= 10 {
			idNumber = depositRecord.Patient.IdNumber[0:3] + "****" + depositRecord.Patient.IdNumber[7:10]
		} else {
			idNumberCount := 0
			for idNumberCount < len(depositRecord.Patient.IdNumber) {
				if idNumberCount >= 4 && idNumberCount <= 8 {
					idNumber += "*"
				} else {
					yo := string([]rune(depositRecord.Patient.IdNumber)[idNumberCount])
					idNumber += yo
				}
				idNumberCount += 1
			}
		}
	}
	fmt.Println("idNumber", idNumber)
	f.SetCellValue("保證金收據", "E3", "身分證字號："+idNumber)
	f.SetCellValue("保證金收據", "G3", "收據單號："+depositRecord.IdNumber)

	f.SetCellValue("保證金收據", "H4", *depositRecord.Organization.TaxIdNumber)
	if depositRecord.Type == "charge" {
		f.SetCellValue("保證金收據", "A5", "入住保證金") // 項目名稱
	} else {
		f.SetCellValue("保證金收據", "A5", "入住保證金 退費") // 項目名稱
	}
	f.SetCellValue("保證金收據", "C5", depositRecord.Price)
	f.SetCellValue("保證金收據", "D5", depositRecord.Note)
	// f.SetCellStyle("保證金收據", "D5", "D5", wrapTextStyle)
	var owner string
	if depositRecord.Organization.Owner != nil {
		owner = *depositRecord.Organization.Owner
	}
	f.SetCellValue("保證金收據", "H5", owner)
	var addressCity string
	if depositRecord.Organization.AddressCity != nil {
		addressCity = *depositRecord.Organization.AddressCity
	}
	var addressDistrict string
	if depositRecord.Organization.AddressDistrict != nil {
		addressDistrict = *depositRecord.Organization.AddressDistrict
	}
	var address string
	if depositRecord.Organization.Address != nil {
		address = *depositRecord.Organization.Address
	}
	f.SetCellValue("保證金收據", "H6", addressCity+addressDistrict+address)
	f.SetCellValue("保證金收據", "H11", depositRecord.User.DisplayName)

	// 合計金額(備註),繳(退)款日期需要function處理
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	f.SetCellFormula("保證金收據", "H12", fmt.Sprintf("=DATE(%s,%s,%s)", year, month, day))
	f.SetCellFormula("保證金收據", "C12", "=SUM(C5:C11)")
	f.SetCellFormula("保證金收據", "D12", "=C12")

	// 機構名稱
	f.SetCellFormula("保證金收據", "A16", "=A1")
	f.SetCellFormula("保證金收據", "A31", "=A1")
	// 床號
	f.SetCellFormula("保證金收據", "A18", "=A3")
	f.SetCellFormula("保證金收據", "A33", "=A3")
	// 姓名
	f.SetCellFormula("保證金收據", "C18", "=C3")
	f.SetCellFormula("保證金收據", "C33", "=C3")
	// 身分證字號
	f.SetCellFormula("保證金收據", "E18", "=E3")
	f.SetCellFormula("保證金收據", "E33", "=E3")
	// 收據單號
	f.SetCellFormula("保證金收據", "G18", "=G3")
	f.SetCellFormula("保證金收據", "G33", "=G3")
	// 統一編號
	f.SetCellFormula("保證金收據", "H19", "=H4")
	f.SetCellFormula("保證金收據", "H34", "=H4")
	// 項目
	f.SetCellFormula("保證金收據", "A20", "=A5")
	f.SetCellFormula("保證金收據", "A35", "=A5")
	// 金額
	f.SetCellFormula("保證金收據", "C20", "=C5")
	f.SetCellFormula("保證金收據", "C35", "=C5")
	// 備註
	f.SetCellFormula("保證金收據", "D20", "=D5")
	f.SetCellFormula("保證金收據", "D35", "=D5")
	// 負責人
	f.SetCellFormula("保證金收據", "H20", "=H5")
	f.SetCellFormula("保證金收據", "H35", "=H5")
	// 機構地址
	f.SetCellFormula("保證金收據", "H21", "=H6")
	f.SetCellFormula("保證金收據", "H36", "=H6")
	// 經辦人
	f.SetCellFormula("保證金收據", "H26", "=H11")
	f.SetCellFormula("保證金收據", "H41", "=H11")
	// 合計
	f.SetCellFormula("保證金收據", "C27", "=SUM(C20:C26)")
	f.SetCellFormula("保證金收據", "C42", "=SUM(C35:C41)")
	// 備註
	f.SetCellFormula("保證金收據", "D27", "=C27")
	f.SetCellFormula("保證金收據", "D42", "=C42")
	// 繳(退)款日期
	f.SetCellFormula("保證金收據", "H27", "=H12")
	f.SetCellFormula("保證金收據", "H42", "=H12")

	fileName := depositRecord.Patient.Branch + depositRecord.Patient.Room + depositRecord.Patient.Bed + depositRecord.Patient.LastName + depositRecord.Patient.FirstName + "-保證金收據-" + depositRecord.IdNumber + ".xlsx"
	if err := f.SaveAs(fileName); err != nil {
		r.Logger.Error("printDepositRecord f.SaveAs", zap.Error(err))
		return "", err
	}
	err = r.store.UploadFile(ctx, fileName, "depositRecord")
	if err != nil {
		r.Logger.Error("printDepositRecord r.store.UploadFile", zap.Error(err))
		return "", err
	}
	err = r.store.SetMetadata(ctx, fileName, "depositRecord")
	if err != nil {
		r.Logger.Error("printDepositRecord r.store.SetMetadata", zap.Error(err))
		return "", err
	}
	publicUrl := r.store.GenPublicLink("depositRecord/" + fileName)
	return publicUrl, nil
}
