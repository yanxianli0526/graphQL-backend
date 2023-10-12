package resolvers_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	config "graphql-go-template/envconfig"

	"gitlab.smart-aging.tech/devops/ms-go-kit/observability"
	"go.uber.org/zap"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/generated"
	"graphql-go-template/internal/gql/resolvers"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	env    config.EnvConfig
	logger *zap.Logger
)

func LoadEnv() {
	if err := config.Process(&env); err != nil {
		fmt.Println(err)
	}
}

func addContext(user models.User) client.Option {
	return func(bd *client.Request) {
		ctx := bd.HTTP.Context()
		ctx = context.WithValue(ctx, resolvers.UserIdCtxKey, user.ID)
		ctx = context.WithValue(ctx, resolvers.OrganizationId, user.OrganizationId)
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}

func TestMutationResolver_AddPatientBillBasicChargeTypeIsChargeToCharge(t *testing.T) {
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		logger, err := observability.SetupLogger(env.Env.Debug)
		if err != nil {
			fmt.Println("observability.SetupLogger")
		}

		LoadEnv()
		gorm, _ := orm.Factory(env.Database)
		resolversData := resolvers.Resolver{ORM: gorm, Logger: logger}
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolversData})))

		//  這邊是在塞值讓之後傳api進去的內容有userId跟organizationId 不然api會找不到這兩個值然後一直死掉
		userId, _ := uuid.Parse("c71f802b-fea7-4278-9062-1181937df0d9")
		organizationId, _ := uuid.Parse("a9994b72-c534-423f-b857-615133ffa248")
		user := models.User{
			ID:             userId,
			OrganizationId: organizationId,
		}

		//  先新增一個全新的住民帳單(統一都拿這個人測試 林可可)
		patientId := "317edd46-1b4e-418c-8395-38f4eff59a8c"

		var patientBillResponse struct {
			CreatePatientBill string
		}

		patientBillCreateRequest := `mutation {
			createPatientBill(
				input: {
					patientId: "` + patientId + `",
					billDate:"2023-06-06T14:46:54+08:00"
				}
			)
		}
		`
		c.MustPost(patientBillCreateRequest, &patientBillResponse, addContext(user))
		// 這邊拿剛剛新增完的住民帳單新增固定月費
		var addPatientBillBasicChargeResponse struct {
			AddPatientBillBasicCharge bool
		}
		patientBillBasicChargePrice := "1500"
		patientBillBasicChargeType := "charge"
		addPatientBillBasicChargeRequest := `
		mutation {
			addPatientBillBasicCharge(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + patientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillBasicChargeRequest, &addPatientBillBasicChargeResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillBasicChargePriceInt, err := strconv.Atoi(patientBillBasicChargePrice)
		if err != nil {
			fmt.Println("patientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, patientBillBasicChargePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillBasicChargeResponse struct {
			UpdatePatientBillBasicCharge bool
		}
		updatePatientBillBasicChargePrice := "1800"
		updatePatientBillBasicChargeType := "charge"
		updatePatientBillBasicChargeRequest := `
		mutation {
			updatePatientBillBasicCharge(
				input: {
					basicChargeId: "` + createdPatientBill.BasicCharges[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + updatePatientBillBasicChargePrice + `
					startDate: "2023-05-01T14:46:54+08:00"
					endDate: "2023-05-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillBasicChargeRequest, &updatePatientBillBasicChargeResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillBasicChargePriceInt, err := strconv.Atoi(updatePatientBillBasicChargePrice)
		if err != nil {
			fmt.Println("updatePatientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-patientBillBasicChargePriceInt+updatePatientBillBasicChargePriceInt)
		// 再來要搞刪除
		var deletePatientBillBasicChargeResponse struct {
			DeletePatientBillBasicCharge bool
		}
		deletePatientBillBasicChargeRequest := `
		mutation {
			deletePatientBillBasicCharge(
				basicChargeId: "` + updatedPatientBill.BasicCharges[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillBasicChargeRequest, &deletePatientBillBasicChargeResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updatedPatientBill.BasicCharges[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillBasicChargeTypeIsChargeToRefund(t *testing.T) {
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		logger, err := observability.SetupLogger(env.Env.Debug)
		if err != nil {
			fmt.Println("observability.SetupLogger")
		}

		LoadEnv()
		gorm, _ := orm.Factory(env.Database)
		resolversData := resolvers.Resolver{ORM: gorm, Logger: logger}
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolversData})))

		//  這邊是在塞值讓之後傳api進去的內容有userId跟organizationId 不然api會找不到這兩個值然後一直死掉
		userId, _ := uuid.Parse("c71f802b-fea7-4278-9062-1181937df0d9")
		organizationId, _ := uuid.Parse("a9994b72-c534-423f-b857-615133ffa248")
		user := models.User{
			ID:             userId,
			OrganizationId: organizationId,
		}

		//  先新增一個全新的住民帳單(統一都拿這個人測試 林可可)
		patientId := "317edd46-1b4e-418c-8395-38f4eff59a8c"

		var patientBillResponse struct {
			CreatePatientBill string
		}

		patientBillCreateRequest := `mutation {
			createPatientBill(
				input: {
					patientId: "` + patientId + `",
					billDate:"2023-06-06T14:46:54+08:00"
				}
			)
		}
		`
		c.MustPost(patientBillCreateRequest, &patientBillResponse, addContext(user))

		// 這邊拿剛剛新增完的住民帳單新增固定月費
		var addPatientBillBasicChargeResponse struct {
			AddPatientBillBasicCharge bool
		}
		patientBillBasicChargePrice := "1500"
		patientBillBasicChargeType := "charge"
		addPatientBillBasicChargeRequest := `
		mutation {
			addPatientBillBasicCharge(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + patientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillBasicChargeRequest, &addPatientBillBasicChargeResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillBasicChargePriceInt, err := strconv.Atoi(patientBillBasicChargePrice)
		if err != nil {
			fmt.Println("patientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, patientBillBasicChargePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillBasicChargeResponse struct {
			UpdatePatientBillBasicCharge bool
		}
		updatePatientBillBasicChargePrice := "1400"
		updatePatientBillBasicChargeType := "refund"
		updatePatientBillBasicChargeRequest := `
		mutation {
			updatePatientBillBasicCharge(
				input: {
					basicChargeId: "` + createdPatientBill.BasicCharges[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + updatePatientBillBasicChargePrice + `
					startDate: "2023-05-01T14:46:54+08:00"
					endDate: "2023-05-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillBasicChargeRequest, &updatePatientBillBasicChargeResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillBasicChargePriceInt, err := strconv.Atoi(updatePatientBillBasicChargePrice)
		if err != nil {
			fmt.Println("updatePatientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		fmt.Println("updatedPatientBill.AmountDue", updatedPatientBill.AmountDue)
		fmt.Println("createdPatientBill.AmountDue-patientBillBasicChargePriceInt-updatePatientBillBasicChargePriceInt", createdPatientBill.AmountDue-patientBillBasicChargePriceInt-updatePatientBillBasicChargePriceInt)

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-patientBillBasicChargePriceInt-updatePatientBillBasicChargePriceInt)
		// 再來要搞刪除
		var deletePatientBillBasicChargeResponse struct {
			DeletePatientBillBasicCharge bool
		}
		deletePatientBillBasicChargeRequest := `
		mutation {
			deletePatientBillBasicCharge(
				basicChargeId: "` + updatedPatientBill.BasicCharges[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillBasicChargeRequest, &deletePatientBillBasicChargeResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updatedPatientBill.BasicCharges[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillBasicChargeTypeIsRefundToCharge(t *testing.T) {
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		logger, err := observability.SetupLogger(env.Env.Debug)
		if err != nil {
			fmt.Println("observability.SetupLogger")
		}

		LoadEnv()
		gorm, _ := orm.Factory(env.Database)
		resolversData := resolvers.Resolver{ORM: gorm, Logger: logger}
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolversData})))

		//  這邊是在塞值讓之後傳api進去的內容有userId跟organizationId 不然api會找不到這兩個值然後一直死掉
		userId, _ := uuid.Parse("c71f802b-fea7-4278-9062-1181937df0d9")
		organizationId, _ := uuid.Parse("a9994b72-c534-423f-b857-615133ffa248")
		user := models.User{
			ID:             userId,
			OrganizationId: organizationId,
		}

		//  先新增一個全新的住民帳單(統一都拿這個人測試 林可可)
		patientId := "317edd46-1b4e-418c-8395-38f4eff59a8c"

		var patientBillResponse struct {
			CreatePatientBill string
		}

		patientBillCreateRequest := `mutation {
			createPatientBill(
				input: {
					patientId: "` + patientId + `",
					billDate:"2023-06-06T14:46:54+08:00"
				}
			)
		}
		`
		c.MustPost(patientBillCreateRequest, &patientBillResponse, addContext(user))

		// 這邊拿剛剛新增完的住民帳單新增固定月費
		var addPatientBillBasicChargeResponse struct {
			AddPatientBillBasicCharge bool
		}
		patientBillBasicChargePrice := "1500"
		patientBillBasicChargeType := "refund"
		addPatientBillBasicChargeRequest := `
		mutation {
			addPatientBillBasicCharge(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillBasicChargeType + `"
					taxType:"stampTax"
					unit: "月"
					price: ` + patientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillBasicChargeRequest, &addPatientBillBasicChargeResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillBasicChargePriceInt, err := strconv.Atoi(patientBillBasicChargePrice)
		if err != nil {
			fmt.Println("patientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -patientBillBasicChargePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillBasicChargeResponse struct {
			UpdatePatientBillBasicCharge bool
		}
		updatePatientBillBasicChargePrice := "1400"
		updatePatientBillBasicChargeType := "charge"
		updatePatientBillBasicChargeRequest := `
		mutation {
			updatePatientBillBasicCharge(
				input: {
					basicChargeId: "` + createdPatientBill.BasicCharges[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + updatePatientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillBasicChargeRequest, &updatePatientBillBasicChargeResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillBasicChargePriceInt, err := strconv.Atoi(updatePatientBillBasicChargePrice)
		if err != nil {
			fmt.Println("updatePatientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+patientBillBasicChargePriceInt+updatePatientBillBasicChargePriceInt)
		// 再來要搞刪除
		var deletePatientBillBasicChargeResponse struct {
			DeletePatientBillBasicCharge bool
		}
		deletePatientBillBasicChargeRequest := `
		mutation {
			deletePatientBillBasicCharge(
				basicChargeId: "` + updatedPatientBill.BasicCharges[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillBasicChargeRequest, &deletePatientBillBasicChargeResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updatedPatientBill.BasicCharges[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillBasicChargeTypeIsRefundToRefund(t *testing.T) {
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		logger, err := observability.SetupLogger(env.Env.Debug)
		if err != nil {
			fmt.Println("observability.SetupLogger")
		}

		LoadEnv()
		gorm, _ := orm.Factory(env.Database)
		resolversData := resolvers.Resolver{ORM: gorm, Logger: logger}
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolversData})))

		//  這邊是在塞值讓之後傳api進去的內容有userId跟organizationId 不然api會找不到這兩個值然後一直死掉
		userId, _ := uuid.Parse("c71f802b-fea7-4278-9062-1181937df0d9")
		organizationId, _ := uuid.Parse("a9994b72-c534-423f-b857-615133ffa248")
		user := models.User{
			ID:             userId,
			OrganizationId: organizationId,
		}

		//  先新增一個全新的住民帳單(統一都拿這個人測試 林可可)
		patientId := "317edd46-1b4e-418c-8395-38f4eff59a8c"

		var patientBillResponse struct {
			CreatePatientBill string
		}

		patientBillCreateRequest := `mutation {
			createPatientBill(
				input: {
					patientId: "` + patientId + `",
					billDate:"2023-06-06T14:46:54+08:00"
				}
			)
		}
		`
		c.MustPost(patientBillCreateRequest, &patientBillResponse, addContext(user))

		// 這邊拿剛剛新增完的住民帳單新增固定月費
		var addPatientBillBasicChargeResponse struct {
			AddPatientBillBasicCharge bool
		}
		patientBillBasicChargePrice := "1500"
		patientBillBasicChargeType := "refund"
		addPatientBillBasicChargeRequest := `
		mutation {
			addPatientBillBasicCharge(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillBasicChargeType + `"
					taxType:"stampTax"
					unit: "月"
					price: ` + patientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillBasicChargeRequest, &addPatientBillBasicChargeResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillBasicChargePriceInt, err := strconv.Atoi(patientBillBasicChargePrice)
		if err != nil {
			fmt.Println("patientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -patientBillBasicChargePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillBasicChargeResponse struct {
			UpdatePatientBillBasicCharge bool
		}
		updatePatientBillBasicChargePrice := "1400"
		updatePatientBillBasicChargeType := "refund"
		updatePatientBillBasicChargeRequest := `
		mutation {
			updatePatientBillBasicCharge(
				input: {
					basicChargeId: "` + createdPatientBill.BasicCharges[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillBasicChargeType + `"
					taxType: "stampTax"
					unit: "月"
					price: ` + updatePatientBillBasicChargePrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillBasicChargeRequest, &updatePatientBillBasicChargeResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillBasicChargePriceInt, err := strconv.Atoi(updatePatientBillBasicChargePrice)
		if err != nil {
			fmt.Println("updatePatientBillBasicChargePriceInt strconv.Atoi(patientBillBasicChargePrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+patientBillBasicChargePriceInt-updatePatientBillBasicChargePriceInt)
		// 再來要搞刪除
		var deletePatientBillBasicChargeResponse struct {
			DeletePatientBillBasicCharge bool
		}
		deletePatientBillBasicChargeRequest := `
		mutation {
			deletePatientBillBasicCharge(
				basicChargeId: "` + updatedPatientBill.BasicCharges[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillBasicChargeRequest, &deletePatientBillBasicChargeResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updatedPatientBill.BasicCharges[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}
