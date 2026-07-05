package services

import (
	"log"
	"sync"
	"time"

	"tunisianet-scraper/internal/repository"
)

func StartCacheWarmer() {

	ticker := time.NewTicker(time.Hour)
	//ticker := time.NewTicker(30 * time.Second)

	go func() {

		for range ticker.C {

			log.Println("🔥 Cache warmer started")

			queries, err := repository.GetQueriesToRefresh(20)

			if err != nil {
				log.Println(err)
				continue
			}

			var wg sync.WaitGroup

			// Limit to 3 queries running simultaneously
			semaphore := make(chan struct{}, 3)

			for _, query := range queries {

				wg.Add(1)

				go func(q string) {

					defer wg.Done()

					semaphore <- struct{}{}
					defer func() {
						<-semaphore
					}()

					log.Println("Refreshing:", q)

					RefreshQuery(q)

					_ = repository.UpdateRefreshTime(q)

				}(query)
			}

			wg.Wait()

			log.Println("✅ Cache warmer finished")
		}

	}()
}
