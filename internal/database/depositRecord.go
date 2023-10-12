package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetDepositRecordLastIdNumber(db *gorm.DB, organizationId uuid.UUID, inputDate time.Time) (string, error) {
	var depositRecord models.DepositRecord

	loc, _ := time.LoadLocation("Asia/Taipei")
	date := inputDate.Format("2006-01-02")
	startDate := date + " 00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate, loc)
	endDate := date + " 23:59:59"
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate, loc)

	result := db.Where("organization_id = ? AND date BETWEEN ? AND ?", organizationId, startTime, endTime).Last(&depositRecord)

	if result.RowsAffected == 0 {
		return "no record", nil
	}

	if result.Error != nil {
		return "", result.Error
	}
	return depositRecord.IdNumber, nil
}

func GetDepositRecordById(db *gorm.DB, depositRecordId, organizationId uuid.UUID, preloadUser, preloadPatient bool) (*models.DepositRecord, error) {
	if preloadUser {
		db = db.Preload("User")
	}
	if preloadPatient {
		db = db.Preload("Patient")
	}

	var depositRecord models.DepositRecord
	err := db.Where("id = ? AND organization_id = ?", depositRecordId, organizationId).First(&depositRecord).Error
	if err != nil {
		return nil, err
	}
	return &depositRecord, nil
}

func GetDepositRecordForPrintById(db *gorm.DB, depositRecordId, organizationId uuid.UUID) (*models.DepositRecord, error) {
	var depositRecord models.DepositRecord
	err := db.Preload("Organization").Preload("Patient").Preload("User").Where("id = ? AND organization_id = ?", depositRecordId, organizationId).First(&depositRecord).Error
	if err != nil {
		return nil, err
	}
	return &depositRecord, nil
}

func GetDepositRecords(db *gorm.DB, organizationId, patientId uuid.UUID, preloadUser, preloadPatient bool) ([]*models.DepositRecord, error) {
	if preloadUser {
		db = db.Preload("User")
	}
	if preloadPatient {
		db = db.Preload("Patient")
	}

	var depositRecords []*models.DepositRecord
	err := db.Where("organization_id = ? AND patient_id = ?", organizationId, patientId).Find(&depositRecords).Error
	if err != nil {
		return nil, err
	}

	return depositRecords, nil
}

func PatientLatestDepositRecords(db *gorm.DB, organizationId uuid.UUID) ([]*models.DepositRecord, []*models.DepositRecord, error) {
	var depositRecordsDescByDate []*models.DepositRecord
	var depositRecordsDescByUpdatedAt []*models.DepositRecord

	result := db.Preload("Patient").Where("organization_id = ?", organizationId).Select("DISTINCT ON (patient_id) *").Order("patient_id, date desc").Find(&depositRecordsDescByDate)
	_ = db.Preload("Patient").Where("organization_id = ?", organizationId).Select("DISTINCT ON (patient_id) *").Order("patient_id, updated_at desc").Find(&depositRecordsDescByUpdatedAt)

	return depositRecordsDescByDate, depositRecordsDescByUpdatedAt, result.Error
}

func CreateDepositRecord(db *gorm.DB, depositRecord *models.DepositRecord) (*uuid.UUID, error) {
	result := db.Create(depositRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return &depositRecord.ID, nil
}

func UpdateDepositRecordById(db *gorm.DB, depositRecordUpdateInput *models.DepositRecord) error {

	err := db.Model(&models.DepositRecord{
		ID: depositRecordUpdateInput.ID,
	}).Updates(map[string]interface{}{
		"note":    depositRecordUpdateInput.Note,
		"user_id": depositRecordUpdateInput.UserId,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func UpdateInvalidDepositRecordById(db *gorm.DB, depositRecordId uuid.UUID) error {
	err := db.Model(&models.DepositRecord{
		ID: depositRecordId,
	}).Update("invalid", true).Error

	if err != nil {
		return err
	}
	return nil
}
