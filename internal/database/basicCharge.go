package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBasicCharge(db *gorm.DB, basicCharge *models.BasicCharge) error {
	return db.Create(basicCharge).Error
}

//這邊是住民帳單開帳時會呼叫
func CreateBasicCharges(db *gorm.DB, basicCharge []*models.BasicCharge) error {
	return db.Create(basicCharge).Error
}

func UpdatePatientBillBasicCharge(db *gorm.DB, patientBill *models.BasicCharge) error {
	return db.Model(&models.BasicCharge{
		ID:             patientBill.ID,
		OrganizationId: patientBill.OrganizationId,
	}).Updates(map[string]interface{}{
		"item_name":  patientBill.ItemName,
		"type":       patientBill.Type,
		"unit":       patientBill.Unit,
		"price":      patientBill.Price,
		"tax_type":   patientBill.TaxType,
		"start_date": patientBill.StartDate,
		"end_date":   patientBill.EndDate,
		"note":       patientBill.Note,
		"user_id":    patientBill.UserId,
	}).Error

}

func AppendAssociationsPatientBillBasicCharge(db *gorm.DB, patientBill *models.PatientBill, basicCharge models.BasicCharge) error {
	return db.Model(&patientBill).Association("BasicCharges").Append(&basicCharge)
}

func GetBasicCharge(db *gorm.DB, organizationId, BasicChargeId uuid.UUID) (*models.BasicCharge, error) {
	var basicCharge *models.BasicCharge

	result := db.Preload("User").Where("organization_id = ? AND id = ?", organizationId, BasicChargeId).Find(&basicCharge)

	return basicCharge, result.Error
}

func UpdateBasicChargesReceiptStatus(db *gorm.DB, needUpdateBasicChargesId []uuid.UUID, status string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.BasicCharge{}).Where("id in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateBasicChargesId, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.BasicCharge{}).Where("id in ?  AND receipt_status = ?", needUpdateBasicChargesId, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.BasicCharge{}).Where("id in ?", needUpdateBasicChargesId).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}
}

func UpdateBasicChargesReceiptStatusInTaxType(db *gorm.DB, needUpdateBasicChargesId []uuid.UUID, status string, taxType []string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.BasicCharge{}).Where("id in ? AND tax_type in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateBasicChargesId, taxType, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.BasicCharge{}).Where("id in ? AND tax_type in ? AND receipt_status = ?", needUpdateBasicChargesId, taxType, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.BasicCharge{}).Where("id in ? AND tax_type in ?", needUpdateBasicChargesId, taxType).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}
}

func DeleteBasicCharge(db *gorm.DB, organizationId, basicChargeId uuid.UUID) error {
	return db.Where("organization_id = ? AND id = ?", organizationId, basicChargeId).Delete(&models.BasicCharge{}).Error
}
