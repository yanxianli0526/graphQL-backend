package services

import (
	"graphql-go-template/internal/models"
	nis "graphql-go-template/pkg/nis/nisApi"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func (svc *Service) GetUserById(userId uuid.UUID) (*models.User, error) {
	user, err := svc.db.GetUserById(userId)

	return user, err
}

func (svc *Service) UserLogout(user *models.User) error {
	// tokenExpiredAt := time.Now()
	// err := svc.db.UserLogout(user, tokenExpiredAt)

	return nil
}

func (svc *Service) GetUserByProviderId(providerId string) (*models.User, bool) {
	user, isExit := svc.db.GetUserByProviderId(providerId)

	return user, isExit
}

func (svc *Service) FirstSyncUser(nisUser *nis.UserInfo, orgId *uuid.UUID, tokenBytes []byte, tokenExpiredAt time.Time) (*models.User, error) {
	user := models.User{
		FirstName:      nisUser.FirstName,
		LastName:       nisUser.LastName,
		DisplayName:    nisUser.DisplayName,
		IdNumber:       nisUser.IdNumber,
		TokenExpiredAt: &tokenExpiredAt,
		ProviderToken:  datatypes.JSON(tokenBytes),
		ProviderId:     nisUser.ProviderId,
		OrganizationId: *orgId,
	}

	userData, err := svc.db.FirstSyncUser(user)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (svc *Service) UpdateUserToken(UserId uuid.UUID, tokenBytes []byte, tokenExpiredAt time.Time) error {

	err := svc.db.UpdateUserToken(UserId, tokenBytes, tokenExpiredAt)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) SyncPeopleInCharges(nisPeopleInCharges []*nis.UserInfo, orgId *uuid.UUID, tokenBytes []byte) ([]models.User, error) {
	peopleInCharges := []models.User{}

	for i := range nisPeopleInCharges {
		user := models.User{
			FirstName:      nisPeopleInCharges[i].FirstName,
			LastName:       nisPeopleInCharges[i].LastName,
			DisplayName:    nisPeopleInCharges[i].DisplayName,
			IdNumber:       nisPeopleInCharges[i].IdNumber,
			ProviderId:     nisPeopleInCharges[i].ProviderId,
			OrganizationId: *orgId,
		}
		peopleInCharges = append(peopleInCharges, user)
	}

	peopleInChargesData, err := svc.db.SyncPeopleInCharges(peopleInCharges)
	if err != nil {
		return nil, err
	}
	return peopleInChargesData, nil
}
