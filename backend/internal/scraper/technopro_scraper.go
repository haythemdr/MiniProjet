package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"tunisianet-scraper/internal/models"
)

func SearchTechnoproPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	searchURL := fmt.Sprintf(
		"https://www.technopro-online.com/module/iqitsearch/searchiqit?s=%s&page=%d",
		url.QueryEscape(search),
		page,
	)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr,en;q=0.9")
	req.Header.Set("Referer", "https://www.technopro-online.com/")
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

	doc.Find("article.product-miniature").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".product-title a").Text(),
		)

		productURL, _ := s.Find(".product-title a").Attr("href")

		imageURL := ""

		img := s.Find(".product-thumbnail img").First()

		if img.Length() > 0 {

			imageURL, _ = img.Attr("data-src")

			if imageURL == "" {
				imageURL, _ = img.Attr("src")
			}

			if imageURL != "" && strings.HasPrefix(imageURL, "/") {
				imageURL = "https://www.technopro-online.com" + imageURL
			}
		}

		price := strings.TrimSpace(
			s.Find(".product-price").First().Text(),
		)

		if name == "" || productURL == "" {
			return
		}

		if seen[productURL] {
			return
		}

		seen[productURL] = true

		products = append(products, models.Product{
			Name:  name,
			Price: price,
			Image: imageURL,
			URL:   productURL,
			Store: "Technopro",
		})
	})

	hasNext := false

	doc.Find(".pagination a").Each(func(i int, s *goquery.Selection) {

		href, _ := s.Attr("href")

		if strings.Contains(href, fmt.Sprintf("page=%d", page+1)) {
			hasNext = true
		}
	})

	return products, hasNext
}
