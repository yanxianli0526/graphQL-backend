package resolvers_test

import (
	"fmt"
	"graphql-go-template/internal/gql/generated"
	"graphql-go-template/internal/gql/resolvers"
	"graphql-go-template/internal/models"
	"strconv"
	"testing"
	"time"

	orm "graphql-go-template/internal/database"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.smart-aging.tech/devops/ms-go-kit/observability"
)

////// 起手為:正數(收費多)在區間內更新
// 新增總額為正數(收費多)在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		fmt.Println("createdPatientBill", createdPatientBill.ID)
		fmt.Println("createdPatientBill.AmountDue", createdPatientBill.AmountDue)

		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-transferRefundLeavePriceInt+updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-transferRefundLeavePriceInt-updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-transferRefundLeavePriceInt)
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-transferRefundLeavePriceInt)
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

////// 起手為:正數(收費多)不在區間內更新
// 新增總額為正數(收費多)不在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveNotInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)不在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveNotInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)不在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveNotInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為正數(收費多)不在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsPositiveNotInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "charge"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

////// 起手為:負數(欠費多)在區間內更新
// 新增總額為負數(欠費多)在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+transferRefundLeavePriceInt+updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+transferRefundLeavePriceInt-updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+transferRefundLeavePriceInt)
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-05-10T10:31:15.000Z"
					endDate: "2023-05-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		transferRefundLeavePriceInt, err := strconv.Atoi(transferRefundLeavePrice)
		if err != nil {
			fmt.Println("patientBillSubsidyPriceInt strconv.Atoi(transferRefundLeavePrice)", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, -transferRefundLeavePriceInt)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+transferRefundLeavePriceInt)
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + createdPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

////// 起手為:負數(欠費多)不在區間內更新
// 新增總額為負數(欠費多)不在區間內更新成在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeNotInIntervalToPositiveInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeNotInIntervalToNegativeInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}

		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-05-09T10:31:15.000Z"
					endDate: "2023-05-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		updateTransferRefundLeavePriceInt, err := strconv.Atoi(updateTransferRefundLeavePrice)
		if err != nil {
			fmt.Println("updatePatientBillSubsidyPriceInt strconv.Atoi(patientBillSubsidyPrice)", err)
		}

		require.Equal(t, updatedPatientBill.AmountDue, createdPatientBill.AmountDue-updateTransferRefundLeavePriceInt)
		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + updatedPatientBill.TransferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue+updateTransferRefundLeavePriceInt)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成不在區間內總額為正數(收費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeNotInIntervalToPositiveNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "charge"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}

// 新增總額為負數(欠費多)不在區間內更新成不在區間內總額為負數(欠費多)
func TestMutationResolver_AddTransferRefundLeave_testPriceIsNegativeNotInIntervalToNegativeNotInInterval(t *testing.T) {
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

		// 這邊拿剛剛新增完的住民帳單新增異動
		var createTransferRefundLeaveResponse struct {
			CreateTransferRefundLeave bool
		}
		transferRefundLeavePrice := "1500"
		transferRefundLeaveType := "refund"
		transferRefundLeaveRequest := `
		mutation {
			createTransferRefundLeave (
				patientId: "` + patientId + `"
				input: {
					startDate: "2023-04-10T10:31:15.000Z"
					endDate: "2023-04-11T10:31:15.000Z"
					reason: "test1"
					isReserveBed: "true"
					note:"test2"
					items:[{itemName:"膳食費",type:"` + transferRefundLeaveType + `",price:` + transferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(transferRefundLeaveRequest, &createTransferRefundLeaveResponse, addContext(user))
		patientBillId, err := uuid.Parse(patientBillResponse.CreatePatientBill)
		if err != nil {
			fmt.Println("uuid.Parse(patientBillResponse.CreatePatientBill)", err)
		}
		createdPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("createdPatientBill orm.GetPatientBillById", err)
		}
		// 確認住民漲單的應繳金額跟新增的異動金額是否相同
		require.Equal(t, createdPatientBill.AmountDue, 0)
		// 再來測試更新還有刪除
		// 先更新
		var updatePatientBillSubsidyResponse struct {
			UpdateTransferRefundLeave bool
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
		transferRefundLeaves, err := orm.GetTransferRefundLeavesByPatientIdBetweenEndDate(gorm.DB, organizationId, patientUUID, startDate, endDate)
		if err != nil {
			fmt.Println("orm.GetTransferRefundLeavesByPatientIdBetweenEndDate", err)
		}
		updateTransferRefundLeavePrice := "1800"
		updateTransferRefundLeaveType := "refund"
		updateTransferRefundLeaveRequest := `
		mutation {
			updateTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
				input: {
					startDate: "2023-04-09T10:31:15.000Z"
					endDate: "2023-04-10T10:31:15.000Z"
					reason: "t1"
					isReserveBed: "false"
					note:"t2"
					items:[{itemName:"膳食費1",type:"` + updateTransferRefundLeaveType + `",price:` + updateTransferRefundLeavePrice + `}]
				}
			)
		}
		`
		c.MustPost(updateTransferRefundLeaveRequest, &updatePatientBillSubsidyResponse, addContext(user))
		updatedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("updatedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, updatedPatientBill.AmountDue, 0)

		// 再來要搞刪除
		var DeleteTransferRefundLeaveResponse struct {
			DeleteTransferRefundLeave bool
		}
		DeleteTransferRefundLeaveRequest := `
		mutation {
			deleteTransferRefundLeave(
				id: "` + transferRefundLeaves[0].ID.String() + `"
			)
		}
		`
		c.MustPost(DeleteTransferRefundLeaveRequest, &DeleteTransferRefundLeaveResponse, addContext(user))
		deletedPatientBill, err := orm.GetPatientBillById(gorm.DB, organizationId, patientBillId)
		if err != nil {
			fmt.Println("deletedPatientBill orm.GetPatientBillById", err)
		}
		require.Equal(t, deletedPatientBill.AmountDue, updatedPatientBill.AmountDue)
		require.Equal(t, deletedPatientBill.AmountDue, 0)
	})
}
