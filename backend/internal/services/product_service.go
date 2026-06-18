package services

import (
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/scraper"
)

func SearchProducts(search string) []models.Product {
	return scraper.SearchTunisianet(search)
}

func GetProductDetails(url string) models.ProductDetails {
	return scraper.GetProductDetails(url)
}
