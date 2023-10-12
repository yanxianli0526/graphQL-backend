package router

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	config "graphql-go-template/envconfig"
	"graphql-go-template/pkg/app"

	orm "graphql-go-template/internal/database"

	ms_kit "gitlab.smart-aging.tech/devops/ms-go-kit/observability"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"graphql-go-template/internal/models"
	"graphql-go-template/internal/restful/services"
	"graphql-go-template/pkg/nis/auth"

	nis "graphql-go-template/pkg/nis/nisApi"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func RegisterNIS(db *orm.GormDatabase, routerGroup *gin.RouterGroup, env config.AuthEndpoint,
) {
	sessionName := "nis_oauth"
	nisHandler := NewNISAPI(db, false, env)

	nisRouter := routerGroup.Group("/nis")
	nisRouter.Use(Session(sessionName))
	{
		nisRouter.GET("/login", nisHandler.LoginHandler)
		nisRouter.POST("/callback", nisHandler.CallbackHandler(env))
	}
}

type NISAPI struct {
	DB          *orm.GormDatabase
	OAuthConfig *auth.OAuth
	debug       bool
}

type frontendType string

const (
	secret frontendType = "jubo-inventory-toll"
)

func NewNISAPI(
	db *orm.GormDatabase,
	debug bool,
	env config.AuthEndpoint,
) *NISAPI {
	secretByte := []byte(secret)

	store = cookie.NewStore(secretByte)
	store.Options(sessions.Options{
		Domain:   "*.jubo.health",
		Path:     "/",
		MaxAge:   43200,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	oAuthConfig := auth.NewOAuth(
		env.AuthURL,
		env.TokenURL,
		env.RedirectURL,
		env.ClientID,
		env.ClientSecret)
	nisService := &NISAPI{
		DB:          db,
		OAuthConfig: oAuthConfig,
		debug:       true,
	}
	return nisService
}

func (u *NISAPI) LoginHandler(ctx *gin.Context) {
	log := ms_kit.GetLoggerFromGinCtx(ctx)
	state, err := generateRandomState()
	if err != nil {
		log.Error("failed to generateRandomState", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generateRandomState",
		})
		return
	}
	var logingURL = u.OAuthConfig.Config.AuthCodeURL(state)
	session := sessions.Default(ctx)

	log.Info("user login", zap.Any("session", session), zap.Any("state", state))
	session.Set("state", state)
	if err := session.Save(); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"loginURL": logingURL,
	})
}

