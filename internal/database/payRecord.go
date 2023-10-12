package orm

import (
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePayRecord(db *gorm.DB, payRecords *models.PayRecord) error {
	return db.Create(payRecords).Error
}

func CreatePayRecords(db *gorm.DB, payRecords []*models.PayRecord) error {
	return db.Create(payRecords).Error
}

func InvalidPayRecord(db *gorm.DB, payRecordId uuid.UUID, input gqlmodels.InvalidPayRecordInput, userId uuid.UUID) error {

	return db.Model(&models.PayRecord{}).Where("id = ? and is_invalid = ?", payRecordId, false).Updates(map[string]interface{}{
		"is_invalid":      true,
		"invalid_date":    input.InvalidDate,
		"invalid_caption": input.InvalidCaption,
		"invalid_user_id": userId,
	}).Error
}

func CancelInvalidPayRecord(db *gorm.DB, payRecordId, userId uuid.UUID) error {
	return db.Model(&models.PayRecord{}).Where("id = ? and is_invalid = ?", payRecordId, true).Updates(map[string]interface{}{
		"is_invalid":      false,
		"invalid_date":    nil,
		"invalid_caption": "",
		"invalid_user_id": userId,
	}).Error
}

func UpdatePayRecordNote(db *gorm.DB, note string, payRecordId, userId uuid.UUID) error {
	return db.Model(&models.PayRecord{}).Where("id = ? and is_invalid = ?", payRecordId, false).Updates(map[string]interface{}{
		"note":    note,
		"user_id": userId,
	}).Error
}

func UpdatePayRecordPaidAmount(db *gorm.DB, paidAmount int, payRecordId uuid.UUID) error {
	return db.Model(&models.PayRecord{}).Where("id = ? and is_invalid = ?", payRecordId, false).Updates(map[string]interface{}{
		"paid_amount": paidAmount,
	}).Error
}

func GetPayRecordCount(db *gorm.DB, organizationId uuid.UUID, payYear, payMonth int) (int64, error) {
	var count int64
	result := db.Model(&models.PayRecord{}).Where("organization_id = ? AND pay_year = ? AND pay_month = ?", organizationId, payYear, payMonth).Count(&count)

	return count, result.Error
}

func GetPayRecordCountByOrganizationId(db *gorm.DB, organizationId uuid.UUID) (int64, error) {
	var count int64
	result := db.Model(&models.PayRecord{}).Where("organization_id = ? ", organizationId).Count(&count)

	return count, result.Error
}

func GetLastestPayRecord(db *gorm.DB, organizationId uuid.UUID, payYear, payMonth int) (*models.PayRecord, error) {
	var payRecords []*models.PayRecord

	result := db.Order("receipt_number desc").Where("organization_id = ? AND pay_year = ? AND pay_month = ?", organizationId, payYear, payMonth).Find(&payRecords)

	return payRecords[0], result.Error
}

func GetLastestPayRecordByOrganizationId(db *gorm.DB, organizationId uuid.UUID) (*models.PayRecord, error) {

	var payRecord *models.PayRecord

	result := db.Order("receipt_number desc").Where("organization_id = ?", organizationId).First(&payRecord)

	return payRecord, result.Error

}

func GetPayRecordById(db *gorm.DB, payRecordId uuid.UUID, preloadOrganization, preloadUser, preloadPayRecordDetails bool) (*models.PayRecord, error) {
	var payRecord *models.PayRecord

	tx := db.Table("pay_records")

	if preloadOrganization {
		tx = tx.Preload("Organization")
	}
	if preloadUser {
		tx = tx.Preload("User").Preload("CreatedUser").Preload("InvalidUser")
	}
	if preloadPayRecordDetails {
		tx = tx.Preload("PayRecordDetails.User")
	}

	result := tx.Preload("Patient").Where("id = ?", payRecordId).First(&payRecord)

	return payRecord, result.Error
}

func GetPayRecordsForPrint(db *gorm.DB, payRecordsId []uuid.UUID, preloadOrganization, preloadUser bool) ([]*models.PayRecord, error) {
	var payRecords []*models.PayRecord

	tx := db.Table("pay_records")

	if preloadOrganization {
		tx = tx.Preload("Organization")
	}
	if preloadUser {
		tx = tx.Preload("User")
	}

	result := tx.Preload("Patient").Order("receipt_number ASC").Where("id in ?", payRecordsId).Find(&payRecords)

	return payRecords, result.Error
}

func GetPayRecords(db *gorm.DB, organizationId uuid.UUID, payYear, payMonth int) ([]*models.PayRecord, error) {
	var payRecords []*models.PayRecord

	result := db.Order("receipt_number desc").Preload("Patient").Where("organization_id = ? AND pay_year = ? AND pay_month = ?", organizationId, payYear, payMonth).Find(&payRecords)

	return payRecords, result.Error
}

func GetPayRecordsByInvalidStatus(db *gorm.DB, organizationId uuid.UUID, payYear, payMonth int, isInvalid bool) ([]*models.PayRecord, error) {
	var payRecords []*models.PayRecord

	result := db.Order("receipt_number asc").Preload("PayRecordDetails").Preload("Patient").Where("organization_id = ? AND pay_year = ? AND pay_month = ? and is_invalid = ?", organizationId, payYear, payMonth, isInvalid).Find(&payRecords)

	return payRecords, result.Error
}

func AppendAssociationsPayRecordPayRecordDetail(db *gorm.DB, payRecord *models.PayRecord, payRecordDetail models.PayRecordDetail) error {
	return db.Model(&payRecord).Association("PayRecordDetails").Append(&payRecordDetail)
}
