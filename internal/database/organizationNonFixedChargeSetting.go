package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOrganizationNonFixedChargeSettingById(db *gorm.DB, organizationNonFixedChargeSettingId uuid.UUID) (*models.OrganizationNonFixedChargeSetting, error) {
	var organizationNonFixedChargeSetting models.OrganizationNonFixedChargeSetting
	err := db.Where("id = ?", organizationNonFixedChargeSettingId).First(&organizationNonFixedChargeSetting).Error
	if err != nil {
		return nil, err
	}
	return &organizationNonFixedChargeSetting, nil
}

func GetOrganizationNonFixedChargeSettings(db *gorm.DB, organizationId uuid.UUID) ([]*models.OrganizationNonFixedChargeSetting, error) {
	var organizationNonFixedChargeSettings []*models.OrganizationNonFixedChargeSetting
	err := db.Where("organization_id = ?", organizationId).Order("item_category,created_at ASC").Find(&organizationNonFixedChargeSettings).Error
	if err != nil {
		return nil, err
	}
	return organizationNonFixedChargeSettings, nil
}

func CreateOrganizationNonFixedChargeSetting(db *gorm.DB, organizationNonFixedChargeSetting *models.OrganizationNonFixedChargeSetting) error {
	return db.Create(organizationNonFixedChargeSetting).Error
}

func UpdateOrganizationNonFixedChargeSettingById(db *gorm.DB, inputOrganizationNonFixedChargeSetting *models.OrganizationNonFixedChargeSetting) error {

	err := db.Model(&models.OrganizationNonFixedChargeSetting{
		ID: inputOrganizationNonFixedChargeSetting.ID,
	}).Updates(map[string]interface{}{
		"item_category": inputOrganizationNonFixedChargeSetting.ItemCategory,
		"item_name":     inputOrganizationNonFixedChargeSetting.ItemName,
		"type":          inputOrganizationNonFixedChargeSetting.Type,
		"unit":          inputOrganizationNonFixedChargeSetting.Unit,
		"price":         inputOrganizationNonFixedChargeSetting.Price,
		"tax_type":      inputOrganizationNonFixedChargeSetting.TaxType,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteOrganizationNonFixedChargeSettingAndFixedChargeSetting(db *gorm.DB, organizationNonFixedChargeSettingId uuid.UUID) error {
	result := db.Delete(&models.OrganizationNonFixedChargeSetting{ID: organizationNonFixedChargeSettingId})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
