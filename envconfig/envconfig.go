package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Env struct {
		HTTPPort int  `envconfig:"HTTP_PORT"    default:"4000"`
		Debug    bool `envconfig:"DEBUG"`
	}
	AuthEndpoint struct {
		BaseURL      string `envconfig:"BASE_URL" default:"http://localhost:8000/api/authenticate/public/preview1"`
		RevokeURL    string `envconfig:"REVOKE_URL" default:"http://localhost:8000/api/auth/revoke"`
		RedirectURL  string `envconfig:"REDIRECT_URL" default:"http://localhost:8010/callback"`
		AuthURL      string `envconfig:"AUTH_URL" default:"http://localhost:8080/auth/login"`
		TokenURL     string `envconfig:"TOKEN_URL" default:"http://localhost:8000/api/auth/token"`
		ClientID     string `envconfig:"CLIENT_ID" default:"xxxxxxxx"`
		ClientSecret string `envconfig:"CLIENT_SECRET" default:"yyyyyyyy"`
	}
	Database struct {
		DBSSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
		DBPort     string `envconfig:"DB_PORT" default:"5432"`
		DBHost     string `envconfig:"DB_HOST" default:"localhost"`
		DBname     string `envconfig:"DB_NAME" default:"test123"`
		DBPassword string `envconfig:"DB_PASSWORD" default:""`
		DBUser     string `envconfig:"DB_USER" default:"postgres"`
		// DBHost     string `envconfig:"DB_HOST" default:"35.185.152.147"`
		// DBPassword string `envconfig:"DB_PASSWORD" default:"jubo50868012"`
		// DBUser     string `envconfig:"DB_USER" default:"inventory-toll-service"`
		// DBname     string `envconfig:"DB_NAME" default:"inventory-toll-dev"`
	}
	PhotoService struct {
		PhotoServiceProtocol string `envconfig:"PHOTO_SERVICE_PROTOCOL" default:"http"`
		PhotoServiceHost     string `envconfig:"PHOTO_SERVICE_HOST" default:"localhost"`
		PhotoServicePort     string `envconfig:"PHOTO_SERVICE_PORT" default:"9999"`
		OriginalBucket       string `envconfig:"ORIGINAL_BUCKET" default:"test-origin-inventory-toll-image"`
		ProcsssedBucket      string `envconfig:"PROCESSED_BUCKET" default:"test-inventory-toll-image"`
	}
	EnvConfig struct {
		Env          Env
		AuthEndpoint AuthEndpoint
		Database     Database
		PhotoService PhotoService
	}
)

func Process(env *EnvConfig) (err error) {
	err = envconfig.Process("", env)
	return
}
