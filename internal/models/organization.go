package models

import "time"

type Organization struct {
	Id               string
	Name             string
	Description      string
	OrganizationType string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TenderList       []Tender
}
