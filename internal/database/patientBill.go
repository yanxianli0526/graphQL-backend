package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPatientBillById(db *gorm.DB, organizationId, patientBillId uuid.UUID) (*models.PatientBill, error) {
	patientBill := models.PatientBill{}
	result := db.Preload("Patient").Preload("Organization").Preload("Subsidies", func(db *gorm.DB) *gorm.DB {
		return db.Order("subsidies.sort_index ASC")
	}).Preload("BasicCharges", func(db *gorm.DB) *gorm.DB {
		return db.Order("basic_charges.sort_index ASC")
	}).Preload("NonFixedChargeRecords", func(db *gorm.DB) *gorm.DB {
		return db.Order("non_fixed_charge_records.item_category,non_fixed_charge_records.non_fixed_charge_date ASC")
	}).Preload("TransferRefundLeaves", func(db *gorm.DB) *gorm.DB {
		return db.Order("transfer_refund_leaves.start_date ASC")
	}).Where("id = ? AND organization_id = ?", patientBillId, organizationId).First(&patientBill)
	if result.Error != nil {
		return nil, result.Error
	}
	return &patientBill, nil
}

func GetPatientBillByPatientIdAndYearMonthHaveData(db *gorm.DB, organizationId, patientBillId uuid.UUID, billYear, billMonth int) (*models.PatientBill, bool) {
	patientBill := models.PatientBill{}
	result := db.Preload("BasicCharges", func(db *gorm.DB) *gorm.DB {
		return db.Order("basic_charges.sort_index ASC")
	}).Preload("Subsidies", func(db *gorm.DB) *gorm.DB {
		return db.Order("subsidies.sort_index ASC")
	}).Preload("NonFixedChargeRecords").Preload("TransferRefundLeaves").
		Where("patient_id = ? AND organization_id = ? AND bill_year = ? AND bill_month = ?", patientBillId, organizationId, billYear, billMonth).First(&patientBill)
	if result.RowsAffected == 0 {
		return nil, false
	}
	return &patientBill, true
}

func GetPatientBillByPatientIdAndYearMonth(db *gorm.DB, organizationId, patientBillId uuid.UUID, billYear, billMonth int) (*models.PatientBill, error) {
	patientBill := models.PatientBill{}
	result := db.Preload("Organization").Preload("Patient").Preload("User").Preload("EditNoteUser").Preload("BasicCharges", func(db *gorm.DB) *gorm.DB {
		return db.Order("basic_charges.sort_index ASC")
	}).Preload("Subsidies", func(db *gorm.DB) *gorm.DB {
		return db.Order("subsidies.sort_index ASC")
	}).Preload("NonFixedChargeRecords", func(db *gorm.DB) *gorm.DB {
		return db.Order("non_fixed_charge_records.item_category,non_fixed_charge_date ASC")
	}).Preload("TransferRefundLeaves", func(db *gorm.DB) *gorm.DB {
		return db.Order("transfer_refund_leaves.start_date ASC")
	}).Where("patient_id = ? AND organization_id = ? AND bill_year = ? AND bill_month = ?", patientBillId, organizationId, billYear, billMonth).First(&patientBill)
	if result.Error != nil {
		return nil, result.Error
	}
	return &patientBill, nil
}

func GetPatientBillsByDate(db *gorm.DB, organizationId uuid.UUID, billYear, billMonth int) ([]*models.PatientBill, error) {
	var patientBills []*models.PatientBill
	err := db.Preload("Patient").Preload("BasicCharges").Preload("Subsidies").Preload("NonFixedChargeRecords").Preload("TransferRefundLeaves").Where("organization_id = ? AND bill_year = ? AND bill_month = ?", organizationId, billYear, billMonth).Find(&patientBills).Error
	if err != nil {
		return nil, err
	}
	return patientBills, nil
}

func GetPatientBills(db *gorm.DB, organizationId uuid.UUID) ([]*models.PatientBill, error) {
	var patientBills []*models.PatientBill
	err := db.Preload("Patient").Where("organization_id = ?", organizationId).Find(&patientBills).Error
	if err != nil {
		return nil, err
	}
	return patientBills, nil
}

func GetPatientBillByMonth(db *gorm.DB, patientId uuid.UUID, billYear, billMonth int, preloadCharge bool) *models.PatientBill {
	var patientBill models.PatientBill
	tx := db.Table("patient_bills")

	if preloadCharge {
		tx.Preload("BasicCharges").Preload("Subsidies").Preload("NonFixedChargeRecords").Preload("TransferRefundLeaves")
	}
	result := tx.Where("patient_id = ? AND bill_year = ? AND bill_month =  ?", patientId, billYear, billMonth).First(&patientBill)
	if result.RowsAffected == 0 {
		return nil
	}

	return &patientBill
}

