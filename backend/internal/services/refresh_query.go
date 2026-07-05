package services

import (
	"log"
	"sync"

	"tunisianet-scraper/internal/repository"
	"tunisianet-scraper/internal/scraper"
)

func RefreshQuery(search string) {

	log.Println("Refreshing query:", search)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		refreshTunisiaNet(search)
	}()

	go func() {
		defer wg.Done()
		refreshMyTek(search)
	}()

	go func() {
		defer wg.Done()
		refreshWiki(search)
	}()

	wg.Wait()

	log.Println("Finished refreshing:", search)
}
func refreshTunisiaNet(search string) {

	page := 1

	for {

		products, hasNext := scraper.SearchTunisianetPage(search, page)

		if len(products) > 0 {

			err := repository.SaveProducts(search, products)

			if err != nil {
				log.Println(err)
			}
		}

		if !hasNext {
			break
		}

		page++
	}
}
func refreshMyTek(search string) {

	page := 1

	for {

		products, hasNext := scraper.SearchMyTekPage(search, page)

		if len(products) > 0 {

			err := repository.SaveProducts(search, products)

			if err != nil {
				log.Println(err)
			}
		}

		if !hasNext {
			break
		}

		page++
	}
}

func refreshWiki(search string) {

	page := 1

	for {

		products, hasNext := scraper.SearchWikiPage(search, page)

		if len(products) > 0 {

			err := repository.SaveProducts(search, products)

			if err != nil {
				log.Println(err)
			}
		}

		if !hasNext {
			break
		}

		page++
	}
}
