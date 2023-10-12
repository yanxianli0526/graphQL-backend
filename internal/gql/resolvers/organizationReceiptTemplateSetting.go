package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Mutations
func (r *mutationResolver) CreateOrganizationReceiptTemplateSetting(ctx context.Context, input gqlmodels.OrganizationReceiptTemplateSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateOrganizationReceiptTemplateSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var organizationPicture string
	if input.OrganizationPicture != nil {
		organizationPicture = *input.OrganizationPicture
	}

	var sealOnePicture string
	if input.SealOnePicture != nil {
		sealOnePicture = *input.SealOnePicture
	}
	var sealTwoPicture string
	if input.SealTwoPicture != nil {
		sealTwoPicture = *input.SealTwoPicture
	}

	var sealThreePicture string
	if input.SealThreePicture != nil {
		sealThreePicture = *input.SealThreePicture
	}

	var sealFourPicture string
	if input.SealFourPicture != nil {
		sealFourPicture = *input.SealFourPicture
	}

	organizationReceiptTemplateSetting := models.OrganizationReceiptTemplateSetting{
		ID:                  uuid.New(),
		Name:                input.Name,
		TaxTypes:            input.TaxTypes,
		OrganizationPicture: organizationPicture,
		TitleName:           input.TitleName,
		PatientInfo:         input.PatientInfo,
		PriceShowType:       input.PriceShowType,
		OrganizationInfoOne: input.OrganizationInfoOne,
		OrganizationInfoTwo: input.OrganizationInfoTwo,
		NoteText:            input.NoteText,
		SealOneName:         input.SealOneName,
		SealOnePicture:      sealOnePicture,
		SealTwoName:         input.SealTwoName,
		SealTwoPicture:      sealTwoPicture,
		SealThreeName:       input.SealThreeName,
		SealThreePicture:    sealThreePicture,
		SealFourName:        input.SealFourName,
		SealFourPicture:     sealFourPicture,
		PartOneName:         input.PartOneName,
		PartTwoName:         input.PartTwoName,
		OrganizationId:      organizationId,
		Organization:        models.Organization{},
	}

	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 新增目前的設定
		err = orm.CreateOrganizationReceiptTemplateSetting(tx, &organizationReceiptTemplateSetting)
		if err != nil {
			r.Logger.Error("CreateOrganizationReceiptTemplateSetting orm.CreateOrganizationReceiptTemplateSetting", zap.Error(err), zap.String("originalUrl", "createOrganizationReceiptTemplateSetting"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		// 	如果有taxTypes代表 有其他的要被更新
		if len(input.TaxTypes) > 0 {
			organizationReceiptTemplateSettingElements := make(map[uuid.UUID]*models.OrganizationReceiptTemplateSetting)
			for i := range input.TaxTypes {
				// 找出各個稅別的內容(再整理資料)
				organizationReceiptTemplateSettingData, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(tx, organizationId, input.TaxTypes[i])
				if err != nil {
					r.Logger.Error("CreateOrganizationReceiptTemplateSetting orm.GetOrganizationReceiptTemplateSettingInTaxType", zap.Error(err), zap.String("originalUrl", "createOrganizationReceiptTemplateSetting"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
				// 整理最後要被更新的資料
				if organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID] == nil {
					var removeCount int
					var taxTypes []string
					for j := range organizationReceiptTemplateSettingData.TaxTypes {
						if input.TaxTypes[i] == organizationReceiptTemplateSettingData.TaxTypes[j] {
							removeCount = j
						}
						taxTypes = append(taxTypes, organizationReceiptTemplateSettingData.TaxTypes[j])
					}
					taxTypes = remove(taxTypes, removeCount)
					organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID] = &models.OrganizationReceiptTemplateSetting{
						ID:             organizationReceiptTemplateSettingData.ID,
						TaxTypes:       taxTypes,
						OrganizationId: organizationId,
					}
				} else {
					var removeCount int
					// 用上面filter後的taxTypes 再去確定有哪些是不要的
					for j := range organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes {
						if input.TaxTypes[i] == organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes[j] {
							removeCount = j
						}
					}
					organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes = remove(organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes, removeCount)
				}
			}

			// 更新項目的taxTypes
			for i := range organizationReceiptTemplateSettingElements {
				err = orm.UpdateOrganizationReceiptTemplateSettingTaxTypes(tx, organizationReceiptTemplateSettingElements[i])
				if err != nil {
					r.Logger.Error("CreateOrganizationReceiptTemplateSetting orm.UpdateOrganizationReceiptTemplateSettingTaxTypes", zap.Error(err), zap.String("originalUrl", "createOrganizationReceiptTemplateSetting"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("CreateOrganizationReceiptTemplateSetting tx.Transaction", zap.Error(err), zap.String("originalUrl", "createOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createOrganizationReceiptTemplateSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createOrganizationReceiptTemplateSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) UpdateOrganizationReceiptTemplateSetting(ctx context.Context, receiptTemplateSettingIdStr string, input gqlmodels.OrganizationReceiptTemplateSettingInput) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationReceiptTemplateSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	receiptTemplateSettingId, err := uuid.Parse(receiptTemplateSettingIdStr)
	if err != nil {
		r.Logger.Warn("UpdateOrganizationReceiptTemplateSetting uuid.Parse(receiptTemplateSettingIdStr)", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	var organizationPicture string
	if input.OrganizationPicture != nil {
		organizationPicture = *input.OrganizationPicture
	}

	var sealOnePicture string
	if input.SealOnePicture != nil {
		sealOnePicture = *input.SealOnePicture
	}
	var sealTwoPicture string
	if input.SealTwoPicture != nil {
		sealTwoPicture = *input.SealTwoPicture
	}

	var sealThreePicture string
	if input.SealThreePicture != nil {
		sealThreePicture = *input.SealThreePicture
	}

	var sealFourPicture string
	if input.SealFourPicture != nil {
		sealFourPicture = *input.SealFourPicture
	}

	organizationReceiptTemplateSetting := models.OrganizationReceiptTemplateSetting{
		ID:                  receiptTemplateSettingId,
		Name:                input.Name,
		TaxTypes:            input.TaxTypes,
		OrganizationPicture: organizationPicture,
		TitleName:           input.TitleName,
		PatientInfo:         input.PatientInfo,
		PriceShowType:       input.PriceShowType,
		OrganizationInfoOne: input.OrganizationInfoOne,
		OrganizationInfoTwo: input.OrganizationInfoTwo,
		NoteText:            input.NoteText,
		SealOneName:         input.SealOneName,
		SealOnePicture:      sealOnePicture,
		SealTwoName:         input.SealTwoName,
		SealTwoPicture:      sealTwoPicture,
		SealThreeName:       input.SealThreeName,
		SealThreePicture:    sealThreePicture,
		SealFourName:        input.SealFourName,
		SealFourPicture:     sealFourPicture,
		PartOneName:         input.PartOneName,
		PartTwoName:         input.PartTwoName,
		OrganizationId:      organizationId,
	}
	tx := r.ORM.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// 新增目前的設定
		err = orm.UpdateOrganizationReceiptTemplateSetting(tx, &organizationReceiptTemplateSetting)
		if err != nil {
			r.Logger.Error("UpdateOrganizationReceiptTemplateSetting orm.UpdateOrganizationReceiptTemplateSetting", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
				zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
			return err
		}
		// 	如果有taxTypes代表 有其他的要被更新
		if len(input.TaxTypes) > 0 {
			organizationReceiptTemplateSettingElements := make(map[uuid.UUID]*models.OrganizationReceiptTemplateSetting)
			for i := range input.TaxTypes {
				// 找出各個稅別的內容(再整理資料)
				organizationReceiptTemplateSettingData, err := orm.GetOrganizationReceiptTemplateSettingInTaxType(tx, organizationId, input.TaxTypes[i])
				if err != nil {
					r.Logger.Error("UpdateOrganizationReceiptTemplateSetting orm.GetOrganizationReceiptTemplateSettingInTaxType", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
				// 不用更新自己(上面已經更新過了)
				// 這邊主要是在更新別的設定的稅別
				if organizationReceiptTemplateSettingData.ID != receiptTemplateSettingId {
					// 整理最後要被更新的資料
					if organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID] == nil {
						var removeCount int
						var taxTypes []string
						for j := range organizationReceiptTemplateSettingData.TaxTypes {
							if input.TaxTypes[i] == organizationReceiptTemplateSettingData.TaxTypes[j] {
								removeCount = j
							}
							taxTypes = append(taxTypes, organizationReceiptTemplateSettingData.TaxTypes[j])
						}
						taxTypes = remove(taxTypes, removeCount)
						organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID] = &models.OrganizationReceiptTemplateSetting{
							ID:             organizationReceiptTemplateSettingData.ID,
							TaxTypes:       taxTypes,
							OrganizationId: organizationId,
						}
					} else {
						var removeCount int
						// 用上面filter後的taxTypes 再去確定有哪些是不要的
						for j := range organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes {
							if input.TaxTypes[i] == organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes[j] {
								removeCount = j
							}
						}
						organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes = remove(organizationReceiptTemplateSettingElements[organizationReceiptTemplateSettingData.ID].TaxTypes, removeCount)
					}
				}
			}

			// 更新項目的taxTypes
			for i := range organizationReceiptTemplateSettingElements {
				err = orm.UpdateOrganizationReceiptTemplateSettingTaxTypes(tx, organizationReceiptTemplateSettingElements[i])
				if err != nil {
					r.Logger.Error("UpdateOrganizationReceiptTemplateSetting orm.UpdateOrganizationReceiptTemplateSettingTaxTypes", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
						zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		r.Logger.Error("UpdateOrganizationReceiptTemplateSetting tx.Transaction", zap.Error(err), zap.String("originalUrl", "updateOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("updateOrganizationReceiptTemplateSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "updateOrganizationReceiptTemplateSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

func (r *mutationResolver) DeleteOrganizationReceiptTemplateSetting(ctx context.Context, receiptTemplateSettingIdStr string) (bool, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("DeleteOrganizationReceiptTemplateSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "deleteOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	receiptTemplateSettingId, err := uuid.Parse(receiptTemplateSettingIdStr)
	if err != nil {
		r.Logger.Warn("DeleteOrganizationReceiptTemplateSetting uuid.Parse(receiptTemplateSettingIdStr)", zap.Error(err), zap.String("originalUrl", "deleteOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}

	err = orm.DeleteOrganizationReceiptTemplateSettingById(r.ORM.DB, receiptTemplateSettingId, organizationId)
	if err != nil {
		r.Logger.Error("DeleteOrganizationReceiptTemplateSetting orm.DeleteOrganizationReceiptTemplateSettingById", zap.Error(err), zap.String("originalUrl", "deleteOrganizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return false, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("deleteOrganizationReceiptTemplateSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "deleteOrganizationReceiptTemplateSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return true, nil
}

// Queries
func (r *queryResolver) OrganizationReceiptTemplateSetting(ctx context.Context, receiptTemplateSettingIdStr string) (*models.OrganizationReceiptTemplateSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationReceiptTemplateSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "organizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	receiptTemplateSettingId, err := uuid.Parse(receiptTemplateSettingIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationReceiptTemplateSetting uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "organizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	organizationReceiptTemplateSetting, err := orm.GetOrganizationReceiptTemplateSettingById(r.ORM.DB, receiptTemplateSettingId, organizationId)
	if err != nil {
		r.Logger.Error("OrganizationReceiptTemplateSetting orm.GetOrganizationReceiptTemplateSettingById", zap.Error(err), zap.String("originalUrl", "organizationReceiptTemplateSetting"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationReceiptTemplateSetting run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationReceiptTemplateSetting"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return organizationReceiptTemplateSetting, nil
}

func (r *queryResolver) OrganizationReceiptTemplateSettings(ctx context.Context) ([]*models.OrganizationReceiptTemplateSetting, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("OrganizationReceiptTemplateSettings uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "organizationReceiptTemplateSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	organizationReceiptTemplateSettings, err := orm.GetOrganizationReceiptTemplateSettings(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("OrganizationReceiptTemplateSettings orm.GetOrganizationReceiptTemplateSettings", zap.Error(err), zap.String("originalUrl", "organizationReceiptTemplateSettings"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("organizationReceiptTemplateSettings run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "organizationReceiptTemplateSettings"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return organizationReceiptTemplateSettings, nil
}

// organizationReceipt resolvers
type organizationReceiptTemplateSettingResolver struct{ *Resolver }

func (r *organizationReceiptTemplateSettingResolver) ID(ctx context.Context, obj *models.OrganizationReceiptTemplateSetting) (string, error) {
	return obj.ID.String(), nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
