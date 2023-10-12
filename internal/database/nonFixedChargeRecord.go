package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateNonFixedChargeRecords(db *gorm.DB, nonFixedChargeRecords []*models.NonFixedChargeRecord) error {
	return db.Create(nonFixedChargeRecords).Error
}

func UpdateAndGetNonFixedChargeRecord(db *gorm.DB, updateNonFixedChargeRecord *models.NonFixedChargeRecord) (*models.NonFixedChargeRecord, error) {
	nonFixedChargeRecord := models.NonFixedChargeRecord{ID: updateNonFixedChargeRecord.ID,
		OrganizationId: updateNonFixedChargeRecord.OrganizationId,
	}
	result := db.Model(&nonFixedChargeRecord).Clauses(clause.Returning{}).Updates(map[string]interface{}{
		"non_fixed_charge_date": updateNonFixedChargeRecord.NonFixedChargeDate,
		"item_category":         updateNonFixedChargeRecord.ItemCategory,
		"item_name":             updateNonFixedChargeRecord.ItemName,
		"type":                  updateNonFixedChargeRecord.Type,
		"unit":                  updateNonFixedChargeRecord.Unit,
		"price":                 updateNonFixedChargeRecord.Price,
		"quantity":              updateNonFixedChargeRecord.Quantity,
		"subtotal":              updateNonFixedChargeRecord.Subtotal,
		"note":                  updateNonFixedChargeRecord.Note,
		"tax_type":              updateNonFixedChargeRecord.TaxType,
		"user_id":               updateNonFixedChargeRecord.UserId,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	return &nonFixedChargeRecord, result.Error
}

func UpdateNonFixedChargesRecordsReceiptStatus(db *gorm.DB, needUpdateNonFixedChargeRecordsId []uuid.UUID, status string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateNonFixedChargeRecordsId, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ?  AND receipt_status = ?", needUpdateNonFixedChargeRecordsId, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ?", needUpdateNonFixedChargeRecordsId).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}
}

func UpdateNonFixedChargesRecordsReceiptStatusInTaxType(db *gorm.DB, needUpdateNonFixedChargeRecordsId []uuid.UUID, status string, taxType []string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ? AND tax_type in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateNonFixedChargeRecordsId, taxType, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ? AND tax_type in ? AND receipt_status = ?", needUpdateNonFixedChargeRecordsId, taxType, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.NonFixedChargeRecord{}).Where("id in ? AND tax_type in ?", needUpdateNonFixedChargeRecordsId, taxType).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}

}

func DeleteNonFixedChargeRecord(db *gorm.DB, organizationId, nonFixedChargeRecordId uuid.UUID) error {
	return db.Where("organization_id = ? AND id = ?", organizationId, nonFixedChargeRecordId).Delete(&models.NonFixedChargeRecord{}).Error
}

func GetNonFixedChargeRecord(db *gorm.DB, nonFixedChargeRecordId uuid.UUID, preloadPatient, preloadUser bool) (*models.NonFixedChargeRecord, error) {
	var nonFixedChargeRecords *models.NonFixedChargeRecord

	tx := db.Table("non_fixed_charge_records")

	if preloadPatient {
		tx = tx.Preload("Patient")
	}
	if preloadUser {
		tx = tx.Preload("User")
	}

	err := tx.Where("id = ?", nonFixedChargeRecordId).Find(&nonFixedChargeRecords).Error
	if err != nil {
		return nil, err
	}
	return nonFixedChargeRecords, nil
}

func GetNonFixedChargeRecordsByPatientId(db *gorm.DB, patientId uuid.UUID) ([]*models.NonFixedChargeRecord, error) {
	var nonFixedChargeRecords []*models.NonFixedChargeRecord
	err := db.Preload("Patient").Preload("User").Where("patient_id = ?", patientId).Find(&nonFixedChargeRecords).Error
	if err != nil {
		return nil, err
	}
	return nonFixedChargeRecords, nil
}

