package services

import (
	"strings"
	"tunisianet-scraper/internal/models"
)

func normalizeSearchText(value string) string {
	value = strings.ToLower(value)
	replacer := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a",
		"ç", "c",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"î", "i", "ï", "i",
		"ô", "o", "ö", "o",
		"ù", "u", "û", "u", "ü", "u",
		"-", " ", "_", " ", "/", " ", "\\", " ", ".", " ", ",", " ", ":", " ", ";", " ", "'", " ", "\"", " ", "(", " ", ")", " ", "[", " ", "]", " ", "+", " ", "&", " ",
	)
	return strings.Join(strings.Fields(replacer.Replace(value)), " ")
}

func GroupProducts(products []models.Product) []models.ProductGroup {

	groups := make(map[string][]models.Product)

	for _, product := range products {

		key := normalizeSearchText(product.Name)

		groups[key] = append(groups[key], product)
	}

	var result []models.ProductGroup

	for key, products := range groups {

		result = append(result, models.ProductGroup{
			Name:     key,
			Products: products,
		})
	}

	return result
}
