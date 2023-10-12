package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSubsidiesSetting(db *gorm.DB, subsidiesSetting []*models.SubsidySetting) error {
	return db.Create(subsidiesSetting).Error
}

func UpdateSubsidiesSetting(db *gorm.DB, createSubsidiesSetting, updateSubsidiesSetting []*models.SubsidySetting, needDeleteSubsidiesSetting []uuid.UUID, patientId uuid.UUID) error {

	var err error
	// 長度皆為0等同於全部刪除
	if len(createSubsidiesSetting) == 0 && len(updateSubsidiesSetting) == 0 {
		err = db.Where("patient_id = ?", patientId).Delete(&models.SubsidySetting{}).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		// 先更新有id的
		for i := range updateSubsidiesSetting {
			err = db.Model(&models.SubsidySetting{
				ID:        updateSubsidiesSetting[i].ID,
				PatientId: patientId,
			}).Updates(map[string]interface{}{
				"item_name":  updateSubsidiesSetting[i].ItemName,
				"type":       updateSubsidiesSetting[i].Type,
				"price":      updateSubsidiesSetting[i].Price,
				"unit":       updateSubsidiesSetting[i].Unit,
				"id_number":  updateSubsidiesSetting[i].IdNumber,
				"note":       updateSubsidiesSetting[i].Note,
				"sort_index": updateSubsidiesSetting[i].SortIndex,
			}).Error
			if err != nil {
				return err
			}
		}
		if len(createSubsidiesSetting) > 0 {
			// 再做新增
			err = db.Create(createSubsidiesSetting).Error
			if err != nil {
				return err
			}
		}
		return db.Where("patient_id = ? AND id in ?", patientId, needDeleteSubsidiesSetting).Delete(&models.SubsidySetting{}).Error
	}
}

func GetSubsidiesSetting(db *gorm.DB, organizationId uuid.UUID) ([]*models.SubsidySetting, error) {
	var subsidiesSetting []*models.SubsidySetting

	result := db.Preload("Patient").Where("organization_id = ?", organizationId).Find(&subsidiesSetting)

	return subsidiesSetting, result.Error

}

func GetSubsidiesSettingByPatientId(db *gorm.DB, organizationId, patientId uuid.UUID) ([]*models.SubsidySetting, error) {
	var subsidiesSetting []*models.SubsidySetting

	result := db.Order("sort_index asc").Where("organization_id = ? AND patient_id = ?", organizationId, patientId).Find(&subsidiesSetting)

	return subsidiesSetting, result.Error
}

func GetSubsidiesSettingInPatientIds(db *gorm.DB, organizationId uuid.UUID, patientsId []uuid.UUID) ([]*models.SubsidySetting, error) {
	var subsidiesSetting []*models.SubsidySetting

	result := db.Where("organization_id = ? AND patient_id in ?", organizationId, patientsId).Find(&subsidiesSetting)

	return subsidiesSetting, result.Error
}

func GetSubsidySetting(db *gorm.DB, organizationId, subsidySettingId uuid.UUID) (*models.SubsidySetting, error) {
	var subsidySetting *models.SubsidySetting

	result := db.Where("organization_id = ? AND id = ?", organizationId, subsidySettingId).Find(&subsidySetting)

	return subsidySetting, result.Error

}
