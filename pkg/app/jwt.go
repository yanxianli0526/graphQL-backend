package app

import (
	"errors"
	"log"
	"time"

	"graphql-go-template/internal/models"

	"github.com/dgrijalva/jwt-go"
)

const issuer = "Jubo Inc."
const expireHours = 360 // 15 days expired
const secret = "jubo"

type UserClaims struct {
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	jwt.StandardClaims
}

type UserClaimsForVitalLink struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(user *models.User) (string, time.Time) {

	expiresAt := time.Now().Add(time.Hour * time.Duration(expireHours))

	claims := UserClaims{
		UserId:      user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return signedToken, expiresAt
}

func ValidateJWT(authToken string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(authToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid signature method")
			return nil, errors.New("errMsg.InvalidToken")
		}

		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err)
		return nil, errors.New("errMsg.InvalidToken")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !(ok && token.Valid) {
		log.Println("invalid authentication token")
		return nil, errors.New("errMsg.InvalidToken")
	}

	return claims, nil
}
