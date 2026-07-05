package services

import (
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/scraper"
)

func GetProductDetails(url string) models.ProductDetails {
	return scraper.GetProductDetails(url)
}
