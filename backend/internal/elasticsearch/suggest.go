package elasticsearch

import (
	"bytes"
	"encoding/json"
)

type suggestResponse struct {
	Suggest map[string][]struct {
		Options []struct {
			Text string `json:"text"`
		} `json:"options"`
	} `json:"suggest"`
}

func SuggestProducts(query string) ([]string, error) {

	body := map[string]interface{}{
		"suggest": map[string]interface{}{
			"product-suggest": map[string]interface{}{
				"prefix": query,
				"completion": map[string]interface{}{
					"field":           "name_completion",
					"size":            10,
					"skip_duplicates": true,
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

	var response suggestResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	suggestions := []string{}

	for _, entry := range response.Suggest["product-suggest"] {
		for _, option := range entry.Options {
			suggestions = append(suggestions, option.Text)
		}
	}

	return suggestions, nil
}
