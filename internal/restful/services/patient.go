package services

import (
	"graphql-go-template/internal/models"
	nis "graphql-go-template/pkg/nis/nisApi"

	"github.com/google/uuid"
)

func (svc *Service) SyncPatients(nisPatient []*nis.Patient, peopleInCharges []models.User, orgId *uuid.UUID) error {
	patients := []models.Patient{}

	// 這邊的住民跟peopleInCharge是多對多的關係 要mapping一下 相符的id
	for i := range nisPatient {
		users := []*models.User{}
		for j := range nisPatient[i].PeopleInCharge {
			for k := range peopleInCharges {
				if peopleInCharges[k].ProviderId == nisPatient[i].PeopleInCharge[j] {
					user := models.User{
						ID:             peopleInCharges[k].ID,
						FirstName:      peopleInCharges[k].FirstName,
						LastName:       peopleInCharges[k].LastName,
						DisplayName:    peopleInCharges[k].DisplayName,
						IdNumber:       peopleInCharges[k].IdNumber,
						ProviderId:     peopleInCharges[k].ProviderId,
						OrganizationId: *orgId,
					}

					users = append(users, &user)
				}
			}
		}

		patient := models.Patient{
			FirstName:      nisPatient[i].FirstName,
			LastName:       nisPatient[i].LastName,
			IdNumber:       nisPatient[i].IdNumber,
			PhotoUrl:       nisPatient[i].PhotoUrl,
			PhotoXPosition: nisPatient[i].PhotoXPosition,
			PhotoYPosition: nisPatient[i].PhotoYPosition,
			ProviderId:     nisPatient[i].ProviderId,
			Status:         nisPatient[i].Status,
			Branch:         nisPatient[i].Branch,
			Room:           nisPatient[i].Room,
			Bed:            nisPatient[i].Bed,
			Sex:            nisPatient[i].Sex,
			Birthday:       nisPatient[i].Birthday,
			CheckInDate:    nisPatient[i].CheckInDate,
			PatientNumber:  nisPatient[i].PatientNumber,
			RecordNumber:   nisPatient[i].RecordNumber,
			Numbering:      nisPatient[i].Numbering,
			OrganizationId: *orgId,
			Users:          users,
		}
		patients = append(patients, patient)
	}

	err := svc.db.SyncPatients(patients)
	if err != nil {
		return err
	}
	return nil
}
