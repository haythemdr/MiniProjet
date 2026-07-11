package services

import (
	"sort"
	"strings"

	"tunisianet-scraper/internal/elasticsearch"
)

func CorrectQuery(query string) string {

	products, err := elasticsearch.SearchProducts(query)
	if err != nil || len(products) == 0 {
		return query
	}

	freq := map[string]int{}

	for _, p := range products {

		name := NormalizeName(p.Name)

		seen := map[string]bool{}

		for _, word := range strings.Fields(name) {

			if seen[word] {
				continue
			}

			seen[word] = true
			freq[word]++
		}
	}

	queryWords := strings.Fields(NormalizeName(query))

	for _, q := range queryWords {
		freq[q] += 100
	}

	type pair struct {
		Word  string
		Count int
	}

	var words []pair

	for w, c := range freq {

		if len(w) < 2 {
			continue
		}

		words = append(words, pair{w, c})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	var corrected []string

	for _, p := range words {

		if len(corrected) >= len(queryWords)+1 {
			break
		}

		corrected = append(corrected, p.Word)
	}

	if len(corrected) == 0 {
		return query
	}

	return strings.Join(corrected, " ")
}
