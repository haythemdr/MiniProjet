package models

type ProductGroup struct {
	Name     string    `json:"name"`
	Products []Product `json:"products"`
}
