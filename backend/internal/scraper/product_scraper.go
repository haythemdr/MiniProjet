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

func firstAttr(selection *goquery.Selection, attrs ...string) string {
	for _, attr := range attrs {
		value, exists := selection.Attr(attr)
		value = strings.TrimSpace(value)

		if exists && value != "" && !strings.HasPrefix(value, "data:image") {
			return normalizeURL(value)
		}
	}

	return ""
}

func normalizeURL(value string) string {
	if strings.HasPrefix(value, "//") {
		return "https:" + value
	}

	if strings.HasPrefix(value, "/") {
		return "https://www.tunisianet.com.tn" + value
	}

	return value
}

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

		image := s.Find(".product-thumbnail img").First()
		imageURL := firstAttr(
			image,
			"data-full-size-image-url",
			"data-src",
			"data-original",
			"src",
		)

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

	productImage := doc.Find(".product-cover img").First()
	if productImage.Length() == 0 {
		productImage = doc.Find("img.center-block.img-responsive").First()
	}

	product.Image = firstAttr(
		productImage,
		"data-image-large-src",
		"data-full-size-image-url",
		"data-src",
		"src",
	)

	product.Availability = strings.TrimSpace(
		doc.Find(".in-stock").First().Text(),
	)

	product.Description = strings.TrimSpace(
		doc.Find(".prodes").First().Text(),
	)

	return product
}
