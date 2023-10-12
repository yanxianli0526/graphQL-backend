package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	config "graphql-go-template/envconfig"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

var (
	env config.EnvConfig
)

type OAuth struct {
	Config *oauth2.Config
}

type CallbackCredential struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func NewOAuth(AuthURL, TokenURL, redirectURL, ClientID, ClientSecret string) *OAuth {
	scopes := []string{"all"}

	endpoint := oauth2.Endpoint{
		AuthURL:  AuthURL,
		TokenURL: TokenURL,
	}

	return &OAuth{
		Config: &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     endpoint,
		},
	}
}

// GetToken will exchange token with code, and create client for further use
func (oAuth *OAuth) GetTokenWithCode(code string) (*oauth2.Token, error) {
	ctx := context.TODO()

	token, err := oAuth.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("GetAccessToken failed: %v", err)
	}

	return token, nil
}

// New the Client when user first time login
func (oAuth *OAuth) NewClient(token *oauth2.Token) *http.Client {
	return oAuth.Config.Client(context.TODO(), token)
}

func RevokeToken(clientId, clientSecret string, token *oauth2.Token) error {
	formData := url.Values{
		"token":           {token.AccessToken},
		"token_type_hint": {"access_token"},
	}
	payload := strings.NewReader(formData.Encode())
	basicAuthStr := fmt.Sprintf("%v:%v", clientId, clientSecret)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthStr))
	basicAuthToken := fmt.Sprintf("Basic %v", basicAuth)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, env.AuthEndpoint.RevokeURL, payload)
	if err != nil {
		return fmt.Errorf("[pkg/nis] http.NewRequest failed: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("authorization", basicAuthToken)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do failed: %v", err)
	}
	defer res.Body.Close()

	if code := res.StatusCode; code < 200 || code > 299 {
		return fmt.Errorf("http.PostForm failed: %v", err)
	}

	return nil
}
