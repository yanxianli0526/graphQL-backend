package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func GetOrgId(token *oauth2.Token) string {
	/* get orgUuid from metadata */

	metadata := token.Extra("metadata")
	orgUUIDInterface := metadata.(map[string]interface{})["orgId"]
	orgId := fmt.Sprintf("%v", orgUUIDInterface)
	return orgId
}

func GetRefreshTokenExpireAt(token *oauth2.Token) (time.Time, error) {
	refreshTokenExpiresInStr, ok := token.Extra("refresh_token_expires_in").(string)
	if !ok {
		log.Info("[nis] GetRefreshTokenExpireAt: time assertion failed")
		return time.Time{}, errors.New("[nis] GetRefreshTokenExpireAt: time assertion failed")
	}

	refreshTokenExpiresIn, err := strconv.Atoi(refreshTokenExpiresInStr)
	if err != nil {
		log.Info("[nis] GetRefreshTokenExpireAt: string convert to integer failed")
		return time.Time{}, err
	}

	refreshTokenExpiresAt := time.Now().Add(time.Duration(refreshTokenExpiresIn) * time.Second)

	return refreshTokenExpiresAt, nil
}

func MarshalToken(token *oauth2.Token) ([]byte, error) {
	data, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("MarshalToken failed: %v", err)
	}

	return data, nil
}
