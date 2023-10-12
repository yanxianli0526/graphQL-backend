package services

import (
	"graphql-go-template/internal/models"
	nis "graphql-go-template/pkg/nis/nisApi"
	"time"

	"github.com/google/uuid"
)

func (svc *Service) GetOrganizationByProviderId(orgId string) (*uuid.UUID, bool) {
	organizationId, isExit := svc.db.GetOrganizationByProviderId(orgId)

	return organizationId, isExit
}

func (svc *Service) FirstSyncOrganization(nisOrganization *nis.OrganizationInfo, orgId string) (*uuid.UUID, error) {

	organization := models.Organization{
		ID:                       [16]byte{},
		Name:                     nisOrganization.Name,
		Address:                  &nisOrganization.Address,
		Phone:                    &nisOrganization.Tel,
		Fax:                      &nisOrganization.Fax,
		Owner:                    &nisOrganization.Owner,
		Email:                    &nisOrganization.Email,
		TaxIdNumber:              &nisOrganization.TaxIdNumber,
		RemittanceIdNumber:       &nisOrganization.RemittanceIdNumber,
		EstablishmentNumber:      &nisOrganization.EstablishmentNumber,
		Solution:                 nisOrganization.Solution,
		FixedChargeStartMonth:    "thisMonth",
		FixedChargeStartDate:     1,
		FixedChargeEndMonth:      "thisMonth",
		FixedChargeEndDate:       31,
		NonFixedChargeStartMonth: "lastMonth",
		NonFixedChargeStartDate:  1,
		NonFixedChargeEndMonth:   "lastMonth",
		NonFixedChargeEndDate:    31,
		TransferRefundStartMonth: "lastMonth",
		TransferRefundStartDate:  1,
		TransferRefundEndMonth:   "lastMonth",
		TransferRefundEndDate:    31,
		Branchs:                  nisOrganization.Branch,
		ProviderOrgId:            orgId,
		TestTime:                 time.Time{},
		Privacy:                  "unmask",
	}

	organizationId, err := svc.db.FirstSyncOrganization(organization)
	if err != nil {
		return nil, err
	}
	return organizationId, nil
}
