package resolvers_test

import (
	"fmt"
	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/generated"
	"graphql-go-template/internal/gql/resolvers"
	"graphql-go-template/internal/models"
	"strconv"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.smart-aging.tech/devops/ms-go-kit/observability"
)

////// 起手為:正數(收費多)在區間內更新
// 新增總額為正數(收費多)在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordSubtotalFloat)

		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}

		updateNonFixedChargeRecordPrice := "1800"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat+updateNonFixedChargeRecordSubtotalFloat)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為正數(收費多)在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordSubtotalFloat)

		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat-updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// // 新增總額為正數(收費多)在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordSubtotalFloat)

		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-(nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat)
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為正數(收費多)在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, nonFixedChargeRecordSubtotalFloat)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-(nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-nonFixedChargeRecordSubtotalFloat)
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// ////// 起手為:正數(收費多)不在區間內更新
// // 新增總額為正數(收費多)不在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveNotInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// // 新增總額為正數(收費多)不在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveNotInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為正數(收費多)不在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveNotInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// // 新增總額為正數(收費多)不在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsPositiveNotInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

////// 起手為:負數(欠費多)在區間內更新
// 新增總額為負數(欠費多)在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`

		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		// 確認住民漲單的應繳金額跟新增的非非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordSubtotalFloat)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat+updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為負數(欠費多)在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordSubtotalFloat)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat-updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var DeleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		DeleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteNonFixedChargeRecordRequest, &DeleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// // 新增總額為負數(欠費多)在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordSubtotalFloat)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+(nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat)
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為負數(欠費多)在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		nonFixedChargeRecordPriceInt, err := strconv.Atoi(nonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordPrice)", err)
		}
		nonFixedChargeRecordQuantityInt, err := strconv.Atoi(nonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		nonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(nonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt)
		require.Equal(t, createdPatientBill.AmountDue, -nonFixedChargeRecordSubtotalFloat)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+(nonFixedChargeRecordPriceInt*nonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+nonFixedChargeRecordSubtotalFloat)
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + createdPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

////// 起手為:負數(欠費多)不在區間內更新
// 新增總額為負數(欠費多)不在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeNotInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "charge"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeNotInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-05-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateNonFixedChargeRecordPriceInt, err := strconv.Atoi(updateNonFixedChargeRecordPrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		updateNonFixedChargeRecordQuantityInt, err := strconv.Atoi(updateNonFixedChargeRecordQuantity)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		updateNonFixedChargeRecordSubtotalFloat, err := strconv.Atoi(updateNonFixedChargeRecordSubtotal)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(nonFixedChargeRecordQuantity)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-updateNonFixedChargeRecordSubtotalFloat)
		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + updatedPatientBill.NonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+(updateNonFixedChargeRecordPriceInt*updateNonFixedChargeRecordQuantityInt))
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateNonFixedChargeRecordSubtotalFloat)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeNotInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "charge"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddNonFixedChargeRecord_testPriceIsNegativeNotInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增非固定月費
		var createNonFixedChargeRecordResponse struct {
			CreateNonFixedChargeRecord bool
		}
		nonFixedChargeRecordPrice := "1500.0"
		nonFixedChargeRecordQuantity := "2"
		nonFixedChargeRecordSubtotal := "3000.0"
		nonFixedChargeRecordType := "refund"
		nonFixedChargeRecordRequest := `
		mutation {
			createNonFixedChargeRecord (
				patientId: "` + patientId + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + nonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + nonFixedChargeRecordPrice + `
        	quantity: ` + nonFixedChargeRecordQuantity + `
        	subtotal: ` + nonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(nonFixedChargeRecordRequest, &createNonFixedChargeRecordResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的非非固定月費金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0.0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateNonFixedChargeRecord bool
		}
		// 因為原本的create不在區間內 所以不能從patientBill去找 要直接找異動
		taipeiZone, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			fmt.Println("time.LoadLocation(Asia/Taipei)", err)
		}
		startDate := time.Date(2023, 4, 5, 1, 1, 1, 0, taipeiZone)
		endDate := time.Date(2023, 4, 20, 1, 1, 1, 0, taipeiZone)
		patientUUID, err := uuid.Parse(patientId)
		if err != nil {
			fmt.Println("uuid.Parse(patientId)", err)
		}
		nonFixedChargeRecords, err := orm.GetNonFixedChargeRecordsByPatientIdAndDate(gorm.DB, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetNonFixedChargeRecordsByPatientIdBetweenEndDate", err)
		}

		updateNonFixedChargeRecordPrice := "1800.0"
		updateNonFixedChargeRecordQuantity := "2"
		updateNonFixedChargeRecordSubtotal := "3600.0"
		updateNonFixedChargeRecordType := "refund"
		updateNonFixedChargeRecordRequest := `
		mutation {
			updateNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
				input: {
					nonFixedChargeDate: "2023-04-10T10:31:15.000Z"
					itemCategory:"耗材"
        	itemName: "項目2"
       	 	type:"` + updateNonFixedChargeRecordType + `"
        	unit: "包"
        	price: ` + updateNonFixedChargeRecordPrice + `
        	quantity: ` + updateNonFixedChargeRecordQuantity + `
        	subtotal: ` + updateNonFixedChargeRecordSubtotal + `
        	note: "note1"
        	taxType:"stampTax"
				}
			)
		}
		`
		c.MustPost(updateNonFixedChargeRecordRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0.0)

		// 再來要搞刪除
		var deleteNonFixedChargeRecordResponse struct {
			DeleteNonFixedChargeRecord bool
		}
		deleteNonFixedChargeRecordRequest := `
		mutation {
			deleteNonFixedChargeRecord(
				id: "` + nonFixedChargeRecords[0].ID.String() + `"
			)
		}
		`
		c.MustPost(deleteNonFixedChargeRecordRequest, &deleteNonFixedChargeRecordResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0.0)
	})
}
