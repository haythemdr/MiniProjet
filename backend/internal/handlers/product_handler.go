package handlers

import (
	"net/http"
	"strings"

	"tunisianet-scraper/internal/services"

	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {

	search := strings.TrimSpace(c.QueryParam("search"))
	if search == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Le paramètre search est obligatoire",
		})
	}

	products := services.SearchProducts(search)

	return c.JSON(http.StatusOK, products)
}

func GetProductDetails(c echo.Context) error {

	url := strings.TrimSpace(c.QueryParam("url"))
	if url == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Le paramètre url est obligatoire",
		})
	}

	product := services.GetProductDetails(url)

	return c.JSON(http.StatusOK, product)
}
