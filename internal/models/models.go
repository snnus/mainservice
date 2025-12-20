package models

import "time"

type NewServicePointRequest struct {
	Name         string `json:"name"`
	ShortName    string `json:"shortName"`
	OfficeNumber string `json:"officeNumber"`
}

type ServicePoint struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ShortName    string    `json:"shortName"`
	OfficeNumber string    `json:"officeNumber"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Ticket struct {
	Ticket string `json:"ticket"`
}
