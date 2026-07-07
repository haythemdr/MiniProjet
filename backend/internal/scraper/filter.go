package scraper

import "strings"

func relevanceScore(productName, query string) int {
	productName = strings.ToLower(productName)
	query = strings.ToLower(query)

	productWords := strings.Fields(productName)
	queryWords := strings.Fields(query)

	score := 0

	for _, q := range queryWords {
		if len(q) <= 2 {
			continue
		}

		for _, p := range productWords {
			if p == q {
				score++
				break
			}
		}
	}

	return score
}
