package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"tunisianet-scraper/internal/models"
)

func SearchScoopPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	search = url.QueryEscape(search)

	pageURL := fmt.Sprintf(
		"https://www.scoop.com.tn/search?controller=search&s=%s&page=%d",
		search,
		page,
	)

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr,en;q=0.9")
	req.Header.Set("Referer", "https://www.scoop.com.tn/")
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
			s.Find(".product-title h6").Text(),
		)

		if name == "" {
			name = strings.TrimSpace(
				s.Find(".product-title").Text(),
			)
		}

		price := strings.TrimSpace(
			s.Find(".price").First().Text(),
		)

		productURL, exists := s.Find("a.thumbnail.product-thumbnail").Attr("href")
		if !exists || productURL == "" {
			productURL, _ = s.Find("a").First().Attr("href")
		}

		if productURL != "" && !strings.HasPrefix(productURL, "http") {
			productURL = "https://www.scoop.com.tn" + productURL
		}

		imageURL := ""

		img := s.Find("img").First()

		if img.Length() > 0 {

			imageURL, _ = img.Attr("src")

			if imageURL == "" {
				imageURL, _ = img.Attr("data-src")
			}

			if imageURL == "" {
				imageURL, _ = img.Attr("data-lazy")
			}

			if imageURL == "" {
				imageURL, _ = img.Attr("data-original")
			}

			if imageURL != "" && strings.HasPrefix(imageURL, "/") {
				imageURL = "https://www.scoop.com.tn" + imageURL
			}
		}

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
			Store: "Scoop",
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
