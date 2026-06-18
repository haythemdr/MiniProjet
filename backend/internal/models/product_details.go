package models

type ProductDetails struct {
	Name         string `json:"name"`
	Price        string `json:"price"`
	Image        string `json:"image"`
	Availability string `json:"availability"`
	Description  string `json:"description"`
}