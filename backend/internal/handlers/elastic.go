package handlers

import (
	"net/http"

	"tunisianet-scraper/internal/elasticsearch"

	"github.com/labstack/echo/v4"
)

func ElasticSearch(c echo.Context) error {

	query := c.QueryParam("q")

	products, err := elasticsearch.SearchProducts(query)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}
