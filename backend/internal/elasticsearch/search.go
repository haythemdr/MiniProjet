package elasticsearch

import (
	"bytes"
	"encoding/json"

	"tunisianet-scraper/internal/models"
)

type searchResponse struct {
	Hits struct {
		Hits []struct {
			Source models.Product `json:"_source"`
			Score  float64        `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

func SearchProducts(query string) ([]models.Product, error) {

	body := map[string]interface{}{
		"size":    10,
		"_source": []string{"name"},
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query":                query,
					"fuzziness":            "AUTO",
					"prefix_length":        1,
					"max_expansions":       50,
					"fuzzy_transpositions": true,
				},
			},
		},
	}

	data, _ := json.Marshal(body)

	res, err := Client.Search(
		Client.Search.WithIndex("products"),
		Client.Search.WithBody(bytes.NewReader(data)),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response searchResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var products []models.Product

	for _, hit := range response.Hits.Hits {
		products = append(products, hit.Source)
	}

	return products, nil
}
