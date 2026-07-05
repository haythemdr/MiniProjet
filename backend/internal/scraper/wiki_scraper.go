package scraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

func SearchWikiPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	searchURL := "https://wiki.tn/?s=" +
		url.QueryEscape(search) +
		"&post_type=product&dgwt_wcas=1" +
		"&paged=" + strconv.Itoa(page)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://wiki.tn/")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, false
	}

	newProducts := 0

	doc.Find(".product-card--grid").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".product-card__title a").Text(),
		)

		productURL, _ := s.Find(".product-card__title a").Attr("href")
		productURL = strings.TrimSpace(productURL)

		img := s.Find(".product-card__image img")

		imageURL := firstImageURL(
			img,
			searchURL,
			"data-large_image",
			"data-large-image",
			"data-o_src",
			"data-lazy-src",
			"data-src",
			"src",
		)

		imageURL = normalizeURLForHost(imageURL, searchURL)

		priceSelection := s.Find(".price ins .woocommerce-Price-amount bdi").First()
		if priceSelection.Length() == 0 {
			priceSelection = s.Find(".woocommerce-Price-amount bdi").First()
		}

		price := strings.TrimSpace(priceSelection.Text())

		if name == "" || productURL == "" || seen[productURL] {
			return
		}

		seen[productURL] = true
		newProducts++

		products = append(products, models.Product{
			Name:  name,
			Price: price,
			Image: imageURL,
			URL:   productURL,
			Store: "Wiki",
		})
	})

	if len(products) < 24 {
		return products, false
	}

	return products, true
}
