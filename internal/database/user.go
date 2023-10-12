package orm

import (
	"graphql-go-template/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserById(db *gorm.DB, userId uuid.UUID) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserLogout(db *gorm.DB, vUser *models.User, expiredTime time.Time) error {
	return db.Model(vUser).Update("token_expired_at", expiredTime).Error
}

func GetUsers(db *gorm.DB, organizationId uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	err := db.Where("organization_id = ?", organizationId).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *GormDatabase) GetUserById(userId uuid.UUID) (*models.User, error) {
	var user models.User
	err := d.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *GormDatabase) GetUserByProviderId(providerId string) (*models.User, bool) {
	user := &models.User{}
	result := d.DB.Where("provider_id = ?", providerId).First(&user)
	if result.RowsAffected < 1 {
		return nil, false
	}

	return user, true

}

func (d *GormDatabase) FirstSyncUser(user models.User) (*models.User, error) {
	result := d.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (d *GormDatabase) UpdateUserToken(UserId uuid.UUID, tokenBytes []byte, tokenExpiredAt time.Time) error {
	return d.DB.Model(&models.User{}).Where("id = ?", UserId).Updates(map[string]interface{}{
		"provider_token":   tokenBytes,
		"token_expired_at": tokenExpiredAt,
	}).Error
}

func (d *GormDatabase) SyncPeopleInCharges(peopleInCharges []models.User) ([]models.User, error) {

	result := d.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: "provider_id",
		}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name"}),
	}).Create(peopleInCharges)
	if result.Error != nil {
		return nil, result.Error
	}
	return peopleInCharges, result.Error
}
