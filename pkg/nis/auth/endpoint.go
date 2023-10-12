package auth

import "golang.org/x/oauth2"

var Endpoint = oauth2.Endpoint{
	AuthURL:  "http://localhost:8080/auth/login",
	TokenURL: "http://localhost:8000/api/auth/token",
}
