package services

import (
	"sort"
	"strconv"
	"strings"

	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/repository"
	"tunisianet-scraper/internal/scraper"
)

func extractPrice(price string) float64 {
	price = strings.TrimSpace(price)

	price = strings.ReplaceAll(price, "DT", "")
	price = strings.ReplaceAll(price, "TND", "")

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

	// Normalize search to avoid duplicates
	search = strings.TrimSpace(strings.ToLower(search))

	// ==========================
	// 1. Check PostgreSQL first
	// ==========================
	products, err := repository.SearchByQuery(search)

	if err == nil && len(products) > 0 {
		println("✅ Products loaded from PostgreSQL")

		sort.Slice(products, func(i, j int) bool {
			return extractPrice(products[i].Price) < extractPrice(products[j].Price)
		})

		return products
	}

	println("🔍 Products not found in DB, scraping websites...")

	// ==========================
	// 2. Scrape websites
	// ==========================
	products = append(products, scraper.SearchTunisianet(search)...)
	products = append(products, scraper.SearchMyTek(search)...)
	products = append(products, scraper.SearchWiki(search)...)
	// products = append(products, scraper.SearchSpaceNet(search)...)

	// ==========================
	// 3. Save into PostgreSQL
	// ==========================
	err = repository.SaveProducts(search, products)
	if err != nil {
		println("Error saving products:", err.Error())
	}

	// ==========================
	// 4. Sort by price
	// ==========================
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