func GetNonFixedChargeRecordsByPatientIdAndDate(db *gorm.DB, patientId uuid.UUID, nonFixedChargeStartDate, nonFixedChargeEndDate time.Time) ([]*models.NonFixedChargeRecord, error) {
	var nonFixedChargeRecords []*models.NonFixedChargeRecord
	err := db.Preload("Patient").Preload("User").Order("non_fixed_charge_date, item_category, item_name ASC").Where("patient_id = ? AND non_fixed_charge_date BETWEEN ? AND ?", patientId, nonFixedChargeStartDate, nonFixedChargeEndDate).Find(&nonFixedChargeRecords).Error
	if err != nil {
		return nil, err
	}
	return nonFixedChargeRecords, nil
}

func GetNonFixedChargeRecordsInPatientIdsAndDate(db *gorm.DB, patientsId []uuid.UUID, nonFixedChargeStartDate, nonFixedChargeEndDate time.Time) ([]*models.NonFixedChargeRecord, error) {
	var nonFixedChargeRecords []*models.NonFixedChargeRecord
	err := db.Preload("Patient").Preload("User").Where("patient_id in ? AND non_fixed_charge_date BETWEEN ? AND ?", patientsId, nonFixedChargeStartDate, nonFixedChargeEndDate).Find(&nonFixedChargeRecords).Error
	if err != nil {
		return nil, err
	}
	return nonFixedChargeRecords, nil
}

func PatientLatestNonFixedChargeRecords(db *gorm.DB, organizationId uuid.UUID) ([]*models.NonFixedChargeRecord, []*models.NonFixedChargeRecord, []*models.NonFixedChargeRecord, error) {
	var thisMonthNonFixedChargeRecords []*models.NonFixedChargeRecord
	var lastMonthNonFixedChargeRecords []*models.NonFixedChargeRecord
	var nonFixedChargeRecordsDescByUpdatedAt []*models.NonFixedChargeRecord

	// 本月
	loc, _ := time.LoadLocation("Asia/Taipei")
	thisMonthFirstDay := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	thisMonthFirstDay, _ = time.ParseInLocation("2006-01-02 15:04:05", thisMonthFirstDay.Format("2006-01-02")+" 00:00:00", loc)
	thisMonthLastDay := thisMonthFirstDay.AddDate(0, 1, -1)
	thisMonthLastDay, _ = time.ParseInLocation("2006-01-02 15:04:05", thisMonthLastDay.Format("2006-01-02")+" 23:59:59", loc)
	result := db.Preload("Patient").
		Where("organization_id = ? AND non_fixed_charge_date BETWEEN ? AND ?", organizationId, thisMonthFirstDay, thisMonthLastDay).Find(&thisMonthNonFixedChargeRecords)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// 上個月
	lastMonthFirstDay := time.Now().AddDate(0, -1, -time.Now().Day()+1)
	lastMonthFirstDay, _ = time.ParseInLocation("2006-01-02 15:04:05", lastMonthFirstDay.Format("2006-01-02")+" 00:00:00", loc)
	lastMonthLasttDay := lastMonthFirstDay.AddDate(0, 1, -1)
	lastMonthLasttDay, _ = time.ParseInLocation("2006-01-02 15:04:05", lastMonthLasttDay.Format("2006-01-02")+" 23:59:59", loc)
	result = db.Preload("Patient").
		Where("organization_id = ? AND non_fixed_charge_date BETWEEN ? AND ?", organizationId, lastMonthFirstDay, lastMonthLasttDay).
		Find(&lastMonthNonFixedChargeRecords)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	result = db.Preload("Patient").Where("organization_id = ?", organizationId).Select("DISTINCT ON (patient_id) *").Order("patient_id, updated_at desc").Find(&nonFixedChargeRecordsDescByUpdatedAt)

	return thisMonthNonFixedChargeRecords, lastMonthNonFixedChargeRecords, nonFixedChargeRecordsDescByUpdatedAt, result.Error
}
