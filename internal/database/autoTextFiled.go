package orm

import (
	"graphql-go-template/internal/models"

	"gorm.io/gorm"
)

func CreateAutoTextField(db *gorm.DB, autoTextField *models.AutoTextField) error {
	return db.Create(autoTextField).Error
}

func DeleteAutoTextField(db *gorm.DB, deleteAutoTextField *models.AutoTextField) error {
	return db.Where("module_name = ? AND item_name = ? AND text = ? AND organization_id = ?",
		deleteAutoTextField.ModuleName, deleteAutoTextField.ItemName, deleteAutoTextField.Text, deleteAutoTextField.OrganizationId).Delete(&models.AutoTextField{}).Error

}

func GetAutoTextFields(db *gorm.DB, autoTextField *models.AutoTextField) ([]*models.AutoTextField, error) {
	var autoTextFields []*models.AutoTextField
	err := db.Where("module_name = ? AND item_name = ? AND organization_id = ?",
		autoTextField.ModuleName, autoTextField.ItemName, autoTextField.OrganizationId).Find(&autoTextFields).Error
	if err != nil {
		return nil, err
	}
	return autoTextFields, nil
}
