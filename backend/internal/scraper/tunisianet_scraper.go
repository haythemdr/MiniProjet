package scraper

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func SearchTunisianet(search string) []models.Product {

	var products []models.Product

	searchURL := "https://www.tunisianet.com.tn/recherche?controller=search&s=" + url.QueryEscape(search)

	resp, err := httpClient.Get(searchURL)
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

	doc.Find(".thumbnail-container").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(
			s.Find(".product-title a").Text(),
		)

		productURL, _ := s.Find(".product-title a").Attr("href")

		imageURL, _ := s.Find(".product-thumbnail img").Attr("src")

		price := strings.TrimSpace(
			s.Find(".price").First().Text(),
		)

		if name != "" {

			products = append(products, models.Product{
				Name:  name,
				Price: price,
				Image: imageURL,
				URL:   productURL,
			})
		}
	})

	return products
}

func GetProductDetails(productURL string) models.ProductDetails {

	var product models.ProductDetails

	resp, err := httpClient.Get(productURL)
	if err != nil {
		return product
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return product
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return product
	}

	product.Name = strings.TrimSpace(
		doc.Find("h1").First().Text(),
	)

	product.Price = strings.TrimSpace(
		doc.Find(".current-price span").First().Text(),
	)

	product.Image, _ = doc.Find("img.center-block.img-responsive").First().Attr("src")

	product.Availability = strings.TrimSpace(
		doc.Find(".in-stock").First().Text(),
	)

	product.Description = strings.TrimSpace(
		doc.Find(".prodes").First().Text(),
	)

	return product
}
