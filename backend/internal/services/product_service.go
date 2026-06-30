package services

import (
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/scraper"
)

func SearchProducts(search string) []models.Product {

	tunisianetProducts := scraper.SearchTunisianet(search)

	mytekProducts := scraper.SearchMyTek(search)

	products := append(tunisianetProducts, mytekProducts...)

	return products
}

func GetProductDetails(url string) models.ProductDetails {
	return scraper.GetProductDetails(url)
}
