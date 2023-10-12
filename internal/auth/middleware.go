package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	orm "graphql-go-template/internal/database"
	"graphql-go-template/internal/gql/resolvers"
	"graphql-go-template/pkg/app"

	"github.com/google/uuid"
)

// Middleware is used to handle auth logic
func Middleware(orm *orm.GormDatabase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			// 取token
			token, err := GetTokenFromHeader(auth)
			if token == "" || err != nil {
				http.Error(w, "GetTokenFromHeader is fail", http.StatusUnauthorized)
				return
			}

			// 驗證 JWT
			userClaims, err := app.ValidateJWT(token)
			if err != nil {
				http.Error(w, "ValidateJWT is fail", http.StatusUnauthorized)
				return
			}

			// 取得 userId
			userId, err := uuid.Parse(userClaims.UserId)
			if err != nil {
				http.Error(w, "userClaims.UserId uuid.Parse fail", http.StatusUnauthorized)
				return
			}

			// 取得 user
			user, err := orm.GetUserById(userId)
			if err != nil {
				http.Error(w, "GetUserById is fail", http.StatusUnauthorized)
				return
			}

			// 比對 database 上 expired 的時間
			isExpired := user.TokenExpiredAt.Before(time.Now())
			if isExpired {
				http.Error(w, "user TokenExpiredAt is Before", http.StatusUnauthorized)
				return
			}

			// Write your fancy token introspection logic here and if valid user then pass appropriate key in header
			// IMPORTANT: DO NOT HANDLE UNAUTHORIZED USER HERE
			ctx = context.WithValue(ctx, resolvers.UserIdCtxKey, user.ID)
			ctx = context.WithValue(ctx, resolvers.OrganizationId, user.OrganizationId)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTokenFromHeader(auth string) (string, error) {

	prefix := "Bearer "
	hasBearerPrefix := strings.HasPrefix(auth, prefix)
	if !hasBearerPrefix {
		return "", errors.New("errormsg.InvalidToken")
	}

	token := strings.TrimPrefix(auth, prefix)
	return token, nil
}
