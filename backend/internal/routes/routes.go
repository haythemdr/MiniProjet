package routes

import (
	"tunisianet-scraper/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {

	e.GET("/products/stream", handlers.StreamProducts)

	e.GET("/product/details", handlers.GetProductDetails)
}
