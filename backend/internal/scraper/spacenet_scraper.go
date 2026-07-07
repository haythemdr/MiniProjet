package scraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

func SearchSpaceNetPage(search string, page int) ([]models.Product, bool) {

	products := []models.Product{}
	seen := make(map[string]bool)

	searchURL := "https://spacenet.tn/module/ambjolisearch/jolisearch?orderby=position&orderway=desc&search_query=" +
		url.QueryEscape(search)

	if page > 1 {
		searchURL += "&page=" + strconv.Itoa(page)
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://spacenet.tn/")
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

	doc.Find(".product-miniature").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".product_name a").Text(),
		)

		productURL, _ := s.Find(".product_name a").Attr("href")
		productURL = strings.TrimSpace(productURL)

		img := s.Find(".cover_image img")

		imageURL := firstImageURL(
			img,
			searchURL,
			"data-src",
			"data-lazy-src",
			"src",
		)

		imageURL = normalizeURLForHost(imageURL, searchURL)

		price := strings.TrimSpace(
			s.Find(".price").First().Text(),
		)

		if name == "" || productURL == "" || seen[productURL] {
			return
		}

		seen[productURL] = true

		products = append(products, models.Product{
			Name:  name,
			Price: price,
			Image: imageURL,
			URL:   productURL,
			Store: "SpaceNet",
		})
	})

	// SpaceNet displays 24 products per page.
	// If you later discover it displays 20 or another number,
	// simply replace 24 with the correct value.
	if len(products) < 32 {
		return products, false
	}

	return products, true
}
