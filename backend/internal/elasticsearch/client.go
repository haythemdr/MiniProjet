package elasticsearch

import (
	"log"
	"os"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

var Client *es8.Client

func Connect() {

	cfg := es8.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
	}

	var client *es8.Client
	var err error

	for i := 0; i < 15; i++ {

		client, err = es8.NewClient(cfg)
		if err == nil {

			res, err := client.Info()

			if err == nil {
				res.Body.Close()

				Client = client

				log.Println("✅ Connected to Elasticsearch")
				return
			}
		}

		log.Println("Waiting for Elasticsearch...")

		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Could not connect to Elasticsearch")
}
