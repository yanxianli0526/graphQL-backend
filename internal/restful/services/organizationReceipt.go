package services

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
)

func (svc *Service) GetOrganizationReceiptById(orgId *uuid.UUID) int64 {

	count := svc.db.GetOrganizationReceiptById(*orgId)

	return count
}

func (svc *Service) FirstSyncOrganizationReceipt(orgId *uuid.UUID) error {

	organizationReceipt := models.OrganizationReceipt{
		ID:                 uuid.New(),
		FirstText:          "",
		Year:               "Christian",
		YearText:           "",
		Month:              "MM",
		MonthText:          "",
		LastText:           "",
		IsResetInNextCycle: false,
		OrganizationId:     *orgId,
	}

	err := svc.db.FirstSyncOrganizationReceipt(organizationReceipt)
	if err != nil {
		return err
	}
	return nil
}