func GetPatientBillsByNonFixedChargeDate(db *gorm.DB, patientId uuid.UUID, nonFixedChargeStartDate time.Time) []*models.PatientBill {
	var patientBills []*models.PatientBill
	result := db.Preload("NonFixedChargeRecords").Where("patient_id = ? AND non_fixed_charge_start_date <= ? AND non_fixed_charge_end_date >=  ?", patientId, nonFixedChargeStartDate, nonFixedChargeStartDate).Find(&patientBills)
	if result.RowsAffected == 0 {
		return nil
	}

	return patientBills
}

func GetPatientBillsByTransferRefundLeaveDate(db *gorm.DB, patientId uuid.UUID, transferRefundLeaveEndDate time.Time) []*models.PatientBill {
	var patientBills []*models.PatientBill
	result := db.Preload("TransferRefundLeaves").Where("patient_id = ? AND transfer_refund_start_date <= ? AND transfer_refund_end_date >=  ?", patientId, transferRefundLeaveEndDate, transferRefundLeaveEndDate).Find(&patientBills)
	if result.RowsAffected == 0 {
		return nil
	}

	return patientBills
}

func GetPatientBillsByMonth(db *gorm.DB, patientIds []uuid.UUID, billYear, billMonth int) []*models.PatientBill {
	var patientBills []*models.PatientBill

	result := db.Preload("BasicCharges").Preload("Subsidies").Where("patient_id in ? AND bill_year = ? AND bill_month =  ?", patientIds, billYear, billMonth).Find(&patientBills)
	if result.RowsAffected == 0 {
		return nil
	}

	return patientBills
}

func CreatePatientBill(db *gorm.DB, patientBill *models.PatientBill) (*uuid.UUID, error) {
	result := db.Create(patientBill)
	if result.Error != nil {
		return nil, result.Error
	}

	return &patientBill.ID, nil
}

func CreatePatientBills(db *gorm.DB, patientBills []*models.PatientBill) error {
	return db.Create(patientBills).Error
}

func AppendAssociationsPatientBillNonFixedChargeRecord(db *gorm.DB, patientBill *models.PatientBill, nonFixedChargeRecord models.NonFixedChargeRecord) error {
	return db.Model(&patientBill).Association("NonFixedChargeRecords").Append(&nonFixedChargeRecord)
}

func AppendAssociationsPatientBillNonFixedChargeRecords(db *gorm.DB, patientBill *models.PatientBill, nonFixedChargeRecords []*models.NonFixedChargeRecord) error {
	return db.Model(&patientBill).Association("NonFixedChargeRecords").Append(&nonFixedChargeRecords)
}

func DeleteAssociationsPatientBillNonFixedChargeRecord(db *gorm.DB, patientBill *models.PatientBill, nonFixedChargeRecords models.NonFixedChargeRecord) error {
	return db.Model(&patientBill).Association("NonFixedChargeRecords").Delete(&nonFixedChargeRecords)
}

func ClearAssociationsPatientBillNonFixedChargeRecords(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&patientBill).Association("NonFixedChargeRecords").Clear()
}

func AppendAssociationsPatientBillTransferRefundLeave(db *gorm.DB, patientBill *models.PatientBill, transferRefundLeave models.TransferRefundLeave) error {
	return db.Model(&patientBill).Association("TransferRefundLeaves").Append(&transferRefundLeave)
}

func AppendAssociationsPatientBillTransferRefundLeaves(db *gorm.DB, patientBill *models.PatientBill, transferRefundLeave []*models.TransferRefundLeave) error {
	return db.Model(&patientBill).Association("TransferRefundLeaves").Append(&transferRefundLeave)
}

func DeleteAssociationsPatientBillTransferRefundLeave(db *gorm.DB, patientBill *models.PatientBill, transferRefundLeave models.TransferRefundLeave) error {
	return db.Model(&patientBill).Association("TransferRefundLeaves").Delete(&transferRefundLeave)
}

func ClearAssociationsPatientBillTransferRefundLeaves(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&patientBill).Association("TransferRefundLeaves").Clear()
}

func UpdatePatientBillNote(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&models.PatientBill{
		ID: patientBill.ID,
	}).Updates(map[string]interface{}{
		"note":              patientBill.Note,
		"edit_note_date":    time.Now(),
		"edit_note_user_id": patientBill.UserId,
	}).Error
}

func UpdatePatientBillAmountReceived(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&models.PatientBill{
		ID: patientBill.ID,
	}).Updates(map[string]interface{}{
		"amount_received": patientBill.AmountReceived,
	}).Error
}

