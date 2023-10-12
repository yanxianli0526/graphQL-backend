package resolvers

import (
	"context"
	"fmt"
	orm "graphql-go-template/internal/database"
	"strconv"
	"time"

	gqlmodels "graphql-go-template/internal/gql/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Queries
func (r *queryResolver) PatientLatestFixedChargeRecords(ctx context.Context) ([]*gqlmodels.FixedChargeRecord, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("PatientLatestFixedChargeRecords uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "patientLatestFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	basicSettings, err := orm.GetBasicChargeSettings(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PatientLatestFixedChargeRecords orm.GetBasicChargeSettings", zap.Error(err), zap.String("fieldName", "patientLatestFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	subsidies, err := orm.GetSubsidiesSetting(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("PatientLatestFixedChargeRecords orm.GetSubsidiesSetting", zap.Error(err), zap.String("fieldName", "patientLatestFixedChargeRecords"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	// 先看固定月費
	elements := make(map[uuid.UUID]*gqlmodels.FixedChargeRecord)
	for i := range basicSettings {
		if elements[basicSettings[i].PatientId] == nil {
			elements[basicSettings[i].PatientId] = &gqlmodels.FixedChargeRecord{
				Patient:   &basicSettings[i].Patient,
				UpdatedAt: basicSettings[i].UpdatedAt,
			}
			elements[basicSettings[i].PatientId].Items = append(elements[basicSettings[i].PatientId].Items, basicSettings[i].OrganizationBasicChargeSetting.ItemName)
		} else {
			elements[basicSettings[i].PatientId].Items = append(elements[basicSettings[i].PatientId].Items, basicSettings[i].OrganizationBasicChargeSetting.ItemName)
		}
	}

	// 再看補助款
	for i := range subsidies {
		if elements[subsidies[i].PatientId] == nil {
			elements[subsidies[i].PatientId] = &gqlmodels.FixedChargeRecord{
				Patient:   &subsidies[i].Patient,
				UpdatedAt: subsidies[i].UpdatedAt,
			}
			elements[subsidies[i].PatientId].Items = append(elements[subsidies[i].PatientId].Items, subsidies[i].ItemName)
		} else {
			elements[subsidies[i].PatientId].Items = append(elements[subsidies[i].PatientId].Items, subsidies[i].ItemName)
		}
	}

	var fixedChargeRecords []*gqlmodels.FixedChargeRecord
	for _, v := range elements {
		fixedChargeRecord := gqlmodels.FixedChargeRecord{
			Patient:   elements[v.Patient.ID].Patient,
			UpdatedAt: elements[v.Patient.ID].UpdatedAt,
			Items:     elements[v.Patient.ID].Items,
		}
		fixedChargeRecords = append(fixedChargeRecords, &fixedChargeRecord)
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patientLatestFixedChargeRecords run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patientLatestFixedChargeRecords"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return fixedChargeRecords, nil
}
