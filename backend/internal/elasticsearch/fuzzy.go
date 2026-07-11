package elasticsearch

import (
	"bytes"
	"encoding/json"
)

type fuzzyResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Name string `json:"name"`
			} `json:"_source"`

			Score float64 `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

func CorrectQuery(query string) (string, error) {

	body := map[string]interface{}{
		"size": 1,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query":     query,
					"fuzziness": "AUTO",
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
		return query, err
	}

	defer res.Body.Close()

	var response fuzzyResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return query, err
	}

	if len(response.Hits.Hits) == 0 {
		return query, nil
	}

	return response.Hits.Hits[0].Source.Name, nil
}
