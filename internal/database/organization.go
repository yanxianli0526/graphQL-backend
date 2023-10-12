package orm

import (
	"encoding/json"
	"errors"
	"graphql-go-template/internal/models"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetOrganizationById(db *gorm.DB, organizationId uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	err := db.Where("id = ?", organizationId).First(&organization).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func UpdateOrganzationById(db *gorm.DB, inputOrganization *models.Organization) error {

	err := db.Model(&models.Organization{
		ID: inputOrganization.ID,
	}).Updates(map[string]interface{}{
		"name":                 inputOrganization.Name,
		"address_city":         inputOrganization.AddressCity,
		"address_district":     inputOrganization.AddressDistrict,
		"address":              inputOrganization.Address,
		"phone":                inputOrganization.Phone,
		"fax":                  inputOrganization.Fax,
		"owner":                inputOrganization.Owner,
		"email":                inputOrganization.Email,
		"tax_id_number":        inputOrganization.TaxIdNumber,
		"remittance_bank":      inputOrganization.RemittanceBank,
		"remittance_id_number": inputOrganization.RemittanceIdNumber,
		"remittance_user_name": inputOrganization.RemittanceUserName,
		"establishment_number": inputOrganization.EstablishmentNumber,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationPrivacy(db *gorm.DB, inputOrganization *models.Organization) error {
	err := db.Model(&models.Organization{
		ID: inputOrganization.ID,
	}).Updates(map[string]interface{}{
		"privacy": inputOrganization.Privacy,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationBillDateRangeSetting(db *gorm.DB, inputOrganization *models.Organization) error {

	err := db.Model(&models.Organization{
		ID: inputOrganization.ID,
	}).Updates(map[string]interface{}{
		"fixed_charge_start_month":     inputOrganization.FixedChargeStartMonth,
		"fixed_charge_start_date":      inputOrganization.FixedChargeStartDate,
		"fixed_charge_end_month":       inputOrganization.FixedChargeEndMonth,
		"fixed_charge_end_date":        inputOrganization.FixedChargeEndDate,
		"non_fixed_charge_start_month": inputOrganization.NonFixedChargeStartMonth,
		"non_fixed_charge_start_date":  inputOrganization.NonFixedChargeStartDate,
		"non_fixed_charge_end_month":   inputOrganization.NonFixedChargeEndMonth,
		"non_fixed_charge_end_date":    inputOrganization.NonFixedChargeEndDate,
		"transfer_refund_start_month":  inputOrganization.TransferRefundStartMonth,
		"transfer_refund_start_date":   inputOrganization.TransferRefundStartDate,
		"transfer_refund_end_month":    inputOrganization.TransferRefundEndMonth,
		"transfer_refund_end_date":     inputOrganization.TransferRefundEndDate,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (d *GormDatabase) GetOrganizationByProviderId(providerId string) (*uuid.UUID, bool) {
	organization := &models.Organization{}
	result := d.DB.Where("provider_org_id = ?", providerId).First(&organization)
	if result.RowsAffected < 1 {
		return nil, false
	}

	return &organization.ID, true
}

func (d *GormDatabase) FirstSyncOrganization(organization models.Organization) (*uuid.UUID, error) {
	result := d.DB.Create(&organization)
	if result.Error != nil {
		return nil, result.Error
	}

	return &organization.ID, result.Error
}

func MarshalStringArray(a pq.StringArray) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		data, _ := json.Marshal(a)
		io.WriteString(w, string(data))
	})
}

func UnmarshalStringArray(v interface{}) (pq.StringArray, error) {
	a, ok := v.(pq.StringArray)
	if !ok {
		return nil, errors.New("failed to cast to pq.StringArray")
	}
	return a, nil
}
