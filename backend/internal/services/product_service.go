package services

import (
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/scraper"
)

func SearchProducts(search string) []models.Product {

	var products []models.Product

	products = append(products, scraper.SearchTunisianet(search)...)
	products = append(products, scraper.SearchMyTek(search)...)
	products = append(products, scraper.SearchWiki(search)...)
	products = append(products, scraper.SearchSpaceNet(search)...)

	return products
}
func GetProductDetails(url string) models.ProductDetails {
	return scraper.GetProductDetails(url)
}
