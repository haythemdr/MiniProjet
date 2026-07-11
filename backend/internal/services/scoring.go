package services

import (
	"fmt"
	"strings"

	"tunisianet-scraper/internal/models"
)

var accessoryWords = []string{
	"coque",
	"case",
	"cover",
	"chargeur",
	"charger",
	"cable",
	"film",
	"housse",
	"protection",
	"adaptateur",
	"verre",
	"support",
	"pochette",
	"etui",
}

var brands = []string{
	"apple",
	"samsung",
	"asus",
	"hp",
	"lenovo",
	"acer",
	"msi",
	"dell",
	"huawei",
	"xiaomi",
	"honor",
	"oppo",
	"realme",
	"infinix",
}

func CalculateScore(query string, product models.Product) int {

	query = ProcessQuery(query)
	name := NormalizeName(product.Name)

	score := 0

	if name == query {
		score += 1200
	}

	if strings.HasPrefix(name, query) {
		score += 400
	}

	if strings.Contains(name, query) {
		score += 700
	}

	words := strings.Fields(query)

	for i, word := range words {

		if strings.Contains(name, word) {

			switch {

			case i == 0:
				score += 250

			case len(word) >= 5:
				score += 180

			default:
				score += 120
			}

			if strings.HasSuffix(word, "gb") ||
				strings.HasSuffix(word, "tb") {

				score += 200
			}

			if len(word) >= 5 &&
				strings.ContainsAny(word, "0123456789") {

				score += 400
			}

		} else {

			score -= 40
		}
	}

	for _, brand := range brands {

		if strings.Contains(query, brand) &&
			strings.Contains(name, brand) {

			score += 300
		}
	}

	for _, accessory := range accessoryWords {

		if strings.Contains(name, accessory) {

			score -= 700
		}
	}

	fmt.Println("--------------------------------------------")
	fmt.Println("Query      :", query)
	fmt.Println("Product    :", product.Name)
	fmt.Println("Normalized :", name)
	fmt.Println("Score      :", score)

	return score
}
