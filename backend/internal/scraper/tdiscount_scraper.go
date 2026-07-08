package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"tunisianet-scraper/internal/models"
)

func SearchTDiscountPage(search string, page int) ([]models.Product, bool) {

	var products []models.Product
	seen := make(map[string]bool)

	var searchURL string

	if page == 1 {
		searchURL = fmt.Sprintf(
			"https://tdiscount.tn/?s=%s&post_type=product",
			url.QueryEscape(search),
		)
	} else {
		searchURL = fmt.Sprintf(
			"https://tdiscount.tn/page/%d/?s=%s&post_type=product",
			page,
			url.QueryEscape(search),
		)
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

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

	doc.Find("li.product").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".woo-loop-product__title a").Text(),
		)

		productURL, _ := s.Find(".woo-loop-product__title a").Attr("href")

		price := strings.TrimSpace(
			s.Find(".price").First().Text(),
		)

		imageURL := ""

		img := s.Find(".mf-product-thumbnail img").First()

		if img.Length() > 0 {

			imageURL, _ = img.Attr("src")

			if imageURL == "" {
				imageURL, _ = img.Attr("data-src")
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
			Store: "TDiscount",
		})
	})

	hasNext := false

	doc.Find(".page-numbers a").Each(func(i int, s *goquery.Selection) {

		if strings.TrimSpace(s.Text()) == fmt.Sprintf("%d", page+1) {
			hasNext = true
		}
	})

	return products, hasNext
}
