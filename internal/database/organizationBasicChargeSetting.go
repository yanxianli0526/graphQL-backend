package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOrganizationBasicChargeSettingById(db *gorm.DB, organizationBasicChargeSettingId uuid.UUID) (*models.OrganizationBasicChargeSetting, error) {
	var organizationBasicChargeSetting models.OrganizationBasicChargeSetting
	err := db.Where("id = ?", organizationBasicChargeSettingId).First(&organizationBasicChargeSetting).Error
	if err != nil {
		return nil, err
	}
	return &organizationBasicChargeSetting, nil
}

func GetOrganizationBasicChargeSettings(db *gorm.DB, organizationId uuid.UUID) ([]*models.OrganizationBasicChargeSetting, error) {
	var organizationBasicChargeSettings []*models.OrganizationBasicChargeSetting
	err := db.Order("created_at ASC").Where("organization_id = ? ", organizationId).Find(&organizationBasicChargeSettings).Error
	if err != nil {
		return nil, err
	}
	return organizationBasicChargeSettings, nil
}

func CreateOrganizationBasicChargeSetting(db *gorm.DB, organizationBasicChargeSetting *models.OrganizationBasicChargeSetting) error {
	return db.Create(organizationBasicChargeSetting).Error
}

func UpdateOrganizationBasicChargeSettingById(db *gorm.DB, inputOrganizationBasicChargeSetting *models.OrganizationBasicChargeSetting) error {

	err := db.Model(&models.OrganizationBasicChargeSetting{
		ID: inputOrganizationBasicChargeSetting.ID,
	}).Updates(map[string]interface{}{
		"item_name": inputOrganizationBasicChargeSetting.ItemName,
		"type":      inputOrganizationBasicChargeSetting.Type,
		"unit":      inputOrganizationBasicChargeSetting.Unit,
		"price":     inputOrganizationBasicChargeSetting.Price,
		"tax_type":  inputOrganizationBasicChargeSetting.TaxType,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteOrganizationBasicChargeSettingAndBasicChargeSetting(db *gorm.DB, organizationBasicChargeSettingId uuid.UUID) error {
	result := db.Delete(&models.OrganizationBasicChargeSetting{ID: organizationBasicChargeSettingId})
	if result.Error != nil {
		return result.Error
	}

	return db.Where("organization_basic_charge_setting_id = ?", organizationBasicChargeSettingId).Delete(&models.BasicChargeSetting{}).Error
}
