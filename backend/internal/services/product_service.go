package services

import (
	"tunisianet-scraper/internal/elasticsearch"
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

	// Index into Elasticsearch
	for _, product := range products {

		if err := elasticsearch.IndexProduct(product); err != nil {
			return err
		}
	}

	return nil
}
