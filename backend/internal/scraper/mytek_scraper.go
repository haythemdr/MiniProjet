package scraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

func normalizeMyTekImageURL(value string) string {
	value = strings.TrimSpace(value)

	if value == "" || strings.HasPrefix(value, "data:image") {
		return ""
	}

	if strings.HasPrefix(value, "//") {
		return "https:" + value
	}

	if strings.HasPrefix(value, "/") {
		return "https://mk-media.mytek.tn/media/catalog/product" + value
	}

	return value
}

func SearchMyTek(search string) []models.Product {
	var products []models.Product
	seen := make(map[string]bool)

	page := 1

	for {

		searchURL := "https://www.mytek.tn/myteksearch/index/productsearch/?q=" +
			url.QueryEscape(search) +
			"&p=" + strconv.Itoa(page)

		req, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			break
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
		req.Header.Set("Referer", "https://www.mytek.tn/")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := httpClient.Do(req)
		if err != nil {
			break
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			break
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()

		if err != nil {
			break
		}

		newProducts := 0

		doc.Find("#seo-product-data > div").Each(func(i int, s *goquery.Selection) {

			name := strings.TrimSpace(s.AttrOr("data-name", ""))
			productURL := strings.TrimSpace(s.AttrOr("data-url", ""))
			imageURL := normalizeMyTekImageURL(s.AttrOr("data-image", ""))
			price := strings.TrimSpace(s.AttrOr("data-final-price", ""))

			if price == "" {
				price = strings.TrimSpace(s.AttrOr("data-price", ""))
			}

			if name == "" {
				name = strings.TrimSpace(s.Find("a").First().Text())
			}

			if productURL == "" {
				productURL, _ = s.Find("a").First().Attr("href")
				productURL = strings.TrimSpace(productURL)
			}

			if price == "" {
				price = strings.TrimSpace(
					s.Find("[itemprop='price']").First().Text(),
				)
			}

			if name == "" || productURL == "" || seen[productURL] {
				return
			}

			seen[productURL] = true
			newProducts++

			if price != "" &&
				!strings.Contains(price, "DT") &&
				!strings.Contains(price, "TND") {
				price += " DT"
			}

			products = append(products, models.Product{
				Name:  name,
				Price: price,
				Image: imageURL,
				URL:   productURL,
				Store: "MyTek",
			})
		})

		if newProducts == 0 {
			break
		}

		page++
	}

	return products
}
