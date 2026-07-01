package services

import (
	"sort"
	"strconv"
	"strings"
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/scraper"
)

func extractPrice(price string) float64 {
	price = strings.TrimSpace(price)

	price = strings.ReplaceAll(price, "DT", "")
	price = strings.ReplaceAll(price, "TND", "")

	// remove every type of space
	price = strings.ReplaceAll(price, " ", "")
	price = strings.ReplaceAll(price, "\u00A0", "")
	price = strings.ReplaceAll(price, "\u202F", "")

	price = strings.ReplaceAll(price, ",", ".")

	value, err := strconv.ParseFloat(price, 64)
	if err != nil {
		println("Parse error:", price)
		return 0
	}

	return value
}
func SearchProducts(search string) []models.Product {

	var products []models.Product

	products = append(products, scraper.SearchTunisianet(search)...)
	products = append(products, scraper.SearchMyTek(search)...)
	products = append(products, scraper.SearchWiki(search)...)
	products = append(products, scraper.SearchSpaceNet(search)...)
	for _, p := range products {
		println(
			p.Store,
			"|",
			p.Price,
			"|",
			strconv.FormatFloat(extractPrice(p.Price), 'f', 3, 64),
		)
	}

	sort.Slice(products, func(i, j int) bool {
		return extractPrice(products[i].Price) < extractPrice(products[j].Price)
	})

	println("========== AFTER SORT ==========")

	for i := 0; i < 20 && i < len(products); i++ {
		println(
			products[i].Store,
			"|",
			products[i].Price,
			"|",
			strconv.FormatFloat(extractPrice(products[i].Price), 'f', 3, 64),
		)
	}

	return products
}
func GetProductDetails(url string) models.ProductDetails {
	return scraper.GetProductDetails(url)
}
