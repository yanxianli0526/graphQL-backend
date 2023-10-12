package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateOrganizationReceiptTemplateSetting(db *gorm.DB, organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting) error {
	return db.Create(organizationReceiptTemplateSetting).Error
}

func UpdateOrganizationReceiptTemplateSetting(db *gorm.DB, organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting) error {
	return db.Model(&models.OrganizationReceiptTemplateSetting{}).Where("id = ? AND organization_id = ?", organizationReceiptTemplateSetting.ID, organizationReceiptTemplateSetting.OrganizationId).Updates(map[string]interface{}{
		"name":                  organizationReceiptTemplateSetting.Name,
		"tax_types":             organizationReceiptTemplateSetting.TaxTypes,
		"title_name":            organizationReceiptTemplateSetting.TitleName,
		"patient_info":          organizationReceiptTemplateSetting.PatientInfo,
		"price_show_type":       organizationReceiptTemplateSetting.PriceShowType,
		"organization_info_one": organizationReceiptTemplateSetting.OrganizationInfoOne,
		"organization_info_two": organizationReceiptTemplateSetting.OrganizationInfoTwo,
		"note_text":             organizationReceiptTemplateSetting.NoteText,
		"seal_one_name":         organizationReceiptTemplateSetting.SealOneName,
		"seal_one_picture":      organizationReceiptTemplateSetting.SealOnePicture,
		"seal_two_name":         organizationReceiptTemplateSetting.SealTwoName,
		"seal_two_picture":      organizationReceiptTemplateSetting.SealTwoPicture,
		"seal_three_name":       organizationReceiptTemplateSetting.SealThreeName,
		"seal_three_picture":    organizationReceiptTemplateSetting.SealThreePicture,
		"seal_four_name":        organizationReceiptTemplateSetting.SealFourName,
		"seal_four_picture":     organizationReceiptTemplateSetting.SealFourPicture,
		"part_one_name":         organizationReceiptTemplateSetting.PartOneName,
		"part_two_name":         organizationReceiptTemplateSetting.PartTwoName,
		"organization_picture":  organizationReceiptTemplateSetting.OrganizationPicture,
	}).Error
}

func UpdateOrganizationReceiptTemplateSettingTaxTypes(db *gorm.DB, organizationReceiptTemplateSetting *models.OrganizationReceiptTemplateSetting) error {
	return db.Model(&models.OrganizationReceiptTemplateSetting{}).Where("id = ? AND organization_id = ?", organizationReceiptTemplateSetting.ID, organizationReceiptTemplateSetting.OrganizationId).Updates(map[string]interface{}{
		"tax_types": organizationReceiptTemplateSetting.TaxTypes,
	}).Error
}

func DeleteOrganizationReceiptTemplateSettingById(db *gorm.DB, receiptTemplateSettingId, organizationId uuid.UUID) error {
	result := db.Where("id = ? AND organization_id = ?", receiptTemplateSettingId, organizationId).Delete(&models.OrganizationReceiptTemplateSetting{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetOrganizationReceiptTemplateSettingById(db *gorm.DB, receiptTemplateSettingId, organizationId uuid.UUID) (*models.OrganizationReceiptTemplateSetting, error) {
	organizationReceiptTemplateSetting := &models.OrganizationReceiptTemplateSetting{}
	result := db.Where("id = ? AND organization_id = ?", receiptTemplateSettingId, organizationId).First(&organizationReceiptTemplateSetting)
	if result.RowsAffected < 1 {
		return nil, result.Error
	}

	return organizationReceiptTemplateSetting, nil
}

func GetOrganizationReceiptTemplateSettingInTaxType(db *gorm.DB, organizationId uuid.UUID, taxType string) (*models.OrganizationReceiptTemplateSetting, error) {
	organizationReceiptTemplateSetting := &models.OrganizationReceiptTemplateSetting{}
	result := db.Order("updated_at asc").Where("organization_id = ? AND ? = ANY (tax_types)", organizationId, taxType).First(&organizationReceiptTemplateSetting)
	if result.RowsAffected < 1 {
		return nil, result.Error
	}

	return organizationReceiptTemplateSetting, nil
}

func GetOrganizationReceiptTemplateSettings(db *gorm.DB, organizationId uuid.UUID) ([]*models.OrganizationReceiptTemplateSetting, error) {
	organizationReceiptTemplateSettings := []*models.OrganizationReceiptTemplateSetting{}
	result := db.Order("created_at asc").Where("organization_id = ?", organizationId).Find(&organizationReceiptTemplateSettings)
	if result.RowsAffected < 1 {
		return nil, result.Error
	}

	return organizationReceiptTemplateSettings, nil
}

func (d *GormDatabase) GetOrganizationReceiptTemplateSettingById(organizationId uuid.UUID) int64 {
	var count int64
	d.DB.Model(&models.OrganizationReceiptTemplateSetting{}).Where("organization_id = ?", organizationId).Count(&count)

	return count
}

func (d *GormDatabase) FirstSyncOrganizationReceiptTemplateSetting(organizationReceiptTemplateSettings []models.OrganizationReceiptTemplateSetting) error {
	result := d.DB.Create(&organizationReceiptTemplateSettings)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