// CallbackHandler will receive State and Code from NIS OAuth Callback
func (u *NISAPI) CallbackHandler(env config.AuthEndpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := ms_kit.GetLoggerFromGinCtx(ctx)
		session := sessions.Default(ctx)
		retrievedState := session.Get("state")
		log.Info("CallbackHandler", zap.Any("session", session), zap.Any("retrievedState", retrievedState))
		svc := services.New(ctx, u.DB)

		// 從前端來的code 和 Status
		callbackCredentials := &auth.CallbackCredential{}
		if err := ctx.Bind(callbackCredentials); err != nil {
			log.Warn("ctx.Bind failed: ", zap.Error(err))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing required fields"))
			return
		}
		if !u.debug && retrievedState != callbackCredentials.State {
			log.Warn("!u.debug && retrievedState != callbackCredentials.State failed: ", zap.Error(fmt.Errorf("retrievedState")))
			log.Warn("!u.debug && retrievedState != callbackCredentials.State failed: ", zap.Error(fmt.Errorf("callbackCredentials.State")))
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state: %s", callbackCredentials.State))
			return
		}
		var token *oauth2.Token
		var err error
		var client *http.Client
		token, err = u.OAuthConfig.GetTokenWithCode(callbackCredentials.Code)

		client = u.OAuthConfig.NewClient(token)
		if err != nil {
			log.Warn("ctx.Bind failed: ", zap.Error(err))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("nis callback fail: %s", err))
			return
		}

		if success := app.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
			log.Warn("GetTokenWithCode failed", zap.Error(err))
			return
		}

		nisAPI := nis.New(client, env)

		// // 來自nis的orgID
		nisOrgId := auth.GetOrgId(token)

		// 先檢查現在的資料庫有沒有organization了
		// 如果沒有就sync nis的資料庫新增一份給so cool
		organizationId, organizationIsExist := svc.GetOrganizationByProviderId(nisOrgId)
		if !organizationIsExist {
			nisOrganization, err := nisAPI.GetNisOrganizationInfo(nisOrgId, token)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("GetOrganizationInfo", zap.Error(err))
				return
			}

			organizationId, err = svc.FirstSyncOrganization(nisOrganization, nisOrgId)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("FirstSyncOrganization", zap.Error(err))
				return
			}
		}

		//  新增機構的收據設定(default)
		organizationReceiptCount := svc.GetOrganizationReceiptById(organizationId)
		if organizationReceiptCount < 1 {
			err = svc.FirstSyncOrganizationReceipt(organizationId)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("FirstSyncOrganizationReceipt", zap.Error(err))
				return
			}
		}

		//  新增機構的列印設定(default)
		organizationReceiptTemplateSettingCount := svc.GetOrganizationReceiptTemplateSettingById(organizationId)
		if organizationReceiptTemplateSettingCount < 1 {
			err = svc.FirstSyncOrganizationReceiptTemplateSetting(organizationId)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("FirstSyncOrganizationReceiptTemplateSetting", zap.Error(err))
				return
			}
		}

		// // 呼叫nis的User Info
		userInfo, err := nisAPI.GetNisUserInfo(nisOrgId, token)
		if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
			log.Error("nisAPI.GetUserInfo failed", zap.Error(err))
			return
		}

		// 包含 access_token refresh_token expiry
		tokenBytes, err := auth.MarshalToken(token)
		if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
			log.Error("nisOAuthThreeLegged.MarshalToken", zap.Error(err))
			return
		}

		// 取得so cool 的 token到期時間
		tokenExpiredAt, err := auth.GetRefreshTokenExpireAt(token)
		if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
			log.Error("nisOAuthThreeLegged.GetRefreshTokenExpireAt failed", zap.Error(err))
			return
		}

		// 先檢查現在的資料庫有沒有user了
		// 如果沒有就sync nis的資料庫新增一份給so cool
		user, userIsExist := svc.GetUserByProviderId(userInfo.ProviderId)
		// 如果找不到這個user表示第一次登入 (sync整份資料) 如果有找到只更新token
		if !userIsExist {
			user, err = svc.FirstSyncUser(userInfo, organizationId, tokenBytes, tokenExpiredAt)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("svc.FirstSyncUser", zap.Error(err))
				return
			}
		} else {
			err = svc.UpdateUserToken(user.ID, tokenBytes, tokenExpiredAt)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("svc.FirstSyncUser", zap.Error(err))
				return
			}
		}

		// 呼叫nis的這邊sync的是負責人員和住民
		usersAndPatients, err := nisAPI.GetNisPeopleInChargesAndPatientsInfo(nisOrgId, token)
		if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
			log.Error("nisAPI.GetNisPeopleInChargesAndPatientsInfo failed", zap.Error(err))
			return
		}
		var peopleInCharges []models.User
		if len(usersAndPatients.PeopleInCharges) > 0 {
			// 要先sync負責人員(和住民有關聯問題)
			peopleInCharges, err = svc.SyncPeopleInCharges(usersAndPatients.PeopleInCharges, organizationId, tokenBytes)
			if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
				log.Error("svc.SyncPeopleInCharges failed", zap.Error(err))
				return
			}
		}

		err = svc.SyncPatients(usersAndPatients.Patients, peopleInCharges, organizationId)
		if success := app.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
			log.Error("svc.SyncPatients failed", zap.Error(err))
			return
		}

		// /* return JWT info */
		signedToken, _ := app.GenerateJWT(user)

		log.Debug("generate JWT signedToken: ", zap.String("signedToken", signedToken))

		ctx.JSON(200, models.TokenExternal{
			Token:     signedToken,
			ExpiredAt: tokenExpiredAt.Unix(),
			// TestToken: token,
		})

	}
}

// 需要靠這個才能拿到正確的state 不要想自己用自己的方法產一組
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var store cookie.Store

func Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}
