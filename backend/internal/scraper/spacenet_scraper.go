package scraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

func SearchSpaceNet(search string) []models.Product {

	var products []models.Product
	seen := make(map[string]bool)

	page := 1

	for {

		searchURL := "https://spacenet.tn/module/ambjolisearch/jolisearch?orderby=position&orderway=desc&search_query=" +
			url.QueryEscape(search)

		if page > 1 {
			searchURL += "&page=" + strconv.Itoa(page)
		}

		println("URL:", searchURL)

		req, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			break
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
		req.Header.Set("Referer", "https://spacenet.tn/")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := httpClient.Do(req)
		if err != nil {
			println("HTTP ERROR:", err.Error())
			break
		}

		println("Status:", resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			break
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			println("GOQUERY ERROR:", err.Error())
			break
		}

		println("Document loaded")
		resp.Body.Close()

		if err != nil {
			break
		}
		println("Status:", resp.StatusCode)
		println("Document loaded")

		newProducts := 0
		count := doc.Find(".product-miniature").Length()
		println("Page:", page, "Products found:", count)

		doc.Find(".product-miniature").Each(func(i int, s *goquery.Selection) {

			name := strings.TrimSpace(
				s.Find(".product_name a").Text(),
			)

			productURL, _ := s.Find(".product_name a").Attr("href")
			productURL = strings.TrimSpace(productURL)

			imageURL, _ := s.Find(".cover_image img").Attr("src")
			imageURL = strings.TrimSpace(imageURL)

			price := strings.TrimSpace(
				s.Find(".price").First().Text(),
			)

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
				Store: "SpaceNet",
			})
		})

		if newProducts == 0 {
			break
		}

		page++
	}

	return products
}
