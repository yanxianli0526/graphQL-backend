package services

import (
	"graphql-go-template/internal/models"

	"github.com/google/uuid"
)

func (svc *Service) GetOrganizationReceiptTemplateSettingById(orgId *uuid.UUID) int64 {

	count := svc.db.GetOrganizationReceiptTemplateSettingById(*orgId)

	return count
}

func (svc *Service) FirstSyncOrganizationReceiptTemplateSetting(orgId *uuid.UUID) error {
	organizationReceiptTemplateSettings := []models.OrganizationReceiptTemplateSetting{}

	organizationReceiptTemplateSettingTamplateOne := models.OrganizationReceiptTemplateSetting{
		ID:                  uuid.New(),
		Name:                "月費收據",
		TaxTypes:            []string{"allTax", "stampTax"},
		TitleName:           "收據",
		PatientInfo:         []string{"areaAndClass", "bedAndRoom", "checkInDate", "idNumber"},
		PriceShowType:       "classAddUp",
		OrganizationInfoOne: []string{"taxIdNumber", "phone"},
		OrganizationInfoTwo: []string{"address", "phone"},
		NoteText:            "",
		SealOneName:         "機構",
		SealOnePicture:      "",
		SealTwoName:         "負責人",
		SealTwoPicture:      "",
		SealThreeName:       "經辦人",
		SealThreePicture:    "",
		SealFourName:        "",
		SealFourPicture:     "",
		PartOneName:         "第一聯 住民收執聯",
		PartTwoName:         "第二聯 機構留存聯",
		OrganizationId:      *orgId,
	}
	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, organizationReceiptTemplateSettingTamplateOne)

	organizationReceiptTemplateSettingTamplateTwo := models.OrganizationReceiptTemplateSetting{
		ID:                  uuid.New(),
		Name:                "雜支收據",
		TaxTypes:            []string{"businessTax", "noTax", "other"},
		TitleName:           "收據",
		PatientInfo:         []string{"areaAndClass", "bedAndRoom"},
		PriceShowType:       "classAddUp",
		OrganizationInfoOne: []string{},
		OrganizationInfoTwo: []string{},
		NoteText:            "",
		SealOneName:         "",
		SealOnePicture:      "",
		SealTwoName:         "",
		SealTwoPicture:      "",
		SealThreeName:       "經辦人",
		SealThreePicture:    "",
		SealFourName:        "",
		SealFourPicture:     "",
		PartOneName:         "第一聯 住民收執聯",
		PartTwoName:         "第二聯 機構留存聯",
		OrganizationId:      *orgId,
	}

	organizationReceiptTemplateSettings = append(organizationReceiptTemplateSettings, organizationReceiptTemplateSettingTamplateTwo)

	err := svc.db.FirstSyncOrganizationReceiptTemplateSetting(organizationReceiptTemplateSettings)
	if err != nil {
		return err
	}
	return nil
}
