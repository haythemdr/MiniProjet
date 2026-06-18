package models

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	URL   string `json:"url"`
}