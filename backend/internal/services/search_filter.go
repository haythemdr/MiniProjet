package services

import (
	"regexp"
	"strings"

	"tunisianet-scraper/internal/models"
)

var numberRegex = regexp.MustCompile(`\d+`)

func ShouldKeepProduct(query string, product models.Product) bool {

	query = strings.ToLower(query)
	name := strings.ToLower(product.Name)

	queryWords := strings.Fields(query)

	// Every important keyword must exist
	for _, word := range queryWords {

		if len(word) <= 2 {
			continue
		}

		if !strings.Contains(name, word) {
			return false
		}
	}

	// Compare numeric values (15,16,5070,512...)
	queryNumbers := numberRegex.FindAllString(query, -1)

	if len(queryNumbers) > 0 {

		productNumbers := numberRegex.FindAllString(name, -1)

		for _, q := range queryNumbers {

			found := false

			for _, p := range productNumbers {

				if q == p {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}
	}

	return true
}
