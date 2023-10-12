package nis

import (
	config "graphql-go-template/envconfig"

	"net/http"
)

type NIS struct {
	client  *http.Client
	baseURL string
}

func New(client *http.Client, env config.AuthEndpoint) *NIS {
	baseURL := env.BaseURL

	return &NIS{
		client:  client,
		baseURL: baseURL,
	}
}

// func (p *NIS) IsThreeLegged() bool {
// 	return p.baseURL == env.AuthEndpoint.BaseURL
// }
