package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"tunisianet-scraper/internal/models"
)

type megaPCProduct struct {
	Title       string  `json:"title"`
	Lien        string  `json:"lien"`
	Price       float64 `json:"price"`
	PrixEnPromo float64 `json:"prixEnPromo"`

	Categorie struct {
		Titre string `json:"titre"`
	} `json:"categorie"`

	Filscateg struct {
		Titre string `json:"titre"`
	} `json:"filscateg"`

	Gallerie struct {
		URLPhoto []string `json:"urlPhoto"`
	} `json:"gallerie"`
}

func SearchMegaPCPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	payload := map[string]interface{}{
		"recordByPage": 12,
		"pageNumber":   fmt.Sprintf("%d", page),
		"key":          search,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, false
	}

	req, err := http.NewRequest(
		"POST",
		"https://apiclt.gi-ga.tech/produit/withQueryForSearchByKey",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://megapc.tn")
	req.Header.Set("Referer", "https://megapc.tn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var result []megaPCProduct

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, false
	}

	for _, item := range result {

		category := url.PathEscape(item.Categorie.Titre)
		subCategory := url.PathEscape(item.Filscateg.Titre)

		productURL :=
			"https://megapc.tn/shop/product/" +
				category +
				"/" +
				subCategory +
				"/" +
				item.Lien

		if seen[productURL] {
			continue
		}

		seen[productURL] = true

		imageURL := ""
		if len(item.Gallerie.URLPhoto) > 0 {
			imageURL = "https://static.gi-ga.tech" + item.Gallerie.URLPhoto[0]
		}

		price := item.Price
		if item.PrixEnPromo > 0 {
			price = item.PrixEnPromo
		}
		products = append(products, models.Product{
			Name:  item.Title,
			Price: fmt.Sprintf("%.0f DT", price),
			Image: imageURL,
			URL:   productURL,
			Store: "MegaPC",
		})
	}

	if len(products) < 12 {
		return products, false
	}

	return products, true
}
