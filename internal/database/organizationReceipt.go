package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *GormDatabase) FirstSyncOrganizationReceipt(organizationReceipt models.OrganizationReceipt) error {
	result := d.DB.Create(&organizationReceipt)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *GormDatabase) GetOrganizationReceiptById(organizationId uuid.UUID) int64 {
	var count int64
	d.DB.Model(&models.OrganizationReceipt{}).Where("organization_id = ?", organizationId).Count(&count)

	return count
}

func GetOrganizationReceiptById(db *gorm.DB, organizationId uuid.UUID) (*models.OrganizationReceipt, error) {
	organizationReceipt := &models.OrganizationReceipt{}
	result := db.Where("organization_id = ?", organizationId).First(&organizationReceipt)
	if result.RowsAffected < 1 {
		return nil, result.Error
	}

	return organizationReceipt, nil
}

func UpdateOrganizationReceiptById(db *gorm.DB, inputOrganizationReceipt *models.OrganizationReceipt) error {
	err := db.Model(&models.OrganizationReceipt{}).Where("organization_id = ?", inputOrganizationReceipt.OrganizationId).Updates(map[string]interface{}{
		"first_text":             inputOrganizationReceipt.FirstText,
		"year":                   inputOrganizationReceipt.Year,
		"year_text":              inputOrganizationReceipt.YearText,
		"month":                  inputOrganizationReceipt.Month,
		"month_text":             inputOrganizationReceipt.MonthText,
		"last_text":              inputOrganizationReceipt.LastText,
		"is_reset_in_next_cycle": inputOrganizationReceipt.IsResetInNextCycle,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
