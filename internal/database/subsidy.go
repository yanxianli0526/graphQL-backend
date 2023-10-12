package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSubsidy(db *gorm.DB, subsidy *models.Subsidy) error {
	return db.Create(subsidy).Error
}

func CreateSubsidies(db *gorm.DB, subsidies []*models.Subsidy) error {
	return db.Create(subsidies).Error
}

func UpdatePatientBillSubsidy(db *gorm.DB, subsidy *models.Subsidy) error {
	return db.Model(&models.Subsidy{
		ID:             subsidy.ID,
		OrganizationId: subsidy.OrganizationId,
	}).Updates(map[string]interface{}{
		"item_name":  subsidy.ItemName,
		"type":       subsidy.Type,
		"price":      subsidy.Price,
		"unit":       subsidy.Unit,
		"id_number":  subsidy.IdNumber,
		"start_date": subsidy.StartDate,
		"end_date":   subsidy.EndDate,
		"note":       subsidy.Note,
		"user_id":    subsidy.UserId,
	}).Error

}

func AppendAssociationsPatientBillSubSidy(db *gorm.DB, patientBill *models.PatientBill, subSidy models.Subsidy) error {
	return db.Model(&patientBill).Association("Subsidies").Append(&subSidy)
}

func GetSubsidy(db *gorm.DB, organizationId, SubsidyId uuid.UUID) (*models.Subsidy, error) {
	var subsidy *models.Subsidy
	result := db.Preload("User").Where("organization_id = ? AND id = ?", organizationId, SubsidyId).Find(&subsidy)
	return subsidy, result.Error
}

func UpdateSubsidiesReceiptStatus(db *gorm.DB, needUpdateSubsidiesId []uuid.UUID, status string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.Subsidy{}).Where("id in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateSubsidiesId, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.Subsidy{}).Where("id in ?  AND receipt_status = ?", needUpdateSubsidiesId, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.Subsidy{}).Where("id in ?", needUpdateSubsidiesId).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}
}

func DeleteSubsidy(db *gorm.DB, organizationId, subsidyId uuid.UUID) error {
	return db.Where("organization_id = ? AND id = ?", organizationId, subsidyId).Delete(&models.Subsidy{}).Error
}
