package nis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type OrganizationInfo struct {
	Name                string   `json:"name"`
	Address             string   `json:"address"`
	Tel                 string   `json:"tel"`
	Fax                 string   `json:"fax"`
	Owner               string   `json:"ownerName"`
	Email               string   `json:"mail"`
	TaxIdNumber         string   `json:"idNumber"`
	RemittanceIdNumber  string   `json:"transferAccount"`
	EstablishmentNumber string   `json:"establishmentNumber"`
	Solution            string   `json:"solution"`
	Branch              []string `json:"branch"`
}

func (p *NIS) GetNisOrganizationInfo(orgId string, token *oauth2.Token) (*OrganizationInfo, error) {
	client := &http.Client{}
	endPoint := fmt.Sprintf("%s/organization/%s/facs/%s", p.baseURL, orgId, orgId)

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
	organization := &OrganizationInfo{}
	if err = json.NewDecoder(res.Body).Decode(organization); err != nil {
		return nil, fmt.Errorf("decode json failed: %w", err)
	}
	return organization, nil
}
