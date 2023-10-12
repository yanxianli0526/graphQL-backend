package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePayRecordDetail(db *gorm.DB, payRecordDetail *models.PayRecordDetail) error {
	return db.Create(payRecordDetail).Error
}

func UpdatePayRecordDetailUser(db *gorm.DB, payRecordDetail *models.PayRecordDetail, organizationId uuid.UUID) error {
	return db.Model(&models.PayRecordDetail{}).Where("id = ? and organization_id = ?", payRecordDetail.ID, organizationId).Updates(map[string]interface{}{
		"user_id": payRecordDetail.UserId,
	}).Error
}

func UpdatePayRecordDetail(db *gorm.DB, payRecordDetail *models.PayRecordDetail, organizationId uuid.UUID) error {
	return db.Model(&models.PayRecordDetail{}).Where("id = ? and organization_id = ?", payRecordDetail.ID, organizationId).Updates(map[string]interface{}{
		"record_date": payRecordDetail.RecordDate,
		"type":        payRecordDetail.Type,
		"price":       payRecordDetail.Price,
		"method":      payRecordDetail.Method,
		"payer":       payRecordDetail.Payer,
		"handler":     payRecordDetail.Handler,
		"note":        payRecordDetail.Note,
		"user_id":     payRecordDetail.UserId,
	}).Error
}

func DeletePayRecordDetail(db *gorm.DB, payRecordDetailId, organizationId uuid.UUID) error {
	return db.Where("id = ? AND organization_id = ?", payRecordDetailId, organizationId).Delete(&models.PayRecordDetail{}).Error

}

func GetPayRecordDetail(db *gorm.DB, payRecrodDetailId, organizationId uuid.UUID) (*models.PayRecordDetail, error) {
	var payRecordDetail models.PayRecordDetail
	err := db.Preload("User").Where("id = ? AND organization_id = ?", payRecrodDetailId, organizationId).First(&payRecordDetail).Error
	if err != nil {
		return nil, err
	}
	return &payRecordDetail, nil
}
