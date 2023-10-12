package resolvers

import (
	"encoding/json"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdatePatientBillAmountDueAndCumulativeUnpaidAmountStruct struct {
	OrganizationId uuid.UUID
	PatientId      uuid.UUID
	PatientBill    *models.PatientBill
	ItemSubtotal   int
	Tx             *gorm.DB
	CRUDType       string
}

type UpdatePatientBillAmountDueStruct struct {
	PatientBill  *models.PatientBill
	ItemSubtotal int
	NewAmountDue int
	Tx           *gorm.DB
}

// 取得多筆住民帳單的應繳金額
func GetPatientBillsAmountDue(patientBills []*models.PatientBill, havePatientConsumptionRecord bool) ([]*models.PatientBill, error) {
	for i := range patientBills {
		var subtotal int
		// 基本月費
		if len(patientBills[i].BasicCharges) > 0 {
			for j := range patientBills[i].BasicCharges {
				if patientBills[i].BasicCharges[j].Type == "charge" {
					subtotal += patientBills[i].BasicCharges[j].Price
				} else {
					subtotal -= patientBills[i].BasicCharges[j].Price
				}
			}
		}
		if len(patientBills[i].Subsidies) > 0 {
			for j := range patientBills[i].Subsidies {
				if patientBills[i].Subsidies[j].Type == "charge" {
					subtotal += patientBills[i].Subsidies[j].Price
				} else {
					subtotal -= patientBills[i].Subsidies[j].Price
				}
			}
		}

		if len(patientBills[i].TransferRefundLeaves) > 0 {
			for j := range patientBills[i].TransferRefundLeaves {
				transferRefundItems := []TransferRefundItem{}
				json.Unmarshal(patientBills[i].TransferRefundLeaves[j].Items, &transferRefundItems)
				// 用item去跑迴圈
				for k := range transferRefundItems {
					// 看看是收費還是退費
					if transferRefundItems[k].Type == "charge" {
						subtotal += transferRefundItems[k].Price
					} else {
						subtotal -= transferRefundItems[k].Price
					}
				}
			}
		}

		if len(patientBills[i].NonFixedChargeRecords) > 0 {
			for j := range patientBills[i].NonFixedChargeRecords {
				if patientBills[i].NonFixedChargeRecords[j].Type == "charge" {
					subtotal += patientBills[i].NonFixedChargeRecords[j].Price * int(patientBills[i].NonFixedChargeRecords[j].Quantity)
				} else {
					subtotal -= patientBills[i].NonFixedChargeRecords[j].Price * int(patientBills[i].NonFixedChargeRecords[j].Quantity)
				}
			}
		}
		patientBills[i].AmountDue = subtotal
	}
	return patientBills, nil
}

// 取得單筆住民帳單的應繳金額
func GetPatientBillAmountDue(patientBill *models.PatientBill, havePatientConsumptionRecord bool) (models.PatientBill, error) {
	var subtotal int
	if len(patientBill.BasicCharges) > 0 {
		for j := range patientBill.BasicCharges {
			if patientBill.BasicCharges[j].Type == "charge" {
				subtotal += patientBill.BasicCharges[j].Price
			} else {
				subtotal -= patientBill.BasicCharges[j].Price
			}
		}
	}
	if len(patientBill.Subsidies) > 0 {
		for j := range patientBill.Subsidies {
			if patientBill.Subsidies[j].Type == "charge" {
				subtotal += patientBill.Subsidies[j].Price
			} else {
				subtotal -= patientBill.Subsidies[j].Price
			}
		}
	}
	if len(patientBill.TransferRefundLeaves) > 0 {
		for j := range patientBill.TransferRefundLeaves {
			transferRefundItems := []TransferRefundItem{}
			json.Unmarshal(patientBill.TransferRefundLeaves[j].Items, &transferRefundItems)
			// 用item去跑迴圈
			for k := range transferRefundItems {
				// 看看是收費還是退費
				if transferRefundItems[k].Type == "charge" {
					subtotal += transferRefundItems[k].Price
				} else {
					subtotal -= transferRefundItems[k].Price
				}
			}
		}
	}
	if len(patientBill.NonFixedChargeRecords) > 0 {
		for j := range patientBill.NonFixedChargeRecords {

			if patientBill.NonFixedChargeRecords[j].Type == "charge" {
				subtotal += patientBill.NonFixedChargeRecords[j].Price * int(patientBill.NonFixedChargeRecords[j].Quantity)
			} else {
				subtotal -= patientBill.NonFixedChargeRecords[j].Price * int(patientBill.NonFixedChargeRecords[j].Quantity)
			}
		}
	}
	patientBill.AmountDue = subtotal
	return *patientBill, nil
}

// 更新住民帳單的應繳金額
func UpdatePatientBillAmountDue(updatePatientBillAmountDueStruct UpdatePatientBillAmountDueStruct) error {
	patientBill := updatePatientBillAmountDueStruct.PatientBill
	tx := updatePatientBillAmountDueStruct.Tx
	// 帳單的應繳金額
	patientBill.AmountDue = updatePatientBillAmountDueStruct.NewAmountDue
	// 更新住民帳單的應繳金額
	err := orm.UpdatePatientBillAmountDue(tx, patientBill)
	if err != nil {
		return err
	}

	return nil
}
