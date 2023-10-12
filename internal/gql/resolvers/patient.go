package resolvers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/models"
	"sort"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Queries
func (r *queryResolver) Patients(ctx context.Context) ([]*models.Patient, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("Patients uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "patients"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patients, err := orm.GetPatients(r.ORM.DB, organizationId, true)
	if err != nil {
		r.Logger.Error("Patients orm.GetPatients", zap.Error(err), zap.String("originalUrl", "patients"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	sort.SliceStable(patients, func(i, j int) bool {
		if patients[i].Branch > patients[j].Branch {
			return false
		}
		if patients[i].Branch < patients[j].Branch {
			return true
		}
		if patients[i].Room > patients[j].Room {
			return false
		}
		if patients[i].Room < patients[j].Room {
			return true
		}
		if patients[i].Bed > patients[j].Bed {
			return false
		}
		if patients[i].Bed < patients[j].Bed {
			return true
		}
		if patients[i].Status > patients[j].Status {
			return false
		}
		if patients[i].Status < patients[j].Status {
			return true
		}
		if patients[i].LastName > patients[j].LastName {
			return false
		}
		if patients[i].LastName < patients[j].LastName {
			return true
		}
		if patients[i].FirstName > patients[j].FirstName {
			return false
		}
		if patients[i].FirstName < patients[j].FirstName {
			return true
		}

		return false
	})
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patients run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patients"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return patients, nil
}

func (r *queryResolver) Patient(ctx context.Context, id string) (*models.Patient, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))

	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("Patient uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "patient"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patientId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("Patient uuid.Parse(id)", zap.Error(err), zap.String("originalUrl", "patient"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	patient, err := orm.GetPatientById(r.ORM.DB, patientId, organizationId)
	if err != nil {
		r.Logger.Error("Patient orm.GetPatientById", zap.Error(err), zap.String("originalUrl", "patient"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("patient run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "patient"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return patient, nil
}

// patient resolvers
type patientResolver struct{ *Resolver }

func (r *patientResolver) ID(ctx context.Context, obj *models.Patient) (string, error) {
	return obj.ID.String(), nil
}
