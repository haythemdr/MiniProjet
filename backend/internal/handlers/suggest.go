package handlers

import (
	"net/http"

	"tunisianet-scraper/internal/elasticsearch"

	"github.com/labstack/echo/v4"
)

func SuggestProducts(c echo.Context) error {

	query := c.QueryParam("q")

	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "query is required",
		})
	}

	suggestions, err := elasticsearch.SuggestProducts(query)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, suggestions)
}
