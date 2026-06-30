package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

func SearchWiki(search string) []models.Product {

	var products []models.Product
	seen := make(map[string]bool)

	searchURL := "https://wiki.tn/?s=" +
		url.QueryEscape(search) +
		"&post_type=product&dgwt_wcas=1"

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return products
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://wiki.tn/")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := httpClient.Do(req)
	if err != nil {
		return products
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return products
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return products
	}

	doc.Find(".product-card--grid").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".product-card__title a").Text(),
		)

		productURL, _ := s.Find(".product-card__title a").Attr("href")
		productURL = strings.TrimSpace(productURL)
		img := s.Find(".product-card__image img")

		imageURL := strings.TrimSpace(img.AttrOr("data-src", ""))

		if imageURL == "" || strings.HasPrefix(imageURL, "data:image") {
			imageURL = strings.TrimSpace(img.AttrOr("src", ""))
		}

		if imageURL == "" || strings.HasPrefix(imageURL, "data:image") {
			srcset := strings.TrimSpace(img.AttrOr("data-srcset", ""))

			if srcset != "" {
				imageURL = strings.Split(srcset, ",")[0]
				imageURL = strings.Fields(imageURL)[0]
			}
		}

		if imageURL == "" || strings.HasPrefix(imageURL, "data:image") {
			srcset := strings.TrimSpace(img.AttrOr("srcset", ""))

			if srcset != "" {
				imageURL = strings.Split(srcset, ",")[0]
				imageURL = strings.Fields(imageURL)[0]
			}
		}

		price := strings.TrimSpace(
			s.Find(".woocommerce-Price-amount bdi").Text(),
		)

		if name == "" || productURL == "" || seen[productURL] {
			return
		}

		seen[productURL] = true
		fmt.Println(name)
		fmt.Println(imageURL)
		fmt.Println("--------------")
		products = append(products, models.Product{
			Name:  name,
			Price: price,
			Image: imageURL,
			URL:   productURL,
			Store: "Wiki",
		})
	})

	return products
}
