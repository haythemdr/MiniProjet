package services

import (
	"strings"
	"sync"

	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/repository"
	"tunisianet-scraper/internal/scraper"
)

func normalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.Join(strings.Fields(name), " ")
	return name
}

func filterDuplicates(
	products []models.Product,
	seen map[string]bool,
) []models.Product {

	var filtered []models.Product

	for _, product := range products {

		key := product.Store + "|" + normalizeName(product.Name)

		if seen[key] {
			continue
		}

		seen[key] = true
		filtered = append(filtered, product)
	}

	return filtered
}
func SearchProductsStream(
	search string,
	out chan<- []models.Product,
) {

	_ = repository.RecordSearch(search)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		page := 1
		liveData := false

		for {

			products, hasNext :=
				scraper.SearchTunisianetPage(search, page)

			if len(products) > 0 {

				liveData = true

				out <- products

				// Refresh cache
				_ = repository.SaveProducts(search, products)
			}

			if !hasNext {
				break
			}

			page++
		}

		// Scraper failed -> use cache
		if !liveData {

			cache, err := repository.SearchByStoreAndQuery(
				"TunisiaNet",
				search,
			)

			if err == nil && len(cache) > 0 {

				out <- cache

				println("📦 TunisiaNet cache used")
			}
		}
	}()

	go func() {
		defer wg.Done()

		page := 1
		liveData := false

		for {

			products, hasNext :=
				scraper.SearchMyTekPage(search, page)

			if len(products) > 0 {

				liveData = true

				out <- products

				_ = repository.SaveProducts(search, products)
			}

			if !hasNext {
				break
			}

			page++
		}

		if !liveData {

			cache, err := repository.SearchByStoreAndQuery(
				"MyTek",
				search,
			)

			if err == nil && len(cache) > 0 {

				println("📦 MyTek cache used")

				out <- cache
			}
		}
	}()
	go func() {
		defer wg.Done()

		page := 1
		liveData := false

		for {

			products, hasNext :=
				scraper.SearchWikiPage(search, page)

			if len(products) > 0 {

				liveData = true

				out <- products

				_ = repository.SaveProducts(search, products)
			}

			if !hasNext {
				break
			}

			page++
		}

		if !liveData {

			cache, err := repository.SearchByStoreAndQuery(
				"Wiki",
				search,
			)

			if err == nil && len(cache) > 0 {

				println("📦 Wiki cache used")

				out <- cache
			}
		}
	}()

	wg.Wait()

	close(out)
}
