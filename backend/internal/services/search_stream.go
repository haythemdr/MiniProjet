package services

import (
	"log"
	"strings"
	"sync"
	"time"

	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/repository"
	"tunisianet-scraper/internal/scraper"
)

const (
	maxSearchPagesPerStore = 20
	maxSearchPageRetries   = 3
	searchRetryDelay       = 500 * time.Millisecond
)

func normalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.Join(strings.Fields(name), " ")
	return name
}

func productKey(product models.Product) string {
	identity := normalizeName(product.Name)
	if identity == "" {
		identity = strings.ToLower(strings.TrimSpace(product.URL))
	}

	return strings.ToLower(strings.TrimSpace(product.Store)) + "|" + identity
}

func filterDuplicates(products []models.Product, seen map[string]bool) []models.Product {
	var filtered []models.Product

	for _, product := range products {
		key := productKey(product)
		if seen[key] {
			continue
		}

		seen[key] = true
		filtered = append(filtered, product)
	}

	return filtered
}

func searchPageWithRetry(
	store string,
	search string,
	page int,
	searchPage func(string, int) ([]models.Product, bool),
) ([]models.Product, bool, bool) {
	var hasNext bool

	for attempt := 1; attempt <= maxSearchPageRetries; attempt++ {
		products, next := searchPage(search, page)
		if products != nil {
			return products, next, true
		}

		hasNext = next
		log.Printf(
			"%s -> page=%d attempt=%d failed",
			store,
			page,
			attempt,
		)

		if attempt < maxSearchPageRetries {
			time.Sleep(searchRetryDelay)
		}
	}

	return nil, hasNext, false
}

func streamStoreProducts(
	search string,
	store string,
	out chan<- models.SearchResponse,
	searchPage func(string, int) ([]models.Product, bool),
) {
	page := 1
	liveData := false
	requestFailed := false
	seen := make(map[string]bool)

	for page <= maxSearchPagesPerStore {
		products, hasNext, ok := searchPageWithRetry(store, search, page, searchPage)
		if !ok {
			requestFailed = true
			break
		}

		products = filterDuplicates(products, seen)

		log.Printf(
			"%s -> page=%d products=%d hasNext=%v",
			store,
			page,
			len(products),
			hasNext,
		)

		if len(products) == 0 {
			break
		}

		liveData = true

		out <- models.SearchResponse{
			Store:       store,
			Source:      "live",
			LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
			Products:    products,
		}

		if err := SaveAndIndexProducts(search, products); err != nil {
			log.Println("save/index products:", err)
		}

		if !hasNext {
			break
		}

		page++
	}

	if page > maxSearchPagesPerStore {
		log.Printf("%s -> stopped after %d pages", store, maxSearchPagesPerStore)
	}

	if liveData || !requestFailed {
		return
	}

	cache, err := repository.SearchByStoreAndQuery(store, search)
	if err != nil || len(cache) == 0 {
		return
	}

	cache = filterDuplicates(cache, make(map[string]bool))
	lastUpdate, _ := repository.GetLastUpdate(store)

	out <- models.SearchResponse{
		Store:       store,
		Source:      "cache",
		LastUpdated: lastUpdate.Format("2006-01-02 15:04:05"),
		Products:    cache,
	}

	log.Printf("%s cache used after live request failed", store)
}

func SearchProductsStream(search string, out chan<- models.SearchResponse) {
	_ = repository.RecordSearch(search)

	var wg sync.WaitGroup
	wg.Add(7)

	go func() {
		defer wg.Done()
		streamStoreProducts(search, "TunisiaNet", out, scraper.SearchTunisianetPage)
	}()

	go func() {
		defer wg.Done()
		streamStoreProducts(search, "MyTek", out, scraper.SearchMyTekPage)
	}()

	go func() {
		defer wg.Done()
		streamStoreProducts(search, "Wiki", out, scraper.SearchWikiPage)
	}()
	go func() {
		defer wg.Done()
		streamStoreProducts(search, "MegaPC", out, scraper.SearchMegaPCPage)
	}()
	go func() {
		defer wg.Done()
		streamStoreProducts(search, "SBS", out, scraper.SearchSBSPage)
	}()
	go func() {
		defer wg.Done()
		streamStoreProducts(search, "Scoop", out, scraper.SearchScoopPage)
	}()
	go func() {
		defer wg.Done()

		page := 1
		liveData := false

		for {

			products, hasNext := scraper.SearchSpaceNetPage(search, page)

			if len(products) > 0 {

				liveData = true

				out <- models.SearchResponse{
					Store:       "SpaceNet",
					Source:      "live",
					LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
					Products:    products,
				}

				if err := SaveAndIndexProducts(search, products); err != nil {
					log.Println("❌", err)
				}
			}

			if !hasNext {
				break
			}

			page++
		}

		if !liveData {

			cache, err := repository.SearchByStoreAndQuery(
				"SpaceNet",
				search,
			)

			if err == nil && len(cache) > 0 {

				lastUpdate, _ := repository.GetLastUpdate("SpaceNet")

				out <- models.SearchResponse{
					Store:       "SpaceNet",
					Source:      "cache",
					LastUpdated: lastUpdate.Format("2006-01-02 15:04:05"),
					Products:    cache,
				}

				log.Println("📦 SpaceNet cache used")
			}
		}
	}()

	wg.Wait()
	close(out)
}
