package orm

import (
	"fmt"
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBasicChargeSettings(db *gorm.DB, basicChargeSettings []*models.BasicChargeSetting) error {
	return db.Create(basicChargeSettings).Error
}

func UpdateBasicChargeSetting(db *gorm.DB, basicChargeSettings []*models.BasicChargeSetting, patientId uuid.UUID, BasicChargeSettingsUUID []uuid.UUID) error {

	// 長度為0等同於全部刪除
	if len(basicChargeSettings) == 0 {
		err := db.Where("patient_id = ?", patientId).Delete(&models.BasicChargeSetting{}).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		for i := range basicChargeSettings {
			var basicChargeSetting models.BasicChargeSetting
			result := db.Where("organization_basic_charge_setting_id = ? AND patient_id = ?", basicChargeSettings[i].OrganizationBasicChargeSettingId, patientId).Find(&basicChargeSetting)
			if result.RowsAffected == 0 {
				err := db.Create(basicChargeSettings[i]).Error
				if err != nil {
					return err
				}
			} else {
				db.Debug().Model(&models.BasicChargeSetting{}).Where("organization_basic_charge_setting_id = ? AND patient_id = ?", basicChargeSettings[i].OrganizationBasicChargeSettingId, patientId).Updates(map[string]interface{}{
					"sort_index": basicChargeSettings[i].SortIndex,
				})
			}
		}

		var needDeleteBasicChargeSetting models.BasicChargeSetting
		err := db.Not(map[string]interface{}{"organization_basic_charge_setting_id": BasicChargeSettingsUUID}).Where("patient_id = ?", patientId).Find(&needDeleteBasicChargeSetting).Error
		if err != nil {
			return err
		}

		err = db.Delete(&needDeleteBasicChargeSetting).Error
		if err != nil {
			return err
		}
		// return db.Create(BasicChargeSetting).Error
		return nil
	}
	return fmt.Errorf("UpdateBasicChargeSetting have error")
}

func GetBasicChargeSettingsByPatientId(db *gorm.DB, organizationId, patientId uuid.UUID) ([]*models.BasicChargeSetting, error) {
	var basicChargeSetting []*models.BasicChargeSetting
	err := db.Order("sort_index asc").Preload("OrganizationBasicChargeSetting").Where("organization_id = ? AND patient_id = ?", organizationId, patientId).Find(&basicChargeSetting).Error
	if err != nil {
		return nil, err
	}
	return basicChargeSetting, nil
}

func GetBasicChargeSettings(db *gorm.DB, organizationId uuid.UUID) ([]*models.BasicChargeSetting, error) {
	var basicChargeSetting []*models.BasicChargeSetting
	err := db.Preload("Patient").Preload("OrganizationBasicChargeSetting").Preload("User").Where("organization_id = ?", organizationId).Find(&basicChargeSetting).Error
	if err != nil {
		return nil, err
	}
	return basicChargeSetting, nil
}

func GetBasicChargeSettingsInPatientIds(db *gorm.DB, organizationId uuid.UUID, patientsId []uuid.UUID) ([]*models.BasicChargeSetting, error) {
	var basicChargeSetting []*models.BasicChargeSetting
	err := db.Order("patient_id desc").Preload("OrganizationBasicChargeSetting").Preload("User").Where("patient_id in ? AND organization_id = ?", patientsId, organizationId).Find(&basicChargeSetting).Error
	if err != nil {
		return nil, err
	}
	return basicChargeSetting, nil
}

func PatientLatestBasicChargeSettings(db *gorm.DB, organizationId uuid.UUID) ([]*models.BasicChargeSetting, error) {
	var basicChargeSettingDescByUpdatedAt []*models.BasicChargeSetting

	result := db.Preload("User").Preload("Patient").Where("organization_id = ?", organizationId).Select("DISTINCT ON (patient_id) *").Order("patient_id, updated_at desc").Find(&basicChargeSettingDescByUpdatedAt)

	return basicChargeSettingDescByUpdatedAt, result.Error
}
