package scraper

import (
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tunisianet-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"

	"strconv"
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

func firstMetaContent(doc *goquery.Document, selector string) string {
	value, _ := doc.Find(selector).First().Attr("content")
	return strings.TrimSpace(html.UnescapeString(value))
}

func normalizeSearchText(value string) string {
	value = strings.ToLower(value)
	replacer := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a",
		"ç", "c",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"î", "i", "ï", "i",
		"ô", "o", "ö", "o",
		"ù", "u", "û", "u", "ü", "u",
		"-", " ", "_", " ", "/", " ", "\\", " ", ".", " ", ",", " ", ":", " ", ";", " ", "'", " ", "\"", " ", "(", " ", ")", " ", "[", " ", "]", " ", "+", " ", "&", " ",
	)
	return strings.Join(strings.Fields(replacer.Replace(value)), " ")
}

func productMatchesSearch(search string, values ...string) bool {
	search = normalizeSearchText(search)
	if search == "" {
		return true
	}

	combined := normalizeSearchText(strings.Join(values, " "))
	if combined == "" {
		return false
	}

	matchedToken := false
	for _, token := range strings.Fields(search) {
		if len(token) > 3 && strings.HasSuffix(token, "s") {
			token = strings.TrimSuffix(token, "s")
		}
		if len(token) < 2 {
			continue
		}

		matchedToken = true
		if !strings.Contains(combined, token) {
			return false
		}
	}

	return matchedToken
}
func SearchTunisianet(search string) []models.Product {

	var products []models.Product
	seen := make(map[string]bool)

	page := 1

	for {

		searchURL := "https://www.tunisianet.com.tn/recherche?controller=search&s=" +
			url.QueryEscape(search) +
			"&submit_search=&page=" + strconv.Itoa(page) +
			"&order=product.price.asc"

		resp, err := httpClient.Get(searchURL)
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

		doc.Find(".thumbnail-container").Each(func(i int, s *goquery.Selection) {

			name := strings.TrimSpace(
				s.Find(".product-title a").Text(),
			)

			if name == "" {
				return
			}

			productURL, _ := s.Find(".product-title a").Attr("href")

			if seen[productURL] {
				return
			}

			seen[productURL] = true
			newProducts++

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

			products = append(products, models.Product{
				Name:  name,
				Price: price,
				Image: imageURL,
				URL:   productURL,
				Store: "TunisiaNet",
			})
		})

		if newProducts == 0 {
			break
		}

		page++
	}

	return products
}

func GetProductDetails(productURL string) models.ProductDetails {
	var product models.ProductDetails

	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		return product
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := httpClient.Do(req)
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

	if strings.Contains(productURL, "mytek.tn") {
		return getMyTekProductDetails(doc)
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

func getMyTekProductDetails(doc *goquery.Document) models.ProductDetails {
	var product models.ProductDetails

	product.Name = strings.TrimSpace(doc.Find(".page-title .base").First().Text())
	if product.Name == "" {
		product.Name = firstMetaContent(doc, "meta[property='og:title']")
	}

	product.Price = firstMetaContent(doc, "meta[itemprop='price']")
	if product.Price == "" {
		product.Price = strings.TrimSpace(doc.Find("[data-price-type='finalPrice']").First().AttrOr("data-price-amount", ""))
	}
	if product.Price != "" {
		product.Price += " TND"
	}

	product.Image = firstMetaContent(doc, "meta[property='og:image']")
	if product.Image == "" {
		product.Image = normalizeMyTekImageURL(doc.Find(".gallery-placeholder img, .fotorama__img").First().AttrOr("src", ""))
	}

	product.Availability = strings.TrimSpace(doc.Find(".stock.available span, .stock-status-placeholder span").First().Text())
	if product.Availability == "" {
		availability, _ := doc.Find("link[itemprop='availability']").First().Attr("href")
		if strings.Contains(strings.ToLower(availability), "instock") {
			product.Availability = "En stock"
		}
	}

	product.Description = strings.TrimSpace(doc.Find(".product.attribute.overview [itemprop='description']").First().Text())
	if product.Description == "" {
		product.Description = strings.TrimSpace(doc.Find(".product.attribute.description .value").First().Text())
	}
	if product.Description == "" {
		product.Description = firstMetaContent(doc, "meta[property='og:description']")
	}

	return product
}
