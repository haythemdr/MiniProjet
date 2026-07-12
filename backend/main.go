package main

import (
	"log"
	"tunisianet-scraper/internal/elasticsearch"
	"tunisianet-scraper/internal/routes"
	"tunisianet-scraper/internal/services"

	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"tunisianet-scraper/internal/database"
)

func main() {

	database.Connect(os.Getenv("DATABASE_URL"))

	if err := database.CreateTables(); err != nil {
		panic(err)
	}
	if os.Getenv("ELASTICSEARCH_URL") != "" {
		elasticsearch.Connect()
		elasticsearch.CreateIndex()
	} else {
		log.Println("⚠️ Elasticsearch disabled")
	}
	//elasticsearch.Connect()
	//elasticsearch.CreateIndex()
	services.LoadSynonyms("internal/config/synonyms.json")
	// Start background cache refresher
	services.StartCacheWarmer()

	e := echo.New()

	e.Use(middleware.CORS())

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
