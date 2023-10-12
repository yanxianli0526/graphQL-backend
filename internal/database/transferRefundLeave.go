package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTransferRefundLeaveById(db *gorm.DB, organizationId, transferRefundLeaveId uuid.UUID) (*models.TransferRefundLeave, error) {
	var transferRefundLeave *models.TransferRefundLeave

	result := db.Where("id = ? AND organization_id = ?", transferRefundLeaveId, organizationId).First(&transferRefundLeave)

	return transferRefundLeave, result.Error
}

func GetTransferRefundLeavesByPatientIdAndDate(db *gorm.DB, organizationId, patientId uuid.UUID, transferRefundLeaveStartDate, transferRefundLeaveEndDate time.Time) ([]*models.TransferRefundLeave, error) {
	var transferRefundLeaves []*models.TransferRefundLeave

	result := db.Preload("User").
		Where("organization_id = ? AND patient_id = ? AND (start_date between ? AND ? OR end_date between ? AND ? OR (start_date < ? AND  end_date > ?))",
			organizationId, patientId, transferRefundLeaveStartDate, transferRefundLeaveEndDate, transferRefundLeaveStartDate, transferRefundLeaveEndDate, transferRefundLeaveStartDate, transferRefundLeaveEndDate).
		Find(&transferRefundLeaves)

	return transferRefundLeaves, result.Error
}

func GetTransferRefundLeavesByPatientIdBetweenEndDate(db *gorm.DB, organizationId, patientId uuid.UUID, transferRefundLeaveStartDate, transferRefundLeaveEndDate time.Time) ([]*models.TransferRefundLeave, error) {
	var transferRefundLeaves []*models.TransferRefundLeave
	result := db.Where("organization_id = ? AND patient_id = ? AND end_date between ? AND ?", organizationId, patientId, transferRefundLeaveStartDate, transferRefundLeaveEndDate).Find(&transferRefundLeaves)
	return transferRefundLeaves, result.Error
}

func GetTransferRefundLeavesInPatientIdsBetweenEndDate(db *gorm.DB, organizationId uuid.UUID, patientsId []uuid.UUID, transferRefundLeaveStartDate, transferRefundLeaveEndDate time.Time) ([]*models.TransferRefundLeave, error) {
	var transferRefundLeaves []*models.TransferRefundLeave

	result := db.Preload("User").Order("start_date ASC").
		Where("organization_id = ? AND patient_id in ? AND (start_date between ? AND ? OR end_date between ? AND ?)",
			organizationId, patientsId, transferRefundLeaveStartDate, transferRefundLeaveEndDate, transferRefundLeaveStartDate, transferRefundLeaveEndDate).
		Find(&transferRefundLeaves)

	return transferRefundLeaves, result.Error
}

func CreateTransferRefundLeave(db *gorm.DB, transferRefundLeave *models.TransferRefundLeave) error {
	return db.Create(transferRefundLeave).Error
}

func UpdateTransferRefundLeave(db *gorm.DB, updateTransferRefundLeave *models.TransferRefundLeave) (*models.TransferRefundLeave, error) {

	transferRefundLeave := models.TransferRefundLeave{
		ID:             updateTransferRefundLeave.ID,
		OrganizationId: updateTransferRefundLeave.OrganizationId,
	}
	result := db.Model(&updateTransferRefundLeave).Clauses(clause.Returning{}).Updates(map[string]interface{}{
		"start_date":     updateTransferRefundLeave.StartDate,
		"end_date":       updateTransferRefundLeave.EndDate,
		"reason":         updateTransferRefundLeave.Reason,
		"is_reserve_bed": updateTransferRefundLeave.IsReserveBed,
		"note":           updateTransferRefundLeave.Note,
		"items":          updateTransferRefundLeave.Items,
		"user_id":        updateTransferRefundLeave.UserId,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	return &transferRefundLeave, result.Error
}

func UpdateTransferRefundLeavesReceiptStatus(db *gorm.DB, needUpdateTransferRefundLeavesId []uuid.UUID, status string, payRecordReceiptDate *time.Time) error {
	// 作廢
	if status == "invalid" {
		return db.Model(&models.TransferRefundLeave{}).Where("id in ? AND (receipt_status = ? OR receipt_status = ?)", needUpdateTransferRefundLeavesId, "issued", "cancelInvalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else if status == "cancelInvalid" {
		// 取消作廢
		return db.Model(&models.TransferRefundLeave{}).Where("id in ?  AND receipt_status = ?", needUpdateTransferRefundLeavesId, "invalid").Updates(map[string]interface{}{
			"receipt_status": status,
		}).Error
	} else {
		// 開帳
		return db.Model(&models.TransferRefundLeave{}).Where("id in ?", needUpdateTransferRefundLeavesId).Updates(map[string]interface{}{
			"receipt_status": status,
			"receipt_date":   *payRecordReceiptDate,
		}).Error
	}
}

func DeleteTransferRefundLeave(db *gorm.DB, transferRefundLeave *models.TransferRefundLeave) error {
	return db.Where("id = ?  AND organization_id = ?", transferRefundLeave.ID, transferRefundLeave.OrganizationId).Delete(&models.TransferRefundLeave{}).Error
}
