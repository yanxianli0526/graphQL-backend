package nis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type UsersAndPatients struct {
	Patients        []*Patient  `json:"patients"`
	PeopleInCharges []*UserInfo `json:"peopleInCharges"`
}

func (p *NIS) GetNisPeopleInChargesAndPatientsInfo(userId string, token *oauth2.Token) (*UsersAndPatients, error) {

	client := &http.Client{}
	endPoint := fmt.Sprintf("%s/user/%s/usersinfo", p.baseURL, userId)

	req, err := http.NewRequest(http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("[pkg/nisApi] http.NewRequest failed: %v", err)
	}

	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[pkg/nisApi] client.Do: %v", err)
	}
	if code := res.StatusCode; code < 200 || code > 299 {
		return nil, fmt.Errorf(endPoint+" api failed: %w", err)
	}
	user := &UsersAndPatients{}
	if err = json.NewDecoder(res.Body).Decode(user); err != nil {
		return nil, fmt.Errorf("decode json failed: %w", err)
	}
	return user, nil
}
