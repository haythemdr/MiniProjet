package scraper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetMyTekCategories() ([]string, error) {

	req, err := http.NewRequest("GET", "https://www.mytek.tn", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var categories []string

	doc.Find("div.title_normal a").Each(func(i int, s *goquery.Selection) {

		link, exists := s.Attr("href")
		if !exists {
			return
		}

		link = strings.TrimSpace(link)

		if link == "" {
			return
		}

		if !strings.HasPrefix(link, "http") {
			link = "https://www.mytek.tn" + link
		}

		if !seen[link] {
			seen[link] = true
			categories = append(categories, link)
		}
	})

	return categories, nil
}
