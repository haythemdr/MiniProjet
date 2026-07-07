package models

type SearchResponse struct {
	Store       string    `json:"store"`
	Source      string    `json:"source"`
	LastUpdated string    `json:"lastUpdated"`
	Products    []Product `json:"products"`
}
