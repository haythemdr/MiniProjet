package services

import (
	"tunisianet-scraper/internal/models"
	"tunisianet-scraper/internal/repository"
)

func SaveAndIndexProducts(
	query string,
	products []models.Product,
) error {

	// Save into PostgreSQL
	if err := repository.SaveProducts(query, products); err != nil {
		return err
	}

	// Elasticsearch disabled for deployment

	return nil
}
