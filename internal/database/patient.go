package orm

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPatientById(db *gorm.DB, patientId, organizationId uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	err := db.Where("id = ? AND organization_id = ?", patientId, organizationId).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func GetPatients(db *gorm.DB, organizationId uuid.UUID, preloadUser bool) ([]*models.Patient, error) {
	var patients []*models.Patient
	tx := db.Table("patients")
	if preloadUser {
		tx.Preload("Users")
	}
	err := tx.Where("organization_id = ?", organizationId).Find(&patients).Error
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (d *GormDatabase) SyncPatients(patients []models.Patient) error {

	return d.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "provider_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "status", "branch", "room", "bed", "sex", "id_number",
			"birthday", "check_in_date", "patient_number", "record_number", "numbering", "photo_url", "photo_x_position", "photo_y_position"}),
	}).Create(patients).Error
}
