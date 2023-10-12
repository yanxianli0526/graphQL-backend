package resolvers

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"graphql-go-template/internal/gql/generated"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"
	"time"
)

type Resolver struct{}

func (r *autoTextFieldResolver) ID(ctx context.Context, obj *models.AutoTextField) (string, error) {
	panic("not implemented")
}

func (r *basicChargeResolver) ID(ctx context.Context, obj *models.BasicCharge) (string, error) {
	panic("not implemented")
}

func (r *basicChargeSettingResolver) ID(ctx context.Context, obj *models.BasicChargeSetting) (string, error) {
	panic("not implemented")
}

func (r *depositRecordResolver) ID(ctx context.Context, obj *models.DepositRecord) (string, error) {
	panic("not implemented")
}

func (r *fileResolver) ID(ctx context.Context, obj *models.File) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateAutoTextField(ctx context.Context, input gqlmodels.AutoTextFieldInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteAutoTextField(ctx context.Context, input gqlmodels.AutoTextFieldInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateBasicCharge(ctx context.Context, input *gqlmodels.BasicChargeInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateBasicCharge(ctx context.Context, input *gqlmodels.BasicChargeInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateDepositRecord(ctx context.Context, input gqlmodels.DepositRecordInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateDepositRecord(ctx context.Context, id string, input gqlmodels.DepositRecordUpdateInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) InvalidDepositRecord(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateFile(ctx context.Context, fileName string) (*gqlmodels.UploadFileResponse, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateNonFixedChargeRecord(ctx context.Context, patientID string, input []*gqlmodels.NonFixedChargeRecordInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateNonFixedChargeRecord(ctx context.Context, id string, input *gqlmodels.NonFixedChargeRecordInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteNonFixedChargeRecord(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganization(ctx context.Context, input *gqlmodels.OrganizationSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationPrivacy(ctx context.Context, input *gqlmodels.OrganizationPrivacyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationBillDateRangeSetting(ctx context.Context, input *gqlmodels.OrganizationBillDateRangeSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateOrganizationBasicChargeSetting(ctx context.Context, input *gqlmodels.OrganizationBasicChargeSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationBasicChargeSetting(ctx context.Context, id string, input *gqlmodels.OrganizationBasicChargeSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteOrganizationBasicChargeSetting(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateOrganizationNonFixedChargeSetting(ctx context.Context, input *gqlmodels.OrganizationNonFixedChargeSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationNonFixedChargeSetting(ctx context.Context, id string, input *gqlmodels.OrganizationNonFixedChargeSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteOrganizationNonFixedChargeSetting(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationReceipt(ctx context.Context, input *gqlmodels.OrganizationReceiptInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateOrganizationReceiptTemplateSetting(ctx context.Context, input gqlmodels.OrganizationReceiptTemplateSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateOrganizationReceiptTemplateSetting(ctx context.Context, id string, input gqlmodels.OrganizationReceiptTemplateSettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteOrganizationReceiptTemplateSetting(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePatientBill(ctx context.Context, input *gqlmodels.PatientBillInput) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePatientBills(ctx context.Context, input *gqlmodels.PatientBillsInput) ([]*models.PatientBill, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePatientBillNote(ctx context.Context, input *gqlmodels.UpdatePatientBillNoteInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePatientBillChargeDates(ctx context.Context, id string, input gqlmodels.UpdatePatientBillChargeDatesInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddPatientBillBasicCharge(ctx context.Context, input *gqlmodels.CreatePatientBillBasicChargeInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePatientBillBasicCharge(ctx context.Context, input *gqlmodels.UpdatePatientBillBasicChargeInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePatientBillBasicCharge(ctx context.Context, basicChargeID string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddPatientBillSubsidy(ctx context.Context, input *gqlmodels.CreatePatientBillSubsidyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePatientBillSubsidy(ctx context.Context, input *gqlmodels.UpdatePatientBillSubsidyInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePatientBillSubsidy(ctx context.Context, subsidyID string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePayRecords(ctx context.Context, input gqlmodels.PayRecordInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) InvalidPayRecord(ctx context.Context, id string, input gqlmodels.InvalidPayRecordInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CancelInvalidPayRecord(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePayRecorNote(ctx context.Context, id string, note string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreatePayRecordDetail(ctx context.Context, payRecrodID string, input gqlmodels.PayRecordDetailInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdatePayRecordDetail(ctx context.Context, id string, input gqlmodels.PayRecordDetailInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeletePayRecordDetail(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateSubsidiesSetting(ctx context.Context, patientID string, input []*gqlmodels.SubsidySettingInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateSubsidiesSetting(ctx context.Context, patientID string, input []*gqlmodels.SubsidySettingUpdateInput, needDeleteID []string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateTransferRefundLeave(ctx context.Context, patientID string, input gqlmodels.TransferRefundLeaveInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateTransferRefundLeave(ctx context.Context, id string, input gqlmodels.TransferRefundLeaveInput) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteTransferRefundLeave(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (r *nonFixedChargeRecordResolver) ID(ctx context.Context, obj *models.NonFixedChargeRecord) (string, error) {
	panic("not implemented")
}

func (r *organizationResolver) ID(ctx context.Context, obj *models.Organization) (string, error) {
	panic("not implemented")
}

func (r *organizationResolver) Privacy(ctx context.Context, obj *models.Organization) (gqlmodels.PrivacyType, error) {
	panic("not implemented")
}

func (r *organizationBasicChargeSettingResolver) ID(ctx context.Context, obj *models.OrganizationBasicChargeSetting) (string, error) {
	panic("not implemented")
}

func (r *organizationNonFixedChargeSettingResolver) ID(ctx context.Context, obj *models.OrganizationNonFixedChargeSetting) (string, error) {
	panic("not implemented")
}

func (r *organizationReceiptResolver) ID(ctx context.Context, obj *models.OrganizationReceipt) (string, error) {
	panic("not implemented")
}

func (r *organizationReceiptResolver) Year(ctx context.Context, obj *models.OrganizationReceipt) (*gqlmodels.YearType, error) {
	panic("not implemented")
}

func (r *organizationReceiptResolver) Month(ctx context.Context, obj *models.OrganizationReceipt) (*gqlmodels.MonthType, error) {
	panic("not implemented")
}

func (r *organizationReceiptTemplateSettingResolver) ID(ctx context.Context, obj *models.OrganizationReceiptTemplateSetting) (string, error) {
	panic("not implemented")
}

func (r *patientResolver) ID(ctx context.Context, obj *models.Patient) (string, error) {
	panic("not implemented")
}

func (r *patientBillResolver) ID(ctx context.Context, obj *models.PatientBill) (string, error) {
	panic("not implemented")
}

func (r *payRecordResolver) ID(ctx context.Context, obj *models.PayRecord) (string, error) {
	panic("not implemented")
}

func (r *payRecordDetailResolver) ID(ctx context.Context, obj *models.PayRecordDetail) (string, error) {
	panic("not implemented")
}

func (r *payRecordDetailResolver) Type(ctx context.Context, obj *models.PayRecordDetail) (gqlmodels.PayRecordDetailType, error) {
	panic("not implemented")
}

func (r *queryResolver) Patient(ctx context.Context, id string) (*models.Patient, error) {
	panic("not implemented")
}

func (r *queryResolver) Patients(ctx context.Context) ([]*models.Patient, error) {
	panic("not implemented")
}

func (r *queryResolver) AutoTextFields(ctx context.Context, field string) ([]*models.AutoTextField, error) {
	panic("not implemented")
}

func (r *queryResolver) BasicChargeSettings(ctx context.Context, patientID string) ([]*models.BasicChargeSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) DepositRecord(ctx context.Context, id string) (*models.DepositRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) DepositRecords(ctx context.Context, patientID string) ([]*models.DepositRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientLatestDepositRecords(ctx context.Context) (*gqlmodels.PatientLatestDepositRecords, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintDepositRecord(ctx context.Context, id string) (*string, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientLatestFixedChargeRecords(ctx context.Context) ([]*gqlmodels.FixedChargeRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) NonFixedChargeRecord(ctx context.Context, id string) (*models.NonFixedChargeRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) NonFixedChargeRecords(ctx context.Context, patientID string, startDate time.Time, endDate time.Time) ([]*models.NonFixedChargeRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientLatestNonFixedChargeRecords(ctx context.Context) (*gqlmodels.PatientLatestNonFixedChargeRecords, error) {
	panic("not implemented")
}

func (r *queryResolver) Organization(ctx context.Context) (*models.Organization, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationBasicChargeSetting(ctx context.Context, id string) (*models.OrganizationBasicChargeSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationBasicChargeSettings(ctx context.Context) ([]*models.OrganizationBasicChargeSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationNonFixedChargeSetting(ctx context.Context, id string) (*models.OrganizationNonFixedChargeSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationNonFixedChargeSettings(ctx context.Context) ([]*models.OrganizationNonFixedChargeSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationReceipt(ctx context.Context) (*models.OrganizationReceipt, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationReceiptTemplateSetting(ctx context.Context, id string) (*models.OrganizationReceiptTemplateSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) OrganizationReceiptTemplateSettings(ctx context.Context) ([]*models.OrganizationReceiptTemplateSetting, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientBill(ctx context.Context, patientID string, billYear int, billMonth int) (*models.PatientBill, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientBills(ctx context.Context, billDate time.Time) ([]*models.PatientBill, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPatientBill(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPatientBillGeneralTable(ctx context.Context, billDate time.Time) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientBillBasicCharge(ctx context.Context, basicChargeID string) (*models.BasicCharge, error) {
	panic("not implemented")
}

func (r *queryResolver) PatientBillSubsidy(ctx context.Context, subsidyID string) (*models.Subsidy, error) {
	panic("not implemented")
}

func (r *queryResolver) PayRecords(ctx context.Context, payDate time.Time) ([]*models.PayRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) PayRecord(ctx context.Context, id string) (*models.PayRecord, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPayRecordDetail(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPayRecordPart(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPayRecordPartByTaxType(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPayRecordsPartByTaxType(ctx context.Context, ids []string) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PrintPayRecordGeneralTable(ctx context.Context, billDate time.Time) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) PayRecordDetail(ctx context.Context, id string) (*models.PayRecordDetail, error) {
	panic("not implemented")
}

func (r *queryResolver) SubsidiesSetting(ctx context.Context, patientID string) ([]*models.SubsidySetting, error) {
	panic("not implemented")
}

func (r *queryResolver) SubsidySetting(ctx context.Context, id string) (*models.SubsidySetting, error) {
	panic("not implemented")
}

func (r *queryResolver) TransferRefundLeave(ctx context.Context, id string) (*models.TransferRefundLeave, error) {
	panic("not implemented")
}

func (r *queryResolver) TransferRefundLeaves(ctx context.Context, patientID string, startDate time.Time, endDate time.Time) ([]*models.TransferRefundLeave, error) {
	panic("not implemented")
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}

func (r *subsidyResolver) ID(ctx context.Context, obj *models.Subsidy) (string, error) {
	panic("not implemented")
}

func (r *subsidySettingResolver) ID(ctx context.Context, obj *models.SubsidySetting) (string, error) {
	panic("not implemented")
}

func (r *transferRefundLeaveResolver) ID(ctx context.Context, obj *models.TransferRefundLeave) (string, error) {
	panic("not implemented")
}

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	panic("not implemented")
}

func (r *userResolver) Preference(ctx context.Context, obj *models.User) (*gqlmodels.UserPreference, error) {
	panic("not implemented")
}

// AutoTextField returns generated.AutoTextFieldResolver implementation.
func (r *Resolver) AutoTextField() generated.AutoTextFieldResolver { return &autoTextFieldResolver{r} }

// BasicCharge returns generated.BasicChargeResolver implementation.
func (r *Resolver) BasicCharge() generated.BasicChargeResolver { return &basicChargeResolver{r} }

// BasicChargeSetting returns generated.BasicChargeSettingResolver implementation.
func (r *Resolver) BasicChargeSetting() generated.BasicChargeSettingResolver {
	return &basicChargeSettingResolver{r}
}

// DepositRecord returns generated.DepositRecordResolver implementation.
func (r *Resolver) DepositRecord() generated.DepositRecordResolver { return &depositRecordResolver{r} }

// File returns generated.FileResolver implementation.
func (r *Resolver) File() generated.FileResolver { return &fileResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// NonFixedChargeRecord returns generated.NonFixedChargeRecordResolver implementation.
func (r *Resolver) NonFixedChargeRecord() generated.NonFixedChargeRecordResolver {
	return &nonFixedChargeRecordResolver{r}
}

// Organization returns generated.OrganizationResolver implementation.
func (r *Resolver) Organization() generated.OrganizationResolver { return &organizationResolver{r} }

// OrganizationBasicChargeSetting returns generated.OrganizationBasicChargeSettingResolver implementation.
func (r *Resolver) OrganizationBasicChargeSetting() generated.OrganizationBasicChargeSettingResolver {
	return &organizationBasicChargeSettingResolver{r}
}

// OrganizationNonFixedChargeSetting returns generated.OrganizationNonFixedChargeSettingResolver implementation.
func (r *Resolver) OrganizationNonFixedChargeSetting() generated.OrganizationNonFixedChargeSettingResolver {
	return &organizationNonFixedChargeSettingResolver{r}
}

// OrganizationReceipt returns generated.OrganizationReceiptResolver implementation.
func (r *Resolver) OrganizationReceipt() generated.OrganizationReceiptResolver {
	return &organizationReceiptResolver{r}
}

// OrganizationReceiptTemplateSetting returns generated.OrganizationReceiptTemplateSettingResolver implementation.
func (r *Resolver) OrganizationReceiptTemplateSetting() generated.OrganizationReceiptTemplateSettingResolver {
	return &organizationReceiptTemplateSettingResolver{r}
}

// Patient returns generated.PatientResolver implementation.
func (r *Resolver) Patient() generated.PatientResolver { return &patientResolver{r} }

// PatientBill returns generated.PatientBillResolver implementation.
func (r *Resolver) PatientBill() generated.PatientBillResolver { return &patientBillResolver{r} }

// PayRecord returns generated.PayRecordResolver implementation.
func (r *Resolver) PayRecord() generated.PayRecordResolver { return &payRecordResolver{r} }

// PayRecordDetail returns generated.PayRecordDetailResolver implementation.
func (r *Resolver) PayRecordDetail() generated.PayRecordDetailResolver {
	return &payRecordDetailResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subsidy returns generated.SubsidyResolver implementation.
func (r *Resolver) Subsidy() generated.SubsidyResolver { return &subsidyResolver{r} }

// SubsidySetting returns generated.SubsidySettingResolver implementation.
func (r *Resolver) SubsidySetting() generated.SubsidySettingResolver {
	return &subsidySettingResolver{r}
}

// TransferRefundLeave returns generated.TransferRefundLeaveResolver implementation.
func (r *Resolver) TransferRefundLeave() generated.TransferRefundLeaveResolver {
	return &transferRefundLeaveResolver{r}
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type autoTextFieldResolver struct{ *Resolver }
type basicChargeResolver struct{ *Resolver }
type basicChargeSettingResolver struct{ *Resolver }
type depositRecordResolver struct{ *Resolver }
type fileResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type nonFixedChargeRecordResolver struct{ *Resolver }
type organizationResolver struct{ *Resolver }
type organizationBasicChargeSettingResolver struct{ *Resolver }
type organizationNonFixedChargeSettingResolver struct{ *Resolver }
type organizationReceiptResolver struct{ *Resolver }
type organizationReceiptTemplateSettingResolver struct{ *Resolver }
type patientResolver struct{ *Resolver }
type patientBillResolver struct{ *Resolver }
type payRecordResolver struct{ *Resolver }
type payRecordDetailResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subsidyResolver struct{ *Resolver }
type subsidySettingResolver struct{ *Resolver }
type transferRefundLeaveResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
