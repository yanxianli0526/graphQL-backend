package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	orm "graphql-go-template/internal/database"
	gqlmodels "graphql-go-template/internal/gql/models"
	"graphql-go-template/internal/models"
	"graphql-go-template/pkg/nis/auth"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Mutations
// func (r *mutationResolver) Login(ctx context.Context, input *gqlmodels.UserInput) (string, error) {
// 	return login(r, input)
// }

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("Logout uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "logout"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("Logout orm.GetUserById", zap.Error(err), zap.String("fieldName", "logout"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}
	tokenExpiredAt := time.Now()

	// 這邊是更新so cool的token過期時間
	err = orm.UserLogout(r.ORM.DB, user, tokenExpiredAt)
	if err != nil {
		r.Logger.Error("Logout orm.UserLogout", zap.Error(err), zap.String("fieldName", "logout"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	token := &oauth2.Token{}
	err = json.Unmarshal(user.ProviderToken, token)
	if err != nil {
		r.Logger.Error("Logout json.Unmarshal", zap.Error(err), zap.String("fieldName", "logout"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return false, err
	}

	// 這邊是刪除nis的token
	auth.RevokeToken(r.ClientID, r.ClientSecret, token)

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("logout run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "logout"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return true, nil
}

// Queries
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("Users uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("fieldName", "users"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	users, err := orm.GetUsers(r.ORM.DB, organizationId)
	if err != nil {
		r.Logger.Error("Users orm.GetUsers", zap.Error(err), zap.Error(err), zap.String("fieldName", "users"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("users run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "users"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(id)
	if err != nil {
		r.Logger.Warn("User uuid.Parse(id)", zap.Error(err), zap.String("fieldName", "user"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("User orm.GetUserById", zap.Error(err), zap.String("fieldName", "user"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("user run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "user"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return user, nil
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		r.Logger.Warn("Me uuid.Parse(userIdStr)", zap.Error(err), zap.String("fieldName", "me"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}

	user, err := orm.GetUserById(r.ORM.DB, userId)
	if err != nil {
		r.Logger.Error("Me orm.GetUserById", zap.Error(err), zap.String("fieldName", "me"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))
		return nil, err
	}
	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("me run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "me"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()))

	return user, nil
}

// user resolvers
type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.ID.String(), nil
}

func (r *userResolver) Preference(ctx context.Context, obj *models.User) (*gqlmodels.UserPreference, error) {
	result := &gqlmodels.UserPreference{}
	err := json.Unmarshal(obj.Preference, &result)
	if err != nil {
		r.Logger.Error("User Preference is inValid", zap.Error(fmt.Errorf("User Preference is inValid ")),
			zap.String("fieldName", "organizationReceipt"), zap.Int64("timestamp", time.Now().Unix()))
	}
	return result, nil
}

// func login(r *mutationResolver, input *gqlmodels.UserInput) (string, error) {
// 	// Implement your login logic here
// 	return "MyFakeToken", nil
// }
