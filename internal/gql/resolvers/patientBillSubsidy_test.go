package resolvers_test

import (
	"fmt"
	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/generated"
	"graphql-go-template/internal/gql/resolvers"
	"graphql-go-template/internal/models"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.smart-aging.tech/devops/ms-go-kit/observability"
)

func TestMutationResolver_AddPatientBillSubsidyTypeIsChargeToCharge(t *testing.T) {
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
		var addPatientBillSubsidyResponse struct {
			AddPatientBillSubsidy bool
		}
		patientBillSubsidyPrice := "1500"
		patientBillSubsidyType := "charge"
		addPatientBillSubsidyRequest := `
		mutation {
			addPatientBillSubsidy(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + patientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillSubsidyRequest, &addPatientBillSubsidyResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillSubsidyPriceInt, err := strconv.Atoi(patientBillSubsidyPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, patientBillSubsidyPriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdatePatientBillSubsidy bool
		}
		updatePatientBillSubsidyPrice := "1800"
		updatePatientBillSubsidyType := "charge"
		updatePatientBillSubsidyRequest := `
		mutation {
			updatePatientBillSubsidy(
				input: {
					subsidyId: "` + createdPatientBill.Subsidies[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + updatePatientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillSubsidyRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillSubsidyPriceInt, err := strconv.Atoi(updatePatientBillSubsidyPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-patientBillSubsidyPriceInt+updatePatientBillSubsidyPriceInt)
		// 再來要搞刪除
		var deletePatientBillSubsidyResponse struct {
			DeletePatientBillSubsidy bool
		}
		deletePatientBillSubsidyRequest := `
		mutation {
			deletePatientBillSubsidy(
				subsidyId: "` + updatedPatientBill.Subsidies[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillSubsidyRequest, &deletePatientBillSubsidyResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updatedPatientBill.Subsidies[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillSubsidyTypeIsChargeToRefund(t *testing.T) {
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
		var addPatientBillSubsidyResponse struct {
			AddPatientBillSubsidy bool
		}
		patientBillSubsidyPrice := "1500"
		patientBillSubsidyType := "charge"
		addPatientBillSubsidyRequest := `
		mutation {
			addPatientBillSubsidy(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + patientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillSubsidyRequest, &addPatientBillSubsidyResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillSubsidyPriceInt, err := strconv.Atoi(patientBillSubsidyPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, patientBillSubsidyPriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdatePatientBillSubsidy bool
		}
		updatePatientBillSubsidyPrice := "1400"
		updatePatientBillSubsidyType := "refund"
		updatePatientBillSubsidyRequest := `
		mutation {
			updatePatientBillSubsidy(
				input: {
					subsidyId: "` + createdPatientBill.Subsidies[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + updatePatientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillSubsidyRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillSubsidyPriceInt, err := strconv.Atoi(updatePatientBillSubsidyPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		fmt.Println("updatedPatientBill.AmountDue", updatedPatientBill.AmountDue)
		fmt.Println("createdPatientBill.AmountDue-patientBillSubsidyPriceInt-updatePatientBillSubsidyPriceInt", createdPatientBill.AmountDue-patientBillSubsidyPriceInt-updatePatientBillSubsidyPriceInt)

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-patientBillSubsidyPriceInt-updatePatientBillSubsidyPriceInt)
		// 再來要搞刪除
		var deletePatientBillSubsidyResponse struct {
			DeletePatientBillSubsidy bool
		}
		deletePatientBillSubsidyRequest := `
		mutation {
			deletePatientBillSubsidy(
				subsidyId: "` + updatedPatientBill.Subsidies[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillSubsidyRequest, &deletePatientBillSubsidyResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updatedPatientBill.Subsidies[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillSubsidyTypeIsRefundToCharge(t *testing.T) {
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
		var addPatientBillSubsidyResponse struct {
			AddPatientBillSubsidy bool
		}
		patientBillSubsidyPrice := "1500"
		patientBillSubsidyType := "refund"
		addPatientBillSubsidyRequest := `
		mutation {
			addPatientBillSubsidy(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + patientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillSubsidyRequest, &addPatientBillSubsidyResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillSubsidyPriceInt, err := strconv.Atoi(patientBillSubsidyPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -patientBillSubsidyPriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdatePatientBillSubsidy bool
		}
		updatePatientBillSubsidyPrice := "1400"
		updatePatientBillSubsidyType := "charge"
		updatePatientBillSubsidyRequest := `
		mutation {
			updatePatientBillSubsidy(
				input: {
					subsidyId: "` + createdPatientBill.Subsidies[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + updatePatientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillSubsidyRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillSubsidyPriceInt, err := strconv.Atoi(updatePatientBillSubsidyPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+patientBillSubsidyPriceInt+updatePatientBillSubsidyPriceInt)
		// 再來要搞刪除
		var deletePatientBillSubsidyResponse struct {
			DeletePatientBillSubsidy bool
		}
		deletePatientBillSubsidyRequest := `
		mutation {
			deletePatientBillSubsidy(
				subsidyId: "` + updatedPatientBill.Subsidies[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillSubsidyRequest, &deletePatientBillSubsidyResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updatedPatientBill.Subsidies[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

func TestMutationResolver_AddPatientBillSubsidyTypeIsRefundToRefund(t *testing.T) {
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
		var addPatientBillSubsidyResponse struct {
			AddPatientBillSubsidy bool
		}
		patientBillSubsidyPrice := "1500"
		patientBillSubsidyType := "refund"
		addPatientBillSubsidyRequest := `
		mutation {
			addPatientBillSubsidy(
				input: {
					id: "` + patientBillResponse.CreatePatientBill + `"
					itemName: "z6"
					type: "` + patientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + patientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "String1"
				}
			)
		}
		`
		c.MustPost(addPatientBillSubsidyRequest, &addPatientBillSubsidyResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		patientBillSubsidyPriceInt, err := strconv.Atoi(patientBillSubsidyPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		// 確認住民漲單的應繳金額跟新增的固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -patientBillSubsidyPriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdatePatientBillSubsidy bool
		}
		updatePatientBillSubsidyPrice := "1400"
		updatePatientBillSubsidyType := "refund"
		updatePatientBillSubsidyRequest := `
		mutation {
			updatePatientBillSubsidy(
				input: {
					subsidyId: "` + createdPatientBill.Subsidies[0].ID.String() + `"
					itemName: "q1"
					type:"` + updatePatientBillSubsidyType + `"
					idNumber: "tt"
					unit: "月"
					price: ` + updatePatientBillSubsidyPrice + `
					startDate: "2023-06-01T14:46:54+08:00"
					endDate: "2023-06-30T14:46:54+08:00"
					note: "aaqqq"
				}
			)
		}
		`
		c.MustPost(updatePatientBillSubsidyRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updatePatientBillSubsidyPriceInt, err := strconv.Atoi(updatePatientBillSubsidyPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+patientBillSubsidyPriceInt-updatePatientBillSubsidyPriceInt)
		// 再來要搞刪除
		var deletePatientBillSubsidyResponse struct {
			DeletePatientBillSubsidy bool
		}
		deletePatientBillSubsidyRequest := `
		mutation {
			deletePatientBillSubsidy(
				subsidyId: "` + updatedPatientBill.Subsidies[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deletePatientBillSubsidyRequest, &deletePatientBillSubsidyResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updatedPatientBill.Subsidies[0].Price)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}
