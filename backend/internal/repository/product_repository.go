package repository

import (
	"context"

	"tunisianet-scraper/internal/database"
	"tunisianet-scraper/internal/models"
)

func SearchByQuery(query string) ([]models.Product, error) {

	rows, err := database.DB.Query(
		context.Background(),
		`
		SELECT name, price, image, url, store
        FROM products
        WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
		`,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.Name,
			&product.Price,
			&product.Image,
			&product.URL,
			&product.Store,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func SaveProducts(query string, products []models.Product) error {

	for _, product := range products {

		_, err := database.DB.Exec(
			context.Background(),
			`
			INSERT INTO products
			(query, name, price, image, url, store)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT(url)
			DO NOTHING
			`,
			query,
			product.Name,
			product.Price,
			product.Image,
			product.URL,
			product.Store,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
