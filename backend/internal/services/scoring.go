package services

import (
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
}

func CalculateScore(query string, product models.Product) int {

	query = strings.ToLower(strings.TrimSpace(query))
	name := strings.ToLower(product.Name)

	score := 0

	// Exact title
	if name == query {
		score += 1000
	}

	// Full query contained
	if strings.Contains(name, query) {
		score += 500
	}

	// Multi-keyword score
	words := strings.Fields(query)

	for _, word := range words {

		if strings.Contains(name, word) {
			score += 120
		} else {
			score -= 10
		}
	}

	// Starts with query
	if strings.HasPrefix(name, query) {
		score += 80
	}

	// Penalize accessories
	for _, word := range accessoryWords {

		if strings.Contains(name, word) {
			score -= 150
		}
	}

	return score
}
