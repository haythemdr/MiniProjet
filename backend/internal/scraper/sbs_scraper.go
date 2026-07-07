package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"tunisianet-scraper/internal/models"
)

func SearchSBSPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	search = url.QueryEscape(search)

	pageURL := fmt.Sprintf(
		"https://www.sbsinformatique.com/recherche?controller=search&s=%s&page=%d",
		search,
		page,
	)

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept-Language", "fr,en;q=0.9")

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

		productURL, _ := s.Find("a.thumbnail.product-thumbnail").Attr("href")

		if productURL == "" {
			productURL, _ = s.Find("a").First().Attr("href")
		}

		if productURL != "" && !strings.HasPrefix(productURL, "http") {
			productURL = "https://www.sbsinformatique.com" + productURL
		}

		imageURL := ""

		if img := s.Find("img").First(); img.Length() > 0 {

			imageURL, _ = img.Attr("src")

			if imageURL == "" {
				imageURL, _ = img.Attr("data-src")
			}

			if imageURL == "" {
				imageURL, _ = img.Attr("data-lazy")
			}

			if imageURL != "" && strings.HasPrefix(imageURL, "/") {
				imageURL = "https://www.sbsinformatique.com" + imageURL
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
			Store: "SBS Informatique",
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
