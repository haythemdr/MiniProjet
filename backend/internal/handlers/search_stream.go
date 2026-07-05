package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/services"

	"github.com/labstack/echo/v4"
)

func StreamProducts(c echo.Context) error {

	search := c.QueryParam("search")

	if search == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "search is required",
		})
	}

	res := c.Response()

	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	channel := make(chan []models.Product)

	go services.SearchProductsStream(search, channel)

	for products := range channel {

		data, err := json.Marshal(products)
		if err != nil {
			continue
		}

		_, _ = res.Write([]byte("data: "))
		_, _ = res.Write(data)
		_, _ = res.Write([]byte("\n\n"))

		flusher.Flush()
	}

	_, _ = res.Write([]byte("event: done\n"))
	_, _ = res.Write([]byte("data: finished\n\n"))

	flusher.Flush()

	return nil
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
