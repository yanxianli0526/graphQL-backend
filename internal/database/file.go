package orm

import (
	"graphql-go-template/internal/models"

	"gorm.io/gorm"
)

func CreateFile(db *gorm.DB, file *models.File) error {
	return db.Create(file).Error
}

// func GetAutoTextFields(db *gorm.DB, autoTextField *models.AutoTextField) ([]*models.AutoTextField, error) {
// 	var autoTextFields []*models.AutoTextField
// 	err := db.Where("module_name = ? AND item_name = ? AND organization_id = ?",
// 		autoTextField.ModuleName, autoTextField.ItemName, autoTextField.OrganizationId).Find(&autoTextFields).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return autoTextFields, nil
// }
