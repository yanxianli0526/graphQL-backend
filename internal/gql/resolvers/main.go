package resolvers

import (
	"context"
	"fmt"

	config "graphql-go-template/envconfig"
	orm "graphql-go-template/internal/database"
	objectStorage "graphql-go-template/pkg/gcp"

	"graphql-go-template/internal/gql/generated"

	"go.uber.org/zap"

	"github.com/99designs/gqlgen/graphql"
)

type contextKey string

var (
	UserIdCtxKey   = contextKey("userId")
	OrganizationId = contextKey("orgId")
)

type Resolver struct {
	ORM                  *orm.GormDatabase
	store                objectStorage.ObjectStorage
	ClientID             string
	ClientSecret         string
	PhotoServiceProtocol string
	PhotoServiceHost     string
	PhotoServicePort     string
	OriginalBucket       string
	ProcsssedBucket      string
	Logger               *zap.Logger
}

func NewRootResolvers(orm *orm.GormDatabase, store objectStorage.ObjectStorage, config config.AuthEndpoint, photoService config.PhotoService, logger *zap.Logger) generated.Config {
	c := generated.Config{
		Resolvers: &Resolver{
			ORM:                  orm, // pass in the ORM instance in the resolvers to be used
			store:                store,
			ClientID:             config.ClientID,
			ClientSecret:         config.ClientSecret,
			PhotoServiceProtocol: photoService.PhotoServiceProtocol,
			PhotoServiceHost:     photoService.PhotoServiceHost,
			PhotoServicePort:     photoService.PhotoServicePort,
			OriginalBucket:       photoService.OriginalBucket,
			ProcsssedBucket:      photoService.ProcsssedBucket,
			Logger:               logger,
		},
	}

	// Schema Directive
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		ctxUserId := ctx.Value(UserIdCtxKey)
		if ctxUserId == nil {
			return nil, fmt.Errorf("You are not authorized to perform this action")
		}
		return next(ctx)
	}

	return c
}

// func (r *Resolver) Test() generated.TestResolver {
// 	return &testResolver{r}
// }

func (r *Resolver) Organization() generated.OrganizationResolver {
	return &organizationResolver{r}
}

func (r *Resolver) AutoTextField() generated.AutoTextFieldResolver {
	return &autoTextFieldResolver{r}
}

func (r *Resolver) File() generated.FileResolver {
	return &fileResolver{r}
}

func (r *Resolver) Patient() generated.PatientResolver {
	return &patientResolver{r}
}

func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}

func (r *Resolver) OrganizationBasicChargeSetting() generated.OrganizationBasicChargeSettingResolver {
	return &organizationBasicChargeSettingResolver{r}
}

func (r *Resolver) OrganizationReceiptTemplateSetting() generated.OrganizationReceiptTemplateSettingResolver {
	return &organizationReceiptTemplateSettingResolver{r}
}

func (r *Resolver) OrganizationReceipt() generated.OrganizationReceiptResolver {
	return &organizationReceiptResolver{r}
}

func (r *Resolver) DepositRecord() generated.DepositRecordResolver {
	return &depositRecordResolver{r}
}

func (r *Resolver) OrganizationNonFixedChargeSetting() generated.OrganizationNonFixedChargeSettingResolver {
	return &organizationNonFixedChargeSettingResolver{r}
}

func (r *Resolver) BasicChargeSetting() generated.BasicChargeSettingResolver {
	return &basicChargeSettingResolver{r}
}

func (r *Resolver) NonFixedChargeRecord() generated.NonFixedChargeRecordResolver {
	return &nonFixedChargeRecordResolver{r}
}

func (r *Resolver) TransferRefundLeave() generated.TransferRefundLeaveResolver {
	return &transferRefundLeaveResolver{r}
}

func (r *Resolver) Subsidy() generated.SubsidyResolver {
	return &subsidyResolver{r}
}

func (r *Resolver) SubsidySetting() generated.SubsidySettingResolver {
	return &subsidySettingResolver{r}
}

func (r *Resolver) BasicCharge() generated.BasicChargeResolver {
	return &basicChargeResolver{r}
}

func (r *Resolver) PatientBill() generated.PatientBillResolver {
	return &patientBillResolver{r}
}

func (r *Resolver) PayRecord() generated.PayRecordResolver {
	return &payRecordResolver{r}
}

func (r *Resolver) PayRecordDetail() generated.PayRecordDetailResolver {
	return &payRecordDetailResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
