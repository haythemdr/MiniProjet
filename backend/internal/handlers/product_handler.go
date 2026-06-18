package handlers

import (
	"net/http"

	"tunisianet-scraper/internal/services"

	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {

	search := c.QueryParam("search")

	products := services.SearchProducts(search)

	return c.JSON(http.StatusOK, products)
}

func GetProductDetails(c echo.Context) error {

	url := c.QueryParam("url")

	product := services.GetProductDetails(url)

	return c.JSON(http.StatusOK, product)
}
