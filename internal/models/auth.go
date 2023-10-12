package models

type TokenExternal struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
	// User      UserExternal `json:"user"`
}
