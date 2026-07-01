package main

import (
	"tunisianet-scraper/internal/routes"

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
	e := echo.New()
	e.Use(middleware.CORS())

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
