package nis

import "time"

type Patient struct {
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Branch         string    `json:"branch"`
	Room           string    `json:"room"`
	Bed            string    `json:"bed"`
	Sex            string    `json:"sex"`
	Birthday       time.Time `json:"birthday"`
	CheckInDate    time.Time `json:"checkInDate"`
	PatientNumber  string    `json:"patientNumber"`
	RecordNumber   string    `json:"recordNumber"`
	Status         string    `json:"status"`
	IdNumber       string    `json:"idNumber"`
	PhotoUrl       string    `json:"photoUrl"`
	PhotoXPosition int       `json:"photoXPosition"`
	PhotoYPosition int       `json:"photoYPosition"`
	ProviderId     string    `json:"providerId"`
	ProviderOrgId  string    `json:"providerOrgId"`
	Numbering      string    `json:"numbering"`
	PeopleInCharge []string  `json:"peopleInCharge"`
}