func UpdatePatientBillChargeDates(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&models.PatientBill{
		ID: patientBill.ID,
	}).Updates(map[string]interface{}{
		"transfer_refund_start_date":  patientBill.TransferRefundStartDate,
		"transfer_refund_end_date":    patientBill.TransferRefundEndDate,
		"non_fixed_charge_start_date": patientBill.NonFixedChargeStartDate,
		"non_fixed_charge_end_date":   patientBill.NonFixedChargeEndDate,
	}).Error
}

func GetPatientBillHaveBasicChargeId(db *gorm.DB, patientId, organization_id uuid.UUID, basicChargeId uuid.UUID) *models.PatientBill {
	var patientBills []*models.PatientBill
	result := db.Preload("BasicCharges", func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", basicChargeId)
	}).Where("patient_id = ? AND organization_id = ? ", patientId, organization_id).Find(&patientBills)
	if result.RowsAffected == 0 {
		return nil
	}
	var patientBill *models.PatientBill
	for i := range patientBills {
		if len(patientBills[i].BasicCharges) > 0 {
			patientBill = patientBills[i]
		}
	}

	return patientBill
}

func GetPatientBillHaveSubsidyId(db *gorm.DB, patientId, organization_id uuid.UUID, subsidyId uuid.UUID) *models.PatientBill {
	var patientBills []*models.PatientBill
	result := db.Preload("Subsidies", func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", subsidyId)
	}).Where("patient_id = ? AND organization_id = ? ", patientId, organization_id).Find(&patientBills)
	if result.RowsAffected == 0 {
		return nil
	}
	var patientBill *models.PatientBill
	for i := range patientBills {
		if len(patientBills[i].Subsidies) > 0 {
			patientBill = patientBills[i]
		}
	}

	return patientBill
}

func UpdatePatientBillAmountDue(db *gorm.DB, patientBill *models.PatientBill) error {
	return db.Model(&models.PatientBill{
		ID:             patientBill.ID,
		OrganizationId: patientBill.OrganizationId,
	}).Updates(map[string]interface{}{
		"amount_due": patientBill.AmountDue,
	}).Error
}

func DeletePatientBillById(db *gorm.DB, patientBillId uuid.UUID, needDeleteBasicChargesID, needDeleteSubsidiesID []uuid.UUID) error {
	// 把多對多的表也刪一下
	db.Exec("delete FROM patient_bill_basic_charges WHERE patient_bill_id = ?", patientBillId)
	db.Exec("delete FROM patient_bill_subsidies WHERE patient_bill_id = ?", patientBillId)
	db.Exec("delete FROM patient_bill_transfer_refund_hospitalizeds WHERE patient_bill_id = ?", patientBillId)
	db.Exec("delete FROM patient_bill_transfer_refund_leaves WHERE patient_bill_id = ?", patientBillId)
	db.Exec("delete FROM patient_bill_non_fixed_charge_records WHERE patient_bill_id = ?", patientBillId)
	// 把多對多得清完後 再把固定費用的資料清一清 (不然會有一堆垃圾)
	if len(needDeleteBasicChargesID) > 0 {
		db.Where("id in ?", needDeleteBasicChargesID).Delete(&models.BasicCharge{})
	}
	if len(needDeleteSubsidiesID) > 0 {
		db.Where("id in ?", needDeleteSubsidiesID).Delete(&models.Subsidy{})
	}

	return db.Delete(&models.PatientBill{ID: patientBillId}).Error
}
func DeletePatientBillsInId(db *gorm.DB, patientsBillId, needDeleteBasicChargesID, needDeleteSubsidiesID []uuid.UUID) error {
	// 把多對多的表也刪一下
	db.Exec("delete FROM patient_bill_basic_charges WHERE patient_bill_id in ?", patientsBillId)
	db.Exec("delete FROM patient_bill_subsidies WHERE patient_bill_id in ?", patientsBillId)
	db.Exec("delete FROM patient_bill_transfer_refund_hospitalizeds WHERE patient_bill_id in ?", patientsBillId)
	db.Exec("delete FROM patient_bill_transfer_refund_leaves WHERE patient_bill_id in ?", patientsBillId)
	db.Exec("delete FROM patient_bill_non_fixed_charge_records WHERE patient_bill_id in ?", patientsBillId)
	// 把多對多得清完後 再把固定費用的資料清一清 (不然會有一堆垃圾)
	if len(needDeleteBasicChargesID) > 0 {
		db.Where("id in ?", needDeleteBasicChargesID).Delete(&models.BasicCharge{})
	}
	if len(needDeleteSubsidiesID) > 0 {
		db.Where("id in ?", needDeleteSubsidiesID).Delete(&models.Subsidy{})
	}
	return db.Where("id in ?", patientsBillId).Delete(&models.PatientBill{}).Error
}
